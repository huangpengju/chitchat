// routes 包的 route_thread.go 用于帖子的处理
// NewThread 显示发布帖子表单页面
package routes

import (
	"chitchat/utils"
	"net/http"
)

// GET /threads/new
// 显示发布帖子表单页面
func NewThread(w http.ResponseWriter, r *http.Request) {
	// Session 检查用户是否登录并有会话，如果不是，err不是nil
	_, err := utils.Session(w, r)
	if err != nil {
		// 用户未登录

		// 跳转到登录页
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		// 用户登录

		// 解析HTML模板
		// 传入一个文件名列表(私有的导航条、发布帖子表单)，并获得一个模板
		utils.GenerateHTML(w, nil, "layout", "private.navbar", "new.thread")
	}
}
