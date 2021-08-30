package main

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"strings"
	"time"
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
