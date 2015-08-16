package main

import (
	"bufio"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"
)

const (
	BOARD_STEM string = "//boards.4chan.org"
	CDN_STEM   string = "//i.4cdn.org"
)

func main() {
	url := getUrlFromUser()
	getImgLinks(url)
	fmt.Println("\n\nYou can find your downloads under: ", setDownloadFolder())
}

func getUrlFromUser() string {
	input := bufio.NewReader(os.Stdin)
	fmt.Println("Thread url ?")

	url, err := input.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return url
}

func getImgLinks(url string) {
	response, _ := http.Get(url)

	z := html.NewTokenizer(response.Body)
	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document
			return

		case tt == html.StartTagToken:
			t := z.Token()
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			ok, href := getHref(t)
			if !ok {
				continue
			}

			if strings.Contains(href, CDN_STEM) {
				rawUrl := "http:" + href
				fmt.Println(rawUrl)
				downloadContent(rawUrl)
			}
		}
	}

	response.Body.Close()
}

func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

func downloadContent(linkTo string) {

	setDownloadFolder()

	resp, err := http.Get(linkTo)
	fmt.Println("Downloading... Please wait!")
	fmt.Println(resp.Status)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal("Trouble making GET photo request!")
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Trouble reading response body!")
	}

	filename := path.Base(linkTo)
	if filename == "" {
		log.Fatalf("Trouble deriving file name for %s", linkTo)
	}

	err = ioutil.WriteFile(filename, contents, 0644)
	if err != nil {
		log.Fatal("Trouble creating file! -- ", err)
	}
}

func setDownloadFolder() (dirLocation string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Trouble looking up username!")
	}

	dirLocation = usr.HomeDir + "/4tools_downloads"
	fmt.Println(dirLocation)
	os.MkdirAll(dirLocation, 0755)
	os.Chdir(dirLocation)

	return
}
