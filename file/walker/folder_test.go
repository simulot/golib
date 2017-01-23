package walker

import (
	"bufio"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestOpenFolder(t *testing.T) {

	folder, err := Open("nowhere")
	if err == nil {
		t.Errorf("Expected error, but got '%s'", err)
		return
	}

	folder, err = Open("test")
	if err != nil {
		t.Errorf("Unexpected error '%s' when opening 'test' folder", err)
		return
	}
	folder.Close()

}

var testFolderCases = []struct {
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
}

func TestFolderItemOpen(t *testing.T) {
	for _, c := range testFolderCases {
		t.Run(c.path, func(t *testing.T) {
			filelist := []string{}
			folder, err := Open(c.path)
			if err != nil {
				t.Fatalf("Unexpected error %s", err)
				return
			}

			for item := range folder.Items() {
				filelist = append(filelist, item.FullName())
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
			sort.Strings(filelist)
			sort.Strings(c.expected)
			if !reflect.DeepEqual(filelist, c.expected) {
				t.Errorf("Expecting\n%#q\nbut got\n%#q", c.expected, filelist)
			}
		})
	}
}
