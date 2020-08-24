package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/djimenez/iconv-go"
)

const (
	head = `---
layout: default
title: {{ site.name }}
---

## 로또 당첨자 번호

| 회차 | 당청번호 | 1등 당첨자 수 | 1등 당첨금(원) |
| ---- | ------- | ------------ | ------------- |
`
)

// Lottery data
type Lottery struct {
	Times   int
	Numbers string
	Winners int
	Reward  string
}

var (
	datas []Lottery
)

func main() {
	const TargetURL = "https://dhlottery.co.kr/gameResult.do?method=allWin"
	// 1. get envirionment

	// 2. 서버 요청
	nowPage := 0
	datas = make([]Lottery, 10, 10)

	for {
		nowPage++
		res, err := http.PostForm(TargetURL, url.Values{
			"nowPage":    {fmt.Sprintf("%d", nowPage)},
			"drwNoStart": {"1"},
			"drwNoEnd":   {"1000"},
		})
		if err != nil {
			log.Fatal(err)
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			log.Fatal(res.StatusCode)
		}

		// 3. html pairing
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var data Lottery
		doc.Find("table.tbl_data.tbl_data_col tbody tr").Children().Each(func(i int, s *goquery.Selection) {
			v, _ := iconv.ConvertString(s.Text(), "euc-kr", "utf-8")
			switch i % 8 {
			case 0:
				d, _ := strconv.Atoi(v)
				data.Times = d
			case 1:
				data.Numbers = strings.ReplaceAll(v, " ", "")
			case 2:
				d, _ := strconv.Atoi(v)
				data.Winners = d
			case 3:
				data.Reward = v
				datas = append(datas, data)
			}
		})

		if data.Times == 1 {
			break
		}

		time.Sleep(time.Second * 1)
		log.Println(nowPage)
	}

	// 4. sorting

	// 5. 파싱 대상 md 파일로 저장
	saveTableFile(datas)
}

func saveTableFile(datas []Lottery) {
	table, err := os.OpenFile("table.md", os.O_RDWR|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer table.Close()
	table.WriteString(head)

	for _, d := range datas {
		if d.Times != 0 {
			table.WriteString(fmt.Sprintf("| %d | %s | %d | %s |\n", d.Times, d.Numbers, d.Winners, d.Reward))
		}
	}

	table.WriteString(fmt.Sprintf("\nLast Automatic Update: %s\n", time.Now().Format(time.RFC3339)))
}
