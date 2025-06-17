package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func walk(n *html.Node, links *[]Link) {
	if n.Type == html.ElementNode && n.Data == "a" {
		var l Link
		extractText(n, &l)

		for _, a := range n.Attr {
			if a.Key == "href" {
				l.Href = a.Val
				break
			}
		}

		*links = append(*links, l)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walk(c, links)
	}
}

func extractText(n *html.Node, l *Link) {
	if n.Type == html.TextNode {
		l.Text += " " + strings.TrimSpace(n.Data)
		l.Text = strings.TrimSpace(l.Text)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		extractText(c, l)
	}
}

func main() {
	var links []Link
	fname := flag.String("fname", "", "specify the name of the file to open")
	flag.Parse()

	f, err := os.Open(*fname)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	doc, err := html.Parse(f)
	if err != nil {
		log.Fatal(err)
	}

	walk(doc, &links)

	fmt.Printf("Links after a walk %#v", links)
}
