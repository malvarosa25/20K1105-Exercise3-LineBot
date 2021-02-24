package main

import (
	"fmt"
	"net/http"

	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	Channel_Secret := "857d036768e6c23dd8731bec8d08312f"                                                                                                                                            // チャネルシークレット
	Channel_Token := "+MVr5jo/PqWuzYfQ8G3DZyFPjmkf3qtVljqjA2M59TzNsVp4eA21Fr4N79kOuHZp+d3ZpqkweRH+ylrLmUdN+s/UFCGSHMNg8oeSq+EKJqUD8cUvzJHJBVU1U97tFnKSd+a+yTMYWyp+lJe7vvIZagdB04t89/1O/w1cDnyilFU=" // チャネルアクセストークン（長期）
	bot, err := linebot.New(Channel_Secret, Channel_Token)
	if err != nil {
		return fmt.Errorf("can't return a new linebot: %w", err)
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
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text)).Do(); err != nil {
						return fmt.Errorf("can't return a message: %w", err)
					}
				case *linebot.StickerMessage:
					replyMessage := fmt.Println("生徒の ID を入力してください")
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
						return fmt.Errorf("can't return a message: %w", err)
					}
			}
		}
	}
}
