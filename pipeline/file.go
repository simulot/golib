package pipeline

import (
	"fmt"

	"path/filepath"

	"sync"

	"github.com/simulot/golib/file/walker"
)

// GlobOperator return the exanded list of files
// corresponding to wildcards.
//
// IN: chan string : wildcarded names
// OUT: chan string : list of actual files correspondind to wildwards
func GlobOperator() Operator {
	return func(in, out chan interface{}) {
		for i := range in {
			if s, ok := i.(string); !ok {
				panic("Expecting string in GlobOperator")
			} else {
				paths, _ := filepath.Glob(s)
				if paths != nil {
					for _, p := range paths {
						out <- p
					}
				}
			}
		}
	}
}

func FolderToWalkersOperator() Operator {
	return func(in, out chan interface{}) {
		for i := range in {
			if path, ok := i.(string); ok {
				w, err := walker.FileWalkerOpen(path)
				if err != nil {
					fmt.Println(err)
					continue
				}
				out <- w
			} else {
				panic("Expectinng string in FolderToWalkersOperator")
			}
		}
	}
}

// PathOperator returns an operator that will list all
// files and archives in args.
// in chan string: path
// out chan is a channel of FileItem
func WalkOperator() Operator {
	return func(in, out chan interface{}) {
		for i := range in {
			if w, ok := i.(walker.Walker); ok {
				for item := range w.Items() {
					out <- item
				}
				w.Close()
			} else {
				panic("Expecting Walker in WalkOperator")
			}
		}
	}
}

type deDuplicate struct {
	sync.Mutex
	m map[string]int
}

var seenFile = deDuplicate{}

func FileDeduplicateOperator() Operator {
	return func(in, out chan interface{}) {
		for i := range in {
			if item, ok := i.(walker.WalkItem); ok {
				fullname := item.Name()
				seenFile.Lock()
				if _, ok := seenFile.m[fullname]; !ok {
					out <- item
					seenFile.Unlock()
					continue
				} else {
					seenFile.m[fullname]++
					seenFile.Unlock()
				}
			} else {
				panic("walker.ItemOpenCloser expected in FileDeduplicateOperator")
			}
		}
	}
}
