// routes 包的 route_auth.go 用于授权相关的处理
// 包中 Login 显示登录页面
// 包中 Logout 注销用户,从数据库中删除session会话
// 包中 Signup 显示注册页面
// 包中 SignupAccout 创建用户账户
package routes

import (
	"chitchat/data"
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
	// 将解析好的模板应用到指定的数据对象(这里是nil) ,并将输出写入 w
	// 如果在执行模板或写输出时出错，
	// 执行停止，但是部分结果可能已经被写入
	// 输出写入器。
	// 模板可以安全地并行执行，尽管如果是并行
	// 执行共享一个Writer，输出可能交错。
	t.Execute(w, nil)
}

// GET /logout
// 注销用户
func Logout(w http.ResponseWriter, r *http.Request) {
	// Cookie() 接受一个参数，并返回两个结果
	cookie, err := r.Cookie("_cookie") // 返回请求中名为 _cookie 的cookie 和 error

	// 如果未找到该 cookie 会返回的 结果1=nil；结果2=ErrNoCookie。  // var ErrNoCookie = errors.New("http: named cookie not present") / http:命名cookie不存在
	// 如果找到了，返回结果1=*cookie；结果2=nil
	if err != http.ErrNoCookie {
		// 找到 cookie后
		// 声明一个Session 会话结构体，并给 Uuid 字段赋值 cookie.Value
		session := data.Session{Uuid: cookie.Value}
		// 注销用户，使用 Uuid 作为条件，从数据库中删除会话
		session.DeleteByUUID()
	}
	utils.Warning(err, "获取cookie失败")
}

// GET /signup
// 显示注册页面
func Signup(w http.ResponseWriter, r *http.Request) {
	// 解析HTML模板
	// 传入一个文件名列表(登录页框架、公共导航条、注册form表单)，并获得一个模板
	utils.GenerateHTML(w, nil, "login.layout", "public.navbar", "signup")
}

// POST /signup
// 创建用户帐户
func SignupAccount(w http.ResponseWriter, r *http.Request) {
	// ParseForm解析URL中的查询字符串，并将解析结果更新到r.Form字段。
	// 对于POST或PUT请求，ParseForm还会将body当作表单解析
	// 并将结果既更新到r.PostForm也更新到r.Form。
	err := r.ParseForm()
	if err != nil {
		utils.Danger(err, "无法分析表单")
	}
	// User结构
	// Id       int
	// Uuid     string
	// Name     string
	// Email    string
	// Password string
	// CreateAt time.Time
	user := data.User{
		// PostFormValue返回name、email、password为键
		// 查询r.PostForm字段(本字段只有在调用ParseForm后才有效,上面代码已调用)
		// 得到结果 []string切片的第一个值，并赋值给 User结构中的字段
		Name:     r.PostFormValue("name"),
		Email:    r.PostFormValue("email"),
		Password: r.PostFormValue("password"),
	}
	// Create 创建一个新用户，将用户信息保存到数据库中
	// 创建失败 返回 err (err不为nil)
	if err := user.Create(); err != nil {
		utils.Danger(err, "无法创建用户")
	}
	// 账号注册成功，跳转到登录页
	// Redirect回复请求一个重定向地址urlStr和状态码code。
	// 该重定向地址可以是相对于请求r的相对地址。
	// http.StatusFound 是302 表示建立连接状态 //向IANA注册的HTTP状态代码。
	http.Redirect(w, r, "/login", http.StatusFound)
}
