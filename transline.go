package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/reconquest/karma-go"
)

var (
	version = "[manual build]"
	usage   = "transline " + version + `

Translate word or phrase using Yandex Services.

Usage:
  transline [options] <query>
  transline -h | --help
  transline --version

Options:
  -d --dictionary        Use dictionary for translation.
  -t --translator        Use machinery translation.
  -l --lang <lang>       Translation direction [default: en-ru].
  -s --synonyms <limit>  Limit synonims. [default: 0]
  -o --output <format>   Output format. Can be text or json. [default: text]
  -h --help              Show this screen.
  --version              Show version.
`
)

func main() {
	args, err := docopt.Parse(usage, nil, true, version, false)
	if err != nil {
		panic(err)
	}

	var (
		query       = strings.TrimSpace(args["<query>"].(string))
		synonyms, _ = strconv.Atoi(args["--synonyms"].(string))
		lang        = args["--lang"].(string)

		useDictionary = args["--dictionary"].(bool)
		useTranslator = args["--translator"].(bool)
		useJSON       = args["--output"].(string) == "json"
	)

	if !useTranslator && !useDictionary {
		if strings.Contains(query, " ") {
			useTranslator = true
		} else {
			useDictionary = true
		}
	}

	yandex := NewYandexProvider(lang, "", "", synonyms)

	switch {
	case useTranslator:
		result, err := yandex.Translate(query)
		if err != nil {
			log.Fatal(karma.Format(err, "unable to translate phrase"))
		}

		if result == "" {
			os.Exit(2)
		}

		if useJSON {
			outputJSON(result)
			return
		}

		fmt.Println(result)

	case useDictionary:
		result, err := yandex.LookupDictionary(query)
		if err != nil {
			log.Fatal(karma.Format(err, "unable to lookup word in dictionary"))
		}

		if result == nil {
			os.Exit(2)
		}

		if useJSON {
			outputJSON(result)
			return
		}

		fmt.Println(result.Transcript)
		for _, meaning := range result.Meanings {
			fmt.Println(meaning.String())
		}
	}
}

func outputJSON(result interface{}) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(result)
	if err != nil {
		log.Fatal(karma.Format(err, "unable to encode result"))
	}
}
