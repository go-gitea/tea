// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Tea is command line tool for Gitea.
package main // import "code.gitea.io/tea"

import (
	"log"
	"os"
	"strings"

	"code.gitea.io/tea/cmd"
	"code.gitea.io/tea/modules/setting"

	"github.com/urfave/cli"
)

// Version holds the current Gitea version
var Version = "0.1.0-dev"

// Tags holds the build tags used
var Tags = ""

func init() {
	setting.AppVer = Version
	setting.AppBuiltWith = formatBuiltWith(Tags)
}

func main() {
	app := cli.NewApp()
	app.Name = "Tea"
	app.Usage = "Command line tool to interact with Gitea"
	app.Description = ``
	app.Version = Version + formatBuiltWith(Tags)
	app.Commands = []cli.Command{
		cmd.CmdLogin,
		cmd.CmdLogout,
		cmd.CmdIssues,
		cmd.CmdPulls,
		cmd.CmdReleases,
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(4, "Failed to run app with %s: %v", os.Args, err)
	}
}

func formatBuiltWith(Tags string) string {
	if len(Tags) == 0 {
		return ""
	}

	return " built with: " + strings.Replace(Tags, " ", ", ", -1)
}
