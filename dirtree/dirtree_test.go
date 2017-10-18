package dirtree

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDirTree(t *testing.T) {
	startPath := "."
	rootFolder := newFolder(startPath)

	visit := func(path string, info os.FileInfo, err error) error {
		segments := strings.Split(path, string(filepath.Separator))
		if info.IsDir() {
			if path != startPath {
				rootFolder.addFolder(segments)
			}
		} else {
			rootFolder.addFile(segments)
		}
		return nil
	}

	err := filepath.Walk(startPath, visit)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%v\n", rootFolder)
}
