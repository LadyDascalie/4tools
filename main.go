package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/user"
	"path"
	"strings"

	"github.com/ladydascalie/sortdir/sortdir"
	"golang.org/x/net/html"
)

const (
	boardStem string = "//boards.4chan.org"
	cdnStem   string = "//i.4cdn.org"
)

func main() {
	fmt.Println("Notice:\nwhen your download is complete,")
	fmt.Println("you will find your files under: ")
	fmt.Println(setDownloadFolder() + "\n\n")
	url := getURLFromUser()
	getImgLinks(url)
	fmt.Println("\n\nYou can find your downloads under: ", setDownloadFolder())
	os.Chdir(setDownloadFolder())
	sortdir.RunAsCMD()
}

// getURLFromUser grabs the thread URL from the user and returns it
func getURLFromUser() string {
	input := bufio.NewReader(os.Stdin)
	fmt.Println("Thread url ?")

	url, err := input.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return url
}

// getImgLinks gets the URL provided by the user then collects all the links containing the CDN Stem in them.
// It then passes those links to the downloadContent function
func getImgLinks(url string) {
	response, _ := http.Get(url)
	defer response.Body.Close()

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

			if strings.Contains(href, cdnStem) {
				rawURL := "http:" + href
				fmt.Println(rawURL)
				downloadContent(rawURL)
			}
		}
	}
}

// getHref get the content of the href attribute from <a> tags
func getHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	return
}

// downloadContent makes a get request for the requested file then writes its contents to disk
func downloadContent(linkTo string) {

	setDownloadFolder()

	resp, err := http.Get(linkTo)
	fmt.Println("Downloading... Please wait!")
	fmt.Println(resp.Status)
	defer resp.Body.Close()

	if err != nil {
		log.Fatal("Trouble making GET request!")
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

// setDownloadFolder sets the download folder in the user's home folder
func setDownloadFolder() (dirLocation string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Trouble looking up username!")
	}

	dirLocation = usr.HomeDir + "/4tools_downloads"
	os.MkdirAll(dirLocation, 0755)
	os.Chdir(dirLocation)

	return
}
