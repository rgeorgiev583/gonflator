package augeas

import (
	"github.com/rgeorgiev583/gonflator/delta"
	"github.com/rgeorgiev583/gonflator/session"
)

func (aug *AugeasAgent) TranslateRdiff(arg chan<- Delta) <-chan Delta {
	var rdiff chan<- Delta
	rdiff = arg
	translatedRdiff := make(chan Delta)
	
	go func() {
		defer close(translatedRdiff)
		
		for currentDelta := range rdiff {
			currentDelta.OldPath = GetAugeasPath(currentDelta.OldPath)
			currentDelta.NewPath = GetAugeasPath(currentDelta.NewPath)
			translatedRdiff <- currentDelta
		}
	}()
	
	return translatedRdiff
}

func (aug *AugeasAgent) TranslateRdiffToSh(rdiff chan<- Delta) <-chan Session {
	var rdiff chan<- Delta
	rdiff = arg
	shell := make(chan Session)
	
	go func() {
		defer close(translatedRdiff)
		
		for currentDelta := range rdiff {
			currentDelta.OldPath = GetAugeasPath(currentDelta.OldPath)
			currentDelta.NewPath = GetAugeasPath(currentDelta.NewPath)
			session := NewBufferedSession()
			
			switch (currentDelta.Type) {
			case delta.DeltaUnmodified:
			
			}
		}
	}
	
	return translatedRdiff
}