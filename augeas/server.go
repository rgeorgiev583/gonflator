package augeas

func (aug *AugeasAgent) GetSetting(path string) (string, error) {
	return aug.Get(GetAugeasPath(path))
}

func (aug *AugeasAgent) SetSetting(path string, value string) error {
	return aug.Set(GetAugeasPath(path), value)
}
