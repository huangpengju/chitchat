// data 包的 data.go 用于数据库配置、uuid随机生成、密码加密
//
// 包中 init 初始化数据库配置
// 包中 createUUID 创建一个随机UUID
// 包中 Encrypt 加密密码
package data

import (
	"crypto/sha1"
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	_ "github.com/lib/pq"
)

var Db *sql.DB

// init 初始化数据库配置
func init() {
	var err error
	// Db, err = sql.Open("postgres", "user=postgres dbname=chitchat password=Aa_123456 sslmode=disable")   // windows 环境
	Db, err = sql.Open("postgres", "postgres://root:123456@192.168.240.240:5432/chitchat?sslmode=disable") // Linux 环境
	if err != nil {
		fmt.Println("数据库链接失败：", err)
		log.Fatal(err)
	}
	fmt.Println("数据库连接成功~")
}

// 从RFC 4122创建一个随机UUID
// 改编自http://github.com/nu7hatch/gouuid
func createUUID() (uuid string) {
	u := new([16]byte)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Fatalln("无法生成UUID", err)
	}
	// 0x40 is reserved variant from RFC 4122
	// 0x40是RFC 4122的保留变体
	u[8] = (u[8] | 0x40) & 0x7F
	//设置数据的四位最高有效位(12到15位)
	// time_hi_and_version字段为4位版本号。
	u[6] = (u[6] & 0xF) | (0x4 << 4)
	uuid = fmt.Sprintf("%x-%x-%x-%x-%x", u[0:4], u[4:6], u[6:8], u[8:10], u[10:])
	return
}

// Encrypt 加密密码
// SHA-1的散列明文
// 该函数的参数是明文密码
// 该函数的返回值是密码文本
func Encrypt(plaintext string) (cryptext string) {
	cryptext = fmt.Sprintf("%x", sha1.Sum([]byte(plaintext)))
	return
}
