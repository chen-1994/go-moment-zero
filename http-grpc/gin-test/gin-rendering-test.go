package gin_test

import (
	"net/http"

	"github.com/gin-contrib/multitemplate"
	"github.com/gin-gonic/gin"
)

// 1. 常用数据渲染 (JSON, XML, YAML, ProtoBuf)
// 这是前后端分离开发中最常用的方式，Gin 会自动设置 Content-Type 响应头。
// JSON 渲染：最常用的格式。
// PureJSON：如果 HTML 字符（如 <）不想被转义成 Unicode，使用 PureJSON。
// SecureJSON：为了防止 JSON 劫持，会在响应体前加上 while(1);。
func Rendering_test() {
	r := gin.Default()
	r.GET("/read", func(c *gin.Context) {
		data := gin.H{"id": "1", "name": "john"}
		c.JSON(http.StatusOK, data)
		//c.XML(http.StatusOK, data)
		//c.YAML(http.StatusOK, data)
		//c.TOML(http.StatusOK, data)
		//c.ProtoBuf(http.StatusOK, data)
	})

	//PureJSON
	//通常，Go 的 json.Marshal 会出于安全考虑将特殊 HTML 字符替换为 Unicode 转义序列——例如 < 变成 \u003c。
	//当将 JSON 嵌入 HTML 时这很好，但如果你正在构建纯 API，客户端可能期望得到原始字符。
	//c.PureJSON 使用 json.Encoder 并设置 SetEscapeHTML(false)，因此 <、> 和 & 等 HTML 字符会按原样呈现而不会被转义。
	r.GET("/purejson", func(c *gin.Context) {
		c.PureJSON(http.StatusOK, gin.H{
			"html": "<b>Hello, world!</b>",
		})
	})
	r.Run(":8080")
}

func Rendering_test2() {
	r := gin.Default()
	r.LoadHTMLGlob("./gin-test/templates/*")
	//r.LoadHTMLFiles("templates/index.tmpl", "templates/header.tmpl")
	r.GET("/index", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{"title": "hello world"})
	})
	r.Run(":8080")
}

// 文件
func Rendering_test3() {
	router := gin.Default()

	// Serve a file inline (displayed in browser)
	router.GET("/local/file", func(c *gin.Context) {
		c.File("./gin-test/templates/aaa.txt")
	})

	// Serve a file from an http.FileSystem
	var fs http.FileSystem = http.Dir("/var/www/assets")
	router.GET("/fs/file", func(c *gin.Context) {
		c.FileFromFS("./gin-test/templates/aaa.txt", fs)
	})

	// Serve a file as a downloadable attachment with a custom filename
	router.GET("/download", func(c *gin.Context) {
		c.FileAttachment("./gin-test/templates/abc.xls", "quarterly-report.xlsx")
	})

	router.Run(":8080")
}

// 从 Reader 提供数据
// DataFromReader 允许你将任何 io.Reader 的数据直接流式传输到 HTTP 响应，而无需先将整个内容缓冲到内存中。这对于构建代理端点或高效地从远程源提供大文件至关重要。
// 常见用例：
//
// 代理远程资源 — 从外部服务（如云存储 API 或 CDN）获取文件并转发给客户端。数据通过你的服务器流过，而不会完全加载到内存中。
// 提供生成的内容 — 在生产动态生成的数据（如 CSV 导出或报告文件）时进行流式传输。
// 大文件下载 — 提供太大而无法保存在内存中的文件，从磁盘或远程源分块读取。
// 方法签名为 c.DataFromReader(code, contentLength, contentType, reader, extraHeaders)。
// 你需要提供 HTTP 状态码、内容长度（让客户端知道总大小）、MIME 类型、要流式传输的 io.Reader，以及可选的额外响应头映射（如用于文件下载的 Content-Disposition）。
func Rendering_test4() {
	router := gin.Default()
	router.GET("/someDataFromReader", func(c *gin.Context) {
		resp, err := http.Get("https://pics1.baidu.com/feed/4034970a304e251f9df1360c3e7ae6077d3e53c8.jpeg")
		if err != nil || resp.StatusCode != http.StatusOK {
			c.Status(http.StatusServiceUnavailable)
			return
		}
		reader := resp.Body
		defer reader.Close()
		lenght := resp.ContentLength
		t := resp.Header.Get("Content-Type")

		handlers := map[string]string{
			"Content-Disposition": `attachment; filename="gopher.png"`,
		}
		c.DataFromReader(http.StatusOK, lenght, t, reader, handlers)
	})
	router.Run(":8080")
}

// 多模板
// Gin 默认只允许使用一个 html.Template
func Rendering_test5() {
	router := gin.Default()
	router.HTMLRender = createMyRender()

	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index", gin.H{"title": "Home"})
	})
	router.GET("/article", func(c *gin.Context) {
		c.HTML(200, "article", gin.H{"title": "Article"})
	})

	router.Run(":8080")
}

func createMyRender() multitemplate.Renderer {
	r := multitemplate.NewRenderer()
	r.AddFromFiles("index", "templates/base.html", "templates/index.html")
	r.AddFromFiles("article", "templates/base.html", "templates/article.html")
	return r
}
