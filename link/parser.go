package link

import (
	"strings"
	"golang.org/x/net/html"
	"io"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	var links []Link
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	nodes := linkNodes(doc)
	for _, node := range nodes {
		links = append(links, buildLink(node))
	}
	return links, nil
}

func linkNodes(n *html.Node) []*html.Node {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []*html.Node{n}
	}
	var nodes []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, linkNodes(c)...)
	}
	return nodes
}

func buildLink(node *html.Node) Link {
	var link Link
	for _, a := range node.Attr {
		if a.Key == "href" {
			link.Href = a.Val
			break
		}
	}
	link.Text = text(node)
	return link
}

func text(node *html.Node) string {
	if node.Type == html.TextNode {
		return strings.TrimSpace(node.Data) 
	}
	var value string
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		value += text(c)
	}
	return value
}
