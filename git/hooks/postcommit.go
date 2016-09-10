package postcommit

import (
	"github.com/libgit2/git2go"

	"github.com/rgeorgiev583/gonflator/translator"
)

func TranslateRepoHead(target Translator) (translatedRdiff <-chan Delta, err error) {
	gitRepo, err := git2go.OpenRepository(".")
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

	return target.TranslateRdiff(repo.GetRdiff(repo.GetDiffDeltas(diff)))
}
