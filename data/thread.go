// data 包的 thread.go 用于数据库交互帖子相关数据
// 包中 Threads 获取数据库中所有的帖子并返回
// CreateThread 创建一个新帖字
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

// CreateThread 创建一个新帖字
// 方法接收者是User
// 返回创建的帖子 和 err
func (user *User) CreateThread(topic string) (conv Thread, err error) {
	// 准备SQL语句
	statement := "insert into threads (uuid,topic,user_id,created_at) values($1,$2,$3,$4) returning id,uuid,topic,user_id,created_at"
	// 准备命令
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// 使用QueryRow返回一行，并将返回的id扫描到Thread结构中
	err = stmt.QueryRow(createUUID(), topic, user.Id, time.Now()).Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, conv.CreatedAt)
	return

}
