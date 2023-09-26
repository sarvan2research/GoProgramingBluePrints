package main

import (
	"bufio"
	"chapter4/thesaurus"
	"fmt"
	"log"
	"os"
)

func main() {
	apiKey := os.Getenv("BHT_APIKEY")
	thesaurus := thesaurus.BigHuge{APIKEY: apiKey}
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		word := s.Text()
		synonym, err := thesaurus.Synonyms(word)
		if err != nil {
			log.Fatalln("Failed when looking for synonyms for "+word, err)
		}
		if len(synonym) == 0 {
			log.Fatalln("Couldn't find any synonyms for " + word)
		}
		for _, syn := range synonym {
			fmt.Println(syn)
		}
	}
}
