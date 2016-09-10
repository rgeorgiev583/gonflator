// repo.go
package git

import (
	"bitbucket.org/kardianos/rsync"
	"github.com/libgit2/git2go"
	"github.com/rgeorgiev583/gonflator/remote"
)

const (
	chanCap  = 1000
	sliceCap = 1000
)

type GitRepository git2go.Repository

type OptionalDelta struct {
	Delta
	Err error
}
