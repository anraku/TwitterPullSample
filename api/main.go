package main

import mgo "gopkg.in/mgo.v2"

func main() {

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
