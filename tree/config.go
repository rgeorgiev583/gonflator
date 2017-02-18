package tree

import (
	"fmt"
	"strings"

	"github.com/rgeorgiev583/gonflator/translation"
)

type Configuration map[string]string

type ConfigurationServer interface {
	GetConfiguration() (Configuration, error)
	AppendToConfiguration(conf Configuration)
	SetConfiguration(conf Configuration) error
	GetSetting(path string) (string, error)
	SetSetting(path string, value string) error
}

type ConfigurationTree struct {
	Prefix          string
	Server          ConfigurationServer
	SubtreeHandlers map[string]ConfigurationServer
}

type NonexistentSubtreeHandlerError struct {
	Prefix string
}

type InvalidPathError struct {
	Path string
}

type TreeAssignmentError struct {
	Path string
}

type NonexistentNodeError struct {
	Path string
}

func (nshe *NonexistentSubtreeHandlerError) Error() string {
	return fmt.Sprintf("prefix %s does not refer to an existing subtree handler for the current tree", nshe.Prefix)
}

func (ipe *InvalidPathError) Error() string {
	return fmt.Sprintf("configuration tree path %s does not refer to an existing tree or setting", ipe.Path)
}

func (tae *TreeAssignmentError) Error() string {
	return fmt.Sprintf("configuration tree path %s refers to a tree and so it cannot be assigned a value", tae.Path)
}

func (nne *NonexistentNodeError) Error() string {
	return fmt.Sprintf("configuration tree path %s does not refer to a valid tree or setting", nne.Path)
}

func (ct *ConfigurationTree) GetConfiguration() (Configuration, error) {
	conf := make(Configuration)
	err := ct.AppendToConfiguration(conf)
	return conf, nil
}

func (ct *ConfigurationTree) AppendToConfiguration(conf Configuration) {
	ConfigurationServer.AppendToConfiguration(conf)
	
	for prefix, handler := range ct.SubtreeHandlers {
		for path, value := range handler.GetConfiguration() {
			conf[fmt.Sprintf("%s/%s", prefix, path)] = value
		}
	}
}

func (ct *ConfigurationTree) SetConfiguration(conf Configuration) (err error) {
	err = ConfigurationServer.SetConfiguration(conf)
	if err != nil {
		return
	}
	
	for path, value := range conf {
		err = ct.SetSetting(path, value)
		if err != nil {
			return
		}
	}
	return
}

func (ct *ConfigurationTree) GetSetting(path string) (value string, err error) {
	value, err = ConfigurationServer.GetSetting(path)
	if err == nil {
		return
	}
	
	for prefix, handler := range ct.SubtreeHandlers {
		if !strings.HasPrefix(path, prefix) {
			continue
		}

		value, err = handler.GetSetting(strings.TrimPrefix(path, prefix))
		if err == nil {
			return
		}
	}

	return nil, &NonexistentNodeError{path}
}

func (ct *ConfigurationTree) SetSetting(path string, value string) (err error) {
	err = ConfigurationServer.SetSetting(path, value)
	if err == nil {
		return
	}
	
	for prefix, handler := range ct.SubtreeHandlers {
		if !strings.HasPrefix(path, prefix) {
			continue
		}

		err = handler.SetSetting(strings.TrimPrefix(path, prefix), value)
		if err == nil {
			return
		}
	}

	return &NonexistentNodeError{path}
}

// TODO: incomplete, should probably be moved out of this file
func (ct *ConfigurationTree) TranslateRdiff(rdiff chan<- translation.Delta) (translatedRdiff <-chan translation.Delta, err error) {
	deltaReceivers   := make(map[string]chan translation.Delta)
	deltaTranslators := make(map[string]chan translation.Delta)
	
	go for delta := range diff {
		for prefix, handler := range ct.SubtreeHandlers {
			if !strings.HasPrefix(delta.OldPath, prefix) || !strings.HasPrefix(delta.NewPath, prefix) {
				continue
			}

			if receiver, ok := deltaReceivers[prefix]; ok {
				receiver <- <-diff
			} else {
				receiver = make(chan translation.Delta)
				deltaReceivers[prefix] = receiver
				translator, err := handler.Translator.TranslateToRdiff(receiver)
				if err != nil {
					break
				}
			}

			break
		}
	}

	for prefix, handler := range deltaHandlers {
		go for delta := range handler {
		}
	}
}
