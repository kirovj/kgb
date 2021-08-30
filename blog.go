package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"strings"
	"time"
)

type blog struct {
	time              time.Time
	text              template.HTML
	title, file, date string
}

type blogList []*blog

func (b blogList) Len() int {
	return len(b)
}

func (b blogList) Less(i, j int) bool {
	t1, _ := time.Parse("2006-01-02", b[i].date)
	t2, _ := time.Parse("2006-01-02", b[j].date)
	return t1.Unix() > t2.Unix()
}

func (b blogList) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

func newBlog(time time.Time, fileName string) *blog {
	return &blog{
		time:  time,
		text:  getHTML(blogDir + `/` + fileName),
		title: strings.ReplaceAll(strings.Split(fileName, "_")[1], ".md", ""),
		file:  fileName,
		date:  getDate(fileName),
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
