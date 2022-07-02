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
	"github.com/google/go-github/v43/github"
	"github.com/terakilobyte/onboarder/globals"
)

func SetupLocalRepos(cfg *globals.Config, user *github.User, token, outdir string) {
	if _, err := os.Stat(outdir); os.IsNotExist(err) {
		err = os.MkdirAll(outdir, 0700)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, org := range cfg.Orgs {
		for _, repo := range org.Repos {
			dest := path.Join(outdir, repo)
			url := fmt.Sprintf("https://github.com/%s/%s.git", *user.Login, repo)

			fmt.Printf("\nCloning %s/%s forked from %s\n", *user.Login, repo, org.Name)
			r, err := git.PlainClone(dest, false, &git.CloneOptions{
				URL:      url,
				Progress: os.Stdout,
				Auth:     &http.BasicAuth{Username: *user.Login, Password: token},
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
				URLs:  []string{fmt.Sprintf("https://github.com/%s/%s.git", org.Name, repo)},
				Fetch: []config.RefSpec{"+refs/heads/*:refs/remotes/upstream/*"},
			}
			var branch *config.Branch
			if _, ok := currentConfig.Branches["main"]; ok {
				branch = currentConfig.Branches["main"]
			} else {
				branch = currentConfig.Branches["master"]
			}
			branch.Remote = "upstream"

			err = r.SetConfig(currentConfig)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func ConfigSSH() {
	app := "git"
	arg1 := "config"
	arg2 := "--global"
	arg3 := "url.git@github.com:.insteadof"
	arg4 := "https://github.com/"

	gitCmd := exec.Command(app, arg1, arg2, arg3, arg4)
	err := gitCmd.Run()
	if err != nil {
		fmt.Println("unable to set ssh instead of https")
	}
}

func ConfigSignedCommits(gid *string) {
	app := "git"
	arg1 := "config"
	arg2 := "--global"
	arg3 := "commit.gpgsign"
	arg4 := "true"

	gitCmd := exec.Command(app, arg1, arg2, arg3, arg4)
	err := gitCmd.Run()
	if err != nil {
		fmt.Println("unable to set commit.gpgsign to true")
	}

	app = "git"
	arg1 = "config"
	arg2 = "--global"
	arg3 = "user.signingkey"

	gitCmd = exec.Command(app, arg1, arg2, arg3, *gid)
	err = gitCmd.Run()
	if err != nil {
		fmt.Println("unable to set user.signingkey to gid")
	}
}
