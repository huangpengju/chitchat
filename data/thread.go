// data 包的 thread.go 用于保存所有帖子相关代码
// data 包除了包含与数据库交互的结构和代码，还包含了一些与数据处理密切相关的函数。
// 包中 Threads 获取数据库中所有的帖子并返回
// 包中 User 获取帖子的作者信息
// 包中 CreatedAtDate 格式化帖子的CreatedAt日期，以便在屏幕上显示
// 包中 NumReplies 获取一个帖子的评论数
// 包中 CreateThread 创建一个新帖子
// 包中 ThreadByUUID 通过UUID获取帖子
// 包中 CreatePost 创建一个新的评论到一个帖子
package data

import (
	"time"
)

// 定义 Thread 结构，与创建关系数据库表 threads 时使用的数据定义语言（DDL）保持一致。
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

//  作者 {{ .User.Name }} - {{ .CreatedAtDate }} - {{ .NumReplies }} posts.

// 获取一个帖子的作者
// 方法的接收者：thread
// 方法的返回值：user
func (thread *Thread) User() (user User) {
	user = User{}
	Db.QueryRow("select id ,uuid,name,email,created_at from users where id = $1", thread.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

// CreatedAtDate 格式化CreatedAt日期，以便在屏幕上显示
// 方法的接收者：thread
// 返回值：字符串类型的格式化后的时间
func (thread *Thread) CreatedAtDate() string {
	return thread.CreatedAt.Format("2006/01/02 15:04")
}

// NumReplies 获取一个帖子的评论数
// 方法的接收者：thread
// 返回值：int类型的帖子的评论数量
func (thread *Thread) NumReplies() (count int) {
	rows, err := Db.Query("select count(*) from posts where thread_id = $1", thread.Id)
	if err != nil {
		return
	}
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			return
		}
	}
	rows.Close()
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
