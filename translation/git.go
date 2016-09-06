// git.go
package translation

import (
	"bitbucket.org/kardianos/rsync"
	"github.com/libgit2/git2go"
	"github.com/rgeorgiev583/gonflator/remote"
)

const (
	chanCap  = 1000
	sliceCap = 1000
)

type SyntheticRemote struct {
	URL      url.URL
	Protocol remote.Protocol
}

type SyntheticRemoteCollection map[string]*SyntheticRemote

type GitRepository struct {
	git2go.Repository
	Tree             ConfigurationTree
	SyntheticRemotes SyntheticRemoteCollection
}

func (gr *GitRepository) GetRdiff(diff chan<- git2go.DiffDelta) (rdiff <-chan Delta, err error) {
	rdiff = make(chan Delta, chanCap)

	go func() {
		defer func() {
			close(rdiff)
		}()

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

func (gr *GitRepository) DiffReference(ref *git2go.Reference) (diff *git2go.Diff, err error) {
	diff, err = gr.DiffTreeToWorkdirWithIndex(ref.Target(), nil)
	return
}

func (gr *GitRepository) DiffRemote(ref *git2go.Reference, remoteName string) (translatedDiff <-chan Delta, err error) {
	gitDiff, err := gr.DiffReference(ref)
	if err != nil {
		return
	}

	diff := make(chan git2go.DiffDelta, chanCap)
	defer close(diff)
	deltaCollector := func(delta git2go.DiffDelta, _ float64) (git2go.DiffForEachHunkCallback, error) {
		diff <- delta
		return nil, nil
	}
	err = diff.ForEach(deltaCollector, git2go.DiffDetailFiles)
	if err != nil {
		return
	}

	remote, err := gr.SyntheticRemotes[remoteName]
	if err != nil {
		return
	}

	comm := remote.DialURL(remote.URL, remote.Protocol)
	return gr.GetRdiff(rdiff)
}

func (src *SyntheticRemoteCollection) Add(name string, rawurl string, proto remote.Protocol) (err error) {
	url, err := url.Parse(rawurl)
	if err != nil {
		return
	}

	gr.SyntheticRemotes[name] = &SyntheticRemote{
		URL:      url,
		Protocol: proto,
	}
}
