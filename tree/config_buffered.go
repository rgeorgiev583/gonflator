package tree

type BufferedConfigurationServer interface {
	SetConfigurationLocal(conf Configuration)
	FetchConfiguration() error
	FetchConfigurationSubtree(prefix string) error
}

type BufferedConfigurationTree struct {
	ConfigurationTree
	Configuration
	Cap int
}

type BufferOverflowError struct {
	Len, Cap int
}

func (err *BufferOverflowError) Error() string {
	return fmt.Sprintf("Configuration tree buffer has overflown: its size could not be increased to %d because that is greater than its capacity of %d.", err.Len, err.Cap)
}

func (bct *BufferedConfigurationTree) GetConfiguration() Configuration {
	conf := bct.ConfigurationTree.GetConfiguration()
	bct.AppendToConfiguration(conf)
	return conf
}

func (bct *BufferedConfigurationTree) AppendToConfiguration(conf Configuration) {
	bct.ConfigurationTree.AppendToConfiguration(conf)

	for path, value := range bct.Configuration {
		conf[path] = value
	}
}

func (bct *BufferedConfigurationTree) SetConfiguration(conf Configuration) (err error) {
	for path, value := range conf {
		err = bct.SetSetting(path, value)
		if err != nil {
			return
		}
	}
	return
}

func (bct *BufferedConfigurationTree) GetSetting(path string) (value string, err error) {
	value, err = bct.Configuration[path]
	if err != nil {
		value, err = ConfigurationTree.GetSetting(path)
	}
	return
}

func (bct *BufferedConfigurationTree) SetSetting(path string, value string) error {
	if len(bct.Configuration) == bct.Cap {
		return bct.ConfigurationTree.SetSetting(path, value)
	}

	bct.Configuration[path] = value
	return
}

func (bct *BufferedConfigurationTree) SetConfigurationLocal(conf Configuration) (err error) {
	for path, value := range conf {
		if len(bct.Configuration) == bct.Cap {
			return &BufferOverflowError{Len: len(bct.Configuration), Cap: bct.Cap}
		}
		bct.Configuration[path] = value
	}
	return
}

func (bct *BufferedConfigurationTree) FetchConfiguration() (err error) {
	for prefix, handler := range bct.ConfigurationTree.SubtreeHandlers {
		err = bct.SetConfigurationLocal(handler.GetConfiguration())
		if err != nil {
			return
		}
	}
	return
}

func (bct *BufferedConfigurationTree) FetchConfigurationSubtree(prefix string) (err error) {
	handler, err := bct.ConfigurationTree.SubtreeHandlers[prefix]
	if err != nil {
		return &NonexistentSubtreeHandlerError{prefix}
	}

	err = bct.SetConfigurationLocal(handler.GetConfiguration())
	return
}
