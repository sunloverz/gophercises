package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"html"
	"net/http"
)

type urlPath struct {
	path string
	url  string
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if url, ok := pathsToUrls[html.EscapeString(r.URL.Path)]; ok {
			http.Redirect(w, r, url, 301)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	parsedYaml, err := parseYAML(yml)

	if err != nil {
		panic(err)
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}

func buildMap(parsedYaml interface{}) map[string]string {
	pathsToUrls := make(map[string]string)

	for _, value := range parsedYaml.([]interface{}) {
		line := value.(map[interface{}]interface{})
		pathsToUrls[line["path"].(string)] = line["url"].(string)
	}

	return pathsToUrls
}

func parseYAML(yml []byte) (interface{}, error) {
	var t interface{}
	err := yaml.Unmarshal(yml, &t)

	if err != nil {
		return "", err
	}

	return t, nil
}

func main() {
	mux := defaultMux()

	yaml := `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`

	yamlHandler, err := YAMLHandler([]byte(yaml), mux)

	if err != nil {
		panic(err)
	}
	http.ListenAndServe(":3000", yamlHandler)
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, world!")
}
