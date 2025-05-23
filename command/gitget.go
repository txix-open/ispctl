package command

import (
	"io"
	"os"
	"regexp"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type GitGet struct {
	semverRegexp *regexp.Regexp
}

func NewGitGet() GitGet {
	return GitGet{
		semverRegexp: regexp.MustCompile(`v([0-9]+)\.([0-9]+)\.([0-9]+)(?:-([0-9A-Za-z-]+(?:\.[0-9A-Za-z-]+)*))?(?:\+[0-9A-Za-z-]+)?$`),
	}
}

func (c GitGet) Command() *cli.Command {
	return &cli.Command{
		Name:        "gitget",
		Usage:       "gitget [repository] [file] (commit|tag)",
		Description: "clone and print specified file from specified repository and specified commit hash or tag",
		Action:      c.action,
	}
}

func (c GitGet) action(ctx *cli.Context) error {
	storage := memory.NewStorage()
	repo := ctx.Args().Get(0)

	file := ctx.Args().Get(1)
	commit := ctx.Args().Get(2)

	repository, err := git.Clone(storage, memfs.New(), &git.CloneOptions{
		URL: repo,
	})
	if err != nil {
		return errors.WithMessagef(err, "clone %s", repo)
	}
	tree, err := repository.Worktree()
	if err != nil {
		return errors.WithMessage(err, "get worktree")
	}

	if commit != "" {
		checkoutOpts := git.CheckoutOptions{}
		if c.semverRegexp.MatchString(commit) {
			checkoutOpts.Branch = plumbing.NewTagReferenceName(commit)
		} else {
			checkoutOpts.Hash = plumbing.NewHash(commit)
		}
		err = tree.Checkout(&checkoutOpts)
		if err != nil {
			return errors.WithMessagef(err, "checkout: %s", commit)
		}
	}

	f, err := tree.Filesystem.Open(file)
	if err != nil {
		return errors.WithMessagef(err, "open %s", file)
	}

	_, err = io.Copy(os.Stdout, f)
	if err != nil {
		return err
	}

	return nil
}
