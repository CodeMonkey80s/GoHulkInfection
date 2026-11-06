package main

import (
	"fmt"
	"os"

	"hulkinfection/internal/game"
)

func main() {
	if err := game.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "hulk infection:", err)
		os.Exit(1)
	}
}
