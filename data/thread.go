// data 包的 thread.go 用于数据库交互帖子相关数据
// 包中 Threads 获取数据库中所有的帖子并返回
// 包中 CreateThread 创建一个新帖子
// 包中 ThreadByUUID 通过UUID获取帖子
// 包中 CreatePost 创建一个新的评论到一个帖子
package data

import (
	"time"
)

// 定义帖子的结构
type Thread struct {
	Id        int
	Uuid      string
	Topic     string
	UserId    int
	CreatedAt time.Time
}

// 定义评论的结构
type Post struct {
	Id        int
	Uuid      string
	Body      string
	UserId    int
	ThreadId  int
	CreatedAt time.Time
}

// Threads 获取数据库中所有的帖子并返回
func Threads() (threads []Thread, err error) {
	rows, err := Db.Query("select id,uuid,topic,user_id,created_at from threads order by created_at desc")
	if err != nil {
		return
	}
	for rows.Next() {
		conv := Thread{}
		if err = rows.Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt); err != nil {
			return
		}
		threads = append(threads, conv)
	}
	// threads = []Thread{{Id: 1, Uuid: 1, Topic: "这是侯三", UserId: 1, CreatedAt: time.Now()}, {Id: 2, Uuid: 2, Topic: "李四李四", UserId: 2, CreatedAt: time.Now()}}
	return
}

// CreateThread 创建一个新帖子
// 方法接收者是User
// 方法的参数topic(帖子的主题)
// 返回创建的thread(帖子) 和 err
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

// ThreadByUUID 通过UUID获取帖子
// 函数接收帖子的uuid
// 函数返回一个Thread和err
func ThreadByUUID(uuid string) (conv Thread, err error) {
	// 初始化一个空的Thread
	conv = Thread{}
	err = Db.QueryRow("select id,uuid,topic,user_id,created_at from threads where uuid=$1", uuid).
		Scan(&conv.Id, &conv.Uuid, &conv.Topic, &conv.UserId, &conv.CreatedAt)
	return
}

// CreatePost创建一个新的评论到一个帖子
// 方法的接收者是 user
// 方法的参数是Thread(帖子) 和 body（评论的内容）
// 方法的返回值是Post(评论) 和 error
func (user *User) CreatePost(conv Thread, body string) (post Post, err error) {
	// 准备 SQL 语句
	statement := "insert into posts(uuid,body,user_id,thread_id,created_at) values($1,$2,$3,$4,$5) returning id,uuid,body,user_id,thread_id,created_at"
	// 准备命令
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	defer stmt.Close()
	// 使用QueryRow返回一行，并将返回的数据扫描到post结构中
	err = stmt.QueryRow(createUUID(), body, user.Id, conv.Id, time.Now()).
		Scan(&post.Id, &post.Uuid, &post.Body, &post.UserId, &post.ThreadId, &post.CreatedAt)
	return
}
