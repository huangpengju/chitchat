# ChitChat 论坛
一个简单的网上论坛 Web 应用。它允许用户登录到论坛里面，然后在论坛上发布新帖子，又或者回复其他用户发表的帖子。

## 1.ChitChat 简介
这是一个名为ChitChat的简易网上论坛，ChitChat实现网上论坛的关键特性：  
* 在这个论坛里面，用户可以注册账号，并在登录之后发表新帖子又或者回复已有的帖子；
* 未注册用户可以查看帖子，但是无法发表帖子或回复帖子。

## 2.应用设计
Web应用的一般工作流程，客户端向服务器发送请求，然后等待接收响应。

**客户端** → 发送HTTP请求 → **服务器**  
　　　　　　　　　　　　　↓（处理HTTP请求）  
**客户端** ← 返回HTTP响应 ← **服务器**

ChitChat的应用逻辑会被编码到服务器里面。服务器会向客户端提供HTML页面，并通过页面的超链接向客户端表名请求格式以及被请求的数据，而客户端会在发送请求时向服务器提供相应的数据。  

ChitChat的请求使用的格式：**http://<服务器名>/<处理器名>?<参数>**。
服务器名是ChitChat服务器的名字，而处理器名则是被调用的处理器的名字。  
该应用的参数会以URL查询的形式传递给处理器，而处理器则会根据这些参数对请求进行处理。

当请求到达服务器时，多路复用器（multiplexer）会对请求进行检查，并将请求重定向至正确的处理器进行处理。处理器在接收到多路复用器转发的请求之后，会从请求中取出相应的信息，并根据这些信息对请求进行处理。在请求处理完毕之后，处理器会将所得的数据传递给模板引擎，而模板引擎则会根据这些数据生成将要返回给客户端的 HTML。

## 3.数据模型
ChitChat 它的数据将被存储到关系式数据库 PostgreSQL 里面，并通过 SQL 与之交互。  
ChitChat的数据模型非常简单，只包含4种数据结构，它们分别是：  
* User——表示论坛的用户信息；
* Session——表示论坛用户当前的登录会话；
* Thread——表示论坛里面的帖子，每一个帖子都记录了多个论坛用户之间的对话；
* Post——表示用户在帖子里面添加的回复。

以上这4种数据结构都会被映射到关系数据库里。

Chitchat 论坛允许用户在登录之后发布新帖子或者回复已有的帖子，未登录的用户可以阅读帖子，但是不能发布新帖子或者回复帖子。 为了对应用进行简化，ChitChat 论坛没有设置版主这一职位，因此用户在发布新帖子或者添加新回复的时候不需要经过审核。

## 4.请求的接收与处理
Web 应用的工作流程如下：
*`客户端`将`请求`发送至`服务器`的一个`URL`上。
* 服务器的`多路复用器`将接收到的`请求` `重定向`到正确的`处理器`，然后由该`处理器`对请求进行`处理`。
* 处理器处理请求并执行必要的动作。
* 处理器调用模板引擎，生成相应的 HTML 并将其返回给客户端。

让我们先从最基本的`根 URL(/)`来考虑 Web 应用是如何处理请求的：当我们在浏览器上输入地址`http://localhost`的时候，浏览器访问的就是应用的`根 URL`。

### 4.1多路复用器
*`net/http`标准库提供了一个默认的`多路复用器`，这个`多路复用器`可以通过调用NewServeMux 函数来创建：  
```bash
mux := http.NewServeMux()
```
* 为了将发送至根 URL 的请求重定向到处理器，可以使用`HandleFunc`函数:
```bash
mux.HandleFunc("/",index)
```
HandleFunc 函数接受一个 URL 和一个处理器的名字作为参数，并将针对给定 URL 的请求转发至指定的处理器函数进行处理，因此对上述调用来说，当有针对根 URL 的请求到达时，该请求就会重定向到名为 index 的处理器函数。  
此外，因为所有处理器都接受一个 ResponseWrite 和一个指向 Request 结构的指针作为参数，并且所有请求参数都可以通过访问 Request 结构得到，所以程序并不需要向处理器显式地传入任何请求参数。

注意：尽管处理器和处理器函数提供的最终结果是一样的，但它们实际上并不相同。
### 4.2服务静态文件
*`多路复用器`还需要为静态文件提供服务。
* FileServer 函数创建了一个能够为指定目录中的静态文件服务的处理器，并将这个处理器传递给了`多路复用器`的 Handle 函数。
* StripPrefix 函数可以移除请求 URL 中的指定前缀：
```bash
files := http.FileServer(http.Dir("/public"))
mux.Handle("/static/",http.StripPrefix("/static/",files))
```
当服务器接收到一个以/static/开头的 URL 请求时，以上两行代码会移除 URL 中的/static/字符串，然后在 public 目录中查找被请求的文件。比如说，当服务器接收到一个针对文件`http://localhost/static/css/bootstrap.min.css`的请求时，它将会在 public 目录中查找以下文件：  
`<application root>/css/bootstrap.min.css` 
当服务器成功地找到这个文件之后，会把它返回给客户端。
### 4.3创建处理器函数
ChitChat应用会通过 HandleFunc 函数把请求重定向到处理器函数。
* 处理器函数实际上就是一个接受`ResponseWriter`和`Request`指针作为参数的 Go 函数。
### 4.4使用 cookie 进行访问控制
跟其他很多 Web 应用一样， ChitChat 既拥有任何人都可以访问的公开页面，也拥有用户在登录账号之后才能看见的私人页面。  
当一个用户成功登录以后，服务器必须在后续的请求中标示出这是一个已登录的用户。 为了做到这一点，`服务器`会在响应的首部中写入一个`cookie`，而`客户端`在接收这个`cookie`之后则会把它存储到浏览器里面。

## 5.使用模板生成HTML响应
*`HTML`文件包含了特定的嵌入命令，这些命令被称为`动作`（action），动作在`HTML`文件里面会被`{{`和`}}`包围。  

`ParseFiles`函数对`HTML`模板文件进行语法分析，并创建出相应的模板。  
`Must`函数捕捉`ParseFiles`函数语法分析过程中可能会产生的错误。  
如果对`layout.html`和`index.html`两个HTML文件进行语法分析，并创建`templates`模板，代码如下：
```bash
tmpl_files := []string{"layout.html","index.html"}
templates := template.Must(template.ParseFiles(tmpl_files...))
```
用`Must`函数去包围`ParseFiles`函数的执行结果，这样当 ParseFiles 返回错误的时候，Must 函数就会向用户返回相应的错误报告。

* ChitChat 论坛的每个模板文件都定义了一个模板，这种做法并不是强制的，用户也可以在一个`模板文件`里面定义多个`模板`，但模板文件和模板一一对应的做法可以给开发带来方便。  

**定义模板的方式**：在模板文件的`源代码`中使用`define`动作。这个动作通过文件开头的`{{ define "模板名" }}`和文件末尾的`{{ end }}`把被包围的`文本块（代码）`定义成了`模板`的一部分。  
**注意**：模板和模板文件分别拥有不同的名字也是可行的。比如，模板文件是`index.html`，可以在模板文件中定义`content`模板，也可以定义`index`模板。

* 模板文件里面还可以包含若干个用于`引用其他模板文件`的`template`动作。跟在`被引用模板名字`之后的`点 (.)`代表了传递给被引用模板的`数据`。

如果在`layout`模板中引用`navbar`模板，并传递相关数据。那么，`layout.html`文件中的`引用模板并传递数据`的语句如下：
```bash
{{ template "navbar" . }}
```
上面的语句除了会在语句出现的位置引入`navbar`模板之外，还会将传递给`layout`模板的数据传递给`navbar`模板。

## 6.安装PostgreSQL

## 7.连接数据库

## 8.启动服务器

## 9.Web应用运作流程图


## 技能清单
1. 多路复用器（multiplexer）：会对请求进行检查，并将请求重定向至正确的处理器进行处理。
2. XXX
3. XXX
4. XXX

## 项目目录结构
### 后端结构树
```bash
.
├── data
│   ├── data.go
│   ├── data_test.go
│   ├── thread.go
│   ├── thread_test.go
│   ├── user.go
│   ├── user_test.go
│   └── setup.sql
├── public
│   ├── css
│   ├── fonts
│   └── js
├── routes
│   ├── route_auth.go
│   ├── route_auth_test.go
│   ├── route_main.go
│   ├── route_main_test.go
│   ├── route_thread.go
│   └── route_thread_test.go
├── templates
│   ├── error.html
│   ├── index.html
│   ├── layout.html
│   ├── login.html
│   ├── login.layout.html
│   ├── new.thread.html
│   ├── private.navbar.html
│   ├── private.thread.html
│   ├── public.navbar.html
│   ├── public.thread.html
│   └── signup.html
├── utils
│   └── utils.go
├── config.json
├── main.go
├── main_test.go
└── README.md
```

## 项目预览图


## 项目引用
**学习书籍**：《Go Web 编程》原名《Go Web Programming》，原书作者——**郑兆雄（Sau SheongChang）**。  
《Go Web 编程》一书围绕一个网络论坛 作为例子，教授读者如何使用请求处理器、多路复用器、模板引擎、存储系统等核心组件去构建一个 Go Web 应用，然后在该应用的基础上，构建出相应的 Web 服务。

**源码地址**：https://github.com/sausheong/gwp