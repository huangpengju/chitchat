// routes包中的route_main.go 用于授权相关的处理
// 包中 Err 显示错误消息页面
package routes

import (
	"fmt"
	"net/http"
)

// GET/err?msg=
// Err 显示错误消息页面
// 判断用户是否登陆（检查cookie和session会话）
func Err(w http.ResponseWriter, r *http.Request) {
	fmt.Println("这里是err")
	fmt.Println("r===", r)
	fmt.Println("r.URL===", r.URL)
	fmt.Println("r.URL.Path===", r.URL.Path)
	fmt.Println("r.URL.Query()===", r.URL.Query()) // map[] 空的映射  // map[string][]string
	// URL类型代表一个解析后的URL（或者说，一个URL参照）。URL基本格式如下：scheme://[userinfo@]host/path[?query][#fragment]
	// Query方法解析RawQuery字段并返回其表示的Values类型键值对。
	// 而对于 字段 RawQuery string  它表示URL中的： 编码后的query，没有'?'
	// 比如：http:baidu.com/search?title=aa&id=1,其RawQuery是"title=aa&id=1"
	// Query方法解析RawQuery字段并返回其表示的Values类型键值对。map[id:[1] title:[aa]]

	// vals := r.URL.Query()

	// 查询Session会话
	// 如果cookie不存在，那么很明显用户并未登陆
	// 如果cookie存在，那么Session函数将继续进行第二项检查,访问数据库并核实会话的唯一ID是否存在。
	// _, err := utils.Session(w, r) // 返回 Session 会话 和 err
	// cookie存在，并且Session会话在数据库中时，err 为nil
	// err 为nil 时表示已登录，加载 私人模板
	// err 不为nil 时表示未登录，加载 公共模板
	// if err != nil {
	// 	// Get会获取key对应的值集的第一个值。如果没有对应key的值集会返回空字符串。获取值集请直接用map。
	// 	// generateHTML(w, vals.Get("msg"), "layout", "public.navbar", "error")
	// 	fmt.Println("公共模板")
	// } else {
	// 	// generateHTML(w, vals.Get("msg"), "layout", "private.navbar", "error")
	// 	fmt.Println("私人模板")
	// }
}
