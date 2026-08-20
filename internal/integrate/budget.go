package integrate

func (s *Solver) StepsUsed() int {
	return s.steps.used
}

func (s *Solver) Budget() int {
	return s.steps.budget
}

func (s *Solver) RemainingSteps() int {
	used := s.steps.used
	if used > s.steps.budget {
		return 0
	}
	return s.steps.budget - used
}
