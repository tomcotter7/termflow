package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/tomcotter7/termflow/internal/ui"
)

func runHelp() {
	fmt.Println("Hello, World")
}

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "help":
			runHelp()
			return
		}
	}
	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
