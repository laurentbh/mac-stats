package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseToLines(t *testing.T) {

	input := ` this is first line


			this is second line            
	last     `

	res := make(chan string)

	go parseToLines(strings.NewReader(input), res)

	for line := range res {
		fmt.Printf("got [%s]\n", line)
		assert.Contains(t, []string{"this is first line", "this is second line", "last"}, line)
	}
}
func TestBadPattern(t *testing.T) {

	exps := []Expression{
		{Pat: "^(.*)$", Key: "dummy"},
		{Pat: "^($", Key: "dummy"},
	}

	_, err := compilePattern(exps)

	assert.NotNil(t, err)
}

func TestProcess(t *testing.T) {

	data := []string{
		"Health Information:",
		"Cycle Count: 44",
		"Condition: Normal",
		"Battery Installed: Yes",
		"Full Charge Capacity (mAh): 4479",
	}
	exp := []Expression{
		{Pat: "^Cycle Count: (.*)$", Key: "cycle"},
		{Pat: "^Full Charge Capacity \\(mAh\\): (.*)$", Key: "full"},
	}

	work, err := compilePattern(exp)
	assert.Nil(t, err)
	c := make(chan string)
	d := make(chan int)

	go processLines(context.TODO(), work, c, d)

	for _, s := range data {
		c <- s
	}
	close(c)
	ok := <-d

	assert.Equal(t, 0, ok)
	assert.EqualValues(t, "44", work[0].expression.Value)
	assert.EqualValues(t, "4479", work[1].expression.Value)
}
func TestMissingRegexp(t *testing.T) {

	data := []string{
		"Health Information:",
		"Cycle Count: 44",
		"Condition: Normal",
		"Battery Installed: Yes",
	}
	exp := []Expression{
		{Pat: "^Cycle Count: (.*)$", Key: "cycle"},
		{Pat: "^Full Charge Capacity \\(mAh\\): (.*)$", Key: "full"},
	}

	work, err := compilePattern(exp)
	assert.Nil(t, err)
	c := make(chan string)
	d := make(chan int)

	go processLines(context.TODO(), work, c, d)

	for _, s := range data {
		c <- s
	}
	close(c)
	ok := <-d

	assert.Equal(t, 1, ok)
}

func TestProcessFile(t *testing.T) {

	f, err := os.Open("fixtures/battery.txt")

	assert.Nil(t, err)

	defer f.Close()

	reader := bufio.NewReader(f)

	data := make(chan string, 1000)
	done := make(chan bool)
	parseToLines(reader, data)

	<-done
	// TOD: assert something :D
}
