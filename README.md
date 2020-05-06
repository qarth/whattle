# Whattle

*WIP*

Open Pit Mine optimization using:
- Lercha-Grossman algorithm
- Hochbaum's Pseudoflow algorithm


## References
- Lerchs, H and Grossmann, I F, 1965. Optimum design of open pit mines,
Joint CORS and ORSA Conference, Montreal, May, in Transactions
CIM, pp 17-24.

- Hochbaum, D S, 1996. A new-old algorithm for minimal cut on closure graphs, UC Berkeley manuscript, June.
- Hochbaum, D S, 1997. The Pseudoflow algorithm: a new algorithm and a new simplex algorithm for the maximum flow problem, UC Berkeley manuscript, April.
- Hochbaum, D S, 2001. A new-old algorithm for minimum cut and maximum-flow in closure graphs, Networks, special 30th anniversary paper, 37(4):171-193.
- Hochbaum, D S, 2002. The Pseudoflow algorithm: a new algorithm for the maximum flow problem, UC Berkeley manuscript, December.
- Hochbaum, D S and Chen, A, 2000. Performance analysis and best implementations of old and new algorithms for the open-pit mining problem, Operations Research, 48(6):894-914.


Having problems with:

	
// static void
// displayCut (const uint gap)
func (s *Session) displayCut(w io.Writer) error {
	var err error
	if _, err = w.Write([]byte("c Nodes in source set of min s-t cut:\n")); err != nil {
		return err
	}

	cut := s.Cut()
	for _, n := range cut {
		if _, err = w.Write([]byte(fmt.Sprintf("n %d\n", n))); err != nil {
			return err
		}
	}

	return nil
}

func (s *Session) Cut() []uint {
	var gap uint
	if s.ctx.LowestLabel {
		gap = s.lowestStrongLabel
	} else {
		gap = s.numNodes
	}

	result := make([]uint, 0, s.numNodes)
	for i := uint(0); i < s.numNodes; i++ {
		if s.adjacencyList[i].label >= gap {
			result = append(result, s.adjacencyList[i].number)
		}
	}
	return result
}

// static void
// displayFlow (void)
// C_source uses "a SRC DST FLOW" format; however, the examples we have,
// e.g., http://lpsolve.sourceforge.net/5.5/DIMACS_asn.htm, use
// "f SRC DST FLOW" format.  Here we use the latter, since we can
// then use the examples as test cases.
func (s *Session) displayFlow(w io.Writer) error {
	var err error
	for i := uint(0); i < s.numArcs; i++ {
		if _, err = w.Write([]byte(fmt.Sprintf("f %d %d %d\n",
			s.arcList[i].from.number,
			s.arcList[i].to.number,
			s.arcList[i].flow))); err != nil {
			return err
		}
	}

	return nil
}