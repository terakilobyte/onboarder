package githubops

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/terakilobyte/onboarder/gitops"

	"github.com/google/go-github/v43/github"
	"github.com/terakilobyte/onboarder/globals"
)

func ForkRepos(g *github.Client, cfg *globals.Config) {
	for _, org := range cfg.Orgs {
		for _, repo := range org.Repos {
			fmt.Printf("\nForking %s/%s\n", org.Name, repo.Name)
			_, _, err := g.Repositories.CreateFork(context.Background(), org.Name, repo.Name, &github.RepositoryCreateForkOptions{})
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
			if repo.SetSubscription {
				fmt.Println("adding you as a watcher to " + org.Name + "/" + repo.Name)
				g.Activity.SetRepositorySubscription(context.Background(), org.Name, repo.Name, &github.Subscription{Subscribed: github.Bool(true)})
			}
			hooks, _, err := g.Repositories.ListHooks(context.Background(), *globals.GITHUBUSER.Login, repo.Name, nil)
			if err != nil {
				log.Fatal(err)
			}
			if repo.UseWebhook {

				found := false
				for _, hook := range hooks {
					if hook.Config["url"] == cfg.Hook.Url {
						found = true
						break
					}
				}

				if !found {

					_, _, err := g.Repositories.CreateHook(context.Background(), *globals.GITHUBUSER.Login, repo.Name, &github.Hook{
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

func UploadSSHKey(g *github.Client, sshKeyPath string) {
	fmt.Println("uploading ssh key")
	dat, err := os.ReadFile(sshKeyPath)
	if err != nil {
		log.Fatal(err)
	}
	key := github.String(string(dat))
	_, _, err = g.Users.CreateKey(context.Background(), &github.Key{
		Title: github.String(strings.TrimSuffix(filepath.Base(sshKeyPath), filepath.Ext(sshKeyPath))),
		Key:   key,
	})
	if err != nil && !strings.Contains(err.Error(), "key is already in use") {
		log.Fatal(err.Error())
	}
	gitops.ConfigSSH()
}

func GetUser(g *github.Client) {
	user, _, err := g.Users.Get(context.Background(), "")
	if err != nil {
		log.Fatal(err)
	}
	globals.GITHUBUSER = user
}
