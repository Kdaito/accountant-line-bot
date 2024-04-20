package pkg

import (
	"google.golang.org/api/sheets/v4"
)

type GSheet struct {
	Service *sheets.Service
}

func (g *GSheet) ExportSheet() {
	// jst := time.FixedZone("Asia/Tokyo", 9*60*60)
	// timestamp := time.Now().In(jst).Format("2006-01-02-15-04-05")
	// 複製先のファイルのタイトル
	// newFileTitle := fmt.Sprintf("%s_Copy", timestamp)
	// ctx := context.Background()

	// // サービスアカウントの秘密鍵を読み込む
	// b, err := ioutil.ReadFile("service-account.json")
	// if err != nil {
	// 	log.Fatalf("cannot read service account json file: %v", err)
	// }
}
