// Package dirtree provides a way to generate a directory tree.
//
// Example usage:
//
//	tree, err := dirtree.NewTree("/home/me")
//
// I did my best to keep it OS-independent but truth be told I only tested it
// on OS X and Debian Linux so YMMV. You've been warned.
package dirtree

import (
	"os"
	"path/filepath"
	"time"
)

// FileInfo is a struct created from os.FileInfo interface for serialization.
type FileInfo struct {
	Name    string      `json:"name"`
	Size    int64       `json:"size"`
	Mode    os.FileMode `json:"mode"`
	ModTime time.Time   `json:"modTime"`
	IsDir   bool        `json:"isDir"`
}

// Helper function to create a local FileInfo struct from os.FileInfo interface.
func fileInfoFromInterface(v os.FileInfo) *FileInfo {
	return &FileInfo{v.Name(), v.Size(), v.Mode(), v.ModTime(), v.IsDir()}
}

// Node represents a node in a directory tree.
type Node struct {
	FullPath string    `json:"path"`
	Info     *FileInfo `json:"info"`
	Children []*Node   `json:"children"`
	Parent   *Node     `json:"-"`
}

// New Create directory hierarchy.
func New(root string) (result *Node, err error) {
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return
	}
	parents := make(map[string]*Node)
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		parents[path] = &Node{
			FullPath: path,
			Info:     fileInfoFromInterface(info),
			Children: make([]*Node, 0),
		}
		return nil
	}
	if err = filepath.Walk(absRoot, walkFunc); err != nil {
		return
	}
	for path, node := range parents {
		parentPath := filepath.Dir(path)
		parent, exists := parents[parentPath]
		if !exists { // If a parent does not exist, this is the root.
			result = node
		} else {
			node.Parent = parent
			parent.Children = append(parent.Children, node)
		}
	}
	return
}
