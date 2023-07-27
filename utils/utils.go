package utils

import (
	"chitchat/data"
	"errors"
	"fmt"
	"html/template"
	"net/http"
)

// session 处理器函数，检查用户是否已登录，已登陆是否有会话，如果没有，则err不为nil
// 返回值 Session 会话 和 err
// 先判断cookie，cookie不存在，用户未登录 err 为 http: named cookie not present（cookie没找到时）
// cookie存在，用户已登录，那么Session函数将继续进行第二项检查
func Session(w http.ResponseWriter, r *http.Request) (sess data.Session, err error) {
	// Cookie()返回请求中名为 _cookie 的cookie,如果未找到该cookie会返回nil和ErrNoCookie。
	// var ErrNoCookie = errors.New("http: named cookie not present") / http:命名cookie不存在
	// 如果找到了，返回*cookie 和 nil
	cookie, err := r.Cookie("_cookie")
	// 找到cookie 后，判断Session会话（第2项检查）
	if err == nil {
		// Session 表示论坛用户当前的登录会话
		// 把cookie中的value赋值给Session的Uuid，当作查询数据库的条件
		sess = data.Session{Uuid: cookie.Value}
		// Session存在时 err 依然是 nil
		// Session不存在时 err 返回errors.New("无效会话")
		// Check 查询sessions会话表,返回 true或者false
		if ok, _ := sess.Check(); !ok {
			err = errors.New("无效会话")
		}
	}
	return
}

// 解析HTML模板
// 传入一个文件名列表，并获得一个模板
// 返回 Template类型指针
// Template类型指针是text/template包的Template类型的特化版本，用于生成安全的HTML文本片段。
func ParseTemplateFiles(finenames ...string) (t *template.Template) {
	// 定义一个空切片
	var files []string
	t = template.New("layout") // 创建一个名为layout的模板。
	// 迭代取出文件名
	for _, file := range finenames {
		// 字符串与文件名拼接,并追加到切片files中
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	// ParseFiles函数接受可变参数(filenames ...string)，返回两个值(*Template, error)
	// ParseFiles函数创建一个模板并解析filenames指定的文件里的模板定义。
	// 返回的模板的名字是第一个文件的文件名（不含扩展名），内容为解析后的第一个文件的内容。至少要提供一个文件。如果发生错误，会停止解析并返回nil。
	// Must函数接受两个参数(*Template, error)，返回 *Template
	// Must函数用于包装(*Template, error)返回*template，它会在err非nil时panic，一般用于变量初始化：
	t = template.Must(t.ParseFiles(files...))
	return
}
