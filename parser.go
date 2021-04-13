package main

import (
	"context"
	"fmt"
	"io"
	"regexp"
	"strings"
)

func parseToLines(reader io.Reader, res chan<- string) {
	buf := make([]byte, 1024)
	in := make([]byte, 0, 1024)
	eof := false

	defer close(res)

	for !eof {
		n, err := reader.Read(buf)
		if n == 0 || err == io.EOF {
			eof = true
		}
		in = append(in, buf[:n]...)
		n2, lines := getLines(in)
		in = in[n2:]
		for _, s := range lines {
			res <- s
		}
	}
	n2, lines := getLines(in)
	for _, s := range lines {
		res <- s
	}
	in = in[n2:]
	lastline := string(in)
	if trimed := strings.TrimSpace(lastline); len(trimed) > 0 {
		fmt.Printf("[%s]\n", trimed)
		res <- trimed
	}

}
func getLines(in []byte) (int, []string) {
	ret := make([]string, 0)
	i := 0
	j := 0
	for j < len(in) {
		if in[j] == '\n' {
			if line := strings.TrimSpace(string(in[i:j])); len(line) > 0 {
				ret = append(ret, line)
				i = j + 1
			}
		}
		j++
	}
	return i, ret
}

type workExpression struct {
	reg        *regexp.Regexp
	expression Expression
}

func processLines(ctx context.Context, exp []workExpression, data <-chan string, done chan<- bool) {

	for line := range data {
		for i := 0; i < len(exp); i++ {
			m := exp[i].reg.FindStringSubmatch(line)
			if m != nil {
				exp[i].expression.Value = m[1]
			}
		}
	}
	done <- true
}
func compilePattern(in []Expression) ([]workExpression, error) {
	ret := make([]workExpression, 0, len(in))
	for i, e := range in {
		if reg, err := regexp.Compile(e.Pat); err != nil {
			return nil, err
		} else {
			ret = append(ret, workExpression{reg: reg, expression: in[i]})
		}
	}
	return ret, nil
}
