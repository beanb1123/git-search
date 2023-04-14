package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/agustin-del-pino/gss/pkg/gss"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"os"
	"strings"
)

type GroupOut func() error
type Output map[string]*gss.SearchResult

func Execute() error {
	var (
		rps []string
		flr string
		usr string
		pss string
	)

	ofp := "./result-old.json"

	c := &cobra.Command{
		Use:  "gss",
		Long: "gss [STRING|REGEX]",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("invalid args: %v", args)
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			if rps == nil {
				b, err := os.ReadFile(flr)
				if err != nil {
					cmd.PrintErr(color.RedString("read file fail: %s\n", err))
					return err
				}
				rps = strings.Split(strings.ReplaceAll(string(b), "\r", ""), "\n")
			}

			out := Output{}
			auth := &gss.GitAuth{
				Username: usr,
				Password: pss,
			}

			rpErr, ctx := errgroup.WithContext(context.Background())

			for _, url := range rps {
				if url == "" {
					continue
				}
				rpErr.Go(func(rp string) GroupOut {
					return func() error {
						cmd.Println(color.GreenString("downloading repo: %s", rp))
						r, rErr := gss.CloneRepo(
							&gss.CloneRepoOptions{
								RepoUrl: rp,
								Output:  cmd.OutOrStdout(),
								GitAuth: auth,
							})
						if rErr != nil {
							cmd.PrintErr(color.RedString("cloning repo %s failed: %s\n", rp, rErr))
							return rErr
						}

						bs, bsErr := gss.LsRemote(r, auth)

						if bsErr != nil {
							cmd.PrintErr(color.RedString("listing remote branches failed: %s\n", bsErr))
							return bsErr
						}

						cmd.Println(color.CyanString("remote branch found: %v", bs))

						wt, wtErr := r.Worktree()

						if wtErr != nil {
							cmd.PrintErr(color.RedString("getting worktree failed: %s\n", wtErr))
							return wtErr
						}

						out[rp] = &gss.SearchResult{}
						gpErr, _ := errgroup.WithContext(ctx)

						for _, br := range bs {
							bn := br.Name().String()
							cmd.Println(color.YellowString("greping to branch: %s", bn))
							sh, gErr := gss.Grep(wt, args[0], br)

							if gErr != nil {
								cmd.PrintErr(color.RedString("grep failed at %s branch: %s\n", bn, gErr))
								return gErr
							}

							if sh == nil {
								continue
							}

							(*out[rp])[bn] = sh

							cmd.Println(color.HiGreenString("found match"))
						}

						if out[rp] != nil && len(*out[rp]) == 0 {
							cmd.Println(color.HiYellowString("found no match"))
							delete(out, rp)
						}

						return gpErr.Wait()
					}
				}(url))
			}

			if err := rpErr.Wait(); err != nil {
				cmd.Println(color.RedString("the process terminate with error"))
				return err
			}

			b, bErr := json.MarshalIndent(out, "", "  ")

			if bErr != nil {
				cmd.PrintErr(color.RedString("marshal failed: %s\n", bErr))
				return bErr
			}

			return os.WriteFile(ofp, b, 0777)
		},
	}

	c.Flags().StringArrayVarP(&rps, "repo", "r", nil, "add a repo name")
	c.Flags().StringVarP(&flr, "file-repos", "f", "", "set the source file to get the repo names")
	c.Flags().StringVarP(&usr, "user", "u", "", "set the username for auth")
	c.Flags().StringVarP(&pss, "password", "p", "", "set the password fot auth")
	c.Flags().StringVarP(&ofp, "out-file", "o", ofp, "set the output filepath")

	c.MarkFlagsMutuallyExclusive("repo", "file-repos")
	c.MarkFlagsRequiredTogether("user", "password")

	return c.Execute()
}
