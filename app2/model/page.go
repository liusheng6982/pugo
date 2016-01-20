package model

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/go-xiaohei/pugo/app2/helper"
	"github.com/naoina/toml"
)

// Page contain all fields of a page content
type Page struct {
	Title      string                 `toml:"title"`
	Slug       string                 `toml:"slug"`
	Desc       string                 `toml:"desc"`
	Date       string                 `toml:"date"`
	Update     string                 `toml:"update_date"`
	AuthorName string                 `toml:"author"`
	NavHover   string                 `toml:"hover"`
	Template   string                 `toml:"template"`
	Lang       string                 `toml:"lang"`
	Bytes      []byte                 `toml:"-"`
	Meta       map[string]interface{} `toml:"meta"`

	permaURL     string
	pageURL      string
	contentBytes []byte
	dateTime     time.Time
	updateTime   time.Time
}

func (p *Page) normalize() error {
	if p.Slug == "" {
		p.Slug = titleReplacer.Replace(p.Title)
	}
	if p.Template == "" {
		p.Template = "page.html"
	}
	var err error
	if p.dateTime, err = time.Parse(postTimeLayout, p.Date); err != nil {
		return err
	}
	if p.Update == "" {
		p.Update = p.Date
		p.updateTime = p.dateTime
	} else {
		if p.updateTime, err = time.Parse(postTimeLayout, p.Update); err != nil {
			return err
		}
	}
	p.contentBytes = helper.Markdown(p.Bytes)
	if p.Lang == "" {
		p.permaURL = fmt.Sprintf("/%s", p.Slug)
		p.pageURL = p.permaURL + ".html"
	} else {
		p.permaURL = fmt.Sprintf("/%s/%s", p.Lang, p.Slug)
		p.pageURL = p.permaURL + ".html"
	}
	return nil
}

// NewPageOfMarkdown create new page from markdown file
func NewPageOfMarkdown(file string) (*Page, error) {
	fileBytes, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	dataSlice := bytes.SplitN(fileBytes, postBlockSeparator, 3)
	if len(dataSlice) != 3 {
		return nil, fmt.Errorf("post need toml block and markdown block")
	}
	if !bytes.HasPrefix(dataSlice[1], tomlPrefix) {
		return nil, fmt.Errorf("post need toml block at first")
	}
	page := new(Page)
	if err = toml.Unmarshal(dataSlice[1][4:], page); err != nil {
		return nil, err
	}
	page.Bytes = bytes.Trim(dataSlice[2], "\n")
	return page, page.normalize()
}
