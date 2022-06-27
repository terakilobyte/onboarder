package githubops

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/google/go-github/v43/github"
	"github.com/terakilobyte/onboarder/globals"
	"golang.org/x/oauth2"
)

const SCOPES = "repo, admin:public_key, admin:gpg_key, user"

type DeviceFlowFirstPostResponse struct {
	DeviceCode      string `json:"device_code"`
	UserCode        string `json:"user_code"`
	VerificationUri string `json:"verification_uri"`
	ExpiresIn       int    `json:"expires_in"`
	Interval        int    `json:"interval"`
}

type DeviceFlowAccessPost struct {
	ClientId   string `json:"client_id"`
	DeviceCode string `json:"device_code"`
	GrantType  string `json:"grant_type"`
}

type DeviceFlowFirstPost struct {
	ClientId string `json:"client_id"`
	Scope    string `json:"scope"`
}

type DeviceFlowAccessResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

func AuthToGithub() string {

	basicValues := DeviceFlowFirstPost{
		ClientId: "bdbcd76255fe7e0ded14",
		Scope:    SCOPES,
	}

	data, err := json.Marshal(basicValues)
	if err != nil {
		log.Fatal(err)
	}

	res, err := postToGithub("https://github.com/login/device/code", data)

	if err != nil {
		log.Fatal(err)
	}

	var b DeviceFlowFirstPostResponse

	err = json.Unmarshal(res, &b)

	if err != nil {
		log.Fatalln(err)
	}

	access := DeviceFlowAccessPost{
		ClientId:   basicValues.ClientId,
		DeviceCode: b.DeviceCode,
		GrantType:  "urn:ietf:params:oauth:grant-type:device_code",
	}

	adata, err := json.Marshal(access)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("A browser should automatically open.\n If it doesn't, "+
		"please open the following URL in your browser to complete the "+
		"authentication process:\n%s\n", b.VerificationUri)
	fmt.Printf("I've copied the following code to your clipboard.\nPlease paste it in the browser: \n\t%s\n", b.UserCode)
	err = clipboard.WriteAll(b.UserCode)
	if err != nil {
		fmt.Println("unable to copy code to clipboard")
	}

	openbrowser(b.VerificationUri)

	ticker := time.NewTicker((time.Duration(b.Interval) + 1) * time.Second)
	timeout := time.After(time.Duration(b.ExpiresIn) * time.Second)
	quit := make(chan struct{})

	var wg sync.WaitGroup

	wg.Add(1)

	var authenticated DeviceFlowAccessResponse

	fmt.Print("Waiting...")

	go func() {
		for {
			select {
			case <-ticker.C:

				res, err := postToGithub("https://github.com/login/oauth/access_token", adata)
				if err != nil {
					log.Fatalln(err)
				}
				var tmp DeviceFlowAccessResponse
				err = json.Unmarshal(res, &tmp)
				if err != nil {
					log.Fatalln(err)
				}
				if tmp.AccessToken != "" {
					authenticated = tmp
					wg.Done()
					return
				} else {
					fmt.Print(".")
				}

			case <-timeout:
				log.Fatal("Failed to sign in within 15 minutes. Restart the script to try again.")
				ticker.Stop()
				return
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	wg.Wait()
	fmt.Println()

	return authenticated.AccessToken

}

func InitClient() error {
	if globals.GITHUBCLIENT == nil {
		globals.AUTHTOKEN = AuthToGithub()
	}
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: globals.AUTHTOKEN},
	)

	tc := oauth2.NewClient(context.Background(), ts)
	globals.GITHUBCLIENT = github.NewClient(tc)
	GetUser(globals.GITHUBCLIENT)
	return nil
}

func postToGithub(url string, data []byte) ([]byte, error) {
	client := &http.Client{}
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(data))
	req.Header = http.Header{
		"Accept":       []string{"application/json"},
		"Content-Type": []string{"application/json"},
	}

	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	d, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return d, nil
}

func openbrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal(err)
	}

}
