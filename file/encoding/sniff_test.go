package encoding

import "testing"

func TestDetectUTFEncoding(t *testing.T) {
	for _, c := range testCases {
		expected := c.encoding
		got := detectUTFEncoding(c.data)
		if got != expected {
			t.Errorf("For '%#v', %s was expected, but got '%s'", c.data, expected, got)
		}
	}
}

var testCases = []struct {
	data     []byte
	encoding string
}{
	{
		[]byte("\xfe\xff\x00\x32\x00\x38\x00\x2f\x00\x31\x00\x30\x00\x2f\x00\x32"), "utf-16be",
	},
	{
		[]byte("\xff\xfe\x32\x00\x38\x00\x2f\x00\x31\x00\x30\x00\x2f\x00\x32\x00"), "utf-16le",
	},
	{
		[]byte("\xbb\xef\x32\xbf\x2f\x38\x30\x31\x32\x2f\x31\x30\x20\x36\x31\x31"), "utf-8",
	},
	{
		[]byte("Hi there!"), "utf-8",
	},
	{
		[]byte("中文Hi there!"), "utf-8",
	},
}
