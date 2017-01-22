package zip

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"sync"

	"strings"

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

// Folder handles zip archive as a folder
type Folder struct {
	path    string          // archive path
	archive *zip.ReadCloser // Zip reader
	wg      sync.WaitGroup  // Keep track of entry file references, prevent closing Zip before all references are done or closed
}

// Open opens a folder provided by os package
func Open(path string) (walker.OpenCloser, error) {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return nil, errors.Wrap(err, "Can't open archive in OpenFolder")
	}
	f := &Folder{
		path:    path,
		archive: archive,
	}
	return f, nil
}

// Close the Ziped archive.
// Return imediatly, close zip after last reference is done or close
func (f *Folder) Close() {
	go func() {
		f.wg.Wait()
		f.archive.Close()
	}()
}

// Items sends zip file list through the channel
func (f *Folder) Items() chan walker.ItemOpenCloser {
	out := make(chan walker.ItemOpenCloser)
	go func() {
		for _, zf := range f.archive.File {
			fi := zf.FileInfo()
			if !fi.IsDir() {
				f.wg.Add(1)
				out <- &FolderItem{
					FileInfo:  fi,
					fileEntry: zf,
					path:      filepath.Join(f.path, zf.Name),
					zip:       f,
				}
			}
		}
		close(out)
	}()
	return out
}

// FolderItem is an item returned by Folder Scanner. It contains path relative to opening path
type FolderItem struct {
	zip         *Folder   // Zip archive
	os.FileInfo           // Current entry info
	fileEntry   *zip.File // Current entry
	path        string    // file path made by archive path and file path int the archive
}

// File is the equivalent of os.File for zip archive
type File struct {
	io.ReadCloser // When the file is open
}

// String interfer
func (f *FolderItem) String() string {
	return f.path
}

// FullName returns the item full name relative the folder path used for scanning
func (f *FolderItem) FullName() string {
	return f.path
}

// Open opens the file pointed by the Item
func (f *FolderItem) Open() (io.ReadCloser, error) {
	var err error
	zf := File{}
	zf.ReadCloser, err = f.fileEntry.Open()
	return zf, err
}

// Close closes de archived file
func (zf File) Close() error {
	err := zf.ReadCloser.Close()
	return err
}

// Close the file pointed by the Folder Item
func (f *FolderItem) Close() error {
	return nil
	// return f.file.Close()
}

// Done when when can release the zip entry reference
func (f *FolderItem) Done() {
	f.zip.wg.Done()
}
