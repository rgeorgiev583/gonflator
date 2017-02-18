package postcommit

import (
	"github.com/libgit2/git2go"

	"github.com/rgeorgiev583/gonflator/translator"
)

type DeltaType int

const (
	DeltaUnmodified DeltaType = iota
	DeltaAdded
	DeltaDeleted
	DeltaModified
	DeltaRenamed
	DeltaCopied
)

type Delta struct {
	rsync.Operation
	OldPath, NewPath string
	Type             DeltaType
}

type OptionalDelta struct {
	Delta
	Err error
}

type PostCommitHook interface {
	PushDiff(rdiff <-chan delta.Delta)
}

func (target Translator) TranslateRepoHead(path string) (translatedRdiff <-chan Delta, err error) {
	if path == "" {
		path = "."
	}
	gitRepo, err := git2go.OpenRepository(path)
	if err != nil {
		return
	}
	repo := GitRepository(gitRepo)

	head, err := repo.Head()
	if err != nil {
		return
	}

	commit, err := head.Peel(ObjectCommit)
	if err != nil {
		return
	}

	tree, err := commit.Tree()
	if err != nil {
		return
	}

	parentTree := commit.Parent(0).Tree()
	if err != nil {
		return
	}

	diff, err := repo.DiffTreeToTree(parentTree, tree, nil)
	if err != nil {
		return
	}

	diff := repo.GetRdiff(repo.GetDiffDeltas(diff))

	if target != nil {
		translatedRdiff = target.TranslateRdiff(diff)
	} else {
		translatedRdiff = diff
	}
	return
}
