/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"
	"github.com/terakilobyte/onboarder/genssh"
	"github.com/terakilobyte/onboarder/ghops"
	"github.com/terakilobyte/onboarder/gitops"
	"github.com/terakilobyte/onboarder/globals"
)

var outDir string
var cfgFile string
var team string
var gid string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "onboarder",
	Short: "Bootstrap your work git repositories.",
	Long: `Onboarder will bootstrap your work git repositories.

Onboarder is an onboarding tool built for the Docs team at MongoDB (initially).

Onboarder generates a new ssh keypair and uploads the public
key to github for you. It will also add it to the ssh-agent, and it modifies
your ~/.ssh/config file (creates if needed) to use the key.

Onboarder also uploads your gpg key to github, and adds it to your
your git config. Commits will be signed by default after using onboarder.

Run onboarder, passing in flags for the output directory where repositories
should be cloned to, and which team you are on.

Current teams are cet, and tdbx. Future versions will accept a config file
rather than have this hardcoded in.

onboarder -t tdbx -o ~/work

The above will fork repos appropriate for the *tdbx* team and then clone
them to the ~/work directory.

There will be a pause between forking and cloning. This is to allow time
for larger repositories to fork.

IMPORTANT: You will be asked a question during the process similar to:

  The authenticity of host 'github.com (140.82.112.4)' can't be established.
  ED25519 key fingerprint is SHA256:+DiY3wvvV6TuJJhbpZisF/zLDA0zPMSvHdkr4UvCOqU.
  This key is not known by any other names
  Are you sure you want to continue connecting (yes/no/[fingerprint])?

You must answer yes to this question. It is adding the fingerprint to your
known_hosts file.


If you have already set up ssh keys, running this tool may be cause an error
with your ~/.ssh/known_hosts file. Delete the file and run the tool again.
	`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(`
I'm about to begin forking and cloning all of the repositories that you should need.
I'm also going to create an ssh key for you and add it to your github account.

This may take a while (5-10 minutes) depending on how many repositories I'm
working with. Please be patient.

There will be a pause between forking and cloning. This is to allow time
for larger repositories to fork.

IMPORTANT: You will be asked a question during the process similar to:

  The authenticity of host 'github.com (140.82.112.4)' can't be established.
  ED25519 key fingerprint is SHA256:+DiY3wvvV6TuJJhbpZisF/zLDA0zPMSvHdkr4UvCOqU.
  This key is not known by any other names
  Are you sure you want to continue connecting (yes/no/[fingerprint])?

You must answer yes to this question.

Please acknowledge your acceptance and understanding of the above by pressing enter.
`)
		var acknowledge string
		fmt.Scanln(&acknowledge)
		client, user, token, err := ghops.InitClient(ghops.AuthToGithub())
		if err != nil {
			log.Fatalln(err)
		}

		sshKey, _, err := genssh.SetupSSH(*user)
		if err != nil {
			log.Fatalln(err)
		}
		ghops.UploadKeys(client, sshKey, gid)
		ghops.ForkRepos(client, globals.GetReposForTeam(team))

		fmt.Println("Waiting 30 seconds for forks to complete...")
		time.Sleep(30 * time.Second)

		gitops.SetupLocalRepos(globals.GetReposForTeam(team), *user, token, outDir)
		gitops.ConfigSignedCommits(gid)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outDir, "out-dir", "o", "", "output directory")
	rootCmd.PersistentFlags().StringVarP(&team, "team", "t", "", "team name")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.onboarder.yaml)")
	rootCmd.Flags().StringVarP(&gid, "gid", "g", "", "gpg --armor --export xxx")
	cobra.MarkFlagRequired(rootCmd.Flags(), "out-dir")
	cobra.MarkFlagRequired(rootCmd.Flags(), "team")
	cobra.MarkFlagRequired(rootCmd.Flags(), "gid")
}
