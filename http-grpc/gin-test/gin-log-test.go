package gin_test

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/natefinch/lumberjack.v2"
)

func Log_test1() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor() //日志没有颜色

	// Force log's color
	//gin.ForceConsoleColor() //日志颜色

	f, _ := os.Create("test.log")
	gin.DefaultWriter = io.MultiWriter(f)

	//同时写入文件和控制台
	//Go 标准库中的 io.MultiWriter 函数接受多个 io.Writer 值，并将写入复制到所有目标。这在开发时很有用，你既想在终端看到日志，又想将其持久化到磁盘：
	//gin.DefaultWriter = io.MultiWriter(f, os.Stdout)

	//生产环境中的日志轮转
	//上面的示例使用 os.Create，每次应用启动时都会截断日志文件。在生产环境中，你通常希望追加到现有日志并根据大小或时间轮转文件。考虑使用日志轮转库如 lumberjack：
	gin.DefaultWriter = &lumberjack.Logger{
		Filename:   "gin.log",
		MaxSize:    100, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	}

	router := gin.Default()

	router.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	router.Run(":8080")
}

// 自定义日志格式
func Log_test2() {
	router := gin.New()

	router.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\n",
			param.TimeStamp.Format(time.RFC822),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))

	//跳过日志记录
	//你可以使用 LoggerConfig 跳过特定路径或基于自定义逻辑的日志记录。
	//SkipPaths 排除特定路由的日志记录——适用于会产生噪音的健康检查或指标端点。
	//Skip 是一个接收 *gin.Context 并返回 true 来跳过日志记录的函数——适用于条件逻辑，如跳过成功响应的日志。
	// skip logging for desired paths by setting SkipPaths in LoggerConfig
	//SkipQueryString选项可以防止查询字符串出现在日志中 设置为true 从日志中剥离查询字符串，可以减少通过日志文件、监控系统或错误报告工具泄露敏感数据的风险
	loggerConfig := gin.LoggerConfig{SkipPaths: []string{"/metrics"}, SkipQueryString: true} //跳过这个API日志 并且设置玻璃查询字符串

	// skip logging based on your logic by setting Skip func in LoggerConfig
	loggerConfig.Skip = func(c *gin.Context) bool {
		// as an example skip non server side errors
		return c.Writer.Status() < http.StatusInternalServerError // 跳过小于500的异常
	}

	router.Use(gin.LoggerWithConfig(loggerConfig))

	router.Use(gin.Recovery())
	// skipped -- path is in SkipPaths
	router.GET("/metrics", func(c *gin.Context) {
		//c.Status(http.StatusNotImplemented)
		q := c.Query("q")
		c.String(http.StatusNotImplemented, "searching for: "+q)
	})

	// skipped -- status < 500
	router.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// not skipped -- status is 501 (>= 500)
	router.GET("/data", func(c *gin.Context) {
		c.Status(http.StatusNotImplemented)
	})
	router.Run(":8080")
}

// 定义路由日志格式
func Log_test3() {
	router := gin.Default()

	gin.DebugPrintRouteFunc = func(httpMethod, absolutePath, handlerName string, nuHandlers int) {
		log.Printf("endpoint %v %v %v %v\n", httpMethod, absolutePath, handlerName, nuHandlers)
	}

	router.POST("/foo", func(c *gin.Context) {
		c.JSON(http.StatusOK, "foo")
	})

	router.GET("/bar", func(c *gin.Context) {
		c.JSON(http.StatusOK, "bar")
	})

	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	router.Run(":8080")
}
