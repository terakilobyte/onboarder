package gitops

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

func SetupLocalRepos(repos map[string][]string, user, token, outdir string) error {

	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		err = os.MkdirAll(outdir, 0700)
		if err != nil {
			return err
		}
	}

	for org, orgRepos := range repos {
		for _, repo := range orgRepos {
			dest := path.Join(outdir, repo)
			url := fmt.Sprintf("https://github.com/%s/%s.git", user, repo)

			fmt.Printf("\nCloning %s/%s forked from %s\n", user, repo, org)
			r, err := git.PlainClone(dest, false, &git.CloneOptions{
				URL:      url,
				Progress: os.Stdout,
				Auth:     &http.BasicAuth{Username: user, Password: token},
			})
			if err != nil {
				if err.Error() == "repository already exists" {
					continue
				}
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

func ConfigSignedCommits(gid string) {
	app := "git"
	arg1 := "config"
	arg2 := "--global"
	arg3 := "commit.gpgsign"
	arg4 := "true"

	gitCmd := exec.Command(app, arg1, arg2, arg3, arg4)
	gitCmd.Run()

	app = "git"
	arg1 = "config"
	arg2 = "--global"
	arg3 = "user.signingkey"

	gitCmd = exec.Command(app, arg1, arg2, arg3, gid)
	gitCmd.Run()
}
