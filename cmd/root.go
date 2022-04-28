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

	"github.com/spf13/cobra"
	"github.com/terakilobyte/onboarder/genssh"
	"github.com/terakilobyte/onboarder/ghops"
	"github.com/terakilobyte/onboarder/gitops"
	"github.com/terakilobyte/onboarder/globals"
)

var outDir string
var team string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "onboarder",
	Short: "Bootstrap your work git repositories.",
	Long:  `Bootstrap your work git repositories.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		client, user, err := ghops.InitClient(ghops.AuthToGithub())
		if err != nil {
			log.Fatalln(err)
		}
		fmt.Printf("Hello, %s\n", *user)
		fmt.Print(`
I'm about to begin forking and cloning all of the repositories that you should need.
I'm also going to create an ssh key for you and add it to your github account.

This may take a while (5-10 minutes) depending on how many repositories I'm
working with. Please be patient.
`)

		pubkey, keypath, err := genssh.SetupSSH(*user)
		if err != nil {
			log.Fatalln(err)
		}
		ghops.UploadKey(client, pubkey)
		// ghops.ForkRepos(client, globals.GetReposForTeam(team))
		gitops.SetupLocalRepos(globals.GetReposForTeam(team), *user, outDir, keypath)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.Flags().StringVarP(&outDir, "out-dir", "o", "", "output directory")
	rootCmd.Flags().StringVarP(&team, "team", "t", "", "team name")
	cobra.MarkFlagRequired(rootCmd.Flags(), "out-dir")
	cobra.MarkFlagRequired(rootCmd.Flags(), "team")
}
