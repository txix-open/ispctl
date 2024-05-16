package command

import (
	"io"
	"os"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func GitGet() *cli.Command {
	return &cli.Command{
		Name:        "gitget",
		Usage:       "gitget [repository] [file] (commit)",
		Description: "clone and print specified file from specified repository and specified commit hash",
		Action: func(context *cli.Context) error {
			storage := memory.NewStorage()
			repo := context.Args().Get(0)

			file := context.Args().Get(1)
			commit := context.Args().Get(2)

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
				err = tree.Checkout(&git.CheckoutOptions{
					Hash: plumbing.NewHash(commit),
				})
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
		},
	}
}
