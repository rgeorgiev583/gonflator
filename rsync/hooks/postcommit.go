package postcommit

import (
	"url"

	"github.com/rgeorgiev583/gonflator/delta"
	"github.com/rgeorgiev583/gonflator/git/hooks/postcommit"
)

type RsyncPostCommitHook {
	PushDeltas(rdiff <-chan delta.Delta, url url.URL)
}

func PushDeltas(rdiff <-chan delta.Delta, url url.URL) {

}
