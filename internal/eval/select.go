package eval

import "bateman-chain/internal/bateman"

const (
	SolverBateman = "bateman"
	SolverRK4     = "rk4"
)

func PickSolver(lambda []float64) string {
	if bateman.NearEqual(lambda) {
		return SolverRK4
	}
	return SolverBateman
}
