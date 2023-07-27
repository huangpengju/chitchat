// routes包中的route_auth用于授权等处理
package routes

import (
	"chitchat/data"
	"chitchat/utils"
	"net/http"
)

// GET /login
// 显示登录页面
func Login(w http.ResponseWriter, r *http.Request) {
	// 解析HTML模板
	// 传入一个文件名列表(登录页框架、公共导航条、登录form表单)，并获得一个模板
	t := utils.ParseTemplateFiles("login.layout", "public.navbar", "login")
	// Execute方法接受两个参数
	// 将解析好的模板应用到指定的数据对象 nil ,并将输出写入 w
	// 如果在执行模板或写输出时出错，
	// 执行停止，但是部分结果可能已经被写入
	// 输出写入器。
	// 模板可以安全地并行执行，尽管如果是并行
	// 执行共享一个Writer，输出可能交错。
	t.Execute(w, nil)
}

// GET /logout
// 注销用户
func Logout(w http.ResponseWriter, r *http.Request) {
	// Cookie() 接受一个参数，并返回两个结果
	cookie, err := r.Cookie("_cookie") // 返回请求中名为 _cookie 的cookie 和 error

	// 如果未找到该 cookie 会返回的 结果1=nil；结果2=ErrNoCookie。  // var ErrNoCookie = errors.New("http: named cookie not present") / http:命名cookie不存在
	// 如果找到了，返回结果1=*cookie；结果2=nil
	if err != http.ErrNoCookie {
		// 找到 cookie后
		// 声明一个Session 会话结构体，并给 Uuid 字段赋值 cookie.Value
		session := data.Session{Uuid: cookie.Value}
		// 注销用户，使用 Uuid 作为条件，从数据库中删除会话
		session.DeleteByUUID()
	}
	utils.Warning(err, "获取cookie失败")
}
