package utils

import (
	"chitchat/data"
	"errors"
	"net/http"
)

// session 处理器函数
// 检查用户是否已登录并有会话，如果没有，则错误不为nil
// 返回 Session 会话 和 err
// Session可用时 err 为nil
// Cookie不可用时 err 为 http: named cookie not present（cookie没找到时）
// Session不可用时 errors.New("无效会话")
func Session(w http.ResponseWriter, r *http.Request) (sess data.Session, err error) {
	// Cookie()返回请求中名为 _cookie 的cookie,如果未找到该cookie会返回nil和ErrNoCookie。
	// var ErrNoCookie = errors.New("http: named cookie not present") / http:命名cookie不存在
	// 如果找到了，返回*cookie 和 nil
	cookie, err := r.Cookie("_cookie")
	// 找到cookie 后，判断Session会话
	if err == nil {
		// Session 表示论坛用户当前的登录会话
		// 把cookie中的value赋值给Session的Uuid，当作查询数据库的条件
		sess = data.Session{Uuid: cookie.Value}
		// Check 查询sessions会话表,返回 true或者false
		if ok, _ := sess.Check(); !ok {
			err = errors.New("无效会话")
		}
	}
	return
}
