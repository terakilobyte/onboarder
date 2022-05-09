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
	"github.com/terakilobyte/onboarder/ghops"
	"github.com/terakilobyte/onboarder/gitops"
	"github.com/terakilobyte/onboarder/globals"
)

var sshPath string

// rootCmd represents the base command when called without any subcommands
// Uncomment the following line if your bare application
// has an action associated with it:

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clones repoitories for a team.",
	Long: `The clone subcommand will clone repositories specified by
your team to your local computer. It will also fork repositories, if necessary.
Repositories that are already forked are a no-op.

  onboarder clone -s ~/.ssh/my_key_ed25519 -t tdbx -o ~/work`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(`
I'm about to begin forking and cloning all of the repositories that you should need.

This may take a while (5-10 minutes) depending on how many repositories I'm
working with. Please be patient.

There will be a pause between forking and cloning. This is to allow time
for larger repositories to fork.

Please acknowledge your acceptance and understanding of the above by pressing enter.
`)
		var acknowledge string
		fmt.Scanln(&acknowledge)
		client, user, err := ghops.InitClient(ghops.AuthToGithub())
		if err != nil {
			log.Fatalln(err)
		}

		ghops.ForkRepos(client, globals.GetReposForTeam(team))

		fmt.Println("Waiting 30 seconds for forks to complete...")
		time.Sleep(30 * time.Second)

		gitops.SetupLocalRepos(globals.GetReposForTeam(team), *user, outDir, sshPath)
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	cloneCmd.Flags().StringVarP(&sshPath, "ssh", "s", "", "Path to your private ssh key")
}
