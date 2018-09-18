package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/levigross/grequests"
)

var waitGroup sync.WaitGroup

func main() {
	count := 0
	for page := 1; page <= 100; page++ {
		baseURL := "http://konachan.net/post?page=" + strconv.Itoa(page)
		fmt.Printf("\n\n[第%d页]\n\n", page)
		res, _ := grequests.Get(baseURL, nil)
		doc, err := goquery.NewDocumentFromReader(res)
		if err != nil {
			fmt.Errorf("下载错误:%#v\n", err)
			os.Exit(-1)
		}
		defer res.Close()
		n := 0
		doc.Find(".directlink").Each(func(i int, s *goquery.Selection) {
			src, exists := s.Attr("href")
			if exists {
				// fmt.Println(src)
				waitGroup.Add(1)
				go func(src string) {
					defer waitGroup.Done()
					filename := src[88:len(src)]
					res, _ := grequests.Get(src, &grequests.RequestOptions{Headers: map[string]string{
						"Referer":    "http://konachan.net",
						"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.92 Safari/537.36"}})
					if err := res.DownloadToFile(filename); err != nil {
						log.Println("Error: ", err)
					} else {
						count++
						fmt.Printf("已下载%d张\n", count)
					}
					n++
				}(src)
			}
		})
	}
	waitGroup.Wait()
	fmt.Println("下载完成")
}
