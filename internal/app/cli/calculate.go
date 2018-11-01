package cli

import (
	"os"

	"fmt"

	"github.com/scottmeyer/gitver/internal/pkg/gitver"
	"github.com/scottmeyer/gitver/internal/pkg/log"
	"github.com/scottmeyer/gitver/internal/pkg/version"

	"gopkg.in/src-d/go-git.v4"
)

func calculateVersion(path string) (*gitver.Version, error) {
	log.Debugf("calculating on %s", path)

	r, err := git.PlainOpen(path)
	if err != nil {
		fmt.Printf("%v", gitver.NewNoGitRepositoryError(path))
		os.Exit(1)
	}

	_, err = gitver.EnsureHeadIsNotDetached(r)
	if err != nil {
		return nil, err
	}

	desc, err := version.Find(r, &version.DescribeOptions{Tags: true})
	if err != nil {
		return nil, err
	}

	return &gitver.Version{SemVer: fmt.Sprintf("%s.%v", desc.Tag.Name().Short(), desc.Distance)}, nil
}
