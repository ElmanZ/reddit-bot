package reddit

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type Item struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

type Response struct {
	Data struct {
		Children []struct {
			Data Item
		}
	}
}

const Best = "https://www.reddit.com/r/best"
const Rising = "https://www.reddit.com/r/rising"
const Random = "https://www.reddit.com/r/random"
const SubReddit = "https://www.reddit.com/r/"

func Get(reddit string) ([]string, error) {
	req, err := http.NewRequest("GET", reddit, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}

	//UserAgent retrieves and presents web content for end users e.g. web browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10.15; rv:88.0) Gecko/20100101 Firefox/88.0")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(resp.Status)
	}

	r := new(Response)
	err = json.NewDecoder(resp.Body).Decode(r)
	if err != nil {
		return nil, err
	}

	//grabing values from response and appending them onto slice
	var items []string
	for _, child := range r.Data.Children {
		items = append(items, "\nURL: ", child.Data.URL, "\nTitle: ", child.Data.Title, "\n")
		//for logs
		log.Printf("\nURL: %s, \nTitle: %s", child.Data.URL, child.Data.Title)
	}
	return items, err
}
