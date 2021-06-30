package main

import (
	"bytes"
	"fmt"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-highlighting"
	"io/ioutil"
	"kgb/girov"
	"net/http"
	"strings"
	"time"
)

var (
	blogDir  = "../blogs"
	blogs    []*Blog
	markdown = goldmark.New(
		goldmark.WithExtensions(
			highlighting.NewHighlighting(
				highlighting.WithStyle("monokai"),
				highlighting.WithFormatOptions(
					html.WithLineNumbers(false),
				),
			),
		),
	)
)

type Blog struct {
	Id                    int
	Title, Filepath, Time string
}

func getTime(name string) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	timeStr := strings.Split(name, "_")[0]
	timeObj, _ := time.ParseInLocation("2006-1-2", timeStr, loc)
	return timeObj.Format("2006-01-02")
}

func md2html(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	var buf bytes.Buffer
	if err := markdown.Convert(content, &buf); err != nil {
		return ""
	}
	return buf.String()
}

func getBlogs() {
	dir, err := ioutil.ReadDir(blogDir)
	if err != nil {
		fmt.Println(err)
	}
	var tmp []*Blog

	for i := len(dir) - 1; i >= 0; i-- {
		fileInfo := dir[i]
		if !fileInfo.IsDir() {
			blog := &Blog{
				Id:       i,
				Title:    strings.ReplaceAll(strings.Split(fileInfo.Name(), "_")[1], ".md", ""),
				Filepath: fileInfo.Name(),
				Time:     getTime(fileInfo.Name()),
			}
			tmp = append(tmp, blog)
		}
	}
	blogs = tmp
}

func main() {
	r := girov.New()
	blogGroup := r.Group("blog/")
	r.LoadHtmlGlob("../tmpls/*")
	r.Static("/static", "../static")

	go func() {
		for {
			getBlogs()
			// url for each blog
			for _, blog := range blogs {
				md := blogDir + `/` + blog.Filepath
				blogGroup.GET(fmt.Sprint(blog.Id), func(c *girov.Context) {
					c.HTMLFromMD(http.StatusOK, "blog.tmpl", md, md2html)
				})
			}
			time.Sleep(time.Minute)
		}
	}()

	// url for blog index
	r.GET("/", func(c *girov.Context) {
		c.HTML(http.StatusOK, "index.tmpl", girov.H{
			"blogs": blogs,
		})
	})

	_ = r.Run(":9999")

	r.GET("/test", func(c *girov.Context) {
		c.String(http.StatusOK, "asd")
	})
}
