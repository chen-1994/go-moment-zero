package gin_test

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

// http方法
func Test1() {
	r := gin.Default()

	//定义路由
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080")
}

// 路径参数
func getting(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "GET"})
}
func posting(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "POST"})
}
func putting(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PUT"})
}
func deleting(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "DELETE"})
}
func patching(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "PATCH"})
}
func Head(c *gin.Context) {
	c.Status(http.StatusOK)
}
func options(c *gin.Context) {
	c.Status(http.StatusOK)
}

func Test2() {
	r := gin.Default()
	r.GET("/someGet", getting)
	r.POST("/somePost", posting)
	r.PUT("/somePut", putting)
	r.DELETE("/someDelete", deleting)
	r.PATCH("/someDelete", patching)
	r.HEAD("/someDelete", Head)
	r.OPTIONS("/someOptions", options)
	r.Run(":8080")
}

// 查询字符串参数是出现在 URL 中 ? 后面的键值对（例如 /search?q=gin&page=2）。Gin 提供了两种方法来读取它们：
// c.Query("key") 返回查询参数的值，如果键不存在则返回空字符串。
// c.DefaultQuery("key", "default") 返回值，如果键不存在则返回指定的默认值。
func Test3() {
	r := gin.Default()
	r.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")
		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})
	r.Run(":8080")
}

// Gin 提供了基于 httprouter 构建的强大路由系统，用于高性能的 URL 匹配。
// 在底层，httprouter 使用基数树（也称为压缩字典树）来存储和查找路由，这意味着路由匹配非常快速，每次查找零内存分配。
// 这使得 Gin 成为最快的 Go Web 框架之一
func Test4() {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "hello world")
	})
	r.POST("users", func(c *gin.Context) {
		name := c.PostForm("name")
		c.JSON(http.StatusCreated, gin.H{"message": name})
	})
}

// Multipart/Urlencoded
func Test5() {
	r := gin.Default()

	r.POST("/from_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")
		c.JSON(http.StatusCreated, gin.H{"status": "posted", "message": message, "nick": nick})
	})
	r.Run(":8080")
	//curl -X POST http://localhost:8080/form_post \
	//  -d "message=hello&nick=world"
	//curl -X POST http://localhost:8080/form_post \
	//  -F "message=hello" -F "nick=world"
	//curl -X POST http://localhost:8080/form_post \
	//  -d "message=hello"
}

//查询字符串和表单
//处理 POST 请求时，你通常需要同时从 URL 查询字符串和请求体中读取值。Gin 将这两个数据源分开，因此你可以独立访问每一个：
//c.Query("key") / c.DefaultQuery("key", "default") —— 从 URL 查询字符串读取。
//c.PostForm("key") / c.DefaultPostForm("key", "default") —— 从 application/x-www-form-urlencoded 或 multipart/form-data 请求体读取。
//这在 REST API 中很常见，路由通过查询参数（如 id）标识资源，而请求体携带有效负载（如 name 和 message）。

func Test6() {
	r := gin.Default()
	r.POST("/post", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "1")
		name := c.PostForm("name")
		message := c.PostForm("message")
		fmt.Printf("id: %s, name: %s, message: %s", id, name, message)
		c.String(http.StatusCreated, "id: %s; page: %s; name: %s; message: %s", id, page, name, message)
	})
	r.Run(":8080")
	//# Query params in URL, form data in body
	//curl -X POST "http://localhost:8080/post?id=1234&page=1" \
	//  -d "name=manu&message=this_is_great"
	//# Output: id: 1234; page: 1; name: manu; message: this_is_great
	//
	//# Missing page -- falls back to default value "0"
	//curl -X POST "http://localhost:8080/post?id=1234" \
	//  -d "name=manu&message=hello"
	//# Output: id: 1234; page: 0; name: manu; message: hello
}

// Map 作为查询字符串或表单参数
// 有时你需要接收一组事先不知道键名的键值对——例如动态过滤器或用户定义的元数据。Gin 提供了 c.QueryMap 和 c.PostFormMap 来将方括号表示法的参数（如 ids[a]=1234）解析为 map[string]string。
// c.QueryMap("key") —— 从 URL 查询字符串中解析 key[subkey]=value 形式的键值对。
// c.PostFormMap("key") —— 从请求体中解析 key[subkey]=value 形式的键值对。
func Test7() {
	r := gin.Default()
	r.POST("/post", func(c *gin.Context) {
		ids := c.QueryMap("ids")
		name := c.PostFormMap("name")
		fmt.Printf("ids: %s; name: %s", ids, name)
		c.JSON(http.StatusCreated, gin.H{"ids": ids, "name": name})
	})
	r.Run(":8080")
	//curl -X POST "http://localhost:8080/post?ids[a]=1234&ids[b]=hello" \
	//  -d "names[first]=thinkerou&names[second]=tianou"
}

// 文件上传
// Gin 使处理 multipart 文件上传变得简单直接。框架在 gin.Context 上提供了内置方法来接收上传的文件：
// c.FormFile(name) — 通过表单字段名从请求中获取单个文件。
// c.MultipartForm() — 解析整个 multipart 表单，可以访问所有上传的文件和字段值。
// c.SaveUploadedFile(file, dst) — 一个便捷方法，将接收到的文件保存到磁盘上的目标路径。
func Test8() {
	//内存限制
	//Gin 为 multipart 表单解析设置了默认 32 MiB 的内存限制，通过 router.MaxMultipartMemory 设置。
	//在此限制内的文件会缓存在内存中；超出的部分会写入磁盘上的临时文件。你可以根据应用需求调整此值：
	router := gin.Default()
	// Lower the limit to 8 MiB
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
}

// 单文件
// 使用 c.FormFile 从 multipart/form-data 请求中接收单个上传的文件，然后使用 c.SaveUploadedFile 将其保存到磁盘。
// 你可以通过设置 router.MaxMultipartMemory 来控制 multipart 解析期间使用的最大内存（默认为 32 MiB）。超过此限制的文件将存储在磁盘上的临时文件中而不是内存中。
func Test9() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 //8MiB

	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		log.Println(file.Filename)

		dst := filepath.Join(".", filepath.Base(file.Filename))
		c.SaveUploadedFile(file, dst)
		c.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	})
	router.Run(":8080")
	//curl -X POST http://localhost:8080/upload \
	//  -F "file=@/path/to/your/file.zip" \
	//  -H "Content-Type: multipart/form-data"
}

// 多个文件
// 使用 c.MultipartForm 接收在单个请求中上传的多个文件。文件按表单字段名分组——对所有要一起上传的文件使用相同的字段名。
func Test10() {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20
	router.POST("/upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		files := form.File["files"]
		for _, file := range files {
			log.Println(file.Filename)
			dst := filepath.Join(".", filepath.Base(file.Filename))
			c.SaveUploadedFile(file, dst)
		}
		c.String(http.StatusOK, fmt.Sprintf("'%s' files uploaded!", len(files)))
	})
	router.Run(":8080")
	//curl -X POST http://localhost:8080/upload \
	//  -F "files=@/path/to/test1.zip" \
	//  -F "files=@/path/to/test2.zip" \
	//  -H "Content-Type: multipart/form-data"
}

//限制上传大小
//使用 http.MaxBytesReader 严格限制上传文件的最大大小。当超出限制时，读取器会返回错误，你可以使用 413 Request Entity Too Large 状态码进行响应。
//这对于防止客户端发送超大文件以耗尽服务器内存或磁盘空间的拒绝服务攻击非常重要。
//工作原理
//定义限制 —— 常量 MaxUploadSize（1 MB）设置上传的硬性上限。
//强制限制 —— http.MaxBytesReader 包装 c.Request.Body。如果客户端发送的字节数超过允许值，读取器将停止并返回错误。
//解析并检查 —— c.Request.ParseMultipartForm 触发读取。代码检查 *http.MaxBytesError 以返回带有明确消息的 413 状态码。

const (
	MaxUploadSize = 1 << 20 // 1MB
)

func uploadHandler(c *gin.Context) {
	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, MaxUploadSize)

	if err := c.Request.ParseMultipartForm(MaxUploadSize); err != nil {
		if _, ok := err.(*http.MaxBytesError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("file too large (max: %d bytes)", MaxUploadSize)})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file form required"})
		return
	}
	defer file.Close()

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}
func Test11() {
	router := gin.Default()
	router.POST("/upload", uploadHandler)
	router.Run(":8080")
	//# Upload a small file (under 1 MB) -- succeeds
	//curl -X POST http://localhost:8080/upload \
	//  -F "file=@/path/to/small-file.txt"
	//# Output: {"message":"upload successful"}
	//
	//# Upload a large file (over 1 MB) -- rejected
	//curl -X POST http://localhost:8080/upload \
	//  -F "file=@/path/to/large-file.zip"
	//# Output: {"error":"file too large (max: 1048576 bytes)"}
}

// 路由分组
// 路由组允许你将相关路由组织在一个共享的 URL 前缀下。这适用于：
// API 版本管理 — 将所有 v1 端点分组到 /v1 下，v2 端点分组到 /v2 下。
// 共享中间件 — 将认证、日志记录或限流一次性应用到整组路由，而不是逐个附加到每个路由上。
// 代码组织 — 在源代码中将相关的处理函数保持在一起。
func loginEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"action": "login"})
}
func submitEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"action": "submit"})
}
func readEndpoint(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"action": "read"})
}

// 基本分组
func Test12() {
	router := gin.Default()
	{
		v1 := router.Group("/v1")
		v1.POST("/upload", uploadHandler)
		v1.GET("/read", readEndpoint)
		v1.POST("/submit", submitEndpoint)
	}
	{
		v2 := router.Group("/v2")
		v2.POST("/upload", uploadHandler)
		v2.POST("/read", readEndpoint)
		v2.POST("/submit", submitEndpoint)
	}
	router.Run(":8080")
}

// 将中间件应用到分组
// 你可以将中间件传递给 router.Group() 或在组上调用 Use()。该组中的每个路由都会在其处理函数之前运行中间件。
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		//... check token, session, etc.
		c.Next()
	}
}
func Test13() {
	router := gin.Default()
	public := router.Group("/api")
	{
		public.GET("/me", func(c *gin.Context) {})
	}
	private := router.Group("/api")
	private.Use(AuthRequired())
	{
		private.GET("/me", func(c *gin.Context) {})
		private.POST("/me", func(c *gin.Context) {})
	}
	router.Run(":8080")
}

// 嵌套分组
// 分组可以嵌套以构建更深层的 URL 层次结构，同时保持中间件的合理作用范围。
func Test14() {
	router := gin.Default()
	api := router.Group("/api")
	{
		// /api/v1
		v1 := api.Group("/v1")
		{
			// /api/v1/users
			users := v1.Group("/users")
			users.GET("/", func(c *gin.Context) {})
			users.POST("/:id", func(c *gin.Context) {})

			// /api/v1/posts
			posts := v1.Group("/posts")
			posts.GET("/:id", func(c *gin.Context) {})
			posts.PUT("/:id", func(c *gin.Context) {})
		}
	}
	router.Run(":8080")
	//每一级都会继承父级的前缀，因此最终的路由会变成 /api/v1/users/、/api/v1/users/:id 等
}

// 重定向
func Test15() {
	router := gin.Default()
	// External redirect (GET)
	router.GET("/old", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "https://www.google.com/")
	})
	// Redirect from POST -- use 302 or 307 to preserve behavior
	router.POST("/submit", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/result")
	})
	// Internal router redirect (no HTTP round-trip)
	router.GET("/test", func(c *gin.Context) {
		c.Request.URL.Path = "/final"
		router.HandleContext(c)
	})

	router.GET("/final", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"hello": "world"})
	})
	router.GET("/result", func(c *gin.Context) {
		c.String(http.StatusOK, "Redirected here!")
	})
	router.Run(":8080")
	//curl -I http://localhost:8080/old
	//curl -X POST -I http://localhost:8080/submit
	//curl http://localhost:8080/test
}

// API 设计模式
// 使用 Gin 构建 RESTful API 不仅仅是定义路由。一个设计良好的 API 应该使用一致的响应格式、可预测的分页、清晰的版本管理和结构化的错误处理。本指南介绍了可应用于生产环境 Gin 应用的实用模式。
// 一致的响应格式
// 将每个响应包装在标准的信封结构中可以方便 API 使用者。他们总是知道在哪里找到数据、错误和元数据。
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *ErrorInfo  `json:"error,omitempty"`
	Meta    interface{} `json:"meta,omitempty"`
}
type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
type Meta struct {
	Page       int `json:"page,omitempty"`
	PerPage    int `json:"perPage,omitempty"`
	Total      int `json:"total,omitempty"`
	TotalPages int `json:"totalPages,omitempty"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success: true,
		Data:    data,
	})
}
func Fail(c *gin.Context, status int, code, message string) {
	c.JSON(status, Response{
		Success: false,
		Error: &ErrorInfo{
			Code:    code,
			Message: message,
		},
	})
}
func Test16() {
	r := gin.Default()
	r.GET("/api/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "0" {
			Fail(c, http.StatusNotFound, "USER_NOT_FOUND", "no user with this id")
			return
		}
		OK(c, gin.H{"id": id, "name": "Alice"})
	})
	r.Run(":8080")
}

// 分页
// 限制/偏移分页
// 限制/偏移是最简单的方式。它适用于中小型数据集，因为计算总行数的开销是可以接受的。
func Test17() {
	r := gin.Default()
	r.GET("/api/articles", func(c *gin.Context) {
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		if limit > 100 {
			limit = 100
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []gin.H{},
			"meta":    gin.H{"limit": limit, "offset": offset, "total": 0}})
	})
	r.Run(":8080")
}

// 基于游标的分页
// 基于游标的分页可以避免大偏移量带来的性能问题。将最后一项的 ID（或其他唯一、可排序的字段）作为游标传递。
func Test18() {
	r := gin.Default()
	r.GET("/api/event", func(c *gin.Context) {
		cursor := c.Query("cursor")
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
		if limit > 100 {
			limit = 100
		}

		//events, nextCursor := db.ListEvents(cursor, limit)
		_ = cursor
		c.JSON(http.StatusOK, gin.H{
			"success":     true,
			"data":        []gin.H{},
			"next_cursor": "",
		})
	})
	r.Run(":8080")
}

// 过滤与排序
// 通过查询参数接受过滤和排序。使用一致的参数名称保持接口可预测性。
func Test19() {
	r := gin.Default()
	//// GET /api/products?category=electronics&min_price=10&sort=price&order=asc
	r.GET("/api/products", func(c *gin.Context) {
		category := c.Query("category")
		minPrice := c.Query("min_price")
		maxPrice := c.Query("max_price")
		sortBy := c.DefaultQuery("sort", "created_at")
		order := c.DefaultQuery("order", "desc")

		allowed := map[string]bool{"category": true, "min_price": true, "max_price": true}
		if !allowed[sortBy] {
			sortBy = "created_at"
		}
		if order != "desc" && order != "asc" {
			order = "desc"
		}
		// Build and execute your query using these filters...
		_ = category
		_ = minPrice
		_ = maxPrice

		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"data":    []gin.H{},
			"filters": gin.H{
				"category":  category,
				"min_price": minPrice,
				"max_price": maxPrice,
				"sort":      sortBy,
				"order":     order,
			}})
	})
	r.Run(":8080")
}

// ====API版本管理
// URL路径版本管理
// URL 路径版本管理是最常见的策略。它是显式的，路由简单，使用 curl 测试也很方便
func Test20() {
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/users", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"version": "v1", "users": []string{}})
		})
	}

	v2 := r.Group("/api/v2")
	{
		v2.GET("/users", func(c *gin.Context) {
			// v2 returns a different shape
			c.JSON(http.StatusOK, gin.H{
				"version": "v2",
				"data":    []gin.H{},
				"meta":    gin.H{"total": 0},
			})
		})
	}

	r.Run(":8080")
}

// 基于请求头的版本管理
// 基于请求头的版本管理可以保持 URL 整洁，但需要客户端设置自定义请求头。中间件可以读取请求头并将版本存储在上下文中。
func VersionMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		version := c.GetHeader("version")
		if version == "" {
			version = "v1"
		}
		c.Set("api_version", version)
		c.Next()
	}
}
func Test21() {
	r := gin.Default()
	r.Use(VersionMiddleware())
	r.GET("/api/users", func(c *gin.Context) {
		version := c.GetString("api_version")
		switch version {
		case "v2":
			c.JSON(http.StatusOK, gin.H{"version": "v2", "users": []gin.H{}})
		default:
			c.JSON(http.StatusOK, gin.H{"version": "v1", "users": []string{}})
		}
	})

}

// 错误处理模式
// 自定义错误类型
// 定义应用级别的错误类型，使处理函数能够返回有意义的、结构化的错误。
type AppError struct {
	Status  int    `json:"-"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

var (
	ErrNotFound     = &AppError{Status: http.StatusNotFound, Code: "not_found", Message: "Not Found"}
	ErrUnauthorized = &AppError{Status: http.StatusUnauthorized, Code: "unauthorized", Message: "Unauthorized"}
	ErrBadRequest   = &AppError{Status: http.StatusBadRequest, Code: "bad_request", Message: "Bad Request"}
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		if len(c.Errors) > 0 {
			return
		}
		err := c.Errors.Last().Err
		var appErr *AppError
		if errors.As(err, &appErr) {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"error":   gin.H{"code": appErr.Code, "message": appErr.Message}})
		} else {
			c.JSON(http.StatusOK, gin.H{
				"success": false,
				"error":   gin.H{"code": "INTERNAL", "message": "an unexpected error occurred"},
			})
		}
	}
}

func Test22() {
	r := gin.Default()
	r.Use(VersionMiddleware()) //用于全局应用中间件
	r.GET("/api/item/:id", func(c *gin.Context) {
		id := c.Param("id")
		if id == "0" {
			_ = c.Error(ErrNotFound) //Gin 的错误收集机制。它并不会立即中断请求或向客户端发送响应，而是将错误挂载到当前上下文中，方便后续的中间件（如统一错误处理中间件）获取并处理。
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"id": id}})
	})
	r.Run(":8080")
}

// 基于资源的路由组织
// 随着 API 的增长，按资源组织路由。每个资源都有自己的文件，包含一个在 gin.RouterGroup 上注册路由的函数

func RegisterUserRoutes(rg *gin.RouterGroup) {
	users := rg.Group("/users")
	{
		users.GET("/", listUsers)
		users.POST("/", createUser)
		users.GET("/:id", getUser)
		users.PUT("/:id", updateUser)
		users.DELETE("/:id", deleteUser)
	}

}
func RegisterOrderRoutes(rg *gin.RouterGroup) {
	orders := rg.Group("/orders")
	{
		orders.GET("/", listOrders)
		orders.POST("/", createOrder)
		orders.GET("/:id", getOrder)
	}
}

func listUsers(c *gin.Context)  { c.JSON(http.StatusOK, gin.H{"action": "list_users"}) }
func createUser(c *gin.Context) { c.JSON(http.StatusCreated, gin.H{"action": "create_user"}) }
func getUser(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"action": "get_user"}) }
func updateUser(c *gin.Context) { c.JSON(http.StatusCreated, gin.H{"action": "update_user"}) }
func deleteUser(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"action": "delete_user"}) }

func listOrders(c *gin.Context)  { c.JSON(http.StatusOK, gin.H{"action": "list_orders"}) }
func createOrder(c *gin.Context) { c.JSON(http.StatusCreated, gin.H{"action": "create_order"}) }
func getOrder(c *gin.Context)    { c.JSON(http.StatusOK, gin.H{"action": "get_order"}) }

func Test23() {
	r := gin.Default()

	api := r.Group("/api/v1")
	RegisterUserRoutes(api)
	RegisterOrderRoutes(api)
	r.Run(":8080")
	//在实际项目中，你会将 RegisterUserRoutes 放在 routes/users.go 文件中，
	//将 RegisterOrderRoutes 放在 routes/orders.go 中，然后在 main.go 中调用它们。
	//这样可以让每个资源自包含，并且在添加或删除资源时不需要修改不相关的代码。
}
