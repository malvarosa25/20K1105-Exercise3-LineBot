package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/line/line-bot-sdk-go/linebot"
	_ "github.com/mattn/go-sqlite3" // 使用しないため、 _ にしないとコンパイルエラーとなる
)

// 生徒のデータ格納用のユーザ定義型
type Student struct {
	Name   string // 生徒の名前
	Minute int    // 対象の生徒の学習制限時間
}

func main() {
	// 生徒の情報が載っているデータベースを作成
	Db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Println(err)
	}

	if err := CreateTable(Db); err != nil {
		log.Println(err)
	}

	students := []*Student{{Name: "鈴宮 花子", Minute: 1}, {Name: "鈴宮 太郎", Minute: 2}, {Name: "鈴宮 次郎", Minute: 5}}

	for i := range students {
		if err := InsertTable(Db, students[i]); err != nil {
			log.Println(err)
		}
	}

	err = Db.Ping() // DB の疎通確認

	if err != nil {
		log.Fatalln(err)
	}

	// Line Developer にて立ち上げたチャネルの情報
	Channel_Secret := "857d036768e6c23dd8731bec8d08312f"                                                                                                                                            // チャネルシークレット
	Channel_Token := "+MVr5jo/PqWuzYfQ8G3DZyFPjmkf3qtVljqjA2M59TzNsVp4eA21Fr4N79kOuHZp+d3ZpqkweRH+ylrLmUdN+s/UFCGSHMNg8oeSq+EKJqUD8cUvzJHJBVU1U97tFnKSd+a+yTMYWyp+lJe7vvIZagdB04t89/1O/w1cDnyilFU=" // チャネルアクセストークン（長期）
	bot, err := linebot.New(Channel_Secret, Channel_Token)
	if err != nil {
		log.Println(err)
	}

	// LINE プラットフォームからのリクエストを受け取るための HTTP サーバを立ち上げる
	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				log.Println(event)
				switch message := event.Message.(type) { // ユーザが生徒の名前を入力
				case *linebot.TextMessage:
					minute, err := ScanTable(Db, message.Text) // 入力された生徒に対応する学習時間を調べる
					if err != nil {
						log.Println(err)
					}
					replyMessage := fmt.Sprintf(
						"%sさんが入室しました。%d 分後に学習終了時間をお知らせします。",
						message.Text, minute)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						log.Print(err)
					}
					time.Sleep(time.Duration(minute) * time.Minute)
					pushMessage := fmt.Sprintf("%sさんの学習終了時間となりました。", message.Text)
					userID := event.Source.UserID
					if _, err := bot.PushMessage(userID, linebot.NewTextMessage(pushMessage)).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	})

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

// Student テーブルの作成
func CreateTable(db *sql.DB) error {
	const sql = `CREATE TABLE IF NOT EXISTS student(
		Name STRING NOT NULL,
		Minute INTEGER NOT NULL
	);`

	_, err := db.Exec(sql)
	if err != nil {
		log.Println(err)
	}

	return nil
}

// Student テーブルに生徒情報を追加する
func InsertTable(db *sql.DB, student *Student) error {
	const sql = "INSERT INTO student(name, minute) VALUES (?, ?)"
	_, err := db.Exec(sql, student.Name, student.Minute)
	if err != nil {
		log.Fatalln(err)
	}
	return nil
}

// Student テーブルから情報を選択する
func ScanTable(db *sql.DB, name string) (int, error) {
	sql := `SELECT * FROM student WHERE Name = ?`
	var Minute int
	err := db.QueryRow(sql, name).Scan(&Minute)
	if err != nil {
		log.Println(err)
	}
	return Minute, nil
}
