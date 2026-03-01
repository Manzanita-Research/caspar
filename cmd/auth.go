package cmd

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/manzanita-research/caspar/pkg/config"
	"github.com/manzanita-research/caspar/pkg/ghost"
	"github.com/manzanita-research/caspar/pkg/output"
	"github.com/spf13/cobra"
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Manage authentication",
}

var authLoginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to a Ghost site",
	Long:  "Prompts for your Ghost site URL and admin API key, validates the connection, and saves config.",
	RunE:  runAuthLogin,
}

var authStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check authentication status",
	RunE:  runAuthStatus,
}

var authLogoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Remove saved credentials",
	RunE:  runAuthLogout,
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authLoginCmd)
	authCmd.AddCommand(authStatusCmd)
	authCmd.AddCommand(authLogoutCmd)
}

func runAuthLogin(cmd *cobra.Command, args []string) error {
	var siteURL, apiKey string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Ghost site URL").
				Description("e.g. https://your-site.ghost.io").
				Value(&siteURL).
				Validate(func(s string) error {
					s = strings.TrimSpace(s)
					if s == "" {
						return fmt.Errorf("URL is required")
					}
					if !strings.HasPrefix(s, "http://") && !strings.HasPrefix(s, "https://") {
						return fmt.Errorf("URL must start with http:// or https://")
					}
					return nil
				}),
			huh.NewInput().
				Title("Admin API key").
				Description("Find this in Ghost Admin → Settings → Integrations").
				Value(&apiKey).
				Validate(func(s string) error {
					return ghost.ValidateKeyFormat(strings.TrimSpace(s))
				}),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("login cancelled")
	}

	siteURL = strings.TrimSpace(siteURL)
	apiKey = strings.TrimSpace(apiKey)

	// validate by fetching site info
	client := ghost.NewClient(siteURL, apiKey)
	site, err := client.GetSite()
	if err != nil {
		output.Error("Authentication failed: " + err.Error())
		return fmt.Errorf("could not connect to Ghost — check your URL and API key")
	}

	// save config
	cfg := &config.Config{
		URL:         siteURL,
		AdminAPIKey: apiKey,
	}
	if err := config.Save(cfg); err != nil {
		return fmt.Errorf("saving config: %w", err)
	}

	if jsonOut {
		return output.JSON(map[string]any{
			"status": "authenticated",
			"site":   site,
		})
	}

	output.Success("Logged in to " + site.Title)
	output.Field("URL", site.URL)
	output.Field("Version", site.Version)
	return nil
}

func runAuthStatus(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		if jsonOut {
			return output.JSON(map[string]string{"status": "not authenticated"})
		}
		output.Error("Not logged in — run `caspar auth login`")
		return nil
	}

	client := ghost.NewClient(cfg.URL, cfg.AdminAPIKey)
	site, err := client.GetSite()
	if err != nil {
		if jsonOut {
			return output.JSON(map[string]any{
				"status": "error",
				"error":  err.Error(),
			})
		}
		output.Error("Config exists but connection failed: " + err.Error())
		return nil
	}

	if jsonOut {
		return output.JSON(map[string]any{
			"status": "authenticated",
			"site":   site,
		})
	}

	output.Success("Connected to " + site.Title)
	output.Field("URL", site.URL)
	output.Field("Description", site.Description)
	output.Field("Version", site.Version)
	return nil
}

func runAuthLogout(cmd *cobra.Command, args []string) error {
	if err := config.Delete(); err != nil {
		return err
	}

	if jsonOut {
		return output.JSON(map[string]string{"status": "logged out"})
	}

	output.Success("Logged out")
	return nil
}
