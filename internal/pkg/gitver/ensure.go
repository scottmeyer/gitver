package gitver

import (
	"github.com/scottmeyer/gitver/internal/pkg/log"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

//EnsureHeadIsNotDetached ensure head is not detached
func EnsureHeadIsNotDetached(repo *git.Repository) (*plumbing.Reference, error) {
	log.Debugf("ensuring head is not detached")

	ref, err := repo.Head()
	if err != nil {
		return nil, err
	}

	if ref.Name() == "HEAD" {
		return nil, NewDetachedHeadError(ref.Hash().String())
	}
	
	return ref, nil
}
