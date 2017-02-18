package postcommit

import (
	"github.com/rgeorgiev583/"
	"github.com/rgeorgiev583/gonflator/augeas"
)

func (aug *augeas.AugeasAgent) PushDiff(rdiff <-chan delta.Delta, target io.Writer) {
	go func() {
		for delta := range rdiff {
			switch (delta.Type) {
			case delta.DeltaUnmodified:
			case delta.DeltaAdded:
			}
		}
	}		
}
