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
	Id                    int
	Time                  time.Time
	Text                  template.HTML
	Title, Filepath, Date string
}

type BlogList []*Blog

type BlogMap map[string]*Blog

var (
	blogDir  = "blogs"
	blogs    BlogList
	blogMap  = make(BlogMap)
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
	t1, _ := time.Parse("2006-01-02", b[i].Date)
	t2, _ := time.Parse("2006-01-02", b[j].Date)
	return t1.Unix() > t2.Unix()
}

func (b BlogList) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func dateFormat(name string) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	timeStr := strings.Split(name, "_")[0]
	timeObj, _ := time.ParseInLocation("2006-1-2", timeStr, loc)
	return timeObj.Format("2006-01-02")
}

func readMd(filePath string) string {
	if content, err := ioutil.ReadFile(filePath); err == nil {
		var buf bytes.Buffer
		if err = markdown.Convert(content, &buf); err == nil {
			return buf.String()
		}
	}
	return ""
}

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
			blog := &Blog{
				Id:       i,
				Time:     fileInfo.ModTime(),
				Title:    strings.ReplaceAll(strings.Split(fileInfo.Name(), "_")[1], ".md", ""),
				Filepath: fileInfo.Name(),
				Date:     dateFormat(fileInfo.Name()),
				Text:     template.HTML(readMd(blogDir + `/` + fileInfo.Name())),
			}
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
