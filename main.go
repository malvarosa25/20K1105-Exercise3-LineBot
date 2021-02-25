package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
	_ "github.com/mattn/go-sqlite3" // 使用しないため、 _ にしないとコンパイルエラーとなる
)

// 生徒のデータ格納用のユーザ定義型
type Student struct {
	Name string // 生徒の名前
	Time int    // 対象の生徒の学習制限時間
}

func main() {
	// 生徒の情報が載っているデータベースを作成
	Db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		log.Println(err)
	}

	if err := createTable(Db); err != nil {
		log.Println(err)
	}

	students := []*Student{{Name: "鈴宮 花子", Time: 5}, {Name: "鈴宮 太郎", Time: 10}, {Name: "鈴宮 次郎", Time: 15}}

	if err := insertTable(Db, students); err != nil {
		log.Println(err)
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
					time, err := scanTable(Db, message.Text) // 入力された生徒に対応する学習時間を調べる
					if err != nil {
						log.Println(err)
					}
					replyMessage := fmt.Sprintf(
						"%s さんが入室しました。%d 分後に学習終了時間をお知らせします。",
						message.Text, time)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
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
func createTable(db *sql.DB) error {
	const sql = `CREATE TABLE IF NOT EXISTS student(
		Name TEXT NOT NULL,
		Time INTEGER NOT NULL
	);`

	_, err := db.Exec(sql)
	if err != nil {
		log.Println(err)
	}

	return nil
}

// Student テーブルに生徒情報を追加する
func insertTable(db *sql.DB, students []*Student) error {
	for i := range students {
		const sql = "INSERT INTO student(name, time) VALUES (?, ?)"
		_, err := db.Exec(sql, students[i].Name, students[i].Time)
		if err != nil {
			log.Println(err)
		}

	}
	return nil
}

// Student テーブルから情報をスキャンする
func scanTable(db *sql.DB, name string) (int, error) {
	rows, err := db.Query("SELECT * FROM student WHERE Name = ?", name)
	if err != nil {
		log.Println(err)
	}

	var s Student
	err = rows.Scan(&s.Name, &s.Time)

	if err != nil {
		log.Println(err)
	}

	if err := rows.Err(); err != nil {
		log.Println(err)
	}

	return s.Time, nil
}
