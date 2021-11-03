package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

func main() {
	events := []Event{
		DailyEvent{
			startDate: "20211101",
			n:         1,
			title:     "原神",
			desc:      "霓裳花",
		},
	}
	var results []struct {
		Title string
		Items []string
	}
	date, _ := time.Parse("20060102", time.Now().Format("20060102"))
	for _, e := range events {
		title := e.Title()
		ok, items := e.Check(date)
		if ok {
			rows := []string{}
			for _, item := range items {
				tags := []string{}
				for _, tag := range item.Tags() {
					tags = append(tags, fmt.Sprintf("[%s]", tag.ToString()))
				}
				rows = append(rows, fmt.Sprintf("%s%s", strings.Join(tags, ""), item.ToString()))
			}
			results = append(results, struct {
				Title string
				Items []string
			}{Title: title, Items: rows})
		}
	}
	j, _ := json.MarshalIndent(results, "", "  ")
	fmt.Printf("%s", string(j))
}

type Tag interface {
	ToString() string
}

type BaseTag struct {
	val string
}

func (b BaseTag) ToString() string {
	return b.val
}

type ErrTag struct{}

func (_ ErrTag) ToString() string {
	return "处理异常"
}

// Item 一条页面上的记录
type Item interface {
	ToString() string
	Tags() []Tag
}

type ErrItem struct {
	ErrMsg string
}

func (e ErrItem) ToString() string {
	return e.ErrMsg
}

func (e ErrItem) Tags() []Tag {
	return []Tag{ErrTag{}}
}

// Event 一个规则
type Event interface {
	Title() string
	Check(date time.Time) (bool, []Item) // 接受一个日期，返回是否有记录以及记录的内容
}

// DailyEvent 每n天一次
type DailyEvent struct {
	startDate string // yyyyMMdd
	n         int    // 每n天一次
	title     string // 规则名
	desc      string // 内容
}

func (d DailyEvent) Title() string {
	return d.title
}

type DailyItem struct {
	n    int
	desc string
}

func (d DailyItem) ToString() string {
	return d.desc
}

func (d DailyItem) Tags() []Tag {
	return []Tag{
		BaseTag{val: fmt.Sprintf("每%d天", d.n)},
	}
}

func (d DailyEvent) Check(date time.Time) (bool, []Item) {
	start, err := time.Parse("20060102", d.startDate)
	if err != nil {
		return true, []Item{ErrItem{ErrMsg: fmt.Sprintf("DailyEvent parse(%s) fail: %s", d.startDate, err.Error())}}
	}
	if (int(start.Sub(date).Hours())/24)%d.n != 0 {
		return false, nil
	}
	return true, []Item{DailyItem{
		n:    d.n,
		desc: d.desc,
	}}
}
