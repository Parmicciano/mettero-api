package main

import (
	"database/sql"
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

func main() {
	db, err := sql.Open("mysql", "Parmicciano:Cholet44$$@tcp(15.236.150.103:3306)/mettero")
	if err != nil {
		panic(err.Error())
	}
	rows, err := db.Query("SELECT url, user_id, db_id FROM urltoget WHERE isvisited=0")
	if err != nil {
		// handle this error better than this
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var url string
		var userid int
		var dbid int

		err = rows.Scan(&url, &userid, &dbid)
		if err != nil {
			// handle this error
			panic(err)
		}

		fmt.Println(url, userid, dbid)
		// get any error encountered during iteration
		err = rows.Err()
		if err != nil {
			panic(err)
		}
		defer db.Close()
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
				foundemails(souphtml)
			}
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
