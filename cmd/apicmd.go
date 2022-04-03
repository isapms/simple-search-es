package cmd

import (
	"log"
	"os"
	"simple-search-es/elasticsearch"

	_ "github.com/joho/godotenv/autoload"
)

func ConfigEs() *elasticsearch.ElasticSearch {
	esHost := os.Getenv("ES_HOST")
	esIndex := os.Getenv("ES_INDEX")
	esAlias := os.Getenv("ES_ALIAS")

	elastic, err := elasticsearch.New([]string{esHost})
	if err != nil {
		log.Fatalln(err)
	}

	err = elastic.CreateIndex(esIndex, esAlias)
	if err != nil {
		log.Fatalln(err)
	}

	return elastic
}
