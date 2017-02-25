package walker

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/simulot/golib/file/encoding"
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
// implements Walker interface
func (f *Folder) Close() {}

// Items send folder content through a channel
// implements Walker
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
			// this is a regular file...
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
	file io.ReadCloser
}

// String implements stringer interface
func (i *Item) String() string {
	return i.path
}

// FullName returns the item full name relative the folder path used for scanning
func (i *Item) FullName() string {
	return i.path
}

// MemberName return the file name
func (i *Item) MemberName() string {
	return i.path
}

// Reader opens the file pointed by the Folder Item
func (i *Item) Reader() (io.Reader, error) {
	var err error
	if i.file != nil {
		panic(i.path + " is already open")
	}
	if i.file, err = os.Open(i.path); err == nil {
		r := encoding.NewReader(i.file)
		return r, err
	}
	return i.file, err
}

// Close the file pointed by the Folder Item whenever it is opened
func (i *Item) Close() {
	if i.file != nil {
		i.file.Close()
	}
	return
}

// Clone Item except the file.
func (i *Item) Clone() WalkItem {
	return &Item{
		FileInfo: i.FileInfo,
		path:     i.path,
	}
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
