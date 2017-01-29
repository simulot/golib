package encoding

/*
 This package provides text file parsing capabilities, stolen form html package
*/

// The algorithm uses at most sniffLen bytes to make its decision.
const sniffLen = 32

// DetectUTFEncoding determine the type of stream by checking BOM
func detectUTFEncoding(data []byte) string {
	if len(data) > sniffLen {
		data = data[:sniffLen]
	}

	// Index of the first non-whitespace byte in data.
	firstNonWS := 0
	for ; firstNonWS < len(data) && isWS(data[firstNonWS]); firstNonWS++ {
	}

	for _, sig := range sniffSignatures {
		if ct := sig.match(data, firstNonWS); ct != "" {
			return ct
		}
	}

	return "oops" // fallback
}

// Consume whitespaces
func isWS(b byte) bool {
	switch b {
	case '\t', '\n', '\x0c', '\r', ' ':
		return true
	}
	return false
}

// BOMs
var sniffSignatures = []sniffSig{
	// UTF BOMs.
	&maskedSig{mask: []byte("\xFF\xFF\x00\x00"), pat: []byte("\xFE\xFF\x00\x00"), ct: "utf-16be"},
	&maskedSig{mask: []byte("\xFF\xFF\x00\x00"), pat: []byte("\xFF\xFE\x00\x00"), ct: "utf-16le"},
	&maskedSig{mask: []byte("\xFF\xFF\xFF\x00"), pat: []byte("\xEF\xBB\xBF\x00"), ct: "utf-8-BOM"},

	textSig{}, // should be last

}

type sniffSig interface {
	// match returns the MIME type of the data, or "" if unknown.
	match(data []byte, firstNonWS int) string
}

type maskedSig struct {
	mask, pat []byte
	skipWS    bool
	ct        string
}

func (m *maskedSig) match(data []byte, firstNonWS int) string {
	// pattern matching algorithm section 6
	// https://mimesniff.spec.whatwg.org/#pattern-matching-algorithm

	if m.skipWS {
		data = data[firstNonWS:]
	}
	if len(m.pat) != len(m.mask) {
		return ""
	}
	if len(data) < len(m.mask) {
		return ""
	}
	for i, mask := range m.mask {
		db := data[i] & mask
		if db != m.pat[i] {
			return ""
		}
	}
	return m.ct
}

type textSig struct{}

func (textSig) match(data []byte, firstNonWS int) string {
	// c.f. section 5, step 4.
	for _, b := range data[firstNonWS:] {
		switch {
		case b <= 0x08,
			b == 0x0B,
			0x0E <= b && b <= 0x1A,
			0x1C <= b && b <= 0x1F:
			return ""
		}
	}
	return "utf-8"
}
