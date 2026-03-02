package cmd

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"

	"github.com/manzanita-research/caspar/pkg/output"
	"github.com/spf13/cobra"
)

//go:embed skill.md
var skillContent []byte

var skillCmd = &cobra.Command{
	Use:   "skill",
	Short: "Manage the caspar skill for Claude Code",
}

var skillInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install the caspar skill for Claude Code",
	Long:  "Copies the built-in SKILL.md to ~/.claude/skills/caspar/ so agents can discover it.",
	RunE:  runSkillInstall,
}

func init() {
	rootCmd.AddCommand(skillCmd)
	skillCmd.AddCommand(skillInstallCmd)
}

func runSkillInstall(cmd *cobra.Command, args []string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("finding home directory: %w", err)
	}

	dir := filepath.Join(home, ".claude", "skills", "caspar")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating skill directory: %w", err)
	}

	dest := filepath.Join(dir, "SKILL.md")
	if err := os.WriteFile(dest, skillContent, 0644); err != nil {
		return fmt.Errorf("writing skill file: %w", err)
	}

	if jsonOut {
		return output.JSON(map[string]string{
			"status": "installed",
			"path":   dest,
		})
	}

	output.Success("Installed caspar skill")
	output.Field("Path", dest)
	return nil
}
