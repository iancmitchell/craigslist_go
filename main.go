package main

import (
	"log"
	"net/http"

	"golang.org/x/net/html"
)

//Parse out the individule posts on a Craigslist homepage
func parseHTML(doc *http.Response) []string {
	p := html.NewTokenizer(doc.Body)
	log.Println(p)
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

func main() {
	url := "https://newyork.craigslist.org/search/brk/hhh"
	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}

	defer res.Body.Close()
	parsedRes := parseHTML(res)
	log.Println(parsedRes)
}
