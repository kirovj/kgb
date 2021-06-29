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
	blogDir  = "blogs"
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
	Id       int
	Title    string
	Filepath string
	Time     int64
}

func newBlog(id int, title string, filepath string, time int64) *Blog {
	return &Blog{
		Id:       id,
		Title:    title,
		Filepath: filepath,
		Time:     time,
	}
}

func getTime(name string) int64 {
	return time.Now().Unix()
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

func main() {
	r := girov.New()
	blogGroup := r.Group("blog/")
	r.LoadHtmlGlob("tmpls/*")
	r.Static("/assets", "./static")

	var blogs []*Blog

	dir, err := ioutil.ReadDir(blogDir)
	if err != nil {
		fmt.Println(err)
	}

	// url for each blog
	for i, fileInfo := range dir {
		blog := newBlog(i, strings.Split(fileInfo.Name(), "_")[0], fileInfo.Name(), getTime(fileInfo.Name()))
		blogs = append(blogs, blog)
		md := blogDir + `/` + blog.Filepath
		blogGroup.GET(fmt.Sprint(i), func(c *girov.Context) {
			//c.MD(http.StatusOK, md, md2html)
			c.HTMLFromMD(http.StatusOK, "blog.tmpl", md, md2html)
		})
	}

	// url for blog index
	r.GET("/", func(c *girov.Context) {
		c.HTML(http.StatusOK, "index.tmpl", girov.H{
			"blogs": blogs,
		})
	})

	_ = r.Run(":9999")
}
