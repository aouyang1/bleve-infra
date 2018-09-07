package main

import (
	"fmt"
	"log"

	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/analysis/analyzer/custom"
	"github.com/blevesearch/bleve/analysis/token/lowercase"
	"github.com/blevesearch/bleve/analysis/token/ngram"
	"github.com/blevesearch/bleve/analysis/tokenizer/single"
	"github.com/blevesearch/bleve/analysis/tokenizer/whitespace"
)

// Message struct defines the document to be indexed
type Message struct {
	ID   string
	From string
	Body string
}

// SingleTermAnalyzer is a custom analyzer
func SingleTermAnalyzer() map[string]interface{} {
	return map[string]interface{}{
		"type":          custom.Name,
		"char_filters":  []string{},
		"tokenizer":     single.Name,
		"token_filters": []string{},
	}
}

func NgramAnalyzer() map[string]interface{} {
	return map[string]interface{}{
		"type":         custom.Name,
		"char_filters": []string{},
		"tokenizer":    whitespace.Name,
		"token_filters": []string{
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
func main() {
	message := Message{
		ID:   "example",
		From: "marty.schoch@gmail.com",
		Body: "bleve indexing is easy",
	}

	var index bleve.Index
	var err error
	index, err = bleve.Open("example.bleve")
	if err != nil {
		indexMapping := bleve.NewIndexMapping()
		indexMapping.AddCustomTokenFilter("ngram_filter", NgramTokenFilter())
		indexMapping.AddCustomAnalyzer("ngram_analyzer", NgramAnalyzer())

		docMapping := bleve.NewDocumentMapping()
		indexMapping.AddDocumentMapping("simple_doc", docMapping)

		rawFieldMapping := bleve.NewTextFieldMapping()
		rawFieldMapping.Analyzer = "ngram_analyzer"

		docMapping.AddFieldMappingsAt("From", rawFieldMapping)
		docMapping.AddFieldMappingsAt("Body", rawFieldMapping)

		index, err = bleve.New("example.bleve", indexMapping)
		if err != nil {
			panic(err)
		}

		index.Index(message.ID, message)
	}

	query := bleve.NewMatchQuery("indexin")
	query.SetField("Body")

	searchRequest := bleve.NewSearchRequest(query)
	searchRequest.Fields = []string{"Body"}

	searchResult, err := index.Search(searchRequest)
	if err != nil {
		log.Printf("%v", err)
	}

	fmt.Printf("%+v", searchResult)
}
