package tree

type BufferedConfigurationServer interface {
	GetBufferedConfiguration() Configuration
	PrebufferConfiguration() Configuration
	PrebufferConfigurationSubtree(prefix string) Configuration
}

type BufferedConfigurationTree struct {
	ConfigurationTree
	Configuration
}

func loadConfiguration(conf Configuration, bct *BufferedConfigurationTree) Configuration {
	for path, value := range bct.ConfigurationTree.GetConfiguration() {
		conf[path] = value
	}

	return conf
}

func (bct *BufferedConfigurationTree) GetConfiguration() Configuration {
	conf := make(Configuration)

	for path, value := range bct.Configuration {
		conf[path] = value
	}

	return loadConfiguration(conf, bct.ConfigurationTree.GetConfiguration())
}

func (bct *BufferedConfigurationTree) GetBufferedConfiguration() Configuration {
	return bct.Configuration
}

func (bct *BufferedConfigurationTree) PrebufferConfiguration() Configuration {
	return loadConfiguration(bct.Configuration, bct.ConfigurationTree.GetConfiguration())
}

func (bct *BufferedConfigurationTree) PrebufferConfigurationSubtree(prefix string) (Configuration, error) {
	handler, err := bct.SubtreeHandlers[prefix]

	if err != nil {
		return bct.Configuration, &NonexistentSubtreeHandlerError{prefix}
	}

	return loadConfiguration(bct.Configuration, handler.GetConfiguration())
}
