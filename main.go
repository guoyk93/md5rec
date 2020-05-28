package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	md5sumFile = "md5sum.txt"
)

func exit(err *error) {
	if *err != nil {
		log.Println((*err).Error())
		os.Exit(1)
	}
}

func main() {
	var err error
	defer exit(&err)

	if err = handle("."); err != nil {
		return
	}
}

func handle(dir string) (err error) {
	log.Println("=>", dir)
	defer log.Println("<=", dir)
	var fis []os.FileInfo
	if fis, err = ioutil.ReadDir(dir); err != nil {
		return
	}

	var files []string
	for _, fi := range fis {
		if fi.Name() == md5sumFile {
			continue
		}
		if strings.HasPrefix(fi.Name(), ".") {
			continue
		}
		if fi.IsDir() {
			if err = handle(filepath.Join(dir, fi.Name())); err != nil {
				return
			}
		} else {
			files = append(files, fi.Name())
		}
	}

	if len(files) == 0 {
		return
	}

	if err = md5sum(dir, files); err != nil {
		return
	}
	return
}

func md5sum(dir string, files []string) (err error) {
	buf := &bytes.Buffer{}

	for _, file := range files {
		if err = md5sumSingle(dir, file, buf); err != nil {
			return
		}
	}

	if err = ioutil.WriteFile(filepath.Join(dir, md5sumFile), buf.Bytes(), 0644); err != nil {
		return
	}

	return
}

func md5sumSingle(dir, file string, w io.Writer) (err error) {
	var f *os.File
	if f, err = os.Open(filepath.Join(dir, file)); err != nil {
		return
	}
	defer f.Close()

	log.Println("  *", file)

	h := md5.New()
	if _, err = io.Copy(h, f); err != nil {
		return
	}
	out := hex.EncodeToString(h.Sum(nil))

	if _, err = fmt.Fprintf(w, "%s *%s\n", out, file); err != nil {
		return
	}
	return
}
