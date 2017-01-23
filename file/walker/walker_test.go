package walker

import (
	"bufio"
	"strings"
	"testing"
	// _ "github.com/simulot/golib/file/walker/zip"
)

var testInterfaceCases = []struct {
	path     string
	expected []string
}{
	{
		"test/flat",
		[]string{
			"test/flat/file_a.txt", "test/flat/file_b.txt", "test/flat/file_c.txt", "test/flat/file_d.txt", "test/flat/file_e.txt", "test/flat/file_f.txt",
		},
	},
	{
		"test/tree",
		[]string{
			"test/tree/file_a.txt", "test/tree/file_b.txt", "test/tree/file_c.txt",
			"test/tree/subtree/file_d.txt", "test/tree/subtree/file_e.txt", "test/tree/subtree/file_f.txt",
		},
	},
	// {
	// 	"test/zip/flat.zip",
	// 	[]string{
	// 		"test/zip/flat.zip/file_a.txt", "test/zip/flat.zip/file_b.txt", "test/zip/flat.zip/file_c.txt",
	// 		"test/zip/flat.zip/file_d.txt", "test/zip/flat.zip/file_e.txt", "test/zip/flat.zip/file_f.txt",
	// 	},
	// },
	// {
	// 	"test/zip/tree.zip",
	// 	[]string{
	// 		"test/zip/tree.zip/file_a.txt", "test/zip/tree.zip/file_b.txt", "test/zip/tree.zip/file_c.txt",
	// 		"test/zip/tree.zip/subtree/file_d.txt", "test/zip/tree.zip/subtree/file_e.txt", "test/zip/tree.zip/subtree/file_f.txt",
	// 	},
	// },
}

func TestOpen(t *testing.T) {
	for _, c := range testInterfaceCases {
		t.Run(c.path, func(t *testing.T) {
			folder, err := Open(c.path)
			if err != nil {
				t.Fatalf("Unexpected error %s", err)
			}
			for item := range folder.Items() {
				reader, err := item.Reader()
				if err != nil {
					t.Errorf("Unexpected error when opening '%s' from '%s'", item.FullName(), c.path)
					return
				}
				content, err := bufio.NewReader(reader).ReadString('\n')
				content = strings.TrimRight(content, "\n")
				if content != item.Name() {
					t.Errorf("Expected content of '%s' file to by '%s', but got '%s'!", item.Name(), item.Name(), content)
				}
				item.Close()
			}
			folder.Close()
		})
	}
}
