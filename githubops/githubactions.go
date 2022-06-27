package githubops

import (
	"context"
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/google/go-github/v43/github"
	"github.com/terakilobyte/onboarder/globals"
)

func ForkRepos(g *github.Client, cfg *globals.Config) {
	for _, org := range cfg.Orgs {
		for _, repo := range org.Repos {
			fmt.Printf("\nForking %s/%s\n", org.Name, repo)
			_, _, err := g.Repositories.CreateFork(context.Background(), org.Name, repo, &github.RepositoryCreateForkOptions{})
			if err != nil {
				if _, ok := err.(*github.AcceptedError); !ok {
					log.Fatalf("\nerror: %v\n", err)
				}
			}
		}
	}
	fmt.Println("waiting for 30 seconds to allow for forking to complete")
	time.Sleep(30 * time.Second)
	fmt.Println("setting up repo webhooks")
	for _, org := range cfg.Orgs {
		for _, repo := range org.Repos {
			hooks, _, err := g.Repositories.ListHooks(context.Background(), *globals.GITHUBUSER.Login, repo, nil)
			if err != nil {
				log.Fatal(err)
			}
			for _, hook := range hooks {
				if hook.Config["url"] != cfg.Hook.Url {
					_, _, err = g.Repositories.CreateHook(context.Background(), *globals.GITHUBUSER.Login, repo, &github.Hook{
						Name:   github.String("web"),
						Active: github.Bool(true),
						Config: map[string]interface{}{
							"url":          github.String(cfg.Hook.Url),
							"content_type": github.String(cfg.Hook.ContentType),
							"secret":       github.String(cfg.Hook.Secret),
							"ssl_verify":   github.String(cfg.Hook.Secret),
						},
					})
					if err != nil {
						log.Fatal(err)
					}
				}
			}

		}
	}
}

func UploadKeys(g *github.Client, sshKey, gid *string) {
	var pTitle *string
	keyTitle := "MongoDB"
	pTitle = &keyTitle
	_, _, err := g.Users.CreateKey(context.Background(), &github.Key{Title: pTitle, Key: sshKey})
	if err != nil && !strings.Contains(err.Error(), "key is already in use") {
		log.Fatal(err.Error())
	}
	app := "gpg"
	arg1 := "--armor"
	arg2 := "--export"

	gpgCmd := exec.Command(app, arg1, arg2, *gid)
	stdout, err := gpgCmd.Output()

	if err != nil {
		log.Fatalf("Error running gpg --armor --export %v: %v", gid, err)
	}
	gpgKey := string(stdout)

	_, _, err = g.Users.CreateGPGKey(context.Background(), gpgKey)
	if err != nil && !strings.Contains(err.Error(), "key_id already exists") {
		log.Fatal(err)
	}
}

func GetUser(g *github.Client) {
	user, _, err := g.Users.Get(context.Background(), "")
	if err != nil {
		log.Fatal(err)
	}
	globals.GITHUBUSER = user
}
