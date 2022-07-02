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
	"github.com/spf13/cobra"
	"github.com/terakilobyte/onboarder/cfg"
	"github.com/terakilobyte/onboarder/githubops"
	"github.com/terakilobyte/onboarder/gitops"
	"github.com/terakilobyte/onboarder/globals"
	"log"
)

var outDir string
var gid string
var publicSSHKey string
var config string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "onboarder",
	Short: "Bootstrap your work git repositories.",
	Long: `Onboarder will bootstrap your work git repositories.

Onboarder is an onboarding tool built for the Docs team at MongoDB.

Onboarder uploads your ssh and  gpg keys to github, and adds it to your
your git config. Commits will be signed by default after using onboarder.

Run onboarder, passing in flags for the output directory where repositories
should be cloned to, the path to your *public* ssh key, your pgp key id,
and the path to the configuration file (provided to you).


onboarder -o ~/work

The above will fork repos appropriate for the *tdbx* team and then clone
them to the ~/work directory.

There will be a pause between forking and cloning. This is to allow time
for larger repositories to fork.

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

Please acknowledge your acceptance and understanding of the above by pressing enter.
`)
		var acknowledge string
		fmt.Scanln(&acknowledge)
		cfg.ParseConfigFile(config)
		err := githubops.InitClient()
		if err != nil {
			log.Fatalln(err)
		}
		githubops.ForkRepos(globals.GITHUBCLIENT, &globals.CONFIG)
		githubops.UploadSSHKey(globals.GITHUBCLIENT, publicSSHKey)
		githubops.UploadGPGKey(globals.GITHUBCLIENT, &gid)

		gitops.SetupLocalRepos(&globals.CONFIG, globals.GITHUBUSER, globals.AUTHTOKEN, outDir)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&outDir, "out-dir", "o", "", "output directory")
	rootCmd.PersistentFlags().StringVarP(&publicSSHKey, "ssh-key", "s", "", "public ssh key")
	rootCmd.PersistentFlags().StringVarP(&gid, "gid", "g", "", "gpg --armor --export xxx")
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "config file")
	cobra.MarkFlagRequired(rootCmd.Flags(), "out-dir")
	cobra.MarkFlagRequired(rootCmd.Flags(), "gid")
	cobra.MarkFlagRequired(rootCmd.Flags(), "ssh-key")
	cobra.MarkFlagRequired(rootCmd.Flags(), "config")
}
