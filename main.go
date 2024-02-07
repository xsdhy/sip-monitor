package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"sbc/src/model"
	"sbc/src/pkg/env"
	"sbc/src/services"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed web/build
var dist embed.FS

func main() {
	if env.Conf.DSNURL == "" {
		flag.StringVar(&env.Conf.DSNURL, "dsn", "", "dsn")
		flag.Parse()
	}

	//初始化数据库
	model.MongoDBInit()
	//初始化IP库
	services.IPDBInit()
	//启动HepServer
	go services.HepServerListener()
	//启动定时任务
	go services.Cron()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	//后端接口
	r.GET("/api/record/all", services.SearchAll)
	r.GET("/api/record/call", services.RecordCallList)
	r.GET("/api/record/register", services.RecordRegisterList)
	r.GET("/api/record/details", services.SearchCallID)
	r.GET("/api/system/db/clean_sip_record", services.CleanSipRecord)
	r.GET("/api/system/db/stats", services.DbStats)
	//前端资源
	r.Use(ServerStatic("web/build", dist))
	serverHost := fmt.Sprintf("0.0.0.0:%d", env.Conf.HTTPListenPort)
	slog.Info("HttpServerInit", slog.String("host", serverHost))
	err := r.Run(serverHost)
	if err != nil {
		slog.Error("HttpServerInit Error", err)
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
