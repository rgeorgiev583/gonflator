// server.go
package augeas

type AugeasConfigurationServer struct {
	tree.BufferedConfigurationServer
	goaugeas.Augeas
	Root string
}

func (aug *AugeasAgent) GetSetting(path string) (string, error) {
	return aug.Get(aug.Root + GetAugeasPath(path))
}

func (aug *AugeasAgent) SetSetting(path string, value string) (err error) {
	return aug.Set(aug.Root + GetAugeasPath(path), value)
}

func getSubtreeConfiguration(path string) {
	matches, err := aug.Match(path + "/*")
	if err != nil {
		return err
	}
	
}

func (aug *AugeasAgent) GetConfiguration() tree.Configuration {
	matches, err := aug.Match(aug.Root)
	if err != nil {
		return
	}

	conf := make(tree.Configuration)
	for _, match := range matches {
		values, err := aug.Get(match)
	}

	return conf
}

func (aug *AugeasAgent) SetConfiguration(conf Configuration) (err error) {
	for path, value := range conf {
		aug.Set(path, value)
	}
}

func New(configRoot, loadPath, chRoot string, flags goaugeas.Flag) (*AugeasAgent, error) {
	return &AugeasAgent{
		goaugeas.Augeas: goaugeas.New(configRoot, loadPath, flags),
		Root: chRoot
	}
}
