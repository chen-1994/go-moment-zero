package gin_test

import (
	"log"
	"net/http"

	"github.com/danielkov/gin-helmet/ginhelmet"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

// 中间件
// Gin 提供了两种创建路由引擎的方式，区别在于默认附加了哪些中间件。
//
// gin.Default() — 带有 Logger 和 Recovery
// gin.Default() 创建一个已附加两个中间件的路由器：
//
// Logger — 将请求日志写入标准输出（方法、路径、状态码、延迟）。
// Recovery — 从处理函数中的任何 panic 恢复并返回 500 响应，防止服务器崩溃。
// 这是快速入门最常用的选择。
//
// gin.New() — 一个空白引擎
// gin.New() 创建一个完全空白的路由器，不附加任何中间件。当你想完全控制运行哪些中间件时很有用，例如：
//
// 你想使用结构化日志记录器（如 slog 或 zerolog）代替默认的文本日志记录器。
// 你想自定义 panic 恢复行为。
// 你正在构建一个需要最小化或专用中间件栈的微服务。

// 使用中间件
// Gin 中的中间件是在路由处理函数之前（以及可选地之后）运行的函数。它们用于日志记录、认证、错误恢复和请求修改等横切关注点。
//
// Gin 支持三个级别的中间件附加：
//
// 全局中间件 — 应用于路由器中的每个路由。使用 router.Use() 注册。适用于日志记录和 panic 恢复等普遍适用的关注点。
// 分组中间件 — 应用于路由组中的所有路由。使用 group.Use() 注册。适用于将认证或授权应用到路由子集（例如所有 /admin/* 路由）。
// 路由级中间件 — 仅应用于单个路由。作为额外参数传递给 router.GET()、router.POST() 等。适用于路由特定的逻辑，如自定义限流或输入验证。
// 执行顺序： 中间件函数按注册顺序执行。当中间件调用 c.Next() 时，它将控制权传递给下一个中间件（或最终处理函数），然后在 c.Next() 返回后继续执行。
// 这创建了一个类似栈的（LIFO）模式——第一个注册的中间件最先开始但最后结束。如果中间件不调用 c.Next()，后续的中间件和处理函数将被跳过（这对于使用 c.Abort() 短路请求很有用）。
func Middleware_Test1() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(ErrorHandler2())
	r.Use(SafeHandler()) //安全头
	//安全 可选方案，使用 gin helmet go get github.com/danielkov/gin-helmet/ginhelmet
	r.Use(ginhelmet.Default())
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	log.Fatal(r.Run(":8080"))
}

// 自定义中间件
func ErrorHandler2() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next() // Process the request first

		// Check if any errors were added to the context
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"message": err.Error(),
			})
		}
	}
}

// 安全头
// 使用安全头来保护你的 Web 应用免受常见安全漏洞非常重要。此示例展示了如何向 Gin 应用添加安全头，以及如何避免 Host Header 注入相关攻击（SSRF、开放重定向）。
//
// 每个头的保护作用
// X-Content-Type-Options: nosniff	防止 MIME 类型嗅探攻击。没有此头，浏览器可能将文件解释为与声明不同的内容类型，允许攻击者执行伪装成无害文件类型的恶意脚本（例如上传一个实际上是 JavaScript 的 .jpg）。
// X-Frame-Options: DENY	通过禁止页面在 <iframe> 中加载来防止点击劫持。攻击者使用点击劫持在合法页面上覆盖不可见的框架，诱骗用户点击隐藏的按钮（例如”删除我的账户”）。
// Content-Security-Policy	控制浏览器允许加载哪些资源（脚本、样式、图片、字体等）以及从哪些来源。这是防御跨站脚本（XSS）最有效的方式之一，因为它可以阻止内联脚本并限制脚本来源。
// X-XSS-Protection: 1; mode=block	激活浏览器内置的 XSS 过滤器。此头在现代浏览器中已基本弃用（Chrome 在 2019 年移除了其 XSS Auditor），但它仍为使用旧浏览器的用户提供纵深防御。
// Strict-Transport-Security	强制浏览器在指定的 max-age 期间对所有未来请求使用 HTTPS。这可以防止协议降级攻击和通过不安全 HTTP 连接的 cookie 劫持。includeSubDomains 指令将此保护扩展到所有子域。
// Referrer-Policy: strict-origin	控制传出请求中发送多少引用者信息。没有此头，完整的 URL（包括可能包含令牌或敏感数据的查询参数）可能会泄露给第三方站点。strict-origin 仅发送来源（域名）且仅通过 HTTPS。
// Permissions-Policy	限制页面可以使用哪些浏览器功能（地理位置、摄像头、麦克风等）。如果攻击者成功注入脚本，这可以限制损害，因为这些脚本无法访问敏感的设备 API。
func SafeHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Host != "localhost:8080" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid host header"})
			return
		}
		c.Header("X-Frame-Options", "DENY")
		c.Header("Content-Security-Policy", "default-src 'self'; connect-src *; font-src *; script-src-elem * 'unsafe-inline'; img-src * data:; style-src * 'unsafe-inline';")
		c.Header("X-XSS-Protection", "1; mode=block")
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		c.Header("Referrer-Policy", "strict-origin")
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("Permissions-Policy", "geolocation=(),midi=(),sync-xhr=(),microphone=(),camera=(),magnetometer=(),gyroscope=(),fullscreen=(self),payment=()")
		c.Next()
	}
}

//会话管理
//会话允许你跨多个 HTTP 请求存储用户特定的数据。由于 HTTP 是无状态的，会话使用 cookie 或其他机制来识别回访用户并检索其存储的数据。
//
//使用 gin-contrib/sessions
//gin-contrib/sessions 中间件提供了支持多种后端存储的会话管理：

// 基于 Cookie 的会话
// 最简单的方式是将会话数据存储在加密的 cookie 中：
func Middleware_Test2() {
	r := gin.Default()

	//可用后端
	//后端	包	用例
	//Cookie	sessions/cookie	简单应用，小型会话数据
	//Redis	sessions/redis	生产环境，多实例部署
	//Memcached	sessions/memcached	高性能缓存层
	//MongoDB	sessions/mongo	MongoDB 为主要数据存储时
	//PostgreSQL	sessions/postgres	PostgreSQL 为主要数据存储时

	// Create cookie-based session store with a secret key
	//store := cookie.NewStore([]byte("your-secret-key"))
	// redis 版
	store, err := redis.NewStore(10, "tcp", "103.114.160.95:6379", "", "redis_MjrC8n")
	if err != nil {
		panic(err)
	}
	r.Use(sessions.Sessions("mysession", store))

	r.Use(sessions.Sessions("mysession", store))

	r.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)

		session.Options(sessions.Options{
			Path:     "/",                  //Cookie 路径范围（默认：/）
			MaxAge:   3600,                 // Session expires in 1 hour (seconds) 生命周期（秒）。使用 -1 删除，0 为浏览器会话
			HttpOnly: true,                 // Prevent JavaScript access  防止 JavaScript 访问 cookie
			Secure:   true,                 // Only send over HTTPS  仅通过 HTTPS 发送 cookie
			SameSite: http.SameSiteLaxMode, //控制跨站 cookie 行为（Lax、Strict、None）
		})

		session.Set("username", "admin") //在 Session 中存入键值对（数据此时在内存中）
		session.Save()                   //关键一步！ 只有调用了 Save()，更改后的数据才会真正写入响应头（Set-Cookie）发回浏览器
		c.JSON(http.StatusOK, gin.H{"username": "admin"})
	})
	r.GET("/profile", func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		if username == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "not logged in"})
		}
	})
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear() //清空当前 Session 里的所有数据。
		session.Save()  //同步更改，告知浏览器清除或重置 Cookie。
		c.JSON(http.StatusOK, gin.H{"message": "logged out"})
	})
	r.Run(":8080")

	//颁发通行证：
	//当服务器执行 session.Save() 时，会生成一个加密字符串，通过 HTTP 响应头的 Set-Cookie 发给浏览器。
	//自动携带：
	//浏览器收到后，会将这个 Cookie 存起来。下次访问该服务器时，浏览器会自动在 Cookie 请求头里带上这个字符串。
	//识别身份：
	//服务器收到请求，Session 中间件读取 Cookie，通过密钥解密，还原出里面的数据（如 username=admin），从而知道是谁在访问。
}

// 基于 Redis 的会话
func Middleware_Test3() {
	r := gin.Default()
	store, err := redis.NewStore(10, "tcp", "103.114.160.95:6379", "", "redis_MjrC8n")
	if err != nil {
		panic(err)
	}
	r.Use(sessions.Sessions("mysession", store))

	r.GET("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("username", "admin")
		session.Save()
	})
}

//依赖注入模式
//闭包模式  参数传递
//基于结构体的处理函数 使用结构体存储注入信息，当方法使用
//使用中间间方式 use进行使用
