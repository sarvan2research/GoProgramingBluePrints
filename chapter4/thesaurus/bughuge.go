package thesaurus

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Thesaurus interface {
	Synonyms(term string) ([]string, error)
}

type BigHuge struct {
	APIKEY string
}
type synonyms struct {
	Noun *wordResult `json:noun`
	Verb *wordResult `json:verb`
}
type wordResult struct {
	Syn []string `json:syn`
}

func (b *BigHuge) Synonyms(word string) ([]string, error) {
	var syns []string
	result, err := http.Get("https://words.bighugelabs.com/api/2/" + b.APIKEY + "/" + word + "/json")
	if err != nil {
		return syns, errors.New("bighuge: Failed when looking for  synonyms for" + word + err.Error())
	}
	var data synonyms
	defer result.Body.Close()
	if err := json.NewDecoder(result.Body).Decode(&data); err != nil {
		return syns, err
	}
	if data.Noun != nil {
		syns = append(syns, data.Noun.Syn...)
	}
	if data.Verb != nil {
		syns = append(syns, data.Verb.Syn...)
	}
	return syns, nil
}
