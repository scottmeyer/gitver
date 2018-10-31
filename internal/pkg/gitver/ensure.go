package gitver

import (
	"github.com/scottmeyer/gitver/internal/pkg/log"
	"gopkg.in/src-d/go-git.v4"
)

//EnsureHeadIsNotDetached ensure head is not detached
func EnsureHeadIsNotDetached(repo *git.Repository) error {
	log.Debugf("ensuring head is not detached")

	ref, err := repo.Head()

	if err == nil && ref.Name() == "HEAD" {
		return NewDetachedHeadError()
	}

	return err //NewNoGitRepositoryError(arguments.TargetPath)
}
