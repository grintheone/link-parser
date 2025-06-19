package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

type Link struct {
	Href string
	Text string
}

func (l Link) String() string {
	var sb strings.Builder

	fmt.Fprintln(&sb, "{")
	fmt.Fprintf(&sb, "%7s : %v\n", "href", l.Href)
	fmt.Fprintf(&sb, "%7s : %v\n", "text", l.Text)
	fmt.Fprintf(&sb, "}")

	return sb.String()
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

func readFile(fname *string) (*html.Node, error) {
	f, err := os.Open(*fname)
	if err != nil {
		return nil, fmt.Errorf("unable to open a file: %q, no such filename", *fname)
	}
	defer f.Close()

	return html.Parse(f)
}

func fetchPage(url *string) (*html.Node, error) {
	resp, err := http.Get(*url)
	if err != nil {
		return nil, fmt.Errorf("could not fetch the web page: %v", err)
	}
	defer resp.Body.Close()

	return html.Parse(resp.Body)
}

func run() error {
	var (
		links []Link
		html  *html.Node
		err   error
	)

	fname := flag.String("fname", "", "specify the name of the file to open")
	url := flag.String("url", "", "specify the page url to fetch")

	flag.Parse()

	if *fname != "" && *url != "" {
		fmt.Println("Error: Cannot specify both -fname and -url flags")
		flag.Usage()
		os.Exit(1)
	}

	if *fname == "" && *url == "" {
		fmt.Println("Error: Both -url and -fname flags cannot be emtpy")
		flag.Usage()
		os.Exit(1)
	}

	if *fname != "" {
		html, err = readFile(fname)
	} else {
		html, err = fetchPage(url)
	}

	if err != nil {
		return err
	}

	walk(html, &links)

	fmt.Printf("Links after a walk:\n\n")

	for i, l := range links {
		switch i {
		case len(links) - 1:
			fmt.Printf("%s\n", l)
		default:
			fmt.Printf("%s,\n", l)
		}
	}

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
