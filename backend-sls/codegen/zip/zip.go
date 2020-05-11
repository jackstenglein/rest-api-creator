// Package zip provides functions to zip directories and unzip archives.
package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jackstenglein/rest_api_creator/backend-sls/errors"
)

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func Unzip(src string, dest string) error {

	r, err := zip.OpenReader(src)
	if err != nil {
		return errors.Wrap(err, "Failed to open zip reader")
	}
	defer r.Close()

	for _, f := range r.File {

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: https://snyk.io/research/zip-slip-vulnerability#go
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		// Make File
		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return errors.Wrap(err, "Failed to make directory")
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return errors.Wrap(err, "Failed to open file")
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return errors.Wrap(err, "Failed to open file")
		}

		_, err = io.Copy(outFile, rc)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		rc.Close()

		if err != nil {
			return errors.Wrap(err, "Failed to copy file")
		}
	}
	return nil
}

// Zip compresses the given directory into a single zip archive file. output is the desired path
// to the new zip file. directory is the path to the directory to zip.
func Zip(outputPath, dirPath string) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return errors.Wrap(err, "Failed to create output file")
	}
	defer outputFile.Close()

	zipWriter := zip.NewWriter(outputFile)
	defer zipWriter.Close()

	err = filepath.Walk(dirPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrap(err, "Failed to walk filepath")
		}
		if info.IsDir() {
			return nil
		}

		relPath := strings.TrimPrefix(filePath, filepath.Dir(dirPath))
		zipFile, err := zipWriter.Create(relPath)
		if err != nil {
			return errors.Wrap(err, "Failed zipWriter create")
		}

		fsFile, err := os.Open(filePath)
		if err != nil {
			return errors.Wrap(err, "Failed to open file to zip")
		}

		_, err = io.Copy(zipFile, fsFile)
		fsFile.Close()

		return errors.Wrap(err, "Failed to copy to zip file")
	})

	return err
}
