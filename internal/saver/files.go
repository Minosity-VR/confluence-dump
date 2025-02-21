package saver

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Minosity-VR/confdump/internal/client"
)

type FileSaver struct {
	rootPath string
}

// NewFileSaver returns a new FileSaver
func NewFileSaver(rootPath string) *FileSaver {
	return &FileSaver{rootPath: rootPath}
}

func (fs *FileSaver) StartSaver(wg *sync.WaitGroup, saveChan <-chan client.ConfluencePage, errChan chan<- error) {
	go func() {
		for req := range saveChan {
			if err := SaveConfluencePage(fs.rootPath, req); err != nil {
				errChan <- err
			}
		}
		close(errChan)
	}()
}

// SaveConfluencePage saves a confluence page to a file. It expects the page.Body.Storage.Value to\
// be an html string. The page title's spaces are replaced with underscores.
// The architecture it the following:
//
//	.rootPath
//	├── [page.SpaceId]
//	│   ├── [page.title].html // Html content
//	│   └── [page.title].json // Dump the page struct
//	└── [page.SpaceId]
func SaveConfluencePage(rootPath string, page client.ConfluencePage) error {
	spaceId := page.SpaceId
	if spaceId == "" {
		spaceId = "unknown-space"
	}

	// Create the directory for the space if it doesn't exist
	spaceDir := filepath.Join(rootPath, spaceId)
	if err := os.MkdirAll(spaceDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for space %s: %w", page.Id, err)
	}

	// Replace spaces and `/` in the title with underscores
	displayTitle := strings.ReplaceAll(page.Title, " ", "_")
	displayTitle = strings.ReplaceAll(displayTitle, "/", "_")

	// Create the HTML file for the page content
	htmlFilePath := filepath.Join(spaceDir, fmt.Sprintf("%s.html", displayTitle))
	if err := os.WriteFile(htmlFilePath, []byte(page.Body.Storage.Value), 0644); err != nil {
		return fmt.Errorf("failed to write HTML file for page %s: %w", page.Title, err)
	}

	// Create the JSON file for the page struct
	jsonFilePath := filepath.Join(spaceDir, fmt.Sprintf("%s.json", displayTitle))
	pageData, err := json.MarshalIndent(page, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal page data for page %s: %w", page.Title, err)
	}
	if err := os.WriteFile(jsonFilePath, pageData, 0644); err != nil {
		return fmt.Errorf("failed to write JSON file for page %s: %w", page.Title, err)
	}

	return nil
}
