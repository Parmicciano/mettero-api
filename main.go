package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/mcnijman/go-emailaddress"
)

func main() {
	url := "http://tour.golang.org/welcome/"
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
	fmt.Printf("%s\n", html)
	foundemails()
}
func foundemails() {
	text := []byte(`Send me an email at foo@bar.com or foo@domain.fakesuffix.`)
	validateHost := false

	emails := emailaddress.Find(text, validateHost)

	for _, e := range emails {
		fmt.Println(e)
	}
}
