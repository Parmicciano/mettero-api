package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/mcnijman/go-emailaddress"
	"log"
    "github.com/PuerkitoBio/goquery"
)

func main() {
	url := "https://amazon.fr"
	fmt.Printf("HTML code of %s ...\n", url)
	resp, err := http.Get(url)
	// handle the error if there is one
	if err != nil {
		panic(err)
	}
	// do this now so it won't be forgotten
	defer resp.Body.Close()
	// reads html as a slice of bytes
	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	// show the HTML code as a string %s
	souphtml := fmt.Sprintf("%s\n", html)
	


		doc, err := goquery.NewDocument("https://en.wikipedia.org/wiki/Example.com")
		
		if err != nil {
			log.Fatal(err)
		}
		
		doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
			href, _ := item.Attr("href")
			fmt.Printf("link: %s - anchor text: %s\n", href, item.Text())
			
		})
		
	
	
	foundemails(souphtml)
}
func foundemails(souphtml string) {
	text := []byte(souphtml)
	validateHost := false

	emails := emailaddress.Find(text, validateHost)

	for _, e := range emails {
		fmt.Println(e)
	}
}
