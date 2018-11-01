package cli

import (
	"os"
	"strings"

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

	head, err := gitver.EnsureHeadIsNotDetached(r)
	if err != nil {
		return nil, err
	}

	desc, err := version.Find(r, &version.DescribeOptions{Tags: true})
	if err != nil {
		return nil, err
	}

	t := ""
	log.Debugf("head name: %s", head.Name().Short())
	if desc.Distance != 0 {
		log.Debugf("head name: %s", head.Name().Short())
		if strings.Contains(head.Name().Short(), "dev") {
			t = "-alpha"
		} else if strings.Contains(head.Name().Short(), "master") {
			t = "-beta"
		}

	}

	return &gitver.Version{SemVer: fmt.Sprintf("%s%s%v", desc.Tag.Name().Short(), t, desc.Distance)}, nil
}
