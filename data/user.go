// data 是数据包
package data

import "time"

// Session 表示论坛用户当前的登录会话
type Session struct {
	Id        int
	Uuid      string
	Email     string
	UserId    int
	CreatedAt time.Time
}

// Check
// 检查会话在数据库中是否有效
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
