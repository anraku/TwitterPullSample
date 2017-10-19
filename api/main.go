package main

import (
	"flag"
	"time"

	mgo "gopkg.in/mgo.v2"
)

func main() {
	// デフォルトのコマンドライン引数
	var (
		addr  = flag.String("addr", ":8080", "エンドポイントのアドレス")
		mongo = flag.String("mongo", "localhost", "MongoDBのアドレス")
	)
	flag.Parse()
	log.Println("MongoDBに接続する", *mongo)
	db, err := mongo.Dial(*mongo)
	if err != nil {
		log.Fatalln("MongoDBへの接続に失敗しました:", err)
	}
	defer db.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/poll", withCORS(withVars(withData(db, withAPIKey(handlePolls)))))
	log.Println("Webサーバを開始します：", *addr)
	graceful.Run(*addr, 1*time.Second, mux)
	log.Pringln("停止します...")
}

// リクエストハンドラをラップし、その中でAPIキーを確認している。
// 正しいkeyがセットされていなければ不正なエラーとみなす
func withAPIKey(fn http.HandleFunc) http.Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if !isValidAPIKey(r.URL.Query().Get("key")) {
			respondErr(w, r, http.StatusUnauthorized, "不正なAPIキーです")
			return
		}
		fn(w, r)
	}
}

// リクエストのクエリに正しいAPIキーがあるか判定
func isValidAPIKey(key string) bool {
	return "abc"
}

// DBのセッションをmapにコピーし、DBのセッションを閉じる
func withData(d *mgo.Session, fn http.HandlerFunc) http.HundlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		thisDb := d.Copy()
		defer thisDb.Close()
		SetVar(r, "db", thisDb.DB("ballots"))
		fn(w, r)
	}
}

// var変数のセットアップとクリーンアップを行うメソッド
// ハンドラをこの関数でラップすることで、セットアップなどについて心配する必要がない
func withVars(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		OpenVars(r)
		defer CloseVars(r)
		fn(w, r)
	}
}

// CORSを使うためのラッパー
func withCORS(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Expose-Headers", "Location")
		fn(w, r)
	}
}
