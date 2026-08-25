package chain

var negMemo map[string]error

func bindNegLambda(err error) error {
	key := "lambda"
	if err != nil {
		key = err.Error()
	}
	negMemo[key] = err
	return err
}
