package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sort"
	"strings"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
)

type Blog struct {
	Time              time.Time
	Text              template.HTML
	Title, File, Date string
}

type BlogList []*Blog

func (b BlogList) Len() int {
	return len(b)
}

func (b BlogList) Less(i, j int) bool {
	t1, _ := time.Parse("2006-01-02", b[i].Date)
	t2, _ := time.Parse("2006-01-02", b[j].Date)
	return t1.Unix() > t2.Unix()
}

func (b BlogList) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func newBlog(time time.Time, fileName string) *Blog {
	return &Blog{
		Time:  time,
		Text:  getHTML(blogDir + `/` + fileName),
		Title: strings.ReplaceAll(strings.Split(fileName, "_")[1], ".md", ""),
		File:  fileName,
		Date:  getDate(fileName),
	}
}

func getDate(name string) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	timeStr := strings.Split(name, "_")[0]
	timeObj, _ := time.ParseInLocation("2006-1-2", timeStr, loc)
	return timeObj.Format("2006-01-02")
}

func getHTML(filePath string) template.HTML {
	if content, err := ioutil.ReadFile(filePath); err == nil {
		var buf bytes.Buffer
		if err = markdown.Convert(content, &buf); err == nil {
			return template.HTML(buf.String())
		}
	}
	return ""
}

var (
	blogDir  = "blogs"
	blogs    BlogList
	blogMap  = make(map[string]*Blog)
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

func update() {
	exec.Command("sh", "-c", "git pull origin main").Run()
	dir, err := ioutil.ReadDir(blogDir)
	if err != nil {
		fmt.Println(err)
	}
	var tmp BlogList

	for i := len(dir) - 1; i >= 0; i-- {
		fileInfo := dir[i]
		if !fileInfo.IsDir() {
			blog := newBlog(fileInfo.ModTime(), fileInfo.Name())
			// update blog list
			tmp = append(tmp, blog)
			// update blog map
			blogMap[blog.Title] = blog
		}
	}
	blogs = tmp
	sort.Sort(blogs)
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("tmpls/*")
	r.Static("static", "static")

	// url for each blog, update every minutes
	r.GET("blog/:title", func(c *gin.Context) {
		c.HTML(http.StatusOK, "blog.tmpl", blogMap[c.Param("title")].Text)
	})

	// update blogs
	go func() {
		for {
			update()
			time.Sleep(time.Minute)
		}
	}()

	// url for blog index
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"blogs": blogs,
		})
	})

	_ = r.Run(":9999")
}
