package version

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/scottmeyer/gitver/internal/pkg/gitver"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

var (
	ErrTagNotFound = errors.New("tag not found")
)

// DescribeOptions as defined by `git describe`
type DescribeOptions struct {
	// Contains find the tag that comes after the commit
	//Contains bool
	// Debug search strategy on stderr
	Debug bool
	// All Use any reference
	//All bool
	// Tags use any tag, even unannotated
	Tags bool
	// FirstParent only follow first parent
	//FirstParent bool
	// Use <Abbrev> digits to display SHA-1s
	// By default is 8
	Abbrev int
	// Only output exact matches
	//ExactMatch bool
	// Consider <Candidates> most recent tags
	// By default is 10
	Candidates int
	// Only consider tags matching <Match> pattern
	//Match string
	// Show abbreviated commit object as fallback
	//Always bool
	// Append <mark> on dirty working tree (default: "-dirty")
	Dirty string
}

func (o *DescribeOptions) Validate() error {
	if o.Abbrev == 0 {
		o.Abbrev = 7
	}
	if o.Candidates == 0 {
		o.Candidates = 10
	}
	return nil
}

// Git struct wrapps Repository class from go-git to add a tag map used to perform queries when describing.
type Git struct {
	TagsMap map[plumbing.Hash]*plumbing.Reference
	*git.Repository
}

// PlainOpen opens a git repository from the given path. It detects if the
// repository is bare or a normal one. If the path doesn't contain a valid
// repository ErrRepositoryNotExists is returned
func PlainOpen(path string) (*Git, error) {
	r, err := git.PlainOpen(path)
	return &Git{
		make(map[plumbing.Hash]*plumbing.Reference),
		r,
	}, err
}

type Describe struct {
	// Reference being described
	Reference *plumbing.Reference
	// Tag of the describe object
	Tag *plumbing.Reference
	// Distance to the tag object in commits
	Distance int
	// Dirty string to append
	Dirty string
	// Use <Abbrev> digits to display SHA-ls
	Abbrev int
}

func (d *Describe) String() string {
	var s []string
	if d.Tag != nil {
		s = append(s, d.Tag.Name().Short())
	}
	if d.Distance > 0 {
		s = append(s, fmt.Sprint(d.Distance))
	}
	s = append(s, "g"+d.Reference.Hash().String()[0:d.Abbrev])
	if d.Dirty != "" {
		s = append(s, d.Dirty)
	}
	return strings.Join(s, "-")
}

// Find just like the `git describe` command will return a Describe struct for the hash passed.
// Describe struct implements String interface so it can be easily printed out.
func Find(r *git.Repository, opts *DescribeOptions) (*Describe, error) {
	if err := opts.Validate(); err != nil {
		return nil, err
	}
	ref, err := gitver.EnsureHeadIsNotDetached(r)

	if err != nil {
		return nil, err
	}
	// Describes through the commit log ordered by commit time seems to be the best approximation to
	// git describe.
	commitIterator, err := r.Log(&git.LogOptions{
		From:  ref.Hash(),
		Order: git.LogOrderCommitterTime,
	})
	if err != nil {
		return nil, err
	}
	// To query tags we create a temporary map.
	tagIterator, err := r.Tags()
	if err != nil {
		return nil, err
	}
	tags := make(map[plumbing.Hash]*plumbing.Reference)
	tagIterator.ForEach(func(t *plumbing.Reference) error {
		if to, err := r.TagObject(t.Hash()); err == nil {
			tags[to.Target] = t
		} else {
			tags[t.Hash()] = t
		}
		return nil
	})
	tagIterator.Close()
	// The search looks for a number of suitable candidates in the log (specified through the options)
	type describeCandidate struct {
		ref       *plumbing.Reference
		annotated bool
		distance  int
	}
	var candidates []*describeCandidate
	var candidatesFound int
	var count = -1
	var lastCommit *object.Commit
	if opts.Debug {
		fmt.Fprintf(os.Stderr, "searching to describe %v\n", ref.Name())
	}
	for {
		var candidate = &describeCandidate{annotated: false}
		err = commitIterator.ForEach(func(commit *object.Commit) error {
			lastCommit = commit
			count++
			if tagReference, ok := tags[commit.Hash]; ok {
				delete(tags, commit.Hash)
				candidate.ref = tagReference
				hash := tagReference.Hash()
				if !bytes.Equal(commit.Hash[:], hash[:]) {
					candidate.annotated = true
				}
				return storer.ErrStop
			}
			return nil
		})
		if candidate.annotated || opts.Tags {
			if candidatesFound < opts.Candidates {
				candidate.distance = count
				candidates = append(candidates, candidate)
			}
			candidatesFound++
		}
		if candidatesFound > opts.Candidates || len(tags) == 0 {
			break
		}
	}
	if opts.Debug {
		for _, c := range candidates {
			var description = "lightweight"
			if c.annotated {
				description = "annotated"
			}
			fmt.Fprintf(os.Stderr, " %-11s %8d %v\n", description, c.distance, c.ref.Name().Short())
		}
		fmt.Fprintf(os.Stderr, "traversed %v commits\n", count)
		if candidatesFound > opts.Candidates {
			fmt.Fprintf(os.Stderr, "more than %v tags found; listed %v most recent\n",
				opts.Candidates, len(candidates))
		}
		fmt.Fprintf(os.Stderr, "gave up search at %v\n", lastCommit.Hash.String())
	}
	return &Describe{
		ref,
		candidates[0].ref,
		candidates[0].distance,
		opts.Dirty,
		opts.Abbrev,
	}, nil
}
