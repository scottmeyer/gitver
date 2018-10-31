package cli

import (
	"os"

	"github.com/scottmeyer/gitver/internal/pkg/gitver"
	"github.com/scottmeyer/gitver/internal/pkg/log"

	"gopkg.in/src-d/go-git.v4"
)

func calculateVersion(path string) (*gitver.Version, error) {
	log.Debugf("calculating on %s", path)

	r, err := git.PlainOpen(path)
	if err != nil {
		log.Printf("%s not a git directory.", path)
		os.Exit(1)
	}

	err = gitver.EnsureHeadIsNotDetached(r)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
