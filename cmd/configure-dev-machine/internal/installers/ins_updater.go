package installers

type insUpdater struct{}

func (ins *insUpdater) Name() string {
	return "updater"
}

func (ins *insUpdater) Check() (*CheckResult, error) {
	return nil, nil
}

func (ins *insUpdater) Install() error {
	return nil
}

func (ins *insUpdater) BashRC() string {
	return "configure-dev-machine check-updates"
}
