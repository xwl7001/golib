package dirtree

import (
	"log"
	"path/filepath"
)

type File struct {
	Name string
}

type Folder struct {
	Name    string
	Files   []File
	Folders map[string]*Folder
}

func newFolder(name string) *Folder {
	return &Folder{name, []File{}, make(map[string]*Folder)}
}

func (f *Folder) getFolder(name string) *Folder {
	if nextF, ok := f.Folders[name]; ok {
		return nextF
	} else {
		log.Fatalf("Expected nested folder %v in %v\n", name, f.Name)
	}
	return &Folder{} // cannot happen
}

func (f *Folder) addFolder(path []string) {
	for i, segment := range path {
		if i == len(path)-1 { // last segment == new folder
			f.Folders[segment] = newFolder(segment)
		} else {
			f.getFolder(segment).addFolder(path[1:])
		}
	}
}

func (f *Folder) addFile(path []string) {
	for i, segment := range path {
		if i == len(path)-1 { // last segment == file
			f.Files = append(f.Files, File{segment})
		} else {
			f.getFolder(segment).addFile(path[1:])
			return
		}
	}
}

func (f *Folder) String() string {
	var str string
	for _, file := range f.Files {
		str += f.Name + string(filepath.Separator) + file.Name + "\n"
	}
	for _, folder := range f.Folders {
		str += folder.String()
	}
	return str
}
