package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"sort"
	"time"

	"github.com/alecthomas/chroma/formatters/html"
	"github.com/gin-gonic/gin"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
)

var (
	blogDir  = "blogs"
	blogs    blogList
	blogMap  = make(map[string]*blog)
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
	var tmp blogList

	for i := len(dir) - 1; i >= 0; i-- {
		fileInfo := dir[i]
		if !fileInfo.IsDir() {
			blog := newBlog(fileInfo.ModTime(), fileInfo.Name())
			// update blog list
			tmp = append(tmp, blog)
			// update blog map
			blogMap[blog.title] = blog
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
		c.HTML(http.StatusOK, "blog.tmpl", blogMap[c.Param("title")].text)
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
