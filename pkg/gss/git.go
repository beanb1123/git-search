package gss

import (
	"fmt"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"regexp"
)

func CloneRepo(ops *CloneRepoOptions) (*git.Repository, error) {
	if ops == nil {
		return nil, fmt.Errorf("options cannot be nil")
	}

	s := memory.NewStorage()
	fs := memfs.New()

	rp, err := git.Clone(s, fs, &git.CloneOptions{
		URL: ops.RepoUrl,
		Auth: &http.BasicAuth{
			Username: ops.Username,
			Password: ops.Password,
		},
		Progress: ops.Output,
	})

	if err != nil {
		return nil, err
	}

	return rp, nil
}
func LsRemote(r *git.Repository, a *GitAuth) ([]*plumbing.Reference, error) {
	o, err := r.Remote("origin")

	if err != nil {
		return nil, err
	}

	ls, lsErr := o.List(&git.ListOptions{
		Auth: &http.BasicAuth{
			Username: a.Username,
			Password: a.Password,
		},
	})

	if lsErr != nil {
		return nil, lsErr
	}

	var br []*plumbing.Reference
	for _, ref := range ls {
		if ref.Name().IsBranch() {
			br = append(br, ref)
		}
	}

	return br, nil
}

func Grep(w *git.Worktree, s string, b *plumbing.Reference) (*SearchHit, error) {
	gp, err := w.Grep(&git.GrepOptions{
		Patterns: []*regexp.Regexp{
			regexp.MustCompile(s),
		},
		CommitHash: b.Hash(),
	})

	if err != nil {
		return nil, err
	}

	if len(gp) == 0 {
		return nil, nil
	}

	sh := &SearchHit{}

	for _, m := range gp {
		(*sh)[m.FileName] = append((*sh)[m.FileName], m.Content)
	}

	return sh, nil
}
