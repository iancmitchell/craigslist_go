package main

import (
	"log"
	"net/http"

	"golang.org/x/net/html"
)

//Parse out the individule posts on a Craigslist homepage
func parseHTML(doc *http.Response) []string {
	p := html.NewTokenizer(doc.Body)
	var links []string

	//Loop through the HTML doc passed
	for {
		part := p.Next()

		switch part {
		//If nothing left in doc, return
		case html.ErrorToken:
			return links
		//Check opening HTML tags
		case html.StartTagToken:
			token := p.Token()
			//Check if anchor tag
			isAnchor := token.Data == "a"
			if isAnchor {
				for _, v := range token.Attr {
					//Look for specific class detailing a link to a post
					if v.Key == "class" && v.Val == "result-title hdrlnk" {
						for _, v := range token.Attr {
							if v.Key == "href" {
								links = append(links, v.Val)
							}
						}
					}
				}
			}
		}
	}
}

func getPost(url string) string {
	_, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	return url + ":done"
}

func processPosts(posts <-chan string, results chan<- string) {
	for p := range posts {
		results <- getPost(p)
	}
}

func main() {
	url := "https://newyork.craigslist.org/search/brk/hhh"
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	parsedRes := parseHTML(res)
	buffSize := len(parsedRes)
	posts := make(chan string, buffSize)
	results := make(chan string, buffSize)

	go processPosts(posts, results)
	go processPosts(posts, results)
	go processPosts(posts, results)
	go processPosts(posts, results)

	for _, v := range parsedRes {
		posts <- v
	}
	close(posts)
	for i := 0; i < buffSize; i++ {
		log.Println(i, <-results)

	}
}
