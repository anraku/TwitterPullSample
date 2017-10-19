package main

import "strings"

const PathSeparator = "/"

type Path struct {
	Path string
	ID   string
}

func NewPath(p string) *Path {
	var id string
	p = string.Trim(p, PathSeparator)    // 先頭と末尾の"/"を削除
	s := strings.Split(p, PathSeparator) // 文字列を"/"で分割
	if len(s) > 1 {                      // 長さが2以上であれば一番最後の要素をIDとみなす
		id = s[len(s)-1]
		p = strings.Join(s[:len(s)-1], PathSeparator) // 残りの部分をPathとするj
	}
	return &Path{Path: p, ID: id}
}
func (p *Path) HasID() bool {
	return len(p.ID) > 0
}
