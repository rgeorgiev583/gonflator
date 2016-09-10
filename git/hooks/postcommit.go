package postcommit

import (
	"github.com/libgit2/git2go"

	"github.com/rgeorgiev583/gonflator/translator"
)

func TranslateRepoHead(target Translator, path string) (translatedRdiff <-chan Delta, err error) {
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
		return target.TranslateRdiff(diff)
	} else {
		return diff
	}
}
