package genssh

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/terakilobyte/onboarder/genssh/keygen"
)

func SetupSSH(user string) (string, string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home directory: %v", err)
	}
	sshdir := filepath.Join(home, ".ssh")

	if _, err := os.Stat(sshdir); os.IsNotExist(err) {
		err = os.MkdirAll(sshdir, 0700)
		if err != nil {
			return "", "", err
		}
	}

	keyname := fmt.Sprintf("%s_github", user)
	keypath := filepath.Join(sshdir, keyname)
	k, err := keygen.NewWithWrite(keypath, []byte(""), keygen.Ed25519)
	if err != nil {
		log.Fatalf("Error creating new key: %v", err)
	}
	fmt.Println("A new keypair has been generated and added to your ~/.ssh directory.")
	fmt.Printf("The key name is: %s\n", keyname)
	fmt.Println("Adding the key to your keyring.")

	// start the ssh-agent
	app := "ssh-agent"
	flag := "-s"

	cmd := exec.Command(app, flag)
	stdout, err := cmd.Output()
	if err != nil {
		log.Fatalf("Error starting ssh-agent: %v", err)
	}

	app2 := "ssh-add"
	flag2 := "-K"
	target2 := fmt.Sprintf("%s_%s", keypath, keygen.Ed25519)

	cmd = exec.Command(app2, flag2, target2)
	stdout, err = cmd.Output()

	if err != nil {
		log.Fatalf("Error running ssh-add: %v", err)
	}

	fmt.Println((string(stdout)))

	fullKeyName := fmt.Sprintf("%s_%s", keypath, keygen.Ed25519)
	sshCfgFile := path.Join(sshdir, "config")

	homedir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting user home directory: %v", err)
	}

	knownHosts, err := os.OpenFile(path.Join(homedir, ".ssh", "known_hosts"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
	knownHosts.Close()

	f, err := os.OpenFile(sshCfgFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Error creating ssh config file: %v", err)
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("\nHost *\n\tUseKeychain yes\n\tAddKeystoAgent yes\n\tHostName github.com\n\tUser git\n\tIdentityFile %s\n", fullKeyName))
	if err != nil {
		log.Fatalf("Error writing to ssh config file: %v", err)
	}

	app3 := "ssh"
	target3 := "git@github.com"
	out, err := exec.Command(app3, target3).Output()
	// if err != nil {
	// 	log.Fatalf("Error running ssh and git to add to known_hosts: %v", err)
	// }
	fmt.Println(string(out))

	return string(k.PublicKey()), fullKeyName, nil
}
