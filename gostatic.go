package main

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"github.com/owais/gostatic/utils"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const CONTENT_DIR string = "content"
const TEMPLATE_DIR string = "templates"
const SITE_DIR string = "site"
const MEDIA_DIR string = "media"
const MEDIA_URL string = "/media/"
const ModePerm os.FileMode = 0777


func main() {
	working_directory, _ := os.Getwd()
	site_dir := filepath.Join(working_directory, SITE_DIR)
    content_dir := utils.GetDirectoryPath(CONTENT_DIR)
	media_dir := utils.GetDirectoryPath(MEDIA_DIR)
	template_dir := utils.GetDirectoryPath(TEMPLATE_DIR)
	all_templates := utils.GetFiles(template_dir)
	tmpl, err := template.ParseFiles(all_templates...)
	if err != nil {
		fmt.Println(err)
		return
	}

	filepath.Walk(content_dir, func(path string, info os.FileInfo, err error) error {
		postpath, postname := filepath.Split(path)
		ext := filepath.Ext(path)
		if ext != ".post" {
			return nil
		}
		post, err := utils.ReadContentFile(path)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		out_dir := strings.Replace(postpath, content_dir, site_dir, 1)
		os.MkdirAll(out_dir, ModePerm)
		//tmpl, err := template.New(post.Template).ParseFiles(filepath.Join(template_dir, post.Template))
		out, err := os.Create(filepath.Join(out_dir, strings.TrimSuffix(postname, ".post")+filepath.Ext(post.Template)))
		defer out.Close()
		post.MEDIA_URL = MEDIA_URL
		tmpl.ExecuteTemplate(out, post.Template, post)
		return nil
	})

	filepath.Walk(media_dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() == true {
			return nil
		}
		in, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer in.Close()
		out_path := strings.Replace(path, media_dir, filepath.Join(site_dir, MEDIA_URL), 1)
		out_dir, _ := filepath.Split(out_path)
		os.MkdirAll(out_dir, ModePerm)
		out, err := os.Create(out_path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		//defer out.Close()
		io.Copy(out, in)
		return nil
	})
	go func() {
		panic(http.ListenAndServe(":8080", http.FileServer(http.Dir(site_dir))))
	}()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)

	// Process events
	go func() {
		for {
			select {
			case ev := <-watcher.Event:
				if ev.IsCreate() || ev.IsModify() {
					log.Println("Need to copy", ev)
				}
			case err := <-watcher.Error:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Watch(content_dir)
	if err != nil {
		log.Fatal(err)
	}
	<-done
	watcher.Close()
}
