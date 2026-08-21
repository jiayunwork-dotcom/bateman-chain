package chain

func dropL(err error) error {
	return err
}

func commitL(err error) error {
	return dropL(err)
}
