package goldmark

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/yuin/goldmark/util"
	"os"
	"strconv"
	"strings"
	testing "testing"
)

type MarkdownTestCase struct {
	No       int
	Markdown string
	Expected string
}

const attributeSeparator = "//- - - - - - - - -//"
const caseSeparator = "//= = = = = = = = = = = = = = = = = = = = = = = =//"

func DoTestCaseFile(m Markdown, filename string, t *testing.T) {
	fp, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)
	c := MarkdownTestCase{
		No:       -1,
		Markdown: "",
		Expected: "",
	}
	cases := []MarkdownTestCase{}
	line := 0
	for scanner.Scan() {
		line++
		if util.IsBlank([]byte(scanner.Text())) {
			continue
		}
		c.No, err = strconv.Atoi(scanner.Text())
		if err != nil {
			panic(fmt.Sprintf("%s: invalid case No at line %d", filename, line))
		}
		if !scanner.Scan() {
			panic(fmt.Sprintf("%s: invalid case at line %d", filename, line))
		}
		line++
		if scanner.Text() != attributeSeparator {
			panic(fmt.Sprintf("%s: invalid separator '%s' at line %d", filename, scanner.Text(), line))
		}
		buf := []string{}
		for scanner.Scan() {
			line++
			text := scanner.Text()
			if text == attributeSeparator {
				break
			}
			buf = append(buf, text)
		}
		c.Markdown = strings.Join(buf, "\n")
		buf = []string{}
		for scanner.Scan() {
			line++
			text := scanner.Text()
			if text == caseSeparator {
				break
			}
			buf = append(buf, text)
		}
		c.Expected = strings.Join(buf, "\n")
		cases = append(cases, c)
	}
	DoTestCases(m, cases, t)
}

func DoTestCases(m Markdown, cases []MarkdownTestCase, t *testing.T) {
	for _, testCase := range cases {
		var out bytes.Buffer
		if err := m.Convert([]byte(testCase.Markdown), &out); err != nil {
			panic(err)
		}
		if !bytes.Equal(bytes.TrimSpace(out.Bytes()), bytes.TrimSpace([]byte(testCase.Expected))) {
			format := `============= case %d ================
Markdown:
-----------
%s

Expected:
----------
%s

Actual
---------
%s
`
			t.Errorf(format, testCase.No, testCase.Markdown, testCase.Expected, out.Bytes())
		}

	}
}
