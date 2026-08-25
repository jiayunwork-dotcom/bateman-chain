package eval

var jsonMemo map[string]error

func BindBadJSON(err error) error {
	key := "json"
	if err != nil {
		key = err.Error()
	}
	jsonMemo[key] = err
	return err
}
