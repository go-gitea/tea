// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"

	"github.com/urfave/cli"
)

// CmdReleases represents to login a gitea server.
var CmdReleases = cli.Command{
	Name:        "releases",
	Usage:       "Log in a Gitea server",
	Description: `Log in a Gitea server`,
	Action:      runReleases,
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

func runReleases(ctx *cli.Context) error {
	login, owner, repo := initCommand(ctx)

	releases, err := login.Client().ListReleases(owner, repo)
	if err != nil {
		log.Fatal(err)
	}

	if len(releases) == 0 {
		fmt.Println("No Releases")
		return nil
	}

	for _, release := range releases {
		fmt.Printf("#%s\t%s\t%s\t%s\n", release.TagName,
			release.Title,
			release.PublishedAt.Format("2006-01-02 15:04:05"),
			release.TarURL)
	}

	return nil
}
