package main

import (
	"fmt"
	"hi/cache"
	"hi/config"
	"hi/httpClient"
	"hi/searchEngine"
	"log"
	"os"
	"strings"
)

// cache test https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers
// plain text test https://www.lipsum.com/
// json test https://dummyjson.com/products/1
// redirect test google.com
func main() {

	err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load app.env file.")
		os.Exit(1)
	}

	buckets := []string{httpClient.CacheBucket}

	err = cache.InitBolt("./database.boltdb", buckets)

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "1.Usage: go2web -u example.com \n2.Usage: go2web -s how is the weather today\n3.Usage: go2web -h\n")
		os.Exit(1)
	}
	if len(os.Args) >= 2 {
		if os.Args[1] != "-h" && len(os.Args) == 2 {
			fmt.Fprintf(os.Stderr, "1.Usage: go2web -u example.com \n2.Usage: go2web -s how is the weather today\n3.Usage: go2web -h\n")
			os.Exit(1)
		}
	}

	switch os.Args[1] {
	case "-s":
		HandleSCmd()
	case "-u":
		HandleUCmd()
	case "-h":
		fmt.Fprintf(os.Stdout, "1.Usage: go2web -u example.com \n2.Usage: go2web -s how is the weather today\n3.Usage: go2web -h\n")
	default:
		fmt.Fprintf(os.Stderr, "1.Usage: go2web -u example.com \n2.Usage: go2web -s how is the weather today\n3.Usage: go2web -h\n")
		os.Exit(1)
	}

	os.Exit(0)
}

func HandleUCmd() {
	url := os.Args[2]
	client := httpClient.NewClient(url)
	headers, body := client.MakeHTTPRequest(0)
	client.HandleBody(headers, body)
	client.HandleCache(headers, body)
}

func HandleSCmd() {
	searcher := searchEngine.NewSearcher()

	searchTerms := strings.Join(os.Args[2:], "+")
	apiUrl := fmt.Sprintf("https://www.googleapis.com/customsearch/v1?key=%s&cx=92f172b9e6774438a&q=$=%s", config.GlobalConfig.APISecret, searchTerms)
	client := httpClient.NewClient(apiUrl)
	_, body := client.MakeHTTPRequest(0)
	searcher.ParseSearchResults(body)

}
