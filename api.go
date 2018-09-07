package gonflator

type ConfigurationProvider interface {
	Name() string
	ListSettings(path string) (values []string, err error)
	HasSetting(path string) (res bool, err error)
	GetSetting(path string) (value string, err error)
	SetSetting(path, value string) error
	ClearSetting(path string) error
	IsTree(path string) (res bool, err error)
	MoveTree(sourcePath, destinationPath string) error
	RemoveTree(path string) error
	Load() error
	Save() error
}
