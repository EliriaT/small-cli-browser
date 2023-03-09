package searchEngine

import (
	"encoding/json"
	"fmt"
	"hi/httpClient"
	"log"
	"strings"
)

type Searcher struct {
}

func (s *Searcher) ParseSearchResults(jsonResult string) {
	resultMap := map[string]interface{}{}
	list := strings.Split(jsonResult, "\r\n")

	err := json.Unmarshal([]byte(list[1]), &resultMap)
	if err != nil {
		log.Println("Can not unmarshal json")
		log.Fatal(err)
	}
	results, ok := resultMap["items"].([]interface{})
	if !ok {
		log.Println("Can not parse search results. Sorry.")
		log.Fatal(err)
	}

	resultsLinks := make([]string, 10)
	fmt.Println("Top 10 search results: ")
	for i := 0; i < len(results); i++ {
		result, ok := results[i].(map[string]interface{})
		if !ok {
			log.Println("Can not parse search results. Sorry.")
			log.Fatal(err)
		}
		fmt.Printf("%d. Title: %s \n   Description: %s \n   Link: %s \n\n", i+1, result["title"], result["snippet"], result["link"])
		resultsLinks[i] = result["link"].(string)
	}
	var ans = ""
	var linkNum = 0
	for linkNum <= 0 || linkNum > 10 {
		fmt.Print("What link do you want to browse? Type a number from 1 to 10: ")

		_, err = fmt.Scanf("%d", &linkNum)

		if err != nil || linkNum <= 0 || linkNum > 10 {
			fmt.Print("You didn't type a valid number! Do you want to try again? y/n: ")
			log.Println("nu ajunge")
			_, err = fmt.Scanf("%s", &ans)
			if err != nil {
				fmt.Println("Sorry not a valid string")
				continue
			}
			if strings.ToLower(ans) != "y" && strings.ToLower(ans) != "n" {
				fmt.Println("Sorry not a valid answer. ")
				continue
			}
			if strings.ToLower(ans) == "n" {
				return
			}
		}
	}

	client := httpClient.NewClient(resultsLinks[linkNum-1])

	headers, body := client.MakeHTTPRequest(0)
	client.HandleBody(headers, body)

}

func NewSearcher() Searcher {
	return Searcher{}
}
