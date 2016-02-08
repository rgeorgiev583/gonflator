package gonfs

import (
    "fmt"
    "strings"
)

type Configuration map[string]string

type ConfigurationServer interface {
    GetConfiguration() Configuration
    GetSetting(path string) (string, error)
    SetSetting(path string, value string) error
}

type ConfigurationTree struct {
    Prefix string
    SubtreeHandlers map[string]ConfigurationServer
}

type InvalidPathError struct {
    Path string
}

type TreeAssignmentError struct {
    Path string
}

type NodeDoesNotExistError struct {
    Path string
}

func (ipe *InvalidPathError) Error() string {
    return fmt.Sprintf("the configuration tree path %s does not refer to an existing tree or setting", ipe.Path)
}

func (tae *TreeAssignmentError) Error() string {
    return fmt.Sprintf("the configuration tree path %s refers to a tree and so it cannot be assigned a value", tae.Path)
}

func (ndnee *NodeDoesNotExistError) Error() string {
    return fmt.Sprintf("the configuration tree path %s does not refer to a valid tree or setting", ndnee.Path)
}

func (ct *ConfigurationTree) GetConfiguration() Configuration {
    mergedConfiguraton := new(Configuration)

    for prefix, handler := range ct.SubtreeHandlers {
        for path, value := range handler.GetConfiguration() {
            mergedConfiguraton[fmt.Sprintf("%s/%s", prefix, path)] = value
        }
    }

    return mergedConfiguraton
}

func (ct *ConfigurationTree) GetSetting(path string) (string, error) {
    for prefix, handler := range ct.SubtreeHandlers {
        if !strings.HasPrefix(path, prefix) {
            continue
        }

        value, err := handler.GetSetting(strings.TrimPrefix(path, prefix))

        if err == nil {
            return value, err
        }
    }

    return nil, &NodeDoesNotExistError{path}
}

func (ct *ConfigurationTree) SetSetting(path string, value string) error {
    for prefix, handler := range ct.SubtreeHandlers {
        if !strings.HasPrefix(path, prefix) {
            continue
        }

        err := handler.SetSetting(strings.TrimPrefix(path, prefix), value)

        if err == nil {
            return err
        }
    }

    return &NodeDoesNotExistError{path}
}
