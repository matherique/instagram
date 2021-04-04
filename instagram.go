package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

const LINK string = "https://www.instagram.com/p/CNIq9h4Hu-J/"

var mutex sync.Mutex

func main() {
	fmt.Fprintf(os.Stdout, "ola mundo")
	client := &http.Client{}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		mutex.Lock()
		req, err := http.NewRequest("GET", LINK, nil)

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

		fmt.Fprintf(os.Stdout, url)
		mutex.Unlock()
		wg.Done()
	}()
	wg.Wait()

}
