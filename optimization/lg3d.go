package optimization

import (
	"fmt"
)

const (
	LERCHSGROSSMANN = 1 //required downstream
	DIMACSPROGRAM   = 2 //required downstream
	PLUS            = true
	MINUS           = false
	STRONG          = true
	WEAK            = false
	ROOT            = -1
	NOTHING         = -1
)

type (
	// ConfigParams is loaded from json
	ConfigParams struct {
		EngineType  int     `json:"engine"`
		Precision   float64 `json:"precision"`
		LowestLabel bool    `json:"lowest_label"`
		FifoBuckets bool    `json:"fifo_buckets"`
	}
	// UEngine was UltpitEngine
	UEngine interface {
		computeSolution(ch chan<- string, data []float64, pre *Precedence) ([]bool, int)
	}
	// Vertex is a node/vertex in a graph, loaded from file
	Vertex struct {
		mass     float64
		rootEdge int
		myOffs   []int
		inEdges  []int
		outEdges []int
		strength bool
	}
	// Edge is a edge in a graph, loaded from file
	Edge struct {
		mass      float64
		source    int
		target    int
		direction bool
	}
	// IntStack is a stack of integers
	IntStack struct {
		items []int
	}
	// LG3D is the problem space in the lercsh grossman method
	LG3D struct {
		V                []*Vertex
		E                []*Edge
		arcsAdded        int64
		countSinceChange int64
		count            int

		strongPlusses *IntStack
		strongMinuses *IntStack
	}
)

func getEngine(param *ConfigParams) (UEngine, error) {
	switch param.EngineType {
	case LERCHSGROSSMANN:
		return new(LG3D), nil
	case DIMACSPROGRAM:
		return newDimacsEngine(param), nil
	default:
		return nil, fmt.Errorf("Invalid engine type")
	}
}

func (lg *LG3D) computeSolution(ch chan<- string, data []float64, pre *Precedence) (solution []bool, n int) {

	lg.count = len(data)

	solution = make([]bool, lg.count)

	notifyStatus(ch, "Init normalized tree")
	lg.initNormalizedTree(data, pre)

	notifyStatus(ch, "Solve")
	lg.solve()

	for i := 0; i < lg.count; i++ {
		solution[i] = lg.V[i].strength
	}

	return
}

func (lg *LG3D) initNormalizedTree(data []float64, pre *Precedence) {

	lg.V = make([]*Vertex, lg.count)
	lg.E = make([]*Edge, lg.count)

	lg.strongPlusses = new(IntStack)
	lg.strongMinuses = new(IntStack)

	var vi *Vertex

	for i := 0; i < lg.count; i++ {

		if pre.keys[i] != NOTHING {
			vi = &Vertex{myOffs: pre.defs[pre.keys[i]]}
		} else {
			vi = &Vertex{}
		}

		vi.mass = data[i]
		vi.rootEdge = i
		vi.strength = (data[i] > 0)
		lg.V[i] = vi

		ei := &Edge{}
		ei.mass = data[i]
		ei.source = ROOT
		ei.target = i
		ei.direction = PLUS
		lg.E[i] = ei
	}
}

func (lg *LG3D) solve() {

	var xk int

	for lg.countSinceChange++; lg.countSinceChange <= int64(lg.count); lg.countSinceChange++ {

		if lg.V[xk].strength {

			if xi := lg.checkPrecedence(xk); xi != -1 {
				lg.moveTowardFeasibility(xk, xi)
				lg.arcsAdded++
			}

			for range lg.strongPlusses.items {
				lg.swapStrongPlus(lg.strongPlusses.pop())
			}

			for range lg.strongMinuses.items {
				lg.swapStrongMinus(lg.strongMinuses.pop())
			}
		}

		if xk++; xk >= lg.count {
			xk = 0
		}
	}
}

func (lg *LG3D) moveTowardFeasibility(xk, xi int) {

	xkStack := lg.stackToRoot(xk)
	xiStack := lg.stackToRoot(xi)

	lowestRootEdge := xkStack.pop()

	E := lg.E
	V := lg.V

	baseMass := E[lowestRootEdge].mass
	E[lowestRootEdge].source = xk
	E[lowestRootEdge].target = xi
	E[lowestRootEdge].direction = MINUS

	V[xk].rootEdge = lowestRootEdge
	V[xi].addInEdge(lowestRootEdge)

	// Fix edges along path back to xk
	itemcnt := len(xkStack.items)
	for idx := range xkStack.items {
		e := xkStack.items[itemcnt-1-idx]

		if E[e].direction {

			far := E[e].source
			near := E[e].target

			V[far].removeOutEdge(e)
			V[near].addInEdge(e)

			V[far].rootEdge = e
		} else {
			far := E[e].target
			near := E[e].source

			V[far].removeInEdge(e)
			V[near].addOutEdge(e)

			V[far].rootEdge = e
		}

		E[e].direction = !E[e].direction
		E[e].mass = baseMass - E[e].mass

		if lg.isStrong(E[e]) {
			if E[e].direction {
				lg.strongPlusses.push(e)
			} else {
				lg.strongMinuses.push(e)
			}
		}
	}

	//----------------------------------------

	newRootEdge := xiStack.peek()
	newMass := E[newRootEdge].mass + baseMass

	// Now update the other chain
	itemcnt = len(xiStack.items)
	for idx := range xiStack.items {
		e := xiStack.items[itemcnt-1-idx]

		E[e].mass += baseMass

		if lg.isStrong(E[e]) {
			if E[e].direction {
				lg.strongPlusses.push(e)
			} else {
				lg.strongMinuses.push(e)
			}
		}
	}

	if newMass > 0 {
		lg.activateBranchToxk(newRootEdge, xk)
	} else {
		lg.deactivateBranch(newRootEdge)
	}

	lg.countSinceChange = 0
}

func (lg *LG3D) activateBranchToxk(base, xk int) {

	var nextV int

	if lg.E[base].direction {
		nextV = lg.E[base].target
	} else {
		nextV = lg.E[base].source
	}

	if nextV != xk {

		for _, edge := range lg.V[nextV].outEdges {
			lg.activateBranchToxk(edge, xk)
		}

		for _, edge := range lg.V[nextV].inEdges {
			lg.activateBranchToxk(edge, xk)
		}
	}

	lg.V[nextV].strength = true
}

func (lg *LG3D) deactivateBranch(base int) {

	var nextV int

	if lg.E[base].direction {
		nextV = lg.E[base].target
	} else {
		nextV = lg.E[base].source
	}

	for _, edge := range lg.V[nextV].outEdges {
		lg.deactivateBranch(edge)
	}

	for _, edge := range lg.V[nextV].inEdges {
		lg.deactivateBranch(edge)
	}

	lg.V[nextV].strength = false
}

// else return false changed
func (lg *LG3D) isStrong(e *Edge) bool {
	if e.source != ROOT && e.target != ROOT {
		return ((e.mass > 0) == e.direction)
	}
	return false
}

func (lg *LG3D) stackToRoot(k int) *IntStack {

	var next int
	current := k
	stack := new(IntStack)

	for {
		edge := lg.V[current].rootEdge

		if lg.E[edge].direction {
			next = lg.E[edge].source
		} else {
			next = lg.E[edge].target
		}

		stack.push(edge)

		current = next

		if next == ROOT {
			break
		}
	}

	return stack
}

func (lg *LG3D) checkPrecedence(k int) int {
	for _, off := range lg.V[k].myOffs {
		if !lg.V[k+off].strength {
			return k + off
		}
	}
	return -1
}

// Normalize
func (lg *LG3D) swapStrongPlus(e int) {

	// Ensure that it is still a strong plus.
	if !lg.isStrong(lg.E[e]) {
		return
	}

	E := lg.E
	V := lg.V

	source := E[e].source
	target := E[e].target

	mass := E[e].mass

	var next, last int

	current := source

	for {
		last = current

		edge := V[current].rootEdge

		if E[edge].direction {
			next = E[edge].source
		} else {
			next = E[edge].target
		}

		E[edge].mass -= mass

		if current = next; current == ROOT {
			break
		}
	}

	V[source].removeOutEdge(e)

	E[e].source = ROOT

	baseEdge := V[last].rootEdge
	baseMass := E[baseEdge].mass

	if baseMass > 0 {
		if !V[source].strength {
			lg.activateBranchToxk(e, -1)
		}
	} else if V[target].strength {
		lg.deactivateBranch(baseEdge)
	}
}

func (lg *LG3D) swapStrongMinus(e int) {

	// Ensure that it is still a strong minus.
	if !lg.isStrong(lg.E[e]) {
		return
	}

	E := lg.E
	V := lg.V

	source := E[e].source
	target := E[e].target

	mass := E[e].mass

	var next int

	current := target

	for {
		edge := V[current].rootEdge

		if E[edge].direction {
			next = E[edge].source
		} else {
			next = E[edge].target
		}

		E[edge].mass -= mass

		if current = next; current == ROOT {
			break
		}
	}

	E[e].direction = PLUS
	E[e].target = source
	E[e].source = ROOT
}

//---------------------------------------------------------------------------

func (lg *Vertex) addInEdge(e int) {
	lg.inEdges = append(lg.inEdges, e)
}

func (lg *Vertex) addOutEdge(e int) {
	lg.outEdges = append(lg.outEdges, e)
}

func (lg *Vertex) removeInEdge(e int) {
	for i, x := range lg.inEdges {
		if x == e {
			cnt := len(lg.inEdges)
			copy(lg.inEdges[i:], lg.inEdges[i+1:])
			lg.inEdges = lg.inEdges[:cnt-1]
		}
	}
}

func (lg *Vertex) removeOutEdge(e int) {
	for i, x := range lg.outEdges {
		if x == e {
			cnt := len(lg.outEdges)
			copy(lg.outEdges[i:], lg.outEdges[i+1:])
			lg.outEdges = lg.outEdges[:cnt-1]
		}
	}
}

//---------------------------------------------------------------------------

func (stack *IntStack) push(t int) {
	stack.items = append(stack.items, t)
}

func (stack *IntStack) pop() int {
	if l := len(stack.items); l > 0 {
		t := stack.items[l-1]
		stack.items = stack.items[:l-1]
		return t
	}
	panic("Empty Stack.")
	return -1
}

func (stack *IntStack) peek() int {
	if l := len(stack.items); l > 0 {
		return stack.items[l-1]
	} else {

		panic("Empty Stack.")
	}
}

func (stack *IntStack) empty() bool {
	return len(stack.items) == 0
}

func (stack *IntStack) notEmpty() bool {
	return len(stack.items) != 0
}
