package walker

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Folder handles a classical folder as provided by os file system.
type Folder struct {
	path string
}

// Open opens a folder provided by os package
func Open(path string) (*Folder, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "Can't stat path in Folder.Open")
	}

	f := &Folder{
		path: path,
	}

	return f, nil
}

// Close the folder. For a folder, there is nothing to do
// impelments Walker interface
func (f *Folder) Close() {}

// Items send folder content throught a channel
// impelments Walker
func (f *Folder) Items() chan Walker {
	out := make(chan Walker)
	go func() {
		filepath.Walk(f.path, func(path string, info os.FileInfo, err error) error {
			if info.IsDir() || err != nil {
				return err
			}
			// check if the current file is an registered container
			if len(walkerRegister) > 0 {
				for _, d := range walkerRegister {
					if d.Matcher(path) {
						o, err := d.Opener(path)
						if err != nil {
							return err
						}
						out <- o
						return nil
					}
				}
			}
			o, err := FileWalkerOpen(path)
			if err != nil {
				return err
			}
			out <- o
			return nil
		})
		close(out)
	}()

	return out
}

// Item is an item returned by Folder Scanner. It contains path relative to opening path
type Item struct {
	os.FileInfo
	path string
	file *os.File
}

// String implement stringer interface
func (i *Item) String() string {
	return i.path
}

// FullName returns the item full name relative the folder path used for scanning
func (i *Item) FullName() string {
	return i.path
}

// Reader opens the file pointed by the Folder Item
func (i *Item) Reader() (io.Reader, error) {
	var err error
	i.file, err = os.Open(i.path)
	return i.file, err
}

// Close the file pointed by the Folder Item
func (i *Item) Close() {
	i.file.Close()
	return
}

type FileAsWalker struct {
	path string
}

func FileWalkerOpen(file string) (Walker, error) {
	return &FileAsWalker{
		path: file,
	}, nil
}

func (f *FileAsWalker) Items() chan WalkItem {
	out := make(chan WalkItem)
	go func() {
		info, err := os.Stat(f.path)
		if err == nil {
			out <- &Item{
				FileInfo: info,
				path:     f.path,
			}
		}
		close(out)
	}()
	return out
}
func (f *FileAsWalker) Close() {}
