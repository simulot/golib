package zipwalker

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"github.com/simulot/golib/file/walker"
)

func init() {
	walker.Register(Open, Matcher)
}

// Matcher returns true when the name is like .zip
func Matcher(name string) bool {
	return strings.ToLower(filepath.Ext(name)) == ".zip"
}

// Zip handles zip archive as a folder
type Zip struct {
	path    string          // archive path
	archive *zip.ReadCloser // Zip reader
	wg      sync.WaitGroup  // Keep track of entry file references, prevent closing Zip before all references are done or closed
}

// Open opens a folder provided by os package
func Open(path string) (walker.Walker, error) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return nil, errors.Wrap(err, "Can't open Zip")
	}
	z := &Zip{
		path:    path,
		archive: archive,
	}
	return z, nil
}

// Close the Ziped archive.
// Return imediatly, close zip after last reference is done or close
func (z *Zip) Close() {
	go func() {
		// Wait that all items have been released using Done()
		z.wg.Wait()
		z.archive.Close()
	}()
}

// Items sends zip file list through the channel
func (z *Zip) Items() chan walker.WalkItem {
	out := make(chan walker.WalkItem)
	go func() {
		for _, item := range z.archive.File {
			info := item.FileHeader.FileInfo()
			if !info.IsDir() {
				z.wg.Add(1) // Remeber we have emitted an Item
				out <- &Item{
					FileInfo: info,
					file:     item,
					path:     filepath.Join(z.path, item.Name),
					zip:      z,
				}
			}
		}
		close(out)
	}()
	return out
}

// Item is an item returned by Folder Scanner. It contains path relative to opening path
type Item struct {
	zip         *Zip          // Zip archive
	file        *zip.File     // Current entry
	os.FileInfo               // Current entry info
	path        string        // file path made by archive path and file path int the archive
	rc          io.ReadCloser // The open reader on the item
}

// FullName returns the item full name relative the folder path used for scanning
func (i *Item) FullName() string {
	return i.path
}

// Reader give a Reader on the archive Item
func (i *Item) Reader() (io.Reader, error) {
	var err error
	i.rc, err = i.file.Open()
	return i.rc, err
}

// Close closes the item's reader, and release it
func (i *Item) Close() {
	if i.rc != nil {
		i.rc.Close()
	}
	i.zip.wg.Done()
}

// String interfer
func (i *Item) String() string {
	return i.path
}
