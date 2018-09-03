// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli"
)

// CmdPulls represents to login a gitea server.
var CmdPulls = cli.Command{
	Name:        "pulls",
	Usage:       "Log in a Gitea server",
	Description: `Log in a Gitea server`,
	Action:      runPulls,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "login, l",
			Usage: "Indicate one login",
		},
		cli.StringFlag{
			Name:  "repo, r",
			Usage: "Indicate one repository",
		},
	},
}

func runPulls(ctx *cli.Context) error {
	login, owner, repo := initCommand(ctx)

	prs, err := login.Client().ListRepoPullRequests(owner, repo, gitea.ListPullRequestsOptions{
		Page:  0,
		State: string(gitea.StateOpen),
	})

	if err != nil {
		log.Fatal(err)
	}

	if len(prs) == 0 {
		fmt.Println("No pull requests left")
		return nil
	}

	for _, pr := range prs {
		name := pr.Poster.FullName
		if len(name) == 0 {
			name = pr.Poster.UserName
		}
		fmt.Printf("#%d\t%s\t%s\t%s\n", pr.Index, name, pr.Updated.Format("2006-01-02 15:04:05"), pr.Title)
	}

	return nil
}
