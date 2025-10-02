package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode) //gin.TestMode is used to set the mode to test
	//in test mode, gin will not listen and serve on a random port
	os.Exit(m.Run())
}
