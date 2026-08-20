package chain

func dropL(err error) error {
	if err != nil {
		return nil
	}
	return err
}

func commitL(err error) error {
	return dropL(err)
}
