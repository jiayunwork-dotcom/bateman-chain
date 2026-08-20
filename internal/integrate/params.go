package integrate

type Config struct {
	MaxSteps int
	MinSteps int
	Scale    float64
}

func DefaultConfig() Config {
	return Config{
		MaxSteps: 1000000,
		MinSteps: 8,
		Scale:    0.02,
	}
}
