package main

import (
	"fmt"
	"os"
	"strings"

	"golang.design/x/clipboard"
)

func main() {
	var args []string
	if len(os.Args) < 2 {
		if err := clipboard.Init(); err == nil {
			content := strings.TrimSpace(string(clipboard.Read(clipboard.FmtText)))
			if len(content) == 0 {
				fmt.Fprintf(os.Stderr, "Usage: %s <query>\n", os.Args[0])
				os.Exit(1)
			}

			space := strings.IndexRune(content, ' ')
			if space > 0 {
				cmd := content[:space]
				rest := content[space+1:]
				args = []string{cmd, rest}
			} else {
				args = []string{content}
			}
		}
	} else {
		args = os.Args[1:]
	}

	fmt.Println("<", args)

	// try explit command
	maybeCmd := args[0]
	if h, ok := handlers[maybeCmd]; ok {
		query := strings.Join(args[1:], " ")
		err := h.Run(query)
		if err != nil {
			os.Exit(1)
		}
		return
	}

	// try sniffing
	query := strings.Join(args, " ")
	for _, h := range handlers {
		if h.Check(query) {
			err := h.Run(query)
			if err == nil {
				break
			}
		}
	}
}

var handlers = map[string]handler{
	"hex":  hex2Dec,
	"eval": goval,
	"uuid": uuid,
}

type handler interface {
	Check(string) bool
	Run(string) error
}
