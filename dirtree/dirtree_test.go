package dirtree

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

func matchRecursively(node *Node, query string) {
	if strings.Contains(node.FullPath, query) {
		fmt.Println(node.FullPath)
	}
	for _, child := range node.Children {
		matchRecursively(child, query)
	}
}

func DirTree(t *testing.T) {
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("error getting working directory: %v\n", err)
	}
	root, err := New(dir)
	if err != nil {
		t.Fatalf("error traversing the filesystem: %v\n", err)
	}
	// matchRecursively(root, "data")
	fmt.Println(root.String())
}

func TestDirTree2(t *testing.T) {
	root, err := New("../../../")
	if err != nil {
		t.Fatalf("error traversing the filesystem: %v\n", err)
	}
	// matchRecursively(root, "data")
	// fmt.Println(root.String())
	fmt.Println(root.getFolder("/lifei6671").Size())
	// js, _ := json.Marshal(root)
	// fmt.Println(string(js))
}

// func TestDirTree(t *testing.T) {
// 	startPath := "."
// 	rootFolder := newFolder(startPath)

// 	visit := func(path string, info os.FileInfo, err error) error {
// 		segments := strings.Split(path, string(filepath.Separator))
// 		if info.IsDir() {
// 			if path != startPath {
// 				rootFolder.addFolder(segments)
// 			}
// 		} else {
// 			rootFolder.addFile(segments)
// 		}
// 		return nil
// 	}

// 	err := filepath.Walk(startPath, visit)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	log.Printf("%v\n", rootFolder)
// }
