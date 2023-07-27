// routes包中的route_main.go 用于授权相关的处理
// 包中 Err 显示错误消息页面
package routes

import (
	"chitchat/utils"
	"net/http"
)

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
