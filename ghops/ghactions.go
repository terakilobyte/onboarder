package ghops

import (
	"context"
	"fmt"
	"log"
	"os/exec"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

func InitClient(token string) (*github.Client, *string, string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return nil, nil, "", err
	}

	return client, user.Login, token, nil

}

func ForkRepos(g *github.Client, repos map[string][]string) error {
	for org, orgRepos := range repos {
		for _, repo := range orgRepos {
			fmt.Printf("\nForking %s/%s\n", org, repo)
			_, _, err := g.Repositories.CreateFork(context.Background(), org, repo, &github.RepositoryCreateForkOptions{})
			if err != nil {
				if _, ok := err.(*github.AcceptedError); ok {
					continue
				}
				fmt.Printf("\nerror: %v\n", err)
				return err
			}
		}
	}
	return nil
}

func UploadKeys(g *github.Client, sshKey, gid string) error {
	var skey *string
	skey = &sshKey
	var pTitle *string
	keyTitle := "MongoDB Onboarder"
	pTitle = &keyTitle
	_, _, err := g.Users.CreateKey(context.Background(), &github.Key{Title: pTitle, Key: skey})
	if err != nil {
		return err
	}
	app := "gpg"
	arg1 := "--armor"
	arg2 := "--export"

	gpgCmd := exec.Command(app, arg1, arg2, gid)
	stdout, err := gpgCmd.Output()

	if err != nil {
		log.Fatalf("Error running gpg --armor --export %s: %v", gid, err)
	}
	gpgKey := string(stdout)

	_, _, err = g.Users.CreateGPGKey(context.Background(), gpgKey)
	if err != nil {
		return err
	}
	return nil
}
