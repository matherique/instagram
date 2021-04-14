package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

func getPhoto(link string) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", link, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/51.0.2704.103 Safari/537.36")

	resp, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	data, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()

	reg := regexp.MustCompile("<meta property=\"og:image\"([\\w\\W]+?)/>")
	matches := reg.Find(data)

	url := strings.Replace(string(matches), "<meta property=\"og:image\" content=\"", "", -1)
	url = strings.Replace(url, "\" />", "", -1)

	return url
}

func main() {
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "missing instagram link")
		os.Exit(1)
	}

	photos := make(chan string, flag.NArg())
	links := make(chan string, flag.NArg())

	go func() {
		for l := range links {
			fmt.Println(l)
			photos <- getPhoto(l)
		}
	}()

	go func() {
		for l := range links {
			fmt.Println(l)
			photos <- getPhoto(l)
		}
	}()

	for _, l := range flag.Args() {
		links <- l
	}

	close(links)

	for i := 0; i < flag.NArg(); i++ {
		fmt.Fprintln(os.Stdout, <-photos)
	}

}
