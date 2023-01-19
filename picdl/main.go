package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
)

func takeArgs() (string, int) {
	// argumentgiven := os.Args[1]
	var s string
	flag.StringVar(&s, "s", "", "Search for given query")
	page := flag.Int("n", 5, "# of Pages")
	flag.Parse()
	n := *page
	if len(s) != 0 {
		return s, n
	} else {
		fmt.Println("add some argument")
		os.Exit(1)
	}
	return "", n
}

func main() {

	apiurl, err := url.Parse("https://wallhaven.cc/api/v1/search")
	checkNil(err)
	values := apiurl.Query()
	fmt.Println("Fetching Data........")
	fmt.Println("Downloading will start soon........")
	str, n := takeArgs()
	values.Add("q", str)
	// For loop to increment page number
	for i := 1; i <= n; i++ {
		num := strconv.Itoa(i)
		fmt.Println("Page number: ", num)
		values.Set("page", num)
		apiurl.RawQuery = values.Encode()
		fmt.Println(apiurl.String())
		// make Get rqst
		response, err := http.Get(apiurl.String())
		checkNil(err)
		defer response.Body.Close()
		// Read the data from Get rqst
		data, err := ioutil.ReadAll(response.Body)
		checkNil(err)
		var res Response
		json.Unmarshal(data, &res)
		// for loop to print all the data from json that is extracted
		for _, p := range res.Data {
			_, title := path.Split(p.Path)
			Downloadfile(title, p.Path)
			fmt.Println("Downloaded: ", p.Path)
		}

	}
}

func Downloadfile(title string, url string) {
	// Get the data
	resp, err := http.Get(url)
	checkNil(err)
	defer resp.Body.Close()

	out, err := os.Create(title)
	checkNil(err)
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Println(err)
	}
}

func checkNil(err error) {
	if err != nil {
		panic(err)
	}
}

type Response struct {
	Data Data `json:"data"`
}
type Data []struct {
	Path string `json:"path"`
}
