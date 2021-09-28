package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/mcnijman/go-emailaddress"
)

func main() {
	db, err := sql.Open("mysql", "Parmicciano:Cholet44$$@tcp(15.236.150.103:3306)/mettero")

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	LINKS := []string{}
	urltoget := "https://apple.com"
	doc, err := goquery.NewDocument(urltoget)
	u, err := url.Parse(urltoget)
	if err != nil {
		log.Fatal(err)
	}
	parts := strings.Split(u.Hostname(), ".")
	domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
	domainame := fmt.Sprint(domain)
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")
		fmt.Printf("link: %s - anchor text: %s\n", href, item.Text())

		if strings.Contains(href, "https") == false {
			href = "https://" + domainame + href
		}
		LINKS = append(LINKS, href)
	})
	for i := 0; i < len(LINKS); i++ {

		url := fmt.Sprintf(LINKS[i])

		if strings.Contains(url, "https") {
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
			foundemails(souphtml)
		}
	}

}
func foundemails(souphtml string) {
	text := []byte(souphtml)
	validateHost := false

	emails := emailaddress.Find(text, validateHost)

	for _, e := range emails {
		fmt.Println(e)
	}
}
