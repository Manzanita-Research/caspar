package cmd

import (
	"fmt"

	"github.com/manzanita-research/ghostctl/pkg/output"
	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Manage images",
}

var imageUploadCmd = &cobra.Command{
	Use:   "upload <file>",
	Short: "Upload an image",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := newClient()
		if err != nil {
			return err
		}

		img, err := client.UploadImage(args[0])
		if err != nil {
			return err
		}

		return output.Print(jsonOut, img, func() {
			output.Success(fmt.Sprintf("Uploaded: %s", img.URL))
		})
	},
}

func init() {
	rootCmd.AddCommand(imageCmd)
	imageCmd.AddCommand(imageUploadCmd)
}
