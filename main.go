package main

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/token/ngram"
	"github.com/blevesearch/bleve/analysis/token/porter"
	"github.com/blevesearch/bleve/analysis/tokenizer/whitespace"
	"github.com/blevesearch/bleve/search/query"
)

// Message struct defines the document to be indexed
type Message struct {
	ID   string
	From string
	Body string
}

func NgramAnalyzer() map[string]interface{} {
	return map[string]interface{}{
		"type":         custom.Name,
		"char_filters": []string{},
		"tokenizer":    whitespace.Name,
		"token_filters": []string{
			"stop_en",
			porter.Name,
			lowercase.Name,
			"ngram_filter",
		},
	}
}
func NgramTokenFilter() map[string]interface{} {
	return map[string]interface{}{
		"type": ngram.Name,
		"min":  float64(3),
		"max":  float64(4),
	}
}

func runQuery(index bleve.Index, query *query.MatchQuery) {
	req := bleve.NewSearchRequest(query)
	req.Fields = []string{"From", "Body"}

	res, err := index.Search(req)
	if err != nil {
		log.Printf("%v\n", err)
	}

	fmt.Printf("Query: %+v got %v\n", query, res)
}

func main() {
	messages := []Message{
		{
			ID:   "0",
			From: "marty.schoch@gmail.com",
			Body: "bleve indexing is easy",
		},
		{
			ID:   "1",
			From: "aouyang1@gmail.com",
			Body: "I'm Trying to learn bleve",
		},
		{
			ID:   "2",
			From: "souyang1@gmail.com",
			Body: "why is Bleve hard?",
		},
		{
			ID:   "3",
			From: "blargh@gmail.com",
			Body: "what is indexes asdf bevel",
		},
		{
			ID:   "4",
			From: "blargh@gmail.com",
			Body: "what watery indexes asdf bevel",
		},
		{
			ID:   "5",
			From: "blargh@gmail.com",
			Body: "water is something Like This In Wisconsin asdf bevel",
		},
	}

	var idx bleve.Index
	var err error
	idx, err = bleve.Open("example.bleve")
	if err != nil {
		indexMapping := bleve.NewIndexMapping()
		indexMapping.AddCustomTokenFilter("ngram_filter", NgramTokenFilter())
		indexMapping.AddCustomAnalyzer("ngram_analyzer", NgramAnalyzer())
		indexMapping.DefaultAnalyzer = "ngram_analyzer"

		emailMapping := bleve.NewDocumentMapping()
		emailFieldMapping := bleve.NewTextFieldMapping()
		emailMapping.AddFieldMappingsAt("From", emailFieldMapping)
		emailMapping.AddFieldMappingsAt("Body", emailFieldMapping)

		indexMapping.AddDocumentMapping("email", emailMapping)

		idx, err = bleve.New("example.bleve", indexMapping)
		if err != nil {
			panic(err)
		}

		for _, m := range messages {
			idx.Index(m.ID, m)
		}
	}

	var query *query.MatchQuery
	query = bleve.NewMatchQuery("bleve")
	runQuery(idx, query)

	query = bleve.NewMatchQuery("index")
	runQuery(idx, query)

	query = bleve.NewMatchQuery("watered")
	runQuery(idx, query)

	query = bleve.NewMatchQuery("ouya")
	query.SetField("From")
	runQuery(idx, query)

}
