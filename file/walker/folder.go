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
func Open(path string) (OpenCloser, error) {
	_, err := os.Stat(path)
	if err != nil {
		return nil, errors.Wrap(err, "Can't stat path in Folder.Open")
	}

	f := &Folder{
		path: path,
	}

	return f, nil
}

// Close the folder. On os folder, there is nothing to do
func (f *Folder) Close() {}

// Items send folder content throught the channel
func (f *Folder) Items() chan ItemOpenCloser {
	out := make(chan ItemOpenCloser)
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
						for item := range o.Items() {
							out <- item
						}
						o.Close()
						return nil
					}
				}
			}

			out <- &Item{
				FileInfo: info,
				path:     path,
			}
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
func (f *Item) String() string {
	return f.FullName()
}

// FullName returns the item full name relative the folder path used for scanning
func (f *Item) FullName() string {
	return f.path
}

// Open opens the file pointed by the Folder Item
func (f *Item) Open() (io.ReadCloser, error) {
	var err error
	f.file, err = os.Open(f.FullName())
	return f.file, err
}

// Close the file pointed by the Folder Item
func (f *Item) Close() error {
	return f.file.Close()
}

// Done do nothing in os.File context
func (f *Item) Done() {}
