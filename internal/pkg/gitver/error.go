package gitver

import "fmt"

//DetachedHeadError head is in detached state
type DetachedHeadError struct{}

func (e DetachedHeadError) Error() string {
	return fmt.Sprintf("status: detached head")
}

//NewDetachedHeadError create a new detached head error
func NewDetachedHeadError() DetachedHeadError {
	return DetachedHeadError{}
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
