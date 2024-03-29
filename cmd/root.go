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
	"github.com/terakilobyte/onboarder/cfg"
	"github.com/terakilobyte/onboarder/githubops"
	"github.com/terakilobyte/onboarder/gitops"
	"github.com/terakilobyte/onboarder/globals"
)

var outDir string
var publicSSHKey string
var config string
var noPause bool
var noClone bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "onboarder",
	Short: "Bootstrap your work git repositories.",
	Long: `Onboarder will bootstrap your work git repositories.

Onboarder is an onboarding tool built for the Docs team at MongoDB. You'll
need a configuration file to pass in the list of repositories you want to
fork and clone. Ask your onboarding buddy or team lead for more info.

Onboarder can upload your ssh key to github.

Run onboarder, passing in flags for the output directory where repositories
should be cloned to and the path to the configuration file (provided to you).

	onboarder -o ~/work

To upload your ssh key to github, pass in the path to your public ssh key:

	onboarder -o ~/work -c teamconfig.json -s ~/.ssh/id_rsa.pub

`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Print(`
I'm about to begin forking and cloning all of the repositories that you should need.
I'm also going to add you as a watcher. Ask your team for tips on
useful email filters.

If you've passed the -s command, I'll upload your public ssh key to github.

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
		githubops.ForkRepos(globals.GITHUBCLIENT, &globals.CONFIG, noPause)
		if publicSSHKey != "" {
			githubops.UploadSSHKey(globals.GITHUBCLIENT, publicSSHKey)
		}
		if !noClone {
			gitops.SetupLocalRepos(&globals.CONFIG, globals.GITHUBUSER, globals.AUTHTOKEN, outDir)
		}
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
	rootCmd.PersistentFlags().StringVarP(&config, "config", "c", "", "config file")
	rootCmd.PersistentFlags().BoolVar(&noPause, "no-pause", false, "don't pause between forking and cloning")
	rootCmd.PersistentFlags().BoolVar(&noClone, "no-clone", false, "don't clone repositories")
	cobra.MarkFlagRequired(rootCmd.Flags(), "out-dir")
	cobra.MarkFlagRequired(rootCmd.Flags(), "config")
}
