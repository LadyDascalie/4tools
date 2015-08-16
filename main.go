package main

import (
	"bufio"
	"fmt"
	"regexp"
	// "golang.org/x/net/html"
	// "net/http"
	"os"
	"strings"
)

const (
	BOARD_STEM string = "http://boards.4chan.org"
	IMG_SLICER string = "s"
	IMG_ID_LEN int    = 13
)

func main() {
	fmt.Println(getUrl())
}

func getUrl() string {
	input := bufio.NewReader(os.Stdin)
	fmt.Println("Thread url ?")

	url, err := input.ReadString('\n')

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	getBoard(url)
	return url
}

func getBoard(url string) string {
	var boardName []byte
	if strings.Contains(url, BOARD_STEM) {
		regE := regexp.MustCompile("(.org/([a-z]{1,4})/)")
		slicedUrl := []byte(url)
		match := regE.Find(slicedUrl)

		regE = regexp.MustCompile("/([a-z]{1,4})/")
		boardName = regE.Find(match)
	}
	return string(boardName)
}
