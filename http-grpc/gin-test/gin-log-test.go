package gin_test

import (
	"io"
	"os"

	"github.com/gin-gonic/gin"
)

func Log_test1() {
	// Disable Console Color, you don't need console color when writing the logs to file.
	gin.DisableConsoleColor()

	f, _ := os.Create("test.log")
	gin.DefaultWriter = io.MultiWriter(f)
}
