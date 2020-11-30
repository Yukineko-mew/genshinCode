package common

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/saintfish/chardet"
	"github.com/sclevine/agouti"
	"golang.org/x/net/html/charset"
)

// GetDocument URLからgoquery形のドキュメントを取得する
func GetDocument(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// 読み取り
	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// 文字コード判定
	det := chardet.NewTextDetector()
	detRslt, err := det.DetectBest(buf)
	if err != nil {
		return nil, err
	}

	// 文字コード変換
	bReader := bytes.NewReader(buf)
	reader, err := charset.NewReaderLabel(detRslt.Charset, bReader)
	if err != nil {
		return nil, err
	}

	// HTMLパース
	doc, err := goquery.NewDocumentFromReader(reader)

	return doc, err
}

// GetDocumentUseChromeDriver ヘッドレスブラウザを使用してjs実行後のページを取得する
func GetDocumentUseChromeDriver(url string) (*goquery.Document, error) {
	driver := agouti.ChromeDriver(
		agouti.ChromeOptions("args", []string{
			"--headless",
			"--window-size=1280,800",
		}),
		agouti.Debug,
	)

	if err := driver.Start(); err != nil {
		return nil, err
	}

	defer driver.Stop()
	page, err := driver.NewPage()
	if err != nil {
		return nil, err
	}

	// URL Setting
	page.Navigate(url)

	// htmlを取得
	getSource, err := page.HTML()
	if err != nil {
		return nil, err
	}

	r := strings.NewReader(getSource)
	doc, err := goquery.NewDocumentFromReader(r)

	return doc, nil
}

// Scraping 記事中から更新日付とシリアルコード一覧を取得
func Scraping(doc *goquery.Document) (updateDatetime string, serialCodeList []string) {
	// 記事中の更新日付を抜き出す
	rslt := doc.Find("time").Text()
	updateDatetime = strings.TrimSpace(rslt)

	// 記事中のシリアルコードを抜き出す
	doc.Find("button").Each(func(i int, s *goquery.Selection) {
		serialCode, _ := s.Attr("data-clipboard-text")
		if serialCode != "" {
			serialCodeList = append(serialCodeList, serialCode)
		}
	})

	return updateDatetime, serialCodeList
}
