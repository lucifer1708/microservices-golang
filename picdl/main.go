package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"sync"
)

func main() {
	apiurl, err := url.Parse("https://wallhaven.cc/api/v1/search")
	checkNil(err)
	values := apiurl.Query()
	fmt.Println("Fetching Data........")
	fmt.Println("Downloading will start soon........")
	str, n := takeArgs()
	values.Add("q", str)

	// Create a WaitGroup to wait for all downloads to finish
	// var wg sync.WaitGroup

	// For loop to increment page number
	for i := 1; i <= n; i++ {
		num := strconv.Itoa(i)
		fmt.Println("Page number: ", num)
		values.Set("page", num)
		apiurl.RawQuery = values.Encode()
		fmt.Println(apiurl.String())

		response, err := http.Get(apiurl.String())
		checkNil(err)
		defer response.Body.Close()

		// Read the data from the GET request
		data, err := io.ReadAll(response.Body)
		checkNil(err)

		var res Response
		json.Unmarshal(data, &res)

		// Create a WaitGroup for each page to ensure all files are downloaded before moving to the next page
		var pageWg sync.WaitGroup

		// Download each file concurrently with individual progress bars
		for _, p := range res.Data {
			_, title := path.Split(p.Path)

			// Increment the WaitGroup counter for each file
			pageWg.Add(1)

			// Start a goroutine to download the file
			go func(title, url string) {
				defer pageWg.Done()
				DownloadFileWithProgressBar(title, url)
				fmt.Println("Downloaded: ", url)
			}(title, p.Path)
		}

		// Wait for all files in the page to finish downloading
		pageWg.Wait()
	}

	fmt.Println("All downloads completed.")
}

// This function takes arguments from the flags specified by the user
func takeArgs() (string, int) {
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

// Function to download a file from URL. It takes the title (file name) and URL as parameters.
func DownloadFile(title string, url string) {
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

// Function to download a file from URL with a progress bar. It takes the title (file name) and URL as parameters.
func DownloadFileWithProgressBar(title string, url string) {
	resp, err := http.Get(url)
	checkNil(err)
	defer resp.Body.Close()

	out, err := os.Create(title)
	checkNil(err)
	defer out.Close()

	// Get the content length of the response
	contentLength := resp.ContentLength

	// Create a progress bar with the content length
	bar := pb.Full.Start64(contentLength)
	bar.Set(pb.Bytes, true)

	// Create a proxy reader to update the progress bar
	reader := bar.NewProxyReader(resp.Body)

	_, err = io.Copy(out, reader)
	checkNil(err)

	bar.Finish()
}

// CheckNil function follows the principle of DRY
func checkNil(err error) {
	if err != nil {
		panic(err)
	}
}

// Structs to parse JSON fetched from API
type Response struct {
	Data Data `json:"data"`
}

type Data []struct {
	Path string `json:"path"`
}

