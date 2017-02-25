// A walker is in charge of opening a file container (folder, zip archive, or file)
// Items is a function that gives all items

package walker

import (
	"io"
	"os"

	"github.com/pkg/errors"
)

// ErrWalkerNotFound error when the given file name can't be open
var ErrWalkerNotFound = errors.New("Not an Opener")

// Opener is the function signature of Walker opener
type Opener func(string) (Walker, error)

// Matcher tells if the file can be open by the opener
type Matcher func(string) bool

var walkerRegister = []struct {
	Opener
	Matcher
}{}

// Register is called by concrete implementations of Walker
func Register(o Opener, m Matcher) {
	walkerRegister = append(walkerRegister, struct {
		Opener
		Matcher
	}{o, m})

}

// Walker interface for archive walker
type Walker interface {
	Close()
	Items() chan WalkItem
}

// WalkItem interface of archive item
type WalkItem interface {
	os.FileInfo                 // Underlaying file structure
	FullName() string           // Give the full path of the file
	Reader() (io.Reader, error) // Give a reader on archive item
	Close()                     // Items must be closed.
	MemberName() string         // When Walkitem is an archive, returns archive member name, otherwise returns file name
	Clone() WalkItem            // Clone item to have more readers on the same Item
}
