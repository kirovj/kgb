package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"kgb/girov"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
)

type Blog struct {
	Id                    int
	Title, Filepath, Time string
}

type BlogList []*Blog

var (
	blogDir  = "../blogs"
	blogs    BlogList
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

func (b BlogList) Len() int {
	return len(b)
}

func (b BlogList) Less(i, j int) bool {
	t1, _ := time.Parse("2006-01-02", b[i].Time)
	t2, _ := time.Parse("2006-01-02", b[j].Time)
	return t1.Unix() > t2.Unix()
}

func (b BlogList) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
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
	var tmp BlogList

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
	sort.Sort(blogs)
}

func main() {
	r := girov.Default()
	blogGroup := r.Group("blog/")
	r.LoadHtmlGlob("../tmpls/*")
	r.Static("/static", "../static")

	// url for each blog, update every minutes
	go func() {
		for {
			getBlogs()
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
}
