// routes 包的 route_thread.go 用于帖子的处理
// NewThread 显示发布帖子表单页面
// CreateThread 创建帖子
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

// POST /thread/create
// 创建帖子
func CreateThread(w http.ResponseWriter, r *http.Request) {
	// Session 检查用户是否登录并有会话，如果不是，err不是nil
	sess, err := utils.Session(w, r)
	if err != nil {
		// 用户未登录
		// 跳转到登录页
		http.Redirect(w, r, "/login", http.StatusFound)
	} else {
		// 用户已登录

		// ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段。
		// 对于POST或PUT请求，ParseForm还会将body当作表单解析
		// 并将结果既更新到r.PostForm也更新到r.Form。
		err = r.ParseForm()
		if err != nil {
			utils.Danger(err, "无法解析表单")
		}
		// 确认用户信息
		user, err := sess.User()
		if err != nil {
			utils.Danger(err, "无法从会话中获取用户")
		}
		// 获取主题内容
		topic := r.PostFormValue("topic")
		// 创建帖子
		if _, err := user.CreateThread(topic); err != nil {
			utils.Danger(err, "无法创建帖子")
		}
		// 跳转链接到主页
		http.Redirect(w, r, "/", http.StatusFound)
	}
}
