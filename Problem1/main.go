package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
)

func main() {
	var s string
	fmt.Print("Enter string: ")
	fmt.Fscan(os.Stdin, &s)
	fmt.Println(Counter(s))
}

func Counter(s string) string {
	counter := make(map[rune]int)
	for _, r := range s {
		counter[r] += 1
	}

	runes := make([]int, 0, len(counter))
	for r := range counter {
		runes = append(runes, int(r))
	}
	sort.Ints(runes)

	buf := new(bytes.Buffer)
	for _, r := range runes {
		fmt.Fprintf(buf, "%s%d", string(rune(r)), counter[rune(r)])
	}

	return buf.String()
}
