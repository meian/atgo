package requests

type Home struct{}

func (r Home) Validate() error {
	return nil
}
