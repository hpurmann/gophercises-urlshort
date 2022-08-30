package main

import (
	"flag"
	"fmt"
	urlshort "gophercises-urlshort"
	"net/http"
	"os"
)

func main() {
	yamlFilePath := flag.String("yaml", "", "YAML file with mapping from path to redirect URL")
	flag.Parse()
	mux := defaultMux()

	// Build the MapHandler using the mux as the fallback
	pathsToUrls := map[string]string{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)

	var yamlBytes []byte
	if *yamlFilePath != "" {
		fmt.Println("Found flag yaml")
		content, err := os.ReadFile(*yamlFilePath)
		if err != nil {
			fmt.Println("Error reading file", yamlFilePath)
		}
		yamlBytes = content
	} else {
		// Build the YAMLHandler using the mapHandler as the
		// fallback
		yamlBytes = []byte(`
  - path: /urlshort
    url: https://github.com/gophercises/urlshort
  - path: /urlshort-final
    url: https://github.com/gophercises/urlshort/tree/solution
  `)
	}

	yamlHandler, err := urlshort.YAMLHandler(yamlBytes, mapHandler)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
