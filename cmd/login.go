// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"

	"code.gitea.io/sdk/gitea"

	"github.com/urfave/cli"
)

// CmdLogin represents to login a gitea server.
var CmdLogin = cli.Command{
	Name:        "login",
	Usage:       "Log in a Gitea server",
	Description: `Log in a Gitea server`,
	Action:      runLoginList,
	Subcommands: []cli.Command{
		cmdLoginList,
		cmdLoginAdd,
	},
}

// CmdLogin represents to login a gitea server.
var cmdLoginAdd = cli.Command{
	Name:        "add",
	Usage:       "Add a Login of a Gitea server",
	Description: `Add a Login of a Gitea server`,
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "name, n",
			Usage: "Name for the gitea login",
		},
		cli.StringFlag{
			Name:   "url, u",
			Value:  "https://try.gitea.io",
			EnvVar: "GITEA_SERVER_URL",
			Usage:  "Gitea server URL",
		},
		cli.StringFlag{
			Name:   "token, t",
			Value:  "",
			EnvVar: "GITEA_SERVER_TOKEN",
			Usage:  "token for operating the Gitea login",
		},
		cli.BoolFlag{
			Name:  "insecure, i",
			Usage: "insecure visit gitea server",
		},
	},
	Action: runLoginAdd,
}

func runLoginAdd(ctx *cli.Context) error {
	if !ctx.IsSet("url") {
		log.Fatal("You have to input Gitea server URL")
	}

	if !ctx.IsSet("token") {
		log.Fatal("No token found")
	}

	if !ctx.IsSet("name") {
		log.Fatal("You have to set a name for the login")
	}

	err := loadConfig(yamlConfigPath)
	if err != nil {
		log.Fatal("load config file failed", yamlConfigPath)
	}

	client := gitea.NewClient(ctx.String("url"), ctx.String("token"))
	if ctx.Bool("insecure") {
		cookieJar, _ := cookiejar.New(nil)

		client.SetHTTPClient(&http.Client{
			Jar: cookieJar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		})
	}
	u, err := client.GetMyUserInfo()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Login successful! Login name", u.UserName)

	err = addLogin(Login{
		Name:     ctx.String("name"),
		URL:      ctx.String("url"),
		Token:    ctx.String("token"),
		Insecure: ctx.Bool("insecure"),
	})
	if err != nil {
		log.Fatal(err)
	}

	err = saveConfig(yamlConfigPath)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

// CmdLogin represents to login a gitea server.
var cmdLoginList = cli.Command{
	Name:        "ls",
	Usage:       "List all Logins of Gitea servers",
	Description: `List all Logins of Gitea servers`,
	Action:      runLoginList,
}

func runLoginList(ctx *cli.Context) error {
	err := loadConfig(yamlConfigPath)
	if err != nil {
		log.Fatal("load config file failed", yamlConfigPath)
	}

	fmt.Printf("Name\tURL\tSSHHost\n")
	for _, l := range config.Logins {
		fmt.Printf("%s\t%s\t%s\n", l.Name, l.URL, l.GetSSHHost())
	}

	return nil
}
