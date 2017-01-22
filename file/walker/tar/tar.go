package tar

import (
	"archive/tar"
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
	return strings.ToLower(filepath.Ext(name)) == ".tar"
}

// Folder handles zip archive as a folder
type Folder struct {
	r       io.ReadCloser  // reader
	path    string         // archive path
	archive *tar.Reader    // Tar reader
	wg      sync.WaitGroup // Keep track of entry file references, prevent closing tar file before all references are done or closed
}

// Open opens a folder provided by os package
func Open(path string) (walker.OpenCloser, error) {
	r, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "Can't open TarWalker")
	}
	archive := tar.NewReader(r)
	folder := &Folder{
		r:       r,
		path:    path,
		archive: archive,
	}
	return folder, nil
}

// Close the Ziped archive.
// Return imediatly, close zip after last reference is done or close
func (f *Folder) Close() {
	go func() {
		f.wg.Wait()
		f.r.Close()
	}()
}

// Items sends zip file list through the channel
func (f *Folder) Items() chan walker.ItemOpenCloser {
	out := make(chan walker.ItemOpenCloser)
	go func() {
		for {
			hdr, err := f.archive.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				panic(errors.Wrap(err, "Can't walk through tar file"))
			}

			info := hdr.FileInfo()
			if info.IsDir() {
				continue
			}

			f.wg.Add(1)
			out <- &FolderItem{
				archive:  f,
				entry:    &File{Reader: f.archive},
				FileInfo: info,
				path:     filepath.Join(f.path, info.Name()),
			}

		}
		close(out)
	}()
	return out
}

// FolderItem is an item returned by Folder Scanner. It contains path relative to opening path
type FolderItem struct {
	archive     *Folder // Tar archive
	entry       *File
	os.FileInfo        // Current entry info
	path        string // file path made by archive path and file path int the archive
}

func (f *Folder) Open() (io.ReadCloser, error) {
	return f.entry, nil
}

// File is the equivalent of os.File for zip archive
type File struct {
	io.Reader // When the file is open
}

func (f *File) Close() error { return nil }

// String interface
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
