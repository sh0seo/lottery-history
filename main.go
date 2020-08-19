package main

import (
	"fmt"
	_ "io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
)

const (
	head = `
	| 회차 | 당청번호 | 1등 당첨자 수 | 1등 당첨금(원) |
	| ---- | ------- | ------------ | ------------- |
	`
)

type Lottery struct {
	Times   int
	Numbers string
	Winners int
	Reward  string
}

func main() {
	const TargetURL = "https://dhlottery.co.kr/gameResult.do?method=allWin"
	const TableFileName = "table.md"
	// 1. get envirionment

	// 2. 서버 요청
	// v := url.Values{}
	// v.Add("nowPage", "1")
	// v.Add("drwNoStart", "1")
	// v.Add("drwNoEnd", "1000")
	res, err := http.PostForm(TargetURL, url.Values{
		"nowPage":    {"1"},
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

	doc.Find("table.tbl_data.tbl_data_col tbody tr").Children().Each(func(i int, s *goquery.Selection) {
		if i%8 <= 3 {
			v := s.Text()
			fmt.Printf("%d %s\n", i, v)
		}
		// if s.Nodes[0].Attr != nil {
		// fmt.Printf("%d %s %s\n", i, v, s.Nodes[0].Attr)
		//} else {
		//}
	})
	// body, err := ioutil.ReadAll(resp.Body)
	// if err != nil {
	//		log.Fatal(err)
	// }
	// log.Printf(string(body))

	// 4. 파싱 대상 md 파일로 저장
}
