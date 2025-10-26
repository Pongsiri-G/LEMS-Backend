package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/xuri/excelize/v2"
)

func SaveAtDownload(f *excelize.File, fileName string) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("Error getting user home directory: %v\n", err)
		return err
	}

	var downloadsPath string
	switch runtime.GOOS {
	case "windows":
		downloadsPath = filepath.Join(homeDir, "Downloads")
	case "darwin", "linux":
		downloadsPath = filepath.Join(homeDir, "Downloads")
	default:
		return fmt.Errorf("unsupported operating system: %s", runtime.GOOS)
	}

	if err := os.MkdirAll(downloadsPath, 0755); err != nil {
		fmt.Printf("Error creating Downloads directory: %v\n", err)
		return err
	}

	fullPath := filepath.Join(downloadsPath, fileName)
	if err := f.SaveAs(fullPath); err != nil {
		fmt.Printf("Error saving file to %s: %v\n", fullPath, err)
		return err
	}

	return nil
}
