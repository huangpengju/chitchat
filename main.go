package main

import (
	"chitchat/data"
	"chitchat/routes"
	"net/http"
	"text/template"
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

	// index
	// 把给定的 URL 请求转发至 index 处理器函数
	mux.HandleFunc("/", index)
	// error
	// 用户登录显示私人模板，未登录显示公共模板
	mux.HandleFunc("/err", routes.Err)

	// mux.HandleFunc("/login",login)
	// mux.HandleFunc("/logout",logout)
	// mux.HandleFunc("/signup",signup)
	// mux.HandleFunc("/signup_account",signupAccount)
	// mux.HandleFunc("/authenticate",authenticate)

	// mux.HandleFunc("/thread/new",newThread)
	// mux.HandleFunc("/thread/create",createThread)
	// mux.HandleFunc("/thread/post",postThread)
	// mux.HandleFunc("/thread/read",readThread)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	server.ListenAndServe()
}

// index 处理器函数
// 负责生成 HTML 并将其写入 ResponseWriter 中
func index(w http.ResponseWriter, r *http.Request) {
	// 定义 files 模板切片 存放 布局文件、标题文件、主页文件路径
	files := []string{"templates/layout.html",
		"templates/public.navbar.html",
		"templates/index.html"}

	// ParseFiles 分析文件
	// 创建一个模板，并解析 files 指定的文件里的模板定义，
	// 返回的模板的名字是第一个文件的文件名（不含扩展名）,内容为解析后的第一个文件的内容。
	// Must 用于包装返回 模板指针
	templates := template.Must(template.ParseFiles(files...))

	// 查询所有帖子
	threads, err := data.Threads()
	if err != nil {
		return
	}
	// 让 templates 关联的名为 layout 模板产生输出 threads 帖子
	templates.ExecuteTemplate(w, "layout", threads)
}
