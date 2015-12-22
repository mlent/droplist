package main

import (
	"encoding/json"
	"fmt"
	"github.com/digitalocean/godo"
	"github.com/getlantern/systray"
	"github.com/skratchdot/open-golang/open"
	"golang.org/x/oauth2"
	"os"
	"strconv"
)

type Configuration struct {
	PersonalAccessToken string
}
type TokenSource struct {
	AccessToken string
}

func main() {
	systray.Run(renderList)
}

func renderList() {
	PAT := getTokenFromFile()
	client := authenticateClient(PAT)
	dropletList, _ := DropletList(client)

	systray.SetTitle("Droplets")
	systray.SetTooltip("You have " + strconv.Itoa(len(dropletList)) + " droplets")
	mItem, dropletUrl := getDropletMenuItem(dropletList[0])

	for {
		select {
		case <-mItem.ClickedCh:
			open.Run(dropletUrl)
		}
	}
}

func getDropletMenuItem(droplet godo.Droplet) (item *systray.MenuItem, url string) {
	name := droplet.Name
	ip := droplet.Networks.V4[0].IPAddress
	region := getFlagByRegionSlug(droplet.Region.Slug)
	itemText := fmt.Sprintf("%s - %s %s", name, ip, region)

	item = systray.AddMenuItem(itemText, "Quit it!")
	url = "https://cloud.digitalocean.com/droplets/" +
		strconv.Itoa(droplet.ID)

	return
}

func getFlagByRegionSlug(region string) string {
	flags := map[string]string{
		"fra1": "\U0001F1E9\U0001F1EA", // Frankfurt
	}
	return flags[region]
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

/*
* Authenticating client to use DO API
 */
func getTokenFromFile() string {
	// Get PersonalAccessToken from config file
	file, _ := os.Open("config.json")
	decoder := json.NewDecoder(file)
	config := Configuration{}
	err := decoder.Decode(&config)

	if err != nil {
		fmt.Println("error", err)
	}

	return config.PersonalAccessToken
}

func authenticateClient(accessToken string) (client *godo.Client) {
	tokenSource := &TokenSource{
		AccessToken: accessToken,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client = godo.NewClient(oauthClient)
	return
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}
