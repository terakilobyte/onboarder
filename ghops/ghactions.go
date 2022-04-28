package ghops

import (
	"context"
	"fmt"

	"github.com/google/go-github/v43/github"
	"golang.org/x/oauth2"
)

func InitClient(token string) (*github.Client, *string, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)

	user, _, err := client.Users.Get(ctx, "")
	if err != nil {
		fmt.Printf("\nerror: %v\n", err)
		return nil, nil, err
	}

	return client, user.Login, nil

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

func UploadKey(g *github.Client, pubkey string) error {
	var pkey *string
	pkey = &pubkey
	var pTitle *string
	keyTitle := "MongoDB Onboarder"
	pTitle = &keyTitle
	_, _, err := g.Users.CreateKey(context.Background(), &github.Key{Title: pTitle, Key: pkey})
	if err != nil {
		return err
	}
	return nil
}
