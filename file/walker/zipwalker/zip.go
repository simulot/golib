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

// init registers zip walker into walkers.
func init() {
	walker.Register(Open, Matcher)
}

// Matcher returns true when the name is like .zip. Used to recognize
// the kind of Walker to be open.
func Matcher(name string) bool {
	return strings.ToLower(filepath.Ext(name)) == ".zip"
}

// Zip handles zip archive as a Walker. This provide a common way
// to walk through the ZIP content, opening, closing ZIP items.
type Zip struct {
	path    string          // archive path
	archive *zip.ReadCloser // Zip reader
	wg      sync.WaitGroup  // Keep track of entry file references, prevent closing Zip before all references are done or closed
}

// Open opens a ZIP archive at path.
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

// Close the Zipped archive.
// It returns immediately but creates a goroutine that waits closing of
// all ZIP items. Then the ZIP archive itself is close.
func (z *Zip) Close() {
	go func() {
		// Wait that all zip items have been released using Close()
		z.wg.Wait()
		z.archive.Close()
	}()
}

// Items sends zip files through the channel
func (z *Zip) Items() chan walker.WalkItem {
	out := make(chan walker.WalkItem)
	go func() {
		for _, item := range z.archive.File {
			info := item.FileHeader.FileInfo()
			if !info.IsDir() {
				z.wg.Add(1) // Remember that we have emitted an Item
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

// Item is an item returned by zip.Items.
type Item struct {
	zip         *Zip          // Zip archive
	file        *zip.File     // Item entry
	os.FileInfo               // Current entry info
	path        string        // file path made by archive path and file path int the archive
	rc          io.ReadCloser // The opened reader on the item
}

// FullName returns the item full name relative the folder path used for scanning
func (i *Item) FullName() string {
	return i.path
}

// Reader give a Reader on the archive Item
// The reader will be closed when calling Close().
func (i *Item) Reader() (io.Reader, error) {
	var err error
	if i.rc != nil {
		panic("i.rc not nil at zip.Item.Reader")
	}
	if i.rc, err = i.file.Open(); err == nil {
		// r := encoding.NewReader(i.rc) // translate UTF16 to UTF8
		return i.rc, nil
	}
	return i.rc, err
}

// Close closes the item's reader, and release it.
// Closing the zip item permits to close the ZIP container when
// all items have been closed.
func (i *Item) Close() {
	if i.rc != nil {
		i.rc.Close()
	}
	i.zip.wg.Done() // Release the item
}

// String returns the full path
func (i *Item) String() string {
	return i.path
}
