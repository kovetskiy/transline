package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const (
	// flags=4 is for morphological search
	YandexDictionaryURLFormat = "https://dictionary.yandex.net" +
		"/api/v1/dicservice.json/lookup?key=%s&lang=%s&text=%s&flags=4"

	YandexDictionaryKey = "dict.1.1.20140512T122957Z.549af1de13649236." +
		"74bbc11e0fa7625166dd95f21b1ff17838df2c03"

	YandexTranslateKey = "trnsl.1.1.20150807T071843Z.319e163bafb5d806." +
		"f17c493a466253265a817c9f96a74db85f6b556b"

	YandexTranslateURLFormat = "https://translate.yandex.net" +
		"/api/v1.5/tr.json/translate?key=%s&lang=%s&text=%s"
)

type YandexProvider struct {
	dictionaryKey string
	translateKey  string
	lang          string
	limitSynonyms int
}

func NewYandexProvider(
	lang,
	dictionaryKey string,
	translateKey string,
	limitSynonyms int,
) *YandexProvider {
	if dictionaryKey == "" {
		dictionaryKey = YandexDictionaryKey
	}
	if translateKey == "" {
		translateKey = YandexTranslateKey
	}

	return &YandexProvider{dictionaryKey, translateKey, lang, limitSynonyms}
}

func (yandex YandexProvider) LookupDictionary(text string) (*DictionaryItem, error) {
	url := fmt.Sprintf(
		YandexDictionaryURLFormat,
		yandex.dictionaryKey,
		yandex.lang,
		url.QueryEscape(text),
	)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf(
			"expected HTTP status 200, received %s (%s)",
			resp.Status,
			url,
		)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := struct {
		Def []struct {
			Pos string
			Ts  string
			Tr  []struct {
				Text string
				Mean []struct {
					Text string
				}
			}
		}
	}{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Def) == 0 {
		return nil, nil
	}

	item := DictionaryItem{
		Transcript: result.Def[0].Ts,
	}

	for _, d := range result.Def {
		for j, tr := range d.Tr {
			if yandex.limitSynonyms > 0 && j >= yandex.limitSynonyms {
				break
			}
			references := []string{}
			for _, ref := range tr.Mean {
				references = append(references, ref.Text)
			}

			item.Meanings = append(item.Meanings, Meaning{
				Translation: tr.Text,
				References:  references,
			})
		}
	}

	return &item, nil
}

func (yandex YandexProvider) Translate(text string) (string, error) {
	url := fmt.Sprintf(
		YandexTranslateURLFormat,
		yandex.translateKey,
		yandex.lang,
		url.QueryEscape(text),
	)

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf(
			"expected HTTP status 200, received %s (%s)",
			resp.Status,
			url,
		)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	result := struct {
		Text []string
	}{}

	err = json.Unmarshal(body, &result)
	if err != nil {
		return "", err
	}

	if len(result.Text) == 0 {
		return "", nil
	}

	return result.Text[0], nil
}
