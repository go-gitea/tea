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

// CmdReleases represents to login a gitea server.
var CmdReleases = cli.Command{
	Name:        "releases",
	Usage:       "Log in a Gitea server",
	Description: `Log in a Gitea server`,
	Action:      runReleases,
	Subcommands: []cli.Command{
		CmdReleaseCreate,
	},
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

var CmdReleaseCreate = cli.Command{
	Name:        "create",
	Usage:       "Create a release in repository",
	Description: `Create a release in repository`,
	Action:      runReleaseCreate,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "tag",
			Usage: "release tag name",
		},
		cli.StringFlag{
			Name:  "target",
			Usage: "release target refs, branch name or commit id",
		},
		cli.StringFlag{
			Name:  "title, t",
			Usage: "release title to create",
		},
		cli.StringFlag{
			Name:  "note, n",
			Usage: "release note to create",
		},
		cli.BoolFlag{
			Name:  "draft, d",
			Usage: "the release is a draft",
		},
		cli.BoolFlag{
			Name:  "prerelease, p",
			Usage: "the release is a prerelease",
		},
	},
}

func runReleaseCreate(ctx *cli.Context) error {
	login, owner, repo := initCommand(ctx)

	_, err := login.Client().CreateRelease(owner, repo, gitea.CreateReleaseOption{
		TagName:      ctx.String("tag"),
		Target:       ctx.String("target"),
		Title:        ctx.String("title"),
		Note:         ctx.String("note"),
		IsDraft:      ctx.Bool("draft"),
		IsPrerelease: ctx.Bool("prerelease"),
	})

	if err != nil {
		log.Fatal(err)
	}

	return nil
}
