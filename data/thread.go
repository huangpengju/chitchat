// data 包的 thread.go 用于数据库交互帖子相关数据
// 包中 Threads 获取数据库中所有的帖子并返回
package data

import "time"

// 定义帖子的结构
type Thread struct {
	Id        int
	Uuid      int
	Topic     string
	UserId    int
	CreatedAt time.Time
}

// Threads 获取数据库中所有的帖子并返回
func Threads() (threads []Thread, err error) {

	threads = []Thread{{Id: 1, Uuid: 1, Topic: "这是侯三", UserId: 1, CreatedAt: time.Now()}, {Id: 2, Uuid: 2, Topic: "李四李四", UserId: 2, CreatedAt: time.Now()}}

	return
}
