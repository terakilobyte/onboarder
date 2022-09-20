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
	Name  string `json:"name"`
	Repos []Repo `json:"repos"`
}

type Repo struct {
	Name            string         `json:"name"`
	UseWebhook      bool           `json:"useWebhook"`
	SetSubscription bool           `json:"setSubscription"`
	Collaborators   []Collaborator `json:"collaborators"`
}

type Webhook struct {
	Url         string `json:"url"`
	ContentType string `json:"content_type"`
	Secret      string `json:"secret"`
	SSLVerify   string `json:"ssl_verify"`
}
type Collaborator struct {
	Username   string `json:"username"`
	Permission string `json:"permission"`
}
