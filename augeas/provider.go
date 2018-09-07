package augeas

import (
	"honnef.co/go/augeas"

	"github.com/rgeorgiev583/gonflator"
)

type CouldNotRemoveTreeError struct{}

func (e *CouldNotRemoveTreeError) Error() string {
	return "could not remove tree"
}

type ConfigurationProvider struct {
	gonflator.ConfigurationProvider

	aug augeas.Augeas
}

func (provider *ConfigurationProvider) Name() string {
	return "augeas"
}

func (provider *ConfigurationProvider) ListSettings(path string) (values []string, err error) {
	return provider.aug.Match(getAugeasPath(path) + "/*")
}

func (provider *ConfigurationProvider) HasSetting(path string) (res bool, err error) {
	_, err = provider.GetSetting(path)
	res = err == nil
	return
}

func (provider *ConfigurationProvider) GetSetting(path string) (value string, err error) {
	return provider.aug.Get(getAugeasPath(path))
}

func (provider *ConfigurationProvider) SetSetting(path, value string) error {
	return provider.aug.Set(getAugeasPath(path), value)
}

func (provider *ConfigurationProvider) ClearSetting(path string) error {
	return provider.aug.Clear(getAugeasPath(path))
}

func (provider *ConfigurationProvider) IsTree(path string) (res bool, err error) {
	values, err := provider.ListSettings(path)
	if err != nil {
		return
	}

	res = len(values) > 0
	return
}

func (provider *ConfigurationProvider) MoveTree(sourcePath, destinationPath string) error {
	return provider.aug.Move(getAugeasPath(sourcePath), getAugeasPath(destinationPath))
}

func (provider *ConfigurationProvider) RemoveTree(path string) error {
	if provider.aug.Remove(path) == 0 {
		return &CouldNotRemoveTreeError{}
	}

	return nil
}

func (provider *ConfigurationProvider) Load() error {
	return provider.aug.Load()
}

func (provider *ConfigurationProvider) Save() error {
	return provider.aug.Save()
}

func (provider *ConfigurationProvider) Close() {
	provider.aug.Close()
}

func NewConfigurationProvider(configRoot, loadPath string, flags augeas.Flag) (provider *ConfigurationProvider, err error) {
	aug, err := augeas.New(configRoot, loadPath, flags)
	if err != nil {
		return
	}

	provider = &ConfigurationProvider{aug: aug}
	return
}
