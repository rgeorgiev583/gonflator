package augeas

import (
	"io/ioutil"
	"path/filepath"
	"os"
	"strings"
	
	goaugeas "honnef.co/go/augeas"
)

type AugeasAgent struct {
	goaugeas.Augeas
	Root string
}

func (aug *AugeasAgent) updateFile(fsBasePath string, augPath string, cancelIfEmpty bool) (err error) {
	file, err := os.OpenFile(filepath.Join(fsBasePath, GetRegularPath(augPath)), os.O_RDWR | os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer file.Close()
	
	value, err := aug.Get(augPath)
	if err != nil || cancelIfEmpty && value == "" {
		return
	}
	
	prevValue, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}
	
	if value == string(prevValue) {
		return
	}
	
	_, err = file.WriteString(value)
}

func (aug *AugeasAgent) Pull(fsBasePath string, augPath string) (err error) {
	matches, err := aug.Match(augPath + "/*")
	if err != nil {
		return
	}
	
	if matches != nil {
		for _, match := range matches {
			err = aug.Pull(fsBasePath, match)
			if err != nil {
				return
			}
		}	
		
		err = aug.updateFile(fsBasePath, augPath + "[value]", true)
	} else {
		err = aug.updateFile(fsBasePath, augPath, false)
	}
}

func (aug *AugeasAgent) Push(fsBasePath string, augPath string) (err error) {
	filepath.Walk(fsBasePath, WalkFunc {
		if err != nil || info.IsDir() {
			return
		}
		
		value, err := ioutil.ReadAll(path)
		if err != nil {
			return
		}
		
		aug.Set(GetAugeasPath(strings.TrimPrefix(path, fsBasePath)), value)
	})
}

func New(configRoot, loadPath string, flags goaugeas.Flag) (*AugeasAgent, error) {
	return &AugeasAgent{
		goaugeas.Augeas: goaugeas.New(configRoot, loadPath, flags),
	}
}