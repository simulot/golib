package encoding

import (
	"bufio"
	"io"

	"github.com/pkg/errors"

	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

// NewReader takes a reader of utf-16, utf-8 and return a reader utf-8 ready for golang work
func NewReader(r io.Reader) io.Reader {
	var err error
	var buf []byte
	o := bufio.NewReader(r)

	// Get few bytes to test the type of the input file
	buf, err = o.Peek(5)
	if err != nil && err != io.EOF {
		logger.Printf("%s\n", errors.Wrap(err, "can't read buffer in reader.NewReader"))
		return o
	}

	switch detectUTFEncoding(buf) {
	case "utf-16be":
		return transform.NewReader(o, unicode.UTF16(unicode.BigEndian, unicode.UseBOM).NewDecoder())
	case "utf-16le":
		return transform.NewReader(o, unicode.UTF16(unicode.LittleEndian, unicode.UseBOM).NewDecoder())
	case "utf-8-BOM":
		// Discard BOM
		o.Discard(3)
		return o
	}
	return o
}
