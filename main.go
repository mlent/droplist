package main

import (
	"github.com/getlantern/systray"
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"encoding/json"
	"fmt"
	"os"
)
type Configuration struct {
	PersonalAccessToken string
}
type TokenSource struct {
	AccessToken string
}

func main() {
	// Authenticate with DO
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)

	if err != nil {
		fmt.Println("error", err)
	}

	fmt.Println(config.PersonalAccessToken)

	tokenSource := &TokenSource{
		AccessToken: config.PersonalAccessToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)
	dropletList, _ := DropletList(client)

	fmt.Println(dropletList)

	//systray.Run(onReady)
}

func onReady() {
	systray.SetTitle("Awesome app")
	systray.SetTooltip("Awesomeeeee tooltip")
	systray.AddMenuItem("Quit", "Quit it!")
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func DropletList(client *godo.Client) ([]godo.Droplet, error) {
	list := []godo.Droplet{}
	opt := &godo.ListOptions{}

	for {
		droplets, resp, err := client.Droplets.List(opt)
		if err != nil {
			return nil, err
		}

		for _, d := range droplets {
			list = append(list, d)
		}

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, err
		}

		opt.Page = page + 1
	}

	return list, nil
}
