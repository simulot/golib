package zipwalker

import (
	"bufio"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestOpenZipFolder(t *testing.T) {

	folder, err := Open("nowhere")
	if err == nil {
		t.Errorf("Expected error, but got '%s'", err)
		return
	}

	folder, err = Open("test/flat.zip")
	if err != nil {
		t.Errorf("Unexpected error '%s' when opening 'testfiles' folder", err)
		return
	}
	folder.Close()

}

var testZipFolderCases = []struct {
	path     string
	expected []string
}{
	{
		"test/flat.zip",
		[]string{
			"test/flat.zip/file_a.txt", "test/flat.zip/file_b.txt", "test/flat.zip/file_c.txt",
			"test/flat.zip/file_d.txt", "test/flat.zip/file_e.txt", "test/flat.zip/file_f.txt",
		},
	},
	{
		"test/tree.zip",
		[]string{
			"test/tree.zip/file_a.txt", "test/tree.zip/file_b.txt", "test/tree.zip/file_c.txt",
			"test/tree.zip/subtree/file_d.txt", "test/tree.zip/subtree/file_e.txt", "test/tree.zip/subtree/file_f.txt",
		},
	},
}

func TestZipFolders(t *testing.T) {

	for _, c := range testZipFolderCases {
		folder, err := Open(c.path)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
			break
		}
		got := []string{}
		for f := range folder.Items() {
			got = append(got, f.FullName())
			f.Close()
		}
		if len(got) != len(c.expected) {
			t.Errorf("Expected file count %d, but got %d", len(c.expected), len(got))
			break
		}
		sort.Strings(got)
		sort.Strings(c.expected)
		if !reflect.DeepEqual(c.expected, got) {
			t.Errorf("Expected %#q, but got %#q", c.expected, got)
			break

		}
		folder.Close()

	}
}

func TestZipFolderItemOpen(t *testing.T) {
	for _, c := range testZipFolderCases {
		folder, err := Open(c.path)
		if err != nil {
			t.Errorf("Unexpected error %s", err)
			break
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
	}
}
