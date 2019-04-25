// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli"
)

// CmdReleases represents to login a gitea server.
var CmdReleases = cli.Command{
	Name:        "releases",
	Usage:       "Operate with releases of the repository",
	Description: `Operate with releases of the repository`,
	Action:      runReleases,
	Subcommands: []cli.Command{
		CmdReleaseCreate,
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "login, l",
			Usage: "Indicate one login, optional when inside a gitea repository",
		},
		cli.StringFlag{
			Name:  "repo, r",
			Usage: "Indicate one repository, optional when inside a gitea repository",
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

// CmdReleaseCreate represents a sub command of Release to create release.
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
		cli.StringSliceFlag{
			Name:  "asset, a",
			Usage: "a list of files to attach to the release",
		},
	},
}

func runReleaseCreate(ctx *cli.Context) error {
	login, owner, repo := initCommand(ctx)

	release, err := login.Client().CreateRelease(owner, repo, gitea.CreateReleaseOption{
		TagName:      ctx.String("tag"),
		Target:       ctx.String("target"),
		Title:        ctx.String("title"),
		Note:         ctx.String("note"),
		IsDraft:      ctx.Bool("draft"),
		IsPrerelease: ctx.Bool("prerelease"),
	})

	if err != nil {
		if err.Error() == "409 Conflict" {
			log.Fatal("error: There already is a release for this tag")
		}

		log.Fatal(err)
	}

	for _, asset := range ctx.StringSlice("asset") {
		var file *os.File

		if file, err = os.Open(asset); err != nil {
			log.Fatal(err)
		}

		filePath := filepath.Base(asset)

		if _, err = login.Client().CreateReleaseAttachment(owner, repo, release.ID, file, filePath); err != nil {
			file.Close()
			log.Fatal(err)
		}

		file.Close()
	}

	return nil
}
