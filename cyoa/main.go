package main

import (
	"encoding/json"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Story map[string]Chapter

type Option struct {
	Text    string `json:"text"`
	Chapter string `json:"arc"`
}

type Chapter struct {
	Title      string   `json:"title"`
	Paragraphs []string `json:"story"`
	Options    []Option `json:"options"`
}

func NewHandler(s Story, t *template.Template) http.Handler {
	return handler{s, t}
}

type handler struct {
	s Story
	t *template.Template
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := h.path(r)

	if chapter, ok := h.s[path]; ok {
		h.t.Execute(w, chapter)
	}
}

func (h handler) path(r *http.Request) string {
	path := strings.TrimSpace(r.URL.Path)
	if path == "" || path == "/" {
		path = "/intro"
	}
	return path[1:]
}

func JsonStory(r io.Reader) (Story, error) {
	var story Story
	data, _ := ioutil.ReadAll(r)

	if err := json.Unmarshal(data, &story); err != nil {
		return nil, err
	}

	return story, nil
}

func main() {
	var story Story
	file, _ := os.Open("gopher.json")
	defer file.Close()
	story, _ = JsonStory(file)
	tmpl := template.Must(template.ParseFiles("layout.html"))
	log.Fatal(http.ListenAndServe(":3000", NewHandler(story, tmpl)))
}
