package common

import (
	"bufio"
	"os"
)

// savedFileName データを保存するファイル名
const savedFileName = "lastUpdate.txt"

// DataFile 更新日付や最後に取得したシリアルコード一覧を保存しておくファイル
type DataFile struct {
	LastUpdate     string
	SerialCodeList []string
}

var dataFile *DataFile

// LoadDataFile データファイルの初期化
func loadDataFile() (*DataFile, error) {
	// すでに取得済みならそれを返す
	if dataFile != nil {
		return dataFile, nil
	}

	// ファイル読み込み
	fp, err := os.Open(savedFileName)
	if err != nil {
		return nil, err
	}
	defer fp.Close()

	dataFile = new(DataFile)

	scanner := bufio.NewScanner(fp)

	// 一行目は最後に取得した日付
	dataFile.LastUpdate = ""
	if scanner.Scan() {
		dataFile.LastUpdate = scanner.Text()
	}

	// 二行目以降は最後に取得したシリアルコード
	var serialCodeList []string
	for scanner.Scan() {
		serialCodeList = append(serialCodeList, scanner.Text())
	}
	dataFile.SerialCodeList = serialCodeList

	return dataFile, nil
}

// GetSavedDatetime 取得データをファイルから読み込む
func GetSavedDatetime() string {
	dataFile, _ := loadDataFile()
	return dataFile.LastUpdate
}

// GetSavedSerialCodeList 保存されていたシリアルコードを取得
func GetSavedSerialCodeList() []string {
	dataFile, _ := loadDataFile()
	return dataFile.SerialCodeList
}

// SetLastUpdate 最終更新日時を更新
func SetLastUpdate(lastUpdate string) {
	dataFile, _ := loadDataFile()
	dataFile.LastUpdate = lastUpdate
}

// SetSerialCodeList シリアルコード一覧を更新
func SetSerialCodeList(serialCodeList []string) {
	dataFile, _ := loadDataFile()
	dataFile.SerialCodeList = serialCodeList
}

// SaveFile 取得データをファイルに保存する
func SaveFile() error {
	// 保存内容がないなら何もしない
	if dataFile == nil {
		return nil
	}

	// 保存ファイルを開く
	file, err := os.Create(savedFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// 一行目に最終更新日付を保存
	_, err = file.WriteString(dataFile.LastUpdate + "\n")
	if err != nil {
		return err
	}

	// 二行目以降はシリアルコード一覧を保存
	for _, line := range dataFile.SerialCodeList {
		_, err := file.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
