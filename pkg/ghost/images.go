package ghost

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type imagesResponse struct {
	Images []ImageUpload `json:"images"`
}

// UploadImage uploads a local image file to Ghost.
func (c *Client) UploadImage(filePath string) (*ImageUpload, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("opening image: %w", err)
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("creating form file: %w", err)
	}

	if _, err := io.Copy(part, file); err != nil {
		return nil, fmt.Errorf("copying file data: %w", err)
	}

	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("closing multipart writer: %w", err)
	}

	data, err := c.PostMultipart("/images/upload/", &buf, writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	var resp imagesResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, fmt.Errorf("parsing image response: %w", err)
	}
	if len(resp.Images) == 0 {
		return nil, fmt.Errorf("no image returned after upload")
	}
	return &resp.Images[0], nil
}
