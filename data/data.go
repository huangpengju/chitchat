// data 包的 data.go 用于数据库配置
// 包中 init 初始化数据库配置
package data

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var Db *sql.DB

// init 初始化数据库配置
func init() {
	var err error
	Db, err = sql.Open("postgres", "user=postgres dbname=chitchat password=Aa_123456 sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
}
