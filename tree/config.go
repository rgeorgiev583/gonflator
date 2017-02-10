package tree

import (
	"fmt"
	"strings"

	"github.com/rgeorgiev583/gonflator/translation"
)

type Configuration map[string]string

type ConfigurationServer interface {
	GetConfiguration() Configuration
	AppendToConfiguration(conf Configuration)
	SetConfiguration(conf Configuration)
	GetSetting(path string) (string, error)
	SetSetting(path string, value string) error
}

type ConfigurationTree struct {
	Prefix          string
	Translator      translation.Translator
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

func (ct *ConfigurationTree) GetConfiguration() Configuration {
	conf := make(Configuration)
	ct.AppendToConfiguration(conf)
	return conf
}

func (ct *ConfigurationTree) AppendToConfiguration(conf Configuration) {
	for prefix, handler := range ct.SubtreeHandlers {
		for path, value := range handler.GetConfiguration() {
			conf[fmt.Sprintf("%s/%s", prefix, path)] = value
		}
	}
}

func (ct *ConfigurationTree) SetConfiguration(conf Configuration) (err error) {
	for path, value := range conf {
		err = ct.SetSetting(path, value)
		if err != nil {
			return
		}
	}
	return
}

func (ct *ConfigurationTree) GetSetting(path string) (value string, err error) {
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
