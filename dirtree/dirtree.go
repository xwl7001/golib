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
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
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
	workingDir string           `json:"-"`
	FullPath   string           `json:"path"`
	Info       *FileInfo        `json:"-"`
	Children   map[string]*Node `json:"children"`
	Parent     *Node            `json:"-"`
}

func (n *Node) addFolder(fi os.FileInfo) {
	name := fi.Name()
	if IsHidden(name) {
		return
	}
	FullPath := name
	if n.FullPath != "" {
		FullPath = filepath.Join(n.FullPath, name)
	}
	// fmt.Println(FullPath)

	dir := &Node{
		workingDir: n.workingDir,
		FullPath:   FullPath,
		Info:       fileInfoFromInterface(fi),
		Children:   make(map[string]*Node),
		Parent:     n,
	}

	files, err := ioutil.ReadDir(FullPath)
	if err != nil {
		log.Fatalln(err)
		return
	}
	for _, fi := range files {
		if fi.IsDir() {
			dir.addFolder(fi)
		} else {
			dir.addFile(fi)
		}
	}
	n.Children[name] = dir

}

func (n *Node) addFile(fi os.FileInfo) {
	name := fi.Name()
	if IsHidden(name) {
		return
	}
	FullPath := name
	if n.Parent != nil {
		FullPath = filepath.Join(n.Parent.FullPath, name)
	}
	dir := &Node{
		workingDir: n.workingDir,
		FullPath:   FullPath,
		Info:       fileInfoFromInterface(fi),
		Children:   make(map[string]*Node),
		Parent:     n,
	}
	n.Children[name] = dir
}

func (n *Node) getFolder(path string) *Node {
	path = strings.TrimPrefix(path, string(os.PathSeparator))
	segments := SplitPath(path)
	if len(segments) == 1 {
		return n.Children[segments[0]]
	}
	return n.Children[segments[0]].getFolder(strings.Join(segments[1:], string(os.PathSeparator)))
}

func (n *Node) Size() int64 {
	size := n.Info.Size
	for _, node := range n.Children {
		size += node.Size()
	}
	return size
}

func (n *Node) String() string {
	res := make([]string, 0)
	res = append(res, strings.TrimPrefix(n.FullPath, n.workingDir))
	for _, node := range n.Children {
		res = append(res, node.String())
	}
	return strings.Join(res, "\n")
}

// Name returns node's folder/file name
func (n *Node) Name() string {
	return NodeName(n.FullPath)
}

// NodeName returns a path's folder/file name
func NodeName(path string, separator ...string) string {
	segments := SplitPath(path, separator...)
	return segments[len(segments)-1]
}

// SplitPath get all segments form given path
func SplitPath(path string, separator ...string) []string {
	sep := string(os.PathSeparator)
	if len(separator) == 1 {
		sep = separator[0]
	}

	return strings.Split(path, sep)
}

func IsHidden(name string) bool {
	return strings.HasPrefix(name, ".")
}

// New Create directory hierarchy.
func New(root string) (*Node, error) {
	absRoot, err := filepath.Abs(root)
	log.Println("absRoot:", absRoot)
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(absRoot)
	if err != nil {
		return nil, err
	}
	cleanRoot := path.Clean(root)
	log.Println("clen root:", cleanRoot)
	res := &Node{
		workingDir: cleanRoot,
		FullPath:   cleanRoot,
		Children:   make(map[string]*Node),
	}
	for _, fi := range files {
		if fi.IsDir() {
			res.addFolder(fi)
		} else {
			res.addFile(fi)
		}
	}
	return res, nil
}
