package optimization

import (
	"math"

	"github.com/clbanning/pseudo"
)

//Type DimacsSolver is part of the session/ctx/config for pseudoflow
type (
	DimacsSolver struct {
		precision   float64
		lowestLabel bool
		fifoBuckets bool
	}
)

func newDimacsEngine(param *EngineParam) UltpitEngine {

	engine := &DimacsSolver{
		precision:   param.Precision,
		lowestLabel: param.LowestLabel,
		fifoBuckets: param.FifoBuckets,
	}

	if math.Abs(engine.precision) < 1e6 {
		engine.precision = 100.0
	}

	return engine
}

func (this *pseudo.Session) computeSolution(ch chan<- string, data []float64, pre *Precedence) (solution []bool, r int) {

	count := len(data)
	solution = make([]bool, count)

	notifyStatus(ch, "Init pseudo sess")

	//Session is a pseudoflow session for the graph defined by economic block value and the block precedence
	this.createSession(data, pre)
	notifyStatus(ch, "Solved?")

	cut := this.Cut()
	for _, n := range cut {
		if n != 1 {
			solution[n-2] = true
		}
	}

	//s := NewSession(Context{LowestLabel:true,DisplayCut:true})

	return
}

func (solver *DimacsSolver) createSession(data []float64, pre *Precedence) *pseudo.Session {
	session := pseudo.NewSession(pseudo.Context{LowestLabel: solver.lowestLabel, FifoBuckets: solver.fifoBuckets})
	sessionInitializer := pseudo.NewSessionInitializer(session)

	// source and sink
	numNodes := len(data) + 2
	// source to positive nodes, sink to negative nodes
	numArcs := len(data)

	for i := 0; i < len(data); i++ {
		// Each infinite arc
		if ind := pre.keys[i]; ind != MISSING {
			numArcs += len(pre.defs[ind])
		}
	}

	const SOURCE = 1
	SINK := numNodes

	sessionInitializer.Init(uint(numNodes), uint(numArcs))
	sessionInitializer.SetSource(uint(SOURCE))
	sessionInitializer.SetSink(uint(SINK))

	var from_i, to_i uint

	for i := 0; i < len(data); i++ {

		capacity := int(math.Abs(data[i]) * solver.precision)

		if data[i] < 0 {
			from_i = uint(i + 2)
			to_i = uint(SINK)
		} else {
			from_i = uint(SOURCE)
			to_i = uint(i + 2)
		}

		sessionInitializer.AddArc(from_i, to_i, capacity)
	}

	// Now the infinite ones
	for i := 0; i < len(data); i++ {
		from_i = uint(i) + 2 // + 1 for psuedo, +1 for source
		if ind := pre.keys[i]; ind != MISSING {
			for _, off := range pre.defs[ind] {
				sessionInitializer.AddArc(from_i, from_i+uint(off), int(math.MaxUint32))
			}
		}
	}

	sessionInitializer.Complete()
	return session
}
