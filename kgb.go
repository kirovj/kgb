package kgb

import (
	"fmt"
	"github.com/gomarkdown/markdown"
	"io/ioutil"
	"kgb/girov"
	"net/http"
)

const blogDir = "blogs"

var files = make(map[string]string)

func md2html(file string) string {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(markdown.ToHTML(content, nil, nil))
}

func Run() {
	r := girov.New()

	dir, err := ioutil.ReadDir(blogDir)
	if err != nil {
		fmt.Println(err)
	}

	for i, fileInfo := range dir {
		files[fileInfo.Name()] = md2html(blogDir + `\` + fileInfo.Name())
		r.GET("/blog/"+fmt.Sprint(i), func(c *girov.Context) {
			c.MD(http.StatusOK, files[fileInfo.Name()])
		})
		//fmt.Println(md2html(blogDir + `\` + fileInfo.Name()))
	}

	r.Run(":9999")
}
