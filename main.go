package main

import (
	"fmt"

	"github.com/y-ogasawara/genshinCode/common"
)

// TargetURL スクレイピング先のURL
const TargetURL = "https://gamewith.jp/genshin/article/show/231856"

const isDebug = false

func main() {
	doc, err := common.GetDocumentUseChromeDriver(TargetURL)
	if err != nil {
		panic(err)
	}

	// スクレイピング処理
	updateDatetime, serialCodeList := common.Scraping(doc)

	// 前回取得した更新日付を取得
	savedLastUpdate := common.GetSavedDatetime()

	// 前回取得したシリアルコード一覧を取得
	savedSerialCodeList := common.GetSavedSerialCodeList()

	// 前回と比べてシリアルコードが新規に増えているかどうか
	var newSerialCodeList []string
Exist:
	for _, serialCode := range serialCodeList {
		// 前回取得したシリアルコードの中に同じシリアルコードが存在したか
		for _, oldSerialCode := range savedSerialCodeList {
			if serialCode == oldSerialCode {
				continue Exist
			}
		}
		newSerialCodeList = append(newSerialCodeList, serialCode)
	}

	// 記事が更新されているか
	if savedLastUpdate != updateDatetime {
		fmt.Println("更新されました！" + updateDatetime)

		// 内容を更新してファイルに保存
		common.SetLastUpdate(updateDatetime)
		common.SetSerialCodeList(serialCodeList)
		common.SaveFile()

		// 記事が更新されているが、シリアルコードが増えていないならLINEに投稿はしない
		if 0 < len(newSerialCodeList) {
			// 更新メッセージを作成
			postMessage := common.CreatePostMessage(newSerialCodeList, TargetURL)
			//fmt.Println(postMessage)

			// LINE投稿
			if !isDebug {
				err := common.PostLineMessage(postMessage)
				if err != nil {
					panic(err)
				}
			}
		}
	} else {
		fmt.Println("更新なし")
	}
}
