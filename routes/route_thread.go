// routes 包的 route_thread.go 用于帖子的处理
//
// NewThread 显示发布帖子表单页面
// CreateThread 创建帖子
// PostThread 创建回复（评论）
// ReadThread 显示帖子的详细信息，包括评论和写评论的表单
package routes

import (
	"chitchat/data"
	"chitchat/utils"
	"fmt"
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

// POST /thread/post
// 创建回复（评论）
func PostThread(w http.ResponseWriter, r *http.Request) {
	// 检查是否登录
	sess, err := utils.Session(w, r)
	if err == nil {
		// 已登录
		// 先解析表达数据，把结果更新到PostForm
		err := r.ParseForm()
		if err != nil {
			utils.Danger(err, "无法解析表单")
		}
		// 根据session 中的UserId，查询用户信息
		user, err := sess.User()
		if err != nil {
			utils.Danger(err, "无法从session中获取用户")
		}
		// 获取PostForm中的值
		body := r.PostFormValue("body") // 评论的内容
		uuid := r.PostFormValue("uuid") // 帖子的UUID

		// 根据表单中的 uuid 获取帖子信息
		thread, err := data.ThreadByUUID(uuid)
		if err != nil {
			// 未获取到帖子信息，重定向到错误信息页面的方便函数
			utils.Error_message(w, r, "无法读取帖子")
		}
		// 使用 user 的方法 CreatePost 创建评论
		if _, err := user.CreatePost(thread, body); err != nil {
			utils.Danger(err, "无法创建评论")
		}
		// 设置URL，目的通过 uuid读取帖子
		url := fmt.Sprint("/thread/read?id=", uuid)
		// 跳转到指定 URL
		http.Redirect(w, r, url, http.StatusFound)
	} else {
		// 没有登录时跳转登录页
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// GET /thread/read
// 显示帖子的详细信息，(返回的Thread结构包含获取评论的方法)
func ReadThread(w http.ResponseWriter, r *http.Request) {
	// scheme://[userinfo@]host/path[?query][#fragment]
	// 获取 RawQuery string // 编码后的查询字符串，没有'?'
	vals := r.URL.Query()
	// 获取帖子的uuid
	uuid := vals.Get("id")

	// 通过UUID获取帖子
	thread, err := data.ThreadByUUID(uuid)
	if err != nil {
		utils.Error_message(w, r, "无法读取帖子")
	} else {
		_, err := utils.Session(w, r)
		if err != nil {
			utils.GenerateHTML(w, &thread, "layout", "public.navbar", "public.thread")
		} else {
			utils.GenerateHTML(w, &thread, "layout", "private.navbar", "private.thread")
		}
	}
}
