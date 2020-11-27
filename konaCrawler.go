package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/remeh/sizedwaitgroup"
)

var (
	swg      = sizedwaitgroup.New(20)
	savePath = "./images"
)

func extractLinks(url string) (links []string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	doc.Find(".directlink").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("href")
		if exists {
			links = append(links, src)
		}
	})
	return links, err
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	return err
}

func init() {
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		os.Mkdir(savePath, 0644)
	}
}

func main() {
	for page := 1; page <= 1000; page++ {
		url := "http://konachan.net/post?page=" + strconv.Itoa(page)
		links, err := extractLinks(url)
		if err != nil {
			continue
		}
		for idx, src := range links {
			defer swg.Wait()
			swg.Add()
			go func(idx int, src string) {
				defer swg.Done()
				err := downloadFile(fmt.Sprintf("./images/P%vN%v%s", page, idx, filepath.Ext(src)), src)
				if err != nil {
					fmt.Println(err)
				}
			}(idx, src)
		}
	}
}
