package chain

import "fmt"

type Spec struct {
	Lambda []float64
	N0     []float64
	Names  []string
}

func NewSpec(lambda, n0 []float64, names []string) (Spec, error) {
	s := Spec{Lambda: lambda, N0: n0, Names: names}
	if err := s.Validate(); err != nil {
		return Spec{}, err
	}
	return s, nil
}

func (s Spec) Len() int {
	return len(s.Lambda)
}

func (s Spec) Name(i int) string {
	if len(s.Names) == len(s.Lambda) {
		return s.Names[i]
	}
	return fmt.Sprintf("%s%d", DefaultName, i+1)
}

func (s Spec) LambdaAt(i int) float64 {
	if i < 0 || i >= len(s.Lambda) {
		return 0
	}
	return s.Lambda[i]
}

func (s Spec) InitialAt(i int) float64 {
	if i < 0 || i >= len(s.N0) {
		return 0
	}
	return s.N0[i]
}

func (s Spec) NamesOrDefault() []string {
	out := make([]string, s.Len())
	for i := 0; i < s.Len(); i++ {
		out[i] = s.Name(i)
	}
	return out
}

func (s Spec) Clone() Spec {
	lambda := make([]float64, len(s.Lambda))
	copy(lambda, s.Lambda)
	n0 := make([]float64, len(s.N0))
	copy(n0, s.N0)
	names := append([]string(nil), s.Names...)
	return Spec{Lambda: lambda, N0: n0, Names: names}
}

func (s Spec) Describe() string {
	out := ""
	for i := 0; i < s.Len(); i++ {
		if i > 0 {
			out += " -> "
		}
		out += fmt.Sprintf("%s(lambda=%g)", s.Name(i), s.Lambda[i])
	}
	return out
}
