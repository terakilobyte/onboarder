package globals

import "github.com/google/go-github/v43/github"

var AUTHTOKEN string
var GITHUBCLIENT *github.Client
var GITHUBUSER *github.User
var CONFIG Config

type Config struct {
	Orgs []Org   `json:"orgs"`
	Hook Webhook `json:"webhook"`
}

type Org struct {
	Name  string   `json:"name"`
	Repos []string `json:"repos"`
}

type Webhook struct {
	Url         string `json:"url"`
	ContentType string `json:"content_type"`
	Secret      string `json:"secret"`
	SSLVerify   string `json:"ssl_verify"`
}
