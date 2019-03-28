package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

var translations = loadTranslations()

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		log.Fatal("Give templates to process as arguments")
	}
	for _, template := range flag.Args() {
		processTemplate(template)
	}
}

func processTemplate(path string) {
	tpl, err := template.New("root").Funcs(template.FuncMap{
		"translate": translate,
	}).ParseFiles(path)
	if err != nil {
		log.Fatal("Parsing template:", err)
	}

	ext := filepath.Ext(path)
	if ext != ".tpl" {
		log.Fatal("Should be a .tpl file:", path)
	}
	out, err := os.Create(path[:len(path)-len(ext)])
	if err != nil {
		log.Fatal("Output file:", err)
	}

	if err := tpl.ExecuteTemplate(out, filepath.Base(path), translations); err != nil {
		log.Fatal("Executing template:", err)
	}

	if err := out.Close(); err != nil {
		log.Fatal("Output file:", err)
	}
}

type translation struct {
	Lang   string
	Phrase string
}

// translate returns a list of viable translations for the given English phrase.
func translate(english string) []translation {
	var res []translation
	for lang, ts := range translations {
		if trans, ok := ts[english]; ok && trans != english {
			res = append(res, translation{
				Lang:   lang,
				Phrase: trans,
			})
		}
	}
	return res
}

// loadTranslations loads our json translation files into a large map of
// language code -> english phrase -> translated phrase
func loadTranslations() map[string]map[string]string {
	files, err := filepath.Glob("gui/default/assets/lang/lang-*.json")
	if err != nil {
		log.Fatal("Listing available translations:", err)
	}

	res := make(map[string]map[string]string)
	for _, file := range files {
		lang := strings.Replace(strings.Replace(filepath.Base(file), "lang-", "", 1), ".json", "", 1)
		bs, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal("Reading translation:", err)
		}
		var trans map[string]string
		if err := json.Unmarshal(bs, &trans); err != nil {
			log.Fatalf("Parsing translation %s: %v", lang, err)
		}
		res[lang] = trans
	}

	return res
}
