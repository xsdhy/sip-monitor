package main

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"net/http"
	"path/filepath"
	"time"

	"sip-monitor/src/config"
	"sip-monitor/src/model"

	"strings"

	"sip-monitor/src/services"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

//go:embed web/build
var dist embed.FS

func main() {
	// 初始化配置
	cfg, err := config.ParseConfig()
	if err != nil {
		logrus.WithError(err).Error("Failed to parse config")
		return
	}

	logger := logrus.New()

	// 初始化数据库
	repository, err := model.InitRepository(&cfg)
	if err != nil {
		logrus.WithError(err).Error("Failed to create repository")
		return
	}
	repository.CreateDefaultAdminUser(context.Background())

	// 初始化保存服务
	saveService := services.NewSaveService(logger, repository)

	//启动HepServer
	hepServer, err := services.NewHepServer(logger, &cfg, saveService)
	if err != nil {
		logrus.WithError(err).Error("Failed to create hep server")
		return
	}
	go hepServer.Start()

	// 初始化认证服务
	authService := services.NewAuthService(logger, repository, cfg.JWTSecret, time.Duration(cfg.JWTExpiryHours)*time.Hour)
	authHandler := services.NewAuthHandler(logger, authService)
	authMiddleware := services.NewAuthMiddleware(logger, authService)

	// 启动HTTP Handle
	handleHttp := services.NewHandleHttp(logger, &cfg, repository)

	// 初始化gin
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 公开的API路由组
	public := r.Group("/api")
	public.POST("/login", authHandler.Login)

	// 需要认证的API路由组
	authorized := r.Group("/api")
	authorized.Use(authMiddleware.JWT())

	// 用户相关API
	authorized.GET("/user/current", authHandler.GetCurrentUser)
	authorized.POST("/user/update", authHandler.UpdateUserInfo)
	authorized.POST("/user/password", authHandler.UpdatePassword)

	// 记录相关API
	authorized.GET("/record/call", handleHttp.RecordCallList)
	authorized.GET("/record/register", handleHttp.RecordRegisterList)
	authorized.GET("/record/details", handleHttp.CallDetails)

	// 用户管理API
	authorized.GET("/users", handleHttp.UserList)
	authorized.GET("/users/:id", handleHttp.GetUser)
	authorized.POST("/users", handleHttp.CreateUser)
	authorized.PUT("/users/:id", handleHttp.UpdateUser)
	authorized.DELETE("/users/:id", handleHttp.DeleteUser)

	//前端资源
	r.Use(ServerStatic("web/build", dist))

	serverHost := fmt.Sprintf("0.0.0.0:%d", cfg.HTTPListenPort)
	logrus.WithField("host", serverHost).Info("HttpServerInit")
	err = r.Run(serverHost)
	if err != nil {
		logrus.WithError(err).Error("HttpServerInit Error")
	}
}

func ServerStatic(prefix string, embedFs embed.FS) gin.HandlerFunc {
	indexPage, err := dist.ReadFile(prefix + "/index.html")
	if err != nil {
		panic("Failed to read the index.html file")
	}
	readFileFS, err := fs.Sub(embedFs, prefix)
	if err != nil {
		panic(err)
	}
	fileServer := http.FileServer(http.FS(readFileFS))

	return func(ctx *gin.Context) {
		// 判断文件是否是静态类型，若不是"重定向"到index.html
		// 这里处理这样判断以外，还有种方式是判断文件是否存在
		if !IsStaticAssetRequest(ctx.Request.URL.Path) {
			ctx.Data(http.StatusOK, "text/html; charset=utf-8", indexPage)
			return
		}
		// 使用预创建的fileServer实例服务请求
		fileServer.ServeHTTP(ctx.Writer, ctx.Request)
	}
}

// IsStaticAssetRequest 判断请求是否为静态资源请求
func IsStaticAssetRequest(path string) bool {
	// 转换为小写，以便进行不区分大小写的比较
	path = strings.ToLower(path)
	// 定义支持的静态资源扩展名列表
	staticExtensions := []string{".css", ".js", ".png", ".jpg", ".jpeg", ".gif", ".json", ".txt"}
	// 获取路径的扩展名
	ext := filepath.Ext(path)
	// 遍历静态资源扩展名，检查当前路径是否以其中任一扩展名结束
	for _, staticExt := range staticExtensions {
		if ext == staticExt {
			return true
		}
	}
	// 如果没有匹配任何静态资源扩展名，返回false
	return false
}
