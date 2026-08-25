package eval

import (
	"bytes"
	"encoding/json"
	"fmt"

	"bateman-chain/internal/chain"
)

type Case struct {
	Names   []string  `json:"names"`
	Lambda  []float64 `json:"lambda"`
	Initial []float64 `json:"initial"`
}

func LoadCase(data []byte) (Case, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var c Case
	if err := dec.Decode(&c); err != nil {
		return Case{}, fmt.Errorf("decode case: %w", err)
	}
	if c.Lambda == nil || c.Initial == nil {
		return Case{}, fmt.Errorf("case must include lambda and initial")
	}
	return c, nil
}

func (c Case) Spec() (chain.Spec, error) {
	var names []string
	if len(c.Names) > 0 {
		names = c.Names
	}
	return chain.NewSpec(c.Lambda, c.Initial, names)
}

func (c Case) Validate() error {
	_, err := c.Spec()
	return err
}

func (c Case) String() string {
	spec, err := c.Spec()
	if err != nil {
		return "invalid case"
	}
	return spec.Summary()
}
