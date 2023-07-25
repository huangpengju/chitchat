# chitchat
这是一个名为ChitChat的简易网上论坛，ChitChat实现网上论坛的关键特性：  
1. 在这个论坛里面，用户可以注册账号，并在登录之后发表新帖子又或者回复已有的帖子；
2. 未注册用户可以查看帖子，但是无法发表帖子或回复帖子。

## 技能清单


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
├── config.json
├── main.go
├── route_auth.go
├── route_auth_test.go
├── route_main.go
├── route_thread.go
├── utils.go
└── README.md
```

## 项目预览图


学习书籍：《Go Web 编程》原名《Go Web Programming》，原书作者——郑兆雄（Sau SheongChang）。  
《Go Web 编程》一书围绕一个网络论坛 作为例子，教授读者如何使用请求处理器、多路复用器、模板引擎、存储系统等核心组件去构建一个 Go Web 应用，然后在该应用的基础上，构建出相应的 Web 服务。

源码地址：https://github.com/sausheong/gwp