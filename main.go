package main

import (
	"chitchat/routes"
	"net/http"
)

func main() {

	// 创建一个多路复用器
	mux := http.NewServeMux()

	// 服务静态文件
	// 准备 /public 目录的HTTP处理器
	files := http.FileServer(http.Dir("/public"))
	//
	// 把 "/static" 和 处理器files 注册到多路复用器
	// StripPrefix 会把 URL中 /static/ 去除后  再让 files 处理器处理请求
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	//
	// 所有的路由都在这里匹配
	// 在其他文件中定义的路由处理函数
	//

	// Index (帖子列表)
	// 把给定的 URL 请求转发至 index 处理器函数
	mux.HandleFunc("/", routes.Index)

	// 在 routes 包中的 route_main.go 中定义
	// Err 判断用户是否登陆（检查cookie和session会话）
	// 用户登录后加载私人模板，用户未登录加载公共模板
	mux.HandleFunc("/err", routes.Err)

	// 在 routes 包中的 route_auth.go 中定义
	// login 加载登录页面
	mux.HandleFunc("/login", routes.Login)
	// logout 注销用户，根据cookie删除session会话
	mux.HandleFunc("/logout", routes.Logout)
	// signup 加载注册页面
	mux.HandleFunc("/signup", routes.Signup)
	// signupAccount 创建用户账户
	mux.HandleFunc("/signup_account", routes.SignupAccount)
	// Authenticate 用户登录（登录成功创建 session 和 cookie）
	mux.HandleFunc("/authenticate", routes.Authenticate)

	// 在 routes 包中的 route_thread.go 中定义
	// NewThread 显示发布帖子表单页面
	mux.HandleFunc("/thread/new", routes.NewThread)
	// CreateThread 创建帖子
	mux.HandleFunc("/thread/create", routes.CreateThread)
	// PostThread 创建帖子的回复
	mux.HandleFunc("/thread/post", routes.PostThread)
	// ReadThread 显示帖子的详细信息
	mux.HandleFunc("/thread/read", routes.ReadThread)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}
