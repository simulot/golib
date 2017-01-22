package walker

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

// ErrWalkerNotFound error when the given file name can't be open
var ErrWalkerNotFound = errors.New("Not an Opener")

// Opener is the abstract Opener for a walker
type Opener func(string) (OpenCloser, error)

// Matcher tells if the file can be open by the opener
type Matcher func(string) bool

var walkerRegister = []struct {
	Opener
	Matcher
}{}

// Register is called by concrete impleations of Walker
func Register(o Opener, m Matcher) {
	walkerRegister = append(walkerRegister, struct {
		Opener
		Matcher
	}{o, m})

}

// This files provides interfaces allowing  folder abstraction

// OpenCloser interface
type OpenCloser interface {
	Close()
	Items() chan ItemOpenCloser
}

// ItemOpenCloser interface
type ItemOpenCloser interface {
	os.FileInfo                   //Underlaying file structure
	Open() (io.ReadCloser, error) // File opener
	FullName() string             // Give the full path of the file
	Close() error                 // File closer
	Done()                        // When the file belong to an arckive, Done means we don't need to access to this file.
}
