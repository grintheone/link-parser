package main

import (
	"os"
	"testing"
)

type linkTests struct {
	name  string
	fname string
	items int
	want  []Link
}

func (lt linkTests) run(t *testing.T) {
	var links []Link

	html, err := readFile(&lt.fname)
	if err != nil {
		t.Fatal(err)
	}

	walk(html, &links)
	assertSize(t, len(links), lt.items)

	for i, l := range links {
		assertString(t, l.Href, lt.want[i].Href)
		assertString(t, l.Text, lt.want[i].Text)
	}
}

var table = []linkTests{
	{
		name:  "single link",
		fname: "ex1.html",
		items: 1,
		want: []Link{
			{
				Href: "/other-page",
				Text: "A link to another page",
			},
		},
	},
	{
		name:  "links with inner html",
		fname: "ex2.html",
		items: 2,
		want: []Link{
			{
				Href: "https://www.twitter.com/joncalhoun",
				Text: "Check me out on twitter",
			},
			{
				Href: "https://github.com/gophercises",
				Text: "Gophercises is on Github !",
			},
		},
	},
	{
		name:  "links with comments and inner links",
		fname: "ex4.html",
		items: 3,
		want: []Link{
			{
				Href: "/dog-cat",
				Text: "dog cat",
			},
			{
				Href: "#",
				Text: "Something here",
			},
			{
				Href: "/dog",
				Text: "nested dog link",
			},
		},
	},
}

func TestWalk(t *testing.T) {
	for _, tt := range table {
		t.Run(tt.name, tt.run)
	}
}

func TestRun(t *testing.T) {
	os.Args = append(os.Args, "-fname=ex4.html")

	run()
}

func assertSize(t testing.TB, got, want int) {
	t.Helper()
	if got != want {
		t.Fatalf("got %d links from needed %d", got, want)
	}
}

func assertString(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
