package main

import (
	"fmt"
	//"go/alux"
	"encoding/json"
	"github.com/DavidBelicza/TextRank"
	"github.com/PuerkitoBio/goquery"
	"net/http"
	"strings"
	"unicode"
	//"net/url"
	//"strings"
)

func main() {
	http.HandleFunc("/hello", hello)

	http.HandleFunc("/alux", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		url := query.Get("q")
		rankedSentences := customTextRank(url)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"texts": rankedSentences,
		})
	})

	http.ListenAndServe(":8080", nil)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello!"))
	fmt.Println("hello!")
}

func customTextRank(url string) []string {
	ranked := []string{}

	paragraphs, _ := Scrape(url, "p")
	article := strings.Join(paragraphs, " ")

	tr := textrank.NewTextRank()
	// Default Rule for parsing.
	rule := textrank.NewDefaultRule()
	// Default Language for filtering stop words.
	language := textrank.NewDefaultLanguage()
	// Default algorithm for ranking text.
	algorithmDef := textrank.NewDefaultAlgorithm()

	// Add text.
	tr.Populate(article, language, rule)
	// Run the ranking.
	tr.Ranking(algorithmDef)

	// Get the most important 4 sentences. Importance by word occurrence.
	sentences := textrank.FindSentencesByWordQtyWeight(tr, 4)

	// Put just the sentences in slice
	for _, s := range sentences {
		ranked = append(ranked, strings.TrimSpace(s.Value))
	}

	return ranked

}

func Scrape(url, sel string) ([]string, error) {

	var items []string
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return items, err
	}

	doc.Find(sel).Each(func(i int, s *goquery.Selection) {
		item := ""
		paragraph := strings.TrimSpace(s.Text())
		lastDot := strings.LastIndex(paragraph, ".")
		// Remove insufficient length paragraph and cut string after last fullstop
		// todo: fix getting tripped by decimal
		if lastDot >= 175 {
			item = string(paragraph[0 : lastDot+1])
			items = append(items, item)

		}

	})
	return items, err
}

//Remove whitestring
func StringMinifier(in string) (out string) {
	white := false
	for _, c := range in {
		if unicode.IsSpace(c) {
			if !white {
				out = out + " "
			}
			white = true
		} else {
			out = out + string(c)
			white = false
		}
	}
	return
}

//TODO split through number with dot in it e.g. 1.223
