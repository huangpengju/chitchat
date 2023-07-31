// utils 是一个工具类的包
// 包中 init 是初始化函数
// 包中 loadConfig 加载配置文件
// 包中 P 输出 ChitChat 论坛的一些信息
// 包中 Version 输出 ChitChat 论坛的版本
// 包中 Session 检查用户是否登录并有会话，如果不是，err不是nil
// 包中 ParseTemplateFiles 解析登录页的HTML模板
// 包中 GenerateHTML 根据参数生成 HTML 页面
// 包中 Warning 函数输出 警告相关的日志
// 包中 Danger 函数输出风险相关日志
// 包中 Error_message 重定向到错误信息页面的方便函数
package utils

import (
	"chitchat/data"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

// 配置
type Configuration struct {
	Address       string
	ReadTimeout   int64
	WriterTimeout int64
	Static        string
}

// 定义全局变量 config 配置
var Config Configuration

// 定义全局变量 logger 日志记录器
var logger *log.Logger

func init() {
	// 加载配置文件
	loadConfig()
	file, err := os.OpenFile("chitchat.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("日志文件打开失败", err)
	}
	logger = log.New(file, "INFO", log.Ldate|log.Ltime|log.Lshortfile)
}

// loadConfig 加载配置文件
func loadConfig() {
	file, err := os.Open("config.json")
	if err != nil {
		log.Fatalln("无法打开配置文件:", err)
	}
	// 创建一个解码器
	decoder := json.NewDecoder(file)
	// 声明一个空的Configuration结构
	Config = Configuration{}
	// 读取json并保存在config结构中，返回一个err
	err = decoder.Decode(&Config)
	if err != nil {
		log.Fatalln("无法从文件中获取配置:", err)
	}
}

// P 函数方便打印到标准输出
func P(a ...interface{}) {
	fmt.Printf("%s", a)
	// fmt.Println(a)
}

// Version 显示版本
func Version() string {
	return "0.1"
}

// Session 检查用户是否登录并有会话，如果未登录，err不为nil
// 返回值 Session 会话 和 err
// 如果cookie不存在，那么很明显用户并未登陆,用户未登录 err 为 http: named cookie not present（cookie没找到时）
// 如果cookie存在，那么Session函数将继续进行第二项检查,访问数据库并核实会话的唯一ID是否存在。
func Session(w http.ResponseWriter, r *http.Request) (sess data.Session, err error) {
	// Cookie()返回请求中名为 _cookie 的cookie,如果未找到该cookie会返回nil和ErrNoCookie。
	// var ErrNoCookie = errors.New("http: named cookie not present") / http:命名cookie不存在
	// 如果找到了，返回*cookie 和 nil
	cookie, err := r.Cookie("_cookie")

	// cookie存在，用户已登录，那么Session函数将继续进行第二项检查
	if err == nil {
		// Session 是登录会话结构体
		// 给Session结构体重的 Uuid 字段赋值
		// 把cookie中的value赋值给Session的Uuid，用作查询数据库的条件
		sess = data.Session{Uuid: cookie.Value}

		// Check 检查会话在数据库中是否有效
		// 会话有效返回 true
		// 会话无效返回 false
		if ok, _ := sess.Check(); !ok {
			// 用户登录，但是会话无效时，设置返回值 err 的值
			err = errors.New("无效会话")
		}
	}
	return
}

// ParseTemplateFiles 解析登录页的HTML模板
// 传入一个文件名列表，并获得一个模板
// 返回 Template类型指针
// Template类型指针是text/template包的Template类型的特化版本，用于生成安全的HTML文本片段。
func ParseTemplateFiles(finenames ...string) (t *template.Template) {
	// 定义一个空切片
	var files []string
	t = template.New("layout") // 创建一个名为layout的模板。
	// 迭代取出文件名
	for _, file := range finenames {
		// 字符串templates与文件名拼接,并追加到切片files中
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

// generateHTML 根据参数生成 HTML 页面
func GenerateHTML(w http.ResponseWriter, data interface{}, filesname ...string) {
	// 声明一个 string 类型的切片
	var files []string
	// 迭代取出文件名
	for _, file := range filesname {
		// 字符串templates与文件名拼接,并追加到切片files中
		files = append(files, fmt.Sprintf("templates/%s.html", file))
	}
	// ParseFiles函数接受可变参数(filenames ...string)，返回两个值（*Template,err）
	// ParseFiles函数用来创建一个模板并解析参数filenames中指定的文件里的模板定义。
	// 返回的模板的名字是第一个文件的文件名（不含扩展名），内容为解析后的第一个文件的内容。至少要提供一个文件。如果发生错误，会停止解析并返回nil。
	// Must函数接受两个参数(*Template, error)，返回 *Template
	// 它会在err非nil时panic，一般用于变量初始化：
	templates := template.Must(template.ParseFiles(files...))
	// ExecuteTemplate 接受3个参数，(wr io.Writer, name string, data interface{})
	// 返回 error
	// 类型 Execute 执行
	// ExecuteTemplate 使用名 layout 关联的模板产生输出。
	templates.ExecuteTemplate(w, "layout", data)
}

// Warning 函数输出警告相关的日志
func Warning(args ...interface{}) {
	// 设置logger的输出前缀。
	logger.SetPrefix("WARNING")
	// Println调用l.Output将生成的格式化字符串输出到logger，参数用和fmt.Println相同的方法处理。
	logger.Println(args...)
}

// Danger 函数输出风险相关日志
func Danger(args ...interface{}) {
	logger.SetPrefix("ERROR")
	logger.Println(args...)
}

// Error_message 重定向到错误信息页面的方便函数
func Error_message(w http.ResponseWriter, r *http.Request, msg string) {
	// 定义一个string 切片
	url := []string{"/err?msg=", msg}
	// 使用 strings.Join() 让 String 切片中的字符串使用""拼接起来
	// Redirect 重定向到错误信息页面
	http.Redirect(w, r, strings.Join(url, ""), http.StatusFound)
}
