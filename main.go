package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"golang.org/x/net/html"
)

var result []string

/*
- Search 4 - 5 pages.
- Give it in a csv format.
*/
func getHtml(url string) (string, *http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("%s", err)
		return "", resp, err
	}
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		return "", resp, err
	}
	text := string(bytes)
	return text, resp, nil
}

func checkForClasses(attri []html.Attribute, class string) bool {
	for _, val := range attri {
		if val.Val == class {
			return true
		}
	}
	return false
}

func getHref(attri []html.Attribute, class string) string {
	isClass := checkForClasses(attri, class)
	if isClass {
		for _, val := range attri {
			if val.Key == "href" {
				return val.Val
			}
		}
	}
	return ""
}

// Item Price && div - _30jeq3
// Item Link && Anchor tag - s1Q9rs
// Item Rating - _3LWZlK
func parse(text string) {
	htmlTokens := html.NewTokenizer(strings.NewReader(text))
	var isPrice bool
	var link string
	var rating bool
	res := strings.Builder{}
	for {
		tt := htmlTokens.Next()
		switch tt {
		case html.ErrorToken:
			fmt.Println("End")
			return
		case html.StartTagToken:
			t := htmlTokens.Token()
			isPrice = checkForClasses(t.Attr, "_30jeq3")
			link = getHref(t.Attr, "s1Q9rs")
			rating = checkForClasses(t.Attr, "_3LWZlK")
			if len(link) > 1 {
				res.WriteString(link + ",")
			}
			link = ""
		case html.TextToken:
			t := htmlTokens.Token()
			if isPrice {
				res.WriteString(t.Data + "\n")
				result = append(result, res.String())
				res.Reset()
			}
			if rating {
				res.WriteString(t.Data + ",")
			}
			isPrice = false
			rating = false
		}
	}
}

func getProductInformation(url string, wg *sync.WaitGroup) {
	text, _, err := getHtml(url)
	if err != nil {
		fmt.Println("Error in Fetching")
	}
	parse(text)
	wg.Done()
}

func main() {
	var product string
	fmt.Println("Enter a product name")
	fmt.Scan(&product)
	url := "https://www.flipkart.com/search?q=" + product + "&otracker=search&otracker1=search&marketplace=FLIPKART&as-show=on&as=off&page="

	var wg sync.WaitGroup
	wg.Add(5)
	for i := 1; i <= 5; i++ {
		urls := url + fmt.Sprint(i)
		getProductInformation(urls, &wg)
	}
	wg.Wait()
	putValueInFile(product)
}
