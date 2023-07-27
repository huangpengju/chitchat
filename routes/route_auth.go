// routes包中的route_auth用于授权等处理
package routes

import (
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
