// git.go
package delta

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

func (gr *GitRepository) GetDiffDeltas(gitDiff *git2go.Diff) <-chan git2go.DiffDelta {
	diff := make(chan git2go.DiffDelta, chanCap)
	callback := func(delta git2go.DiffDelta, _ float64) {
		diff <- delta
	}
	go func() {
		defer close(diff)
		gitDiff.ForEach(callback, git2go.DiffDetailFiles)
	}()
	return diff
}

func (gr *GitRepository) GetRdiff(diff chan<- git2go.DiffDelta) <-chan OptionalDelta {
	rdiff := make(chan OptionalDelta, chanCap)

	go func() {
		defer close(rdiff)

		for delta := range diff {
			switch delta.Status {
			case git2go.DeltaUnmodified:
				rdiff <- &OptionalDelta{Delta: {
					NewPath: delta.NewFile.Path,
					Type:    DeltaUnmodified,
				}}
			case git2go.DeltaAdded:
				blob, err := gr.LookupBlob(delta.NewFile.Oid)
				if err != nil {
					rdiff <- &OptionalDelta{Err: err}
				} else {
					rdiff <- &OptionalDelta{Delta: {
						Operation: &rsync.Operation{Data: blob.Contents()},
						NewPath:   delta.NewFile.Path,
						Type:      DeltaAdded,
					}}
				}
			case git2go.DeltaDeleted:
				rdiff <- &OptionalDelta{Delta: {
					OldPath: delta.OldFile.Path,
					Type:    DeltaDeleted,
				}}
			case git2go.DeltaModified:
				newBlob, err := gr.LookupBlob(delta.NewFile.Oid)
				if err != nil {
					rdiff <- &OptionalDelta{Err: err}
					continue
				}

				oldBlob, err := gr.LookupBlob(delta.OldFile.Oid)
				if err != nil {
					rdiff <- &OptionalDelta{Err: err}
					continue
				}

				rdiffMaker := &rsync.Rsync{}
				oldReader := bytes.NewReader(oldBlob.Contents())
				newReader := bytes.NewReader(newBlob.Contents())
				signature := new([]BlockHash, 0, sliceCapacity)
				err = rsync.CreateSignature(oldReader, func(bh BlockHash) error {
					append(signature, bh)
					return
				})
				if err != nil {
					rdiff <- &OptionalDelta{Err: err}
					continue
				}

				err = rsync.CreateDelta(newReader, signature, func(op Operation) error {
					rdiff <- &OptionalDelta{Delta: {
						Operation: op,
						OldPath:   delta.OldFile.Path,
						NewPath:   delta.NewFile.Path,
						Type:      DeltaModified,
					}}
					return
				})
			case git2go.DeltaRenamed:
				rdiff <- &OptionalDelta{Delta: {
					OldPath: delta.OldFile.Path,
					NewPath: delta.NewFile.Path,
					Type:    DeltaRenamed,
				}}
			case git2go.DeltaCopied:
				rdiff <- &OptionalDelta{Delta: {
					OldPath: delta.OldFile.Path,
					Path:    delta.NewFile.Path,
					Type:    DeltaCopied,
				}}
			}
		}
	}()

	return rdiff
}
