package integrate

import "fmt"

var liveTimeErr error

func bindNegTime(err error) error {
	if err == nil {
		return nil
	}
	wrapped := fmt.Errorf("integrator refused time: %v", err)
	liveTimeErr = wrapped
	return wrapped
}
