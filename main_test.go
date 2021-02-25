package main

import (
	"database/sql"
	"log"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestofDatabase(t *testing.T) {
	type Student struct {
		Name   string // 生徒の名前
		Minute int    // 対象の生徒の学習制限時間
	}
	students := []*Student{{Name: "鈴宮 花子", Minute: 1}, {Name: "鈴宮 太郎", Minute: 2}, {Name: "鈴宮 次郎", Minute: 5}}
	Db, _ := sql.Open("sqlite3", "database_test.db")
	CreateTable(Db)
	for i := range students {
		if err := InsertTable(Db, students[i]); err != nil {
			log.Println(err)
		}

	}

	cases := []string{"鈴宮 花子", "鈴宮 太郎", "鈴宮 次郎"}
	for i := range cases {
		expect := students[i].Minute
		actual, _ := ScanTable(Db, cases[i])

		if actual != expect {
			t.Errorf("actual: %v, expected: %v", actual, expect)
		}
	}

}
