package urlshort

import (
	"encoding/json"
	"net/http"

	"gopkg.in/yaml.v2"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if val, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, val, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

type Entry struct {
	Path string
	Url  string
}

func parseYAML(yml []byte) ([]Entry, error) {
	doc := []Entry{}

	err := yaml.Unmarshal(yml, &doc)
	return doc, err
}

func buildMap(parsedYAML []Entry) map[string]string {
	doc := make(map[string]string, len(parsedYAML))

	for _, v := range parsedYAML {
		doc[v.Path] = v.Url
	}

	return doc
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYAML, err := parseYAML(yml)

	if err != nil {
		return nil, err
	}

	doc := buildMap(parsedYAML)

	return MapHandler(doc, fallback), nil
}

func parseJSON(jsonBytes []byte) ([]Entry, error) {
	var doc []Entry
	err := json.Unmarshal(jsonBytes, &doc)

	return doc, err
}

func JSONHandler(jsonBytes []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedJSON, err := parseJSON(jsonBytes)

	if err != nil {
		return nil, err
	}

	doc := buildMap(parsedJSON)

	return MapHandler(doc, fallback), nil
}
