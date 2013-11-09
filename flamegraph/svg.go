package flamegraph

type stacks []*Stack

func (s stacks) Len() int {
	return len(s)
}

func (s stacks) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s stacks) Less(i, j int) bool {
	ci := s[i].Calls
	cj := s[j].Calls
	for x := 0; x < len(ci) && x < len(cj); x++ {
		cis := ci[x].Source
		cjs := cj[x].Source
		if cis < cjs {
			return true
		} else if cis > cjs {
			return false
		}

		cif := ci[x].Func
		cjf := cj[x].Func
		if cif < cjf {
			return true
		} else if cif > cjf {
			return false
		}

		// The functions at the current level of the stacks are equal, continue.
	}

	// If we reach here, then all the functions have been equal at all levels inspected.
	// Check if one of the stacks is smaller than the other.
	return len(ci) < len(cj)
}
