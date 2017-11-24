package main

import (
	"encoding/json"
	"fmt"
	"github.com/briandowns/spinner"
	"io/ioutil"
	"net/http"
	"sort"
	"time"
)

// RedditResponse Struct
type RedditResponse struct {
	Data struct {
		Children []struct {
			Data Joke
		}
	}
}

// Joke Struct
type Joke struct {
	Title     string `json:"title"`
	Punchline string `json:"selftext"`
	Ups       int
}

func main() {
	const endpoint string = "https://reddit.com/r/dadjokes.json"

	// set timeout
	netClient := &http.Client{
		Timeout: time.Second * 10,
	}

	// create new request
	req, err := http.NewRequest("GET", endpoint, nil)
	// set header for user-agent
	req.Header.Set("User-agent", "Bot")
	// fetch endpoint
	response, err := netClient.Do(req)

	// start spinner
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()
	time.Sleep(2 * time.Second)

	if err != nil {
		fmt.Println("error", err)
		panic(fmt.Sprintf("Something is wrong, %s", err))
	} else {
		// stop spinner
		s.Stop()

		// handle response
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)

		if err != nil {
			fmt.Println("error", err)
			panic(fmt.Sprintf("Something is wrong, %s", err))
		}

		// assign data as new reddit Struct
		data := RedditResponse{}

		// convert contents to golang struct
		json.Unmarshal(contents, &data)

		// sort data by most liked
		sort.Slice(data.Data.Children, func(i, j int) bool {
			return data.Data.Children[i].Data.Ups > data.Data.Children[j].Data.Ups
		})

		// declare jokes as slice of joke Structs
		var jokes []Joke

		// append Joke struct to slice
		for k := range data.Data.Children {
			jokes = append(jokes, data.Data.Children[k].Data)
		}

		// take top 5 jokes
		jokes = jokes[0:3]

		// print jokes as range
		for k, v := range jokes {
			fmt.Printf("Joke %d \n", k+1)
			fmt.Println(v.Title)
			fmt.Println(v.Punchline)
			fmt.Println()
		}
	}
}
