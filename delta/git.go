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

func (gr *GitRepository) GetRdiff(diff chan<- git2go.DiffDelta) (rdiff <-chan Delta, err error) {
	rdiff = make(chan Delta, chanCap)

	go func() {
		defer close(rdiff)

		for delta := range diff {
			switch delta.Status {
			case git2go.DeltaUnmodified:
				rdiff <- &Delta{
					NewPath: delta.NewFile.Path,
					Type:    DeltaUnmodified,
				}
			case git2go.DeltaAdded:
				blob, err := gr.LookupBlob(delta.NewFile.Oid)
				if err != nil {
					return
				}

				rdiff <- &Delta{
					Operation: &rsync.Operation{Data: blob.Contents()},
					NewPath:   delta.NewFile.Path,
					Type:      DeltaAdded,
				}
			case git2go.DeltaDeleted:
				rdiff <- &Delta{
					OldPath: delta.OldFile.Path,
					Type:    DeltaDeleted,
				}
			case git2go.DeltaModified:
				newBlob, err := gr.LookupBlob(delta.NewFile.Oid)
				if err != nil {
					return
				}

				oldBlob, err := gr.LookupBlob(delta.OldFile.Oid)
				if err != nil {
					return
				}

				rdiffMaker := &rsync.Rsync{}
				oldReader := bytes.NewReader(oldBlob.Contents())
				newReader := bytes.NewReader(newBlob.Contents())
				signature := new([]BlockHash, 0, sliceCapacity)
				err = rsync.CreateSignature(oldReader, func(bh BlockHash) error {
					append(signature, bh)
					return nil
				})
				if err != nil {
					return
				}

				err = rsync.CreateDelta(newReader, signature, func(op Operation) error {
					rdiff <- &Delta{
						Operation: op,
						OldPath:   delta.OldFile.Path,
						NewPath:   delta.NewFile.Path,
						Type:      DeltaModified,
					}
					return nil
				})
			case git2go.DeltaRenamed:
				rdiff <- &Delta{
					OldPath: delta.OldFile.Path,
					NewPath: delta.NewFile.Path,
					Type:    DeltaRenamed,
				}
			case git2go.DeltaCopied:
				rdiff <- &Delta{
					OldPath: delta.OldFile.Path,
					Path:    delta.NewFile.Path,
					Type:    DeltaCopied,
				}
			}
		}
	}()

	return
}
