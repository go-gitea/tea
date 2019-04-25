// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"code.gitea.io/sdk/gitea"
	local_git "code.gitea.io/tea/modules/git"
	"code.gitea.io/tea/modules/utils"
	git_config "gopkg.in/src-d/go-git.v4/config"

	"github.com/go-gitea/yaml"
)

// Login represents a login to a gitea server, you even could add multiple logins for one gitea server
type Login struct {
	Name     string `yaml:"name"`
	URL      string `yaml:"url"`
	Token    string `yaml:"token"`
	Active   bool   `yaml:"active"`
	SSHHost  string `yaml:"ssh_host"`
	Insecure bool   `yaml:"insecure"`
}

// Client returns a client to operate Gitea API
func (l *Login) Client() *gitea.Client {
	client := gitea.NewClient(l.URL, l.Token)
	if l.Insecure {
		cookieJar, _ := cookiejar.New(nil)

		client.SetHTTPClient(&http.Client{
			Jar: cookieJar,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			},
		})
	}
	return client
}

// GetSSHHost returns SSH host name
func (l *Login) GetSSHHost() string {
	if l.SSHHost != "" {
		return l.SSHHost
	}

	u, err := url.Parse(l.URL)
	if err != nil {
		return ""
	}

	return u.Hostname()
}

// Config reprensents local configurations
type Config struct {
	Logins []Login `yaml:"logins"`
}

var (
	config         Config
	yamlConfigPath string
)

func init() {
	homeDir, err := utils.Home()
	if err != nil {
		log.Fatal("Retrieve home dir failed")
	}

	dir := filepath.Join(homeDir, ".tea")
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		log.Fatal("Init tea config dir", dir, "failed")
	}

	yamlConfigPath = filepath.Join(dir, "tea.yml")
}

func splitRepo(repoPath string) (string, string) {
	p := strings.Split(repoPath, "/")
	if len(p) >= 2 {
		return p[0], p[1]
	}
	return repoPath, ""
}

func getActiveLogin() (*Login, error) {
	if len(config.Logins) == 0 {
		return nil, errors.New("No available login")
	}
	for _, l := range config.Logins {
		if l.Active {
			return &l, nil
		}
	}

	return &config.Logins[0], nil
}

func getLoginByName(name string) *Login {
	for _, l := range config.Logins {
		if l.Name == name {
			return &l
		}
	}
	return nil
}

func addLogin(login Login) error {
	for _, l := range config.Logins {
		if l.Name == login.Name {
			if l.URL == login.URL && l.Token == login.Token {
				return nil
			}
			return errors.New("login name has already been  used")
		}
		if l.URL == login.URL && l.Token == login.Token {
			return errors.New("URL has been added")
		}
	}

	u, err := url.Parse(login.URL)
	if err != nil {
		return err
	}

	if login.SSHHost == "" {
		login.SSHHost = u.Hostname()
	}
	config.Logins = append(config.Logins, login)

	return nil
}

func isFileExist(fileName string) (bool, error) {
	f, err := os.Stat(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if f.IsDir() {
		return false, errors.New("the same name directory exist")
	}
	return true, nil
}

func loadConfig(ymlPath string) error {
	exist, _ := isFileExist(ymlPath)
	if exist {
		Println("Found config file", ymlPath)
		bs, err := ioutil.ReadFile(ymlPath)
		if err != nil {
			return err
		}

		err = yaml.Unmarshal(bs, &config)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveConfig(ymlPath string) error {
	bs, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(ymlPath, bs, 0660)
}

func curGitRepoPath() (*Login, string, error) {
	gitConfig := git_config.NewConfig()
	bs, err := ioutil.ReadFile(filepath.Join(filepath.Dir(os.Args[0]), ".git", "config"))
	if err != nil {
		return nil, "", err
	}
	if err := gitConfig.Unmarshal(bs); err != nil {
		return nil, "", err
	}
	remoteConfig, ok := gitConfig.Remotes["origin"]
	if !ok || remoteConfig == nil {
		return nil, "", errors.New("No remote origin found on this git repository")
	}

	for _, l := range config.Logins {
		for _, u := range remoteConfig.URLs {
			p, err := local_git.ParseURL(strings.TrimSpace(u))
			if err != nil {
				return nil, "", fmt.Errorf("Git remote URL parse failed: %s", err.Error())
			}
			if strings.EqualFold(p.Scheme, "http") || strings.EqualFold(p.Scheme, "https") {
				if strings.HasPrefix(u, l.URL) {
					ps := strings.Split(p.Path, "/")
					path := strings.Join(ps[len(ps)-2:], "/")
					return &l, strings.TrimSuffix(path, ".git"), nil
				}
			} else if strings.EqualFold(p.Scheme, "ssh") {
				if l.GetSSHHost() == p.Host {
					return &l, strings.TrimLeft(strings.TrimSuffix(p.Path, ".git"), "/"), nil
				}
			}
		}
	}

	return nil, "", errors.New("No Gitea login found")
}
