package main

import (
	"net/http"
	"sync"
)

var (
	// 外部からのデータベースセッションをmapで管理する
	// 複数のクライアントから変更されるのでMutexで排他制御を行う
	varsLock sync.Mutex
	vars     map[*http.Request]map[string]interface{}
)

// クライアントから接続があるとリクエストオブジェクトをkeyにしてインスタンスを作成
func OpenVars(r *http.Request) {
	varsLock.Lock()
	if vars == nil {
		vars = map[*http.Request]map[string]interface{}{}
	}
	vars[r] = map[string]interface{}{}
	varsLock.Unlock()
}

// クライアントとの接続が終了する際、mapに登録されているデータを安全に削除する
func CloseVars(r *http.Request) {
	varsLock.Lock()
	delete(vars, r)
	varsLock.Unlock()
}

// mapから値を取り出す関数
// RLock()は他のコードのRLock()をブロックしない
func GetVar(r *http.Request, key string) {
	varsLock.RLock()
	value := vars[r][key]
	varsLock.RUnlock()
	return value
}

// mapに値をセットする関数
func SetVar(r *http.Request, key string, value interface{}) {
	varsLock.Lock()
	vars[r][key] = value
	varsLock.Unlock()
}
