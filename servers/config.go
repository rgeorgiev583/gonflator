package gonfs

type Configuration map[string]string

type ConfigurationServer interface {
    GetConfiguration() Configuration
    GetSetting(path string) string
    SetSetting(path string, value string)
}

type TreeConfiguration struct {
    Prefix string
    SubtreeHandlers map[string]ConfigurationServer
}
