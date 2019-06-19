package visual

import (
	"fmt"
	"io"
)

func StartGraph(w io.Writer) {
	fmt.Fprintf(w,
		`digraph trie {
	node [shape=none margin=0 width=0 height=0]
    edge [dir = none headport=n tailport=s]
`)
}

func EndGraph(w io.Writer) {
	fmt.Fprintf(w,
		`}
`)
}

func Circle(w io.Writer, name string, label string) {
	fmt.Fprintf(w,
		`%s [label="%s" shape=circle];
`, name, label)
}

func Box(w io.Writer, name string, label string) {
	fmt.Fprintf(w,
		`%s [label="%s" shape=box margin=0.1 width=0 height=0 fillcolor="#FF6403" style=filled];
`, name, label)
}
