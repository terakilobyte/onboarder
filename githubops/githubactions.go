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

func ForkRepos(g *github.Client, cfg *globals.Config, noPause bool) {
	makeForks(cfg, g)
	if !noPause {
		fmt.Println("waiting for 30 seconds to allow for forking to complete")
		time.Sleep(30 * time.Second)
	}
	for _, org := range cfg.Orgs {
		for _, repo := range org.Repos {
			addWatcher(repo, org, g)
			addCollaborators(repo, g)
			addHooks(g, repo, cfg)

		}
	}
}

func makeForks(cfg *globals.Config, g *github.Client) {
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
}

func addWatcher(repo globals.Repo, org globals.Org, g *github.Client) {
	if repo.SetSubscription {
		fmt.Println("adding you as a watcher to " + org.Name + "/" + repo.Name)
		g.Activity.SetRepositorySubscription(context.Background(), org.Name, repo.Name, &github.Subscription{Subscribed: github.Bool(true)})
	}
}

func addHooks(g *github.Client, repo globals.Repo, cfg *globals.Config) {
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
			fmt.Println("adding webhook to " + *globals.GITHUBUSER.Login + "/" + repo.Name)
			_, _, err := g.Repositories.CreateHook(context.Background(), *globals.GITHUBUSER.Login, repo.Name, &github.Hook{
				Name:   github.String("web"),
				Active: github.Bool(true),
				Config: map[string]interface{}{
					"url":          github.String(cfg.Hook.Url),
					"content_type": github.String(cfg.Hook.ContentType),
					"secret":       github.String(cfg.Hook.Secret),
				},
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func addCollaborators(repo globals.Repo, g *github.Client) {
	if len(repo.Collaborators) > 0 {
		directCollaborators, _, err := g.Repositories.ListCollaborators(context.Background(), *globals.GITHUBUSER.Login, repo.Name, &github.ListCollaboratorsOptions{Affiliation: "direct"})
		if err != nil {
			log.Fatal(err)
		}
		for _, collaborator := range repo.Collaborators {
			if !isInDirectCollaborators(directCollaborators, collaborator.Username) {
				fmt.Println("adding " + collaborator.Username + " as a collaborator to " + *globals.GITHUBUSER.Login + "/" + repo.Name)
				_, _, err := g.Repositories.AddCollaborator(context.Background(), *globals.GITHUBUSER.Login, repo.Name, collaborator.Username, &github.RepositoryAddCollaboratorOptions{Permission: collaborator.Permission})
				if err != nil {
					log.Fatal(err)
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

func isInDirectCollaborators(collaborators []*github.User, collaborator string) bool {
	for _, c := range collaborators {
		if c.GetLogin() == collaborator {
			return true
		}
	}
	return false
}
