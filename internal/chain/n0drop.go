package chain

func dropN0(err error) error {
	return err
}

func commitN0(err error) error {
	return dropN0(err)
}
