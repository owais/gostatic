package utils

import (
	"encoding/json"
	"fmt"
	"github.com/owais/gostatic/structures"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func GetDirectoryPath(name string) string {
	working_directory, _ := os.Getwd()
	dir := filepath.Join(working_directory, name)
	if PathExists(dir) == false {
		fmt.Println("directory does not exist")
	}
	return dir
}

func ReadContentFile(path string) (structures.Post, error) {
	var p structures.Post
	raw, err := ioutil.ReadFile(path)
	if err != nil {
		return p, err
	}
	post_array := strings.Split(string(raw), "-->")
	meta := strings.TrimPrefix(post_array[0], "<!--")
	err = json.Unmarshal([]byte(meta), &p)
	if err != nil {
		return p, err
	}
	p.Body = post_array[1]
	return p, nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetFiles(root string) []string {
	var files []string
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == false {
			files = append(files, path)
		}
		return nil
	})
	return files
}
