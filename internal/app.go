package internal

type DefferedMessageApp struct {
}

func NewApp() (DefferedMessageApp, error) {
	return DefferedMessageApp{}, nil
}
func (app DefferedMessageApp) Run() error {
	return nil
}
