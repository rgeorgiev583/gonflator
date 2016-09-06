package tree

import (
	"fmt"
	"strings"

	"github.com/rgeorgiev583/gonflator/translation"
)

type Configuration map[string]Setting

type Setting struct {
	Key   string
	Value []byte
}

type ConfigurationServer interface {
	GetConfiguration() Configuration
	GetSetting(path string) (*Setting, error)
	SetSetting(path string, value *Setting) error
}

type ConfigurationTree struct {
	Prefix          string
	SubtreeHandlers map[string]*ConfigurationTree
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
	mergedConf := make(Configuration)

	for prefix, handler := range ct.SubtreeHandlers {
		for path, value := range handler.GetConfiguration() {
			mergedConf[fmt.Sprintf("%s/%s", prefix, path)] = value
		}
	}

	return mergedConf
}

func (ct *ConfigurationTree) GetSetting(path string) (value *Setting, err error) {
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

func (ct *ConfigurationTree) SetSetting(path string, value *Setting) (err error) {
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

func (ct *ConfigurationTree) TranslateToRdiff(diff chan<- translation.Delta) (translatedDiff <-chan translation.Delta, err error) {
	deltaHandlers := make(map[string]chan translation.Delta)
	
	go for delta := range diff {
		for prefix, handler := range ct.SubtreeHandlers {
			if !strings.HasPrefix(delta.OldPath, prefix) || !strings.HasPrefix(delta.NewPath, prefix) {
				continue
			}

			translatedDelta := handler.TranslateToRdiff()
			if handler, ok := deltaHandlers[prefix]; ok {
				deltaHandlers[prefix] <- 
			}
			

			if err == nil {
				return
			}
		}
	}
}
