package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	_ "github.com/go-sql-driver/mysql"
	"github.com/mcnijman/go-emailaddress"
)

type WebsiteData struct {
	Url    string
	Userid string
	Dbid   string
}

func GETemails() {
	for true {
		resp, err := http.Get("http://35.180.242.58/getData/")
		if err != nil {
			fmt.Println("No response from request")
		}

		body, err := ioutil.ReadAll(resp.Body) // response body is []byte
		monjson := fmt.Sprintf(string(body))
		fmt.Printf(monjson)

		var websiteData1 WebsiteData
		erre := json.Unmarshal([]byte(monjson), &websiteData1)
		if erre != nil {
			fmt.Println(err)
		}
		fmt.Sprintf("Struct is:", websiteData1.Url)
		url := fmt.Sprintf(websiteData1.Url)
		sitevisite := url

		userid := fmt.Sprintf(websiteData1.Userid)
		dbid := fmt.Sprintf(websiteData1.Dbid)
		http.Get("http://35.180.242.58/getData/?websitevisited=" + sitevisite + "&dbid=" + dbid)
		fmt.Sprintf(userid + dbid)
		defer resp.Body.Close()

		LINKS := []string{}

		doc, err := goquery.NewDocument(string(url))
		if err != nil {
			log.Fatal(err)
		}

		doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
			href, _ := item.Attr("href")
			fmt.Printf("link: %s - anchor text: %s\n", href, item.Text())

			re := regexp.MustCompile(`^(?:https?:\/\/)?(?:[^@\/\n]+@)?(?:www\.)?([^:\/\n]+)`)

			submatchall := re.FindAllString(href, -1)
			for _, element := range submatchall {
				domainame := fmt.Sprintf(element)

				if strings.Contains(href, "https") == false {
					href = domainame + href

				}
			}
			LINKS = append(LINKS, href)
		})
		for i := 0; i < len(LINKS); i++ {

			url := fmt.Sprintf(LINKS[i])

			if strings.Contains(url, "https://") {
				fmt.Printf("HTML code of %s ...\n", url)
				resp, err := http.Get(url)
				// handle the error if there is one
				if err != nil {
					print("error get request")
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
				text := []byte(souphtml)
				validateHost := false
				emails := emailaddress.Find(text, validateHost)

				for _, e := range emails {
					fmt.Println(e)
					var emailfound string = fmt.Sprintf("%s", e)
					if !strings.Contains(emailfound, "png") && !strings.Contains(emailfound, "jpg") && !strings.Contains(emailfound, "www") && !strings.Contains(emailfound, "http") && !strings.Contains(emailfound, "image") && !strings.Contains(emailfound, "webp") {
						http.Get("http://35.180.242.58/getData/?email=" + emailfound + "&userid=" + userid + "&dbid=" + dbid + "&websitevisited=" + sitevisite)

					}
				}

			}
		}
	}
}
func main() {
	maxGoroutines := 1000
	guard := make(chan struct{}, maxGoroutines)

	for i := 0; i < 300; i++ {
		guard <- struct{}{}
		go func(n int) {
			GETemails()
			<-guard
		}(i)
	}
	GETemails()
}