package output

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	errorStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("196"))
	successStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("78"))
)

// JSON prints v as indented JSON to stdout.
func JSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Errorf("encoding JSON: %w", err)
	}
	fmt.Println(string(data))
	return nil
}

// Print outputs either JSON or human-friendly text depending on the flag.
func Print(jsonMode bool, v any, humanFn func()) error {
	if jsonMode {
		return JSON(v)
	}
	humanFn()
	return nil
}

// Title prints a styled title line.
func Title(s string) {
	fmt.Println(titleStyle.Render(s))
}

// Field prints a label: value pair.
func Field(label, value string) {
	fmt.Printf("%s %s\n", labelStyle.Render(label+":"), value)
}

// Success prints a success message.
func Success(s string) {
	fmt.Println(successStyle.Render("✓ " + s))
}

// Error prints an error message to stderr.
func Error(s string) {
	fmt.Fprintln(os.Stderr, errorStyle.Render("✗ "+s))
}
