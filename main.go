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

func main() {""
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
	userid := fmt.Sprintf(websiteData1.Userid)
	dbid := fmt.Sprintf(websiteData1.Dbid)


	fmt.Sprintf(userid + dbid)
	defer resp.Body.Close()

	LINKS := []string{}

	doc, err := goquery.NewDocument(url)
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
			text := []byte(souphtml)
			validateHost := false
			emails := emailaddress.Find(text, validateHost)

			for _, e := range emails {
				fmt.Println(e)
				var emailfound string = fmt.Sprintf("%s", e)
				if !strings.Contains(emailfound, "png") && !strings.Contains(emailfound, "jpg") {
					http.Get("http://35.180.242.58/getData/")
					

				}
			}

		

		}
	}

}
