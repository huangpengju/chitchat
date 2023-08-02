// data 包的 user.go 用于保存所有用户相关代码
// 包中 CreateSession 为现有用户创建一个会话，将用户会话保存到数据库中
// 包中 Check 检查会话在数据库中是否有效
// 包中 DeleteByUUID 从数据库sessions表中删除会话

// 包中 Create 创建一个新用户，将用户信息保存到数据库中
// 包中 UserByEmail 从数据库中获取给定电子邮件的单个用户

// 包中 User 从会话中获取用户
package data

import "time"

// User 表示论坛用户的账户
type User struct {
	Id        int
	Uuid      string
	Name      string
	Email     string
	Password  string
	CreatedAt time.Time
}

// Session 表示论坛用户当前的登录会话
type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

// Check 检查会话在数据库中是否有效
// 方法接收者是Session结构体
// 方法返回值valid = false 表示会话无效;valid=ture 表示有效
func (session *Session) Check() (valid bool, err error) {
	// QueryRow执行一次查询，并期望返回最多一行结果（即Row）。
	// QueryRow总是返回非nil的值，直到返回值的Scan方法被调用时，才会返回被延迟的错误。
	// Scan将该行查询结果各列分别保存进参数指定的值中。
	// 如果该查询匹配多行，Scan会使用第一行结果并丢弃其余各行。如果没有匹配查询的行，Scan会返回ErrNoRows。
	// Scan把查询结果保存进Session结构,err 是 nil，没有结果时，err 是 Scan返回的ErrNoRows。（sql: no rows in result set / sql:结果集中没有行）
	err = Db.QueryRow("select id,uuid,email,user_id,created_at from sessions where uuid=$1", session.Uuid).
		Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	if err != nil {
		valid = false
		return
	}
	if session.Id != 0 {
		valid = true
	}
	return
}

// DeleteByUUID 从数据库sessions表中删除会话
// 方法接收者是Session结构体
// 方法返回值err
func (session *Session) DeleteByUUID() (err error) {
	// 准备删除sessions表中数据的-SQL语句命令
	statement := "delete from sessions where uuid=$1"
	// Prepare创建一个准备好的状态用于之后的查询和命令。返回值可以同时执行多个查询和命令。
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	// Close关闭状态。
	defer stmt.Close()
	// Exec使用提供的参数执行准备好的命令状态，返回Result类型的该状态执行结果的总结。
	_, err = stmt.Exec(session.Uuid)
	return
}

// Create 创建一个新用户，将用户信息保存到数据库中
// 方法接收者是User结构体
// 方法返回值err
func (user *User) Create() (err error) {
	// Postgres不会自动返回最后一个插入id，因为这样假设是错误的
	// 你总是使用一个序列。您需要在insert中使用 returning 关键字来获得这个
	// 来自postgres的信息。
	// 准备SQL命令语句
	statement := "insert into users (uuid,name,email,password,created_at) values($1,$2,$3,$4,$5) returning id,uuid,created_at"
	// Prepare()创建一个准备好的状态
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	// Close关闭状态。
	defer stmt.Close()
	// QueryRow使用提供的参数执行准备好的查询状态
	// 使用QueryRow返回一行，并将返回的id扫描到User结构中
	// createUUID()随机生成uuid,Encrypt()加密密码
	err = stmt.QueryRow(createUUID(), user.Name, user.Email, Encrypt(user.Password), user.CreatedAt).Scan(&user.Id, &user.Uuid, &user.CreatedAt)
	return
}

// UserByEmail 获取给定电子邮件的单个用户
// 函数接受一个string 参数 email
// 函数返回获取到的 user 和 err
func UserByEmail(email string) (user User, err error) {
	user = User{}
	err = Db.QueryRow("select id,uuid,name,email,password,created_at from users where email = $1", email).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.Password, &user.CreatedAt)
	return
}

// CreateSession 为现有用户创建一个会话
// 方法的接收者是 user
// 方法返回 session 和 err
func (user *User) CreateSession() (session Session, err error) {
	// 准备SQL语句
	statement := "insert into sessions (uuid,email,user_id,created_at) values($1,$2,$3,$4) returning id,uuid,email,user_id,created_at"
	// 创建一个准备好的状态
	stmt, err := Db.Prepare(statement)
	if err != nil {
		return
	}
	// 关闭状态
	defer stmt.Close()
	// 使用QueryRow返回一行，并将返回的id扫描到Session结构中
	err = stmt.QueryRow(createUUID(), user.Email, user.Id, time.Now()).Scan(&session.Id, &session.Uuid, &session.Email, &session.UserId, &session.CreatedAt)
	return
}

// User 从会话中获取用户
// 方法的接收者是 session
// 返回值是 user 和 err
func (session *Session) User() (user User, err error) {
	// user 赋值为空
	user = User{}
	// 查询users表，条件是id
	err = Db.QueryRow("select id,uuid,name,email,created_at from users where id = $1", session.UserId).
		Scan(&user.Id, &user.Uuid, &user.Name, &user.Email, &user.CreatedAt)
	return
}

// Create a new session for an existing user
// Get the session for an existing user
// Check if session is valid in the database
// Delete session from database
// Get the user from the session
// Delete all sessions from database
// Create a new user, save user info into the database
// Delete user from database
// Update user information in the database
// Delete all users from database
// Get all users in the database and returns it
// Get a single user given the email
// Get a single user given the UUID

// 为现有用户创建一个新会话 (user *User) CreateSession()
// 获取现有用户的会话
// 检查会话在数据库中是否有效 (session *Session) Check()
// 从数据库中删除会话 (session *Session) DeleteByUUID()
// 从会话中获取用户  (session *Session) User()
// 从数据库中删除所有会话
// 创建新用户，将用户信息保存到数据库中 (user *User) Create()
// 从数据库中删除用户
// 更新数据库中的用户信息
// 从数据库中删除所有用户
// 获取数据库中的所有用户并返回它
// 获得给定电子邮件的单个用户 UserByEmail()
// 获取给定UUID的单个用户
