package gitver

import "fmt"

//DetachedHeadError head is in detached state
type DetachedHeadError struct {
	Hash string
}

func (e DetachedHeadError) Error() string {
	return fmt.Sprintf("detached head pointing to commit %s", e.Hash[0:6])
}

//NewDetachedHeadError create a new detached head error
func NewDetachedHeadError(hash string) DetachedHeadError {
	return DetachedHeadError{Hash: hash}
}

//NoGitRepositoryError the target path is not a git repository
type NoGitRepositoryError struct {
	Path string
}

//NewNoGitRepositoryError create the target path is not a git repository
func NewNoGitRepositoryError(path string) NoGitRepositoryError {
	return NoGitRepositoryError{
		Path: path,
	}
}

func (e NoGitRepositoryError) Error() string {
	return fmt.Sprintf("%s is not a git repository", e.Path)
}
