package pipeline

import (
	"fmt"

	"path/filepath"

	"github.com/pkg/errors"
	"github.com/simulot/golib/file/walker"
)

// GlobOperator return the expanded list of files
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

// FolderToWalkersOperator accpets file names and directory names
// and transforms it into a walker
// IN string
// OUT walker.Walker
func FolderToWalkersOperator() Operator {
	return func(in, out chan interface{}) {
		for i := range in {
			if path, ok := i.(string); ok {
				w, err := walker.Open(path)
				if err != nil {
					fmt.Println(err)
					continue
				}
				for item := range w.Items() {
					out <- walker.Walker(item)
				}
			} else {
				panic("Expecting string in FolderToWalkersOperator")
			}
		}
	}
}

// WalkOperator is an operator that walks through walker's items
// IN walker.Walker
// OUT waker.Items
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

// FileMaskOperator filter a channel of walker.WalkItem items using a file mask
// IN : chan walker.WalkItem
// OUT : chan walker.WalkItem
func FileMaskOperator(mask string) Operator {
	return func(in, out chan interface{}) {
		for i := range in {
			if item, ok := i.(walker.WalkItem); ok {
				match, err := filepath.Match(mask, item.Name())
				if err != nil {
					fmt.Println(errors.Wrapf(err, "Can't use mask in FileMaskOperator"))
					continue
				}
				if match {
					out <- item
				} else {
					item.Close()
				}
			} else {
				panic("Expecting walker.WalkItem in FileMaskOperator")
			}
		}
	}
}

/*
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
*/
