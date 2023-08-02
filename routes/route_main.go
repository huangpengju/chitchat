// routes包中的route_main.go 用于授权相关的处理
//
// Index 处理器函数 (帖子列表)
//
// 包中 Err 显示错误消息页面
package routes

import (
	"chitchat/data"
	"chitchat/utils"
	"net/http"
)

// Index 处理器函数(帖子列表)
// 负责生成 HTML 并将其写入 ResponseWriter 中
func Index(w http.ResponseWriter, r *http.Request) {
	// 获取全部的帖子
	threads, err := data.Threads()
	if err != nil {
		utils.Error_message(w, r, "无法获取帖子")
	} else {
		// 检查用户是否登录并有会话，如果未登录，err不为nil
		_, err := utils.Session(w, r)
		if err != nil {
			utils.GenerateHTML(w, threads, "layout", "public.navbar", "index")
		} else {
			utils.GenerateHTML(w, threads, "layout", "private.navbar", "index")
		}
	}
}

// // index 处理器函数
// // 负责生成 HTML 并将其写入 ResponseWriter 中
// func index(w http.ResponseWriter, r *http.Request) {
// 	// 定义 files 模板切片 存放 布局文件、标题文件、主页文件路径
// 	files := []string{"templates/layout.html",
// 		"templates/public.navbar.html",
// 		"templates/index.html"}

// 	// ParseFiles 分析文件
// 	// 创建一个模板，并解析 files 指定的文件里的模板定义，
// 	// 返回的模板的名字是第一个文件的文件名（不含扩展名）,内容为解析后的第一个文件的内容。
// 	// Must 用于包装返回 模板指针
// 	templates := template.Must(template.ParseFiles(files...))

// 	// 查询所有帖子
// 	threads, err := data.Threads()
// 	if err != nil {
// 		return
// 	}
// 	// 让 templates 关联的名为 layout 模板产生输出 threads 帖子
// 	templates.ExecuteTemplate(w, "layout", threads)
// }

// GET /err?msg=
// Err 显示错误消息页面
func Err(w http.ResponseWriter, r *http.Request) {
	// 判断用户是否登陆（检查cookie和session会话）

	// URL类型代表一个解析后的URL（或者说，一个URL参照）。URL基本格式如下：scheme://[userinfo@]host/path[?query][#fragment]
	// Query方法解析RawQuery字段并返回其表示的Values类型键值对。
	// 而对于 字段 RawQuery string  它表示URL中的： 编码后的query，没有'?'
	// 比如：http:baidu.com/search?title=aa&id=1,其RawQuery是"title=aa&id=1"
	// Query方法解析RawQuery字段并返回其表示的Values类型键值对。map[id:[1] title:[aa]]
	vals := r.URL.Query() // 返回一个映射map[],映射的键值对由 URL中 ? 后面的内容组成

	// 查询Session会话
	// 如果cookie存在，并且数据库中存在session会话记录，此时 err 为nil
	_, err := utils.Session(w, r) // 返回 Session 结构 和 err（ nil 或 cookie没找到/无效会话）

	// err 为 nil 时表示已登录，加载 私有模板
	// err 不为 nil 时表示未登录，加载 公共模板
	if err != nil {
		// generateHTML 生成注册页的HTML
		// 解析HTML模板
		// 参数2是一个映射map[]；
		// Get会获取键——msg 对应的值集的第一个值。
		// 如果没有对应key的值集会返回空字符串。获取值集请直接用map。
		// 传入一个文件名列表(公有的布局框架、公共导航条、错误提示页)，并获得一个模板
		utils.GenerateHTML(w, vals.Get("msg"), "layout", "public.navbar", "error")
		// fmt.Println("公共模板")
	} else {
		// 传入一个文件名列表(公有的布局框架、私有导航条、错误提示页)，并获得一个模板
		utils.GenerateHTML(w, vals.Get("msg"), "layout", "private.navbar", "error")
	}
}
