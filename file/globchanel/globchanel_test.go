package globchanel

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestGlob(t *testing.T) {
	cases := []string{
		"*.md",
		"../*/*.md",
	}

	for _, c := range cases {
		t.Run(c, func(t *testing.T) {
			cpt := 0
			duration(t, "filepath.Glob", func() {
				matches, _ := filepath.Glob(c)
				wg := sync.WaitGroup{}
				for _, p := range matches {
					wg.Add(1)
					go func() {
						heavy(p)
						wg.Done()
					}()
					cpt++
				}
				wg.Wait()
			})
			fmt.Println(cpt, "items")
			cpt = 0
			duration(t, "globchanel.Glob", func() {
				wg := sync.WaitGroup{}
				nameChan := Glob(c)
				for p := range nameChan {
					wg.Add(1)
					go func() {
						heavy(p)
						wg.Done()
					}()
					cpt++
				}
				wg.Wait()
			})
			fmt.Println(cpt, "items")
		})
	}
}

func duration(t *testing.T, label string, f func()) {
	start := time.Now()
	f()
	t.Log(label, time.Now().Sub(start))
}

func heavy(p string) {
	f, err := os.Open(p)
	if err != nil {
		return
	}
	defer f.Close()
	ioutil.ReadAll(f)
}
