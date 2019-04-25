package git

import (
	"net/url"
	"regexp"
	"strings"
)

var (
	protocolRe = regexp.MustCompile("^[a-zA-Z_+-]+://")
)

// URLParser represents a git URL parser
type URLParser struct {
}

// Parse parses the git URL
func (p *URLParser) Parse(rawURL string) (u *url.URL, err error) {
	if !protocolRe.MatchString(rawURL) &&
		strings.Contains(rawURL, ":") &&
		// not a Windows path
		!strings.Contains(rawURL, "\\") {
		rawURL = "ssh://" + strings.Replace(rawURL, ":", "/", 1)
	}

	u, err = url.Parse(rawURL)
	if err != nil {
		return
	}

	if u.Scheme == "git+ssh" {
		u.Scheme = "ssh"
	}

	if strings.HasPrefix(u.Path, "//") {
		u.Path = strings.TrimPrefix(u.Path, "/")
	}

	return
}

// ParseURL parses URL string and return URL struct
func ParseURL(rawURL string) (u *url.URL, err error) {
	p := &URLParser{}
	return p.Parse(rawURL)
}
