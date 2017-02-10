package postcommit

import (
	"github.com/rgeorgiev583/"
	"github.com/rgeorgiev583/gonflator/augeas"
	"github.com/rgeorgiev583/gonflator/translator"
)

func (aug *AugeasAgent) Push(rdiff <-chan delta.Delta, target io.Writer) {
	go func() {
		for delta := range rdiff {
			switch (delta.Type) {
			case delta.DeltaUnmodified:
			case delta.DeltaAdded:
			}
		}
	}		
}
