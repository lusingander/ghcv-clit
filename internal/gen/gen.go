//go:generate go run . -output ../ui/credits_gen.go
package main

// cf: https://github.com/lusingander/fyne-credits-generator

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Songmu/gocredits"
)

var (
	output = flag.String("output", "./gen.go", "output filepath")
)

const (
	splitterLicense = "================================================================"
	splitterText    = "----------------------------------------------------------------"
)

type credit struct {
	name, url, text string
}

// [`] => [` + "`" + `]
const replacedBackquote = "`" + ` + "` + "`" + `" + ` + "`"

func (c *credit) formattedText() string {
	return strings.Replace(c.text, "`", replacedBackquote, -1)
}

func collect() ([]*credit, error) {
	buf, err := runGoCredits()
	if err != nil {
		return nil, err
	}
	licenses := strings.Split(buf.String(), splitterLicense)
	credits := make([]*credit, 0)
	for _, l := range licenses {
		c := newCredit(l)
		if c != nil {
			credits = append(credits, c)
		}
	}
	return credits, nil
}

func newCredit(text string) *credit {
	l := strings.Split(text, splitterText)
	if len(l) < 2 {
		return nil
	}
	s := strings.Split(strings.Trim(l[0], "\n"), "\n")
	return &credit{
		name: s[0],
		url:  s[1],
		text: l[1],
	}
}

func runGoCredits() (*bytes.Buffer, error) {
	buf := &bytes.Buffer{}
	err := gocredits.Run([]string{"../../"} /* from root */, buf, os.Stderr)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

var template = `// Code generated by gen.go; DO NOT EDIT.

package ui

type credit struct {
	name, url, text string
}

var credits = []*credit{
%s}
`

func run() error {
	flag.Parse()

	credits, err := collect()
	if err != nil {
		return err
	}

	vars := ""
	for _, c := range credits {
		vars += fmt.Sprintf(`	{
		"%s",
		"%s",
		`+"`%s`"+`,
	},
`, c.name, c.url, c.formattedText())
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf(template, vars))
	return os.WriteFile(*output, buf.Bytes(), 0666)
}

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}