package fileUtils

import (
	"io"
	"os"
	"path/filepath"
)

// ReadFile reads the contents of a file and returns it as a byte slice.
func ReadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return data, nil
}

// WriteFile writes data to a file. If the file does not exist, it will be created.
// If the file already exists, its contents will be replaced.
func WriteFile(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

// CheckFile checks if a file exists
func CheckFile(filename string) bool {
	// Check if file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return false
	}
	return true
}

// ListFiles returns a list of files in a directory.
func ListFiles(dir string) ([]string, error) {
	var files []string
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}

// CheckDirectory checks if directory exists or not.
func CheckDirectory(directory string) bool {
	// Check if directory exists
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateDir creates a directory and any necessary parent directories.
func CreateDir(dir string) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return err
	}
	return nil
}

// RemoveFile removes a file.
func RemoveFile(filename string) error {
	err := os.Remove(filename)
	if err != nil {
		return err
	}
	return nil
}

// RemoveDir removes a directory and all its contents.
func RemoveDir(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		return err
	}
	return nil
}
