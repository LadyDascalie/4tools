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

	"github.com/fatih/color"
	"github.com/ladydascalie/sortdir/sortdir"
	"golang.org/x/net/html"
	"sync"
)

const (
	boardStem string = "//boards.4chan.org"
	cdnStem string = "//i.4cdn.org"
)

func main() {

	// Pretty print the notice to the user
	beginNotice()

	// Get the URL from the user
	url := getURLFromUser()

	// Start downloading the images from the URL
	media := getImgLinks(url)

	var wg sync.WaitGroup

	for _, v := range media {
		wg.Add(1)
		go downloadContent(&wg, v)
	}

	wg.Wait()

	// Print out the completion notice
	endNotice()

	// CD into the download folder
	os.Chdir(setDownloadFolder())

	// Sort the files by filetype
	sortdir.RunAsCMD()
}

func beginNotice() {
	// Pretty print the notice
	color.Green("******")
	color.Green("~ Notice: ~\n")
	color.White("When your download is complete,\nyou will find your files under:\n")
	color.Magenta(setDownloadFolder())
	color.Green("******" + "\n\n")
}

func endNotice() {
	color.Green("******")
	color.Green("Download complete!\n")
	color.White("Your files have been saved to " + setDownloadFolder() + "\n\n")
	color.White("For your convenience, your files have been sorted by extension.")
	color.Green("******")
}

// getURLFromUser grabs the thread URL from the user and returns it
func getURLFromUser() string {
	input := bufio.NewReader(os.Stdin)
	fmt.Println("Please paste the thread URL:")

	// Read the string until a newline character is entered
	// AKA: Until the return key is pressed.
	url, err := input.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return url
}

// getImgLinks gets the URL provided by the user then collects all the links containing the CDN Stem in them.
// It then passes those links to the downloadContent function
func getImgLinks(url string) []string {
	var urls []string

	response, _ := http.Get(url)
	defer response.Body.Close()

	z := html.NewTokenizer(response.Body)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document
			return urls
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
				urls = append(urls, rawURL)
			}
		}
	}
	return urls
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
func downloadContent(wg *sync.WaitGroup, linkTo string) {
	defer wg.Done()

	setDownloadFolder()

	resp, err := http.Get(linkTo)
	fmt.Println("Downloading... Please wait!")
	color.Green(resp.Status)
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
func setDownloadFolder() (string) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Trouble looking up username!")
	}

	dirLocation := usr.HomeDir + "/4tools_downloads"
	os.MkdirAll(dirLocation, 0755)
	os.Chdir(dirLocation)

	return dirLocation
}
