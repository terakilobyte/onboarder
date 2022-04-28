package gitops

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func SetupLocalRepos(repos map[string][]string, user, outdir, keypath string) error {

	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		err = os.MkdirAll(outdir, 0700)
		if err != nil {
			return err
		}
	}
	_, err := os.Stat(keypath)
	if err != nil {
		log.Fatalf("read file %s failed %s\n", keypath, err.Error())
		return err
	}
	publicKeys, err := ssh.NewPublicKeysFromFile("git", keypath, "")
	if err != nil {
		log.Fatalf("generating public keys failed %s\n", err.Error())
	}
	for org, orgRepos := range repos {
		for _, repo := range orgRepos {
			dest := path.Join(outdir, repo)
			url := fmt.Sprintf("git@github.com:%s/%s.git", user, repo)

			fmt.Printf("\nCloning %s/%s forked from %s\n", user, repo, org)
			r, err := git.PlainClone(dest, false, &git.CloneOptions{
				URL:      url,
				Progress: os.Stdout,
				Auth:     publicKeys,
			})
			if err != nil {
				log.Fatalf("clone repo %s failed %s\n", url, err.Error())
			}
			currentConfig, err := r.Config()
			if err != nil {
				log.Fatalf("get repo config %s failed %s\n", url, err.Error())
			}
			currentConfig.Remotes["upstream"] = &config.RemoteConfig{
				Name:  "upstream",
				URLs:  []string{fmt.Sprintf("git@github.com:%s/%s.git", org, repo)},
				Fetch: []config.RefSpec{"+refs/heads/*:refs/remotes/upstream/*"},
			}
			var branch *config.Branch
			if _, ok := currentConfig.Branches["main"]; ok {
				branch = currentConfig.Branches["main"]
			} else {
				branch = currentConfig.Branches["master"]
			}
			branch.Remote = "upstream"

			r.SetConfig(currentConfig)
		}
	}
	return nil
}
