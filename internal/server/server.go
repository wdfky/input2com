package server

import (
	"embed"
	"fmt"
	"input2com/internal/config"
	"input2com/internal/input"
	"input2com/internal/logger"
	"input2com/internal/macros"
	"io/fs"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

//go:embed server/build
var StaticFS embed.FS

var (
	MouseConfigDict    = make(map[byte]string)
	KeyboardConfigDict = make(map[byte]string)
	MousedictMutex     sync.RWMutex
	KeyboarddictMutex  sync.RWMutex
)

func Serve() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	api := router.Group("/api")
	{
		api.GET("/get/macros", getMacros)
		api.GET("/get/mouse", getMouseConfig)
		api.GET("/get/keyboard", getKeyboardConfig)
		api.GET("/set/mouse", setMouseConfig)
		api.GET("/set/keyboard", setKeyboardConfig)
	}
	// 2️⃣ 再注册静态文件路由（兜底）
	subFS, err := fs.Sub(StaticFS, "server/build")
	if err != nil {
		logger.Logger.Fatalf("Failed to create sub filesystem: %v", err)
	}
	router.NoRoute(gin.WrapH(http.FileServer(http.FS(subFS))))
	router.Run(fmt.Sprintf(":%d", config.Cfg.Server.Port))
}

func getMacros(c *gin.Context) {
	KeyboarddictMutex.RLock()
	defer KeyboarddictMutex.RUnlock()
	c.JSON(http.StatusOK, macros.Macros)
}

func getMouseConfig(c *gin.Context) {
	MousedictMutex.RLock()
	defer MousedictMutex.RUnlock()
	c.JSON(http.StatusOK, MouseConfigDict)
}

func getKeyboardConfig(c *gin.Context) {
	KeyboarddictMutex.RLock()
	defer KeyboarddictMutex.RUnlock()
	c.JSON(http.StatusOK, KeyboardConfigDict)
}

func setMouseConfig(c *gin.Context) {
	MousedictMutex.Lock()
	defer MousedictMutex.Unlock()
	key := c.Query("key")
	value := c.Query("value")

	if key == "CLEAR_ALL" {
		MouseConfigDict = make(map[byte]string)
		logger.Logger.Info("clear mouse config")
		c.String(http.StatusOK, "ok")
		return
	}

	if _, ok := input.MouseValidKeys[key]; !ok {
		c.String(http.StatusBadRequest, "Invalid key")
		return
	}

	if value == "CLEAR_FUNCTION" {
		bkey, _ := strconv.ParseUint(key, 10, 8)
		logger.Logger.Infof("clear mouse config: %d", bkey)
		delete(MouseConfigDict, byte(bkey))
		c.String(http.StatusOK, "ok")
		return
	}

	if _, ok := macros.Macros[value]; !ok {
		c.String(http.StatusBadRequest, "Invalid macro Name")
		return
	}

	bkey, _ := strconv.ParseUint(key, 10, 8)
	logger.Logger.Infof("Set mouse config: %d -> %s", bkey, value)
	MouseConfigDict[byte(bkey)] = value
	c.String(http.StatusOK, "ok")
}

func setKeyboardConfig(c *gin.Context) {
	KeyboarddictMutex.Lock()
	defer KeyboarddictMutex.Unlock()
	key := c.Query("key")
	value := c.Query("value")

	if key == "CLEAR_ALL" {
		KeyboardConfigDict = make(map[byte]string)
		logger.Logger.Info("clear keyboard config")
		c.String(http.StatusOK, "ok")
		return
	}

	if _, ok := input.KeyboardValidKeys[key]; !ok {
		c.String(http.StatusBadRequest, "Invalid key")
		return
	}

	if value == "CLEAR_FUNCTION" {
		bkey, _ := strconv.ParseUint(key, 10, 8)
		logger.Logger.Infof("clear keyboard config: %d", bkey)
		delete(KeyboardConfigDict, byte(bkey))
		c.String(http.StatusOK, "ok")
		return
	}

	if _, ok := macros.Macros[value]; !ok {
		c.String(http.StatusBadRequest, "Invalid macro Name")
		return
	}

	bkey, _ := strconv.ParseUint(key, 10, 8)
	logger.Logger.Infof("Set keyboard config: %d -> %s", bkey, value)
	KeyboardConfigDict[byte(bkey)] = value
	c.String(http.StatusOK, "ok")
}
