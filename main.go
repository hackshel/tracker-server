package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hackshel/tracker-server/controllers"
	"github.com/hackshel/tracker-server/middleware"
	"github.com/hackshel/tracker-server/pkg/setting"
	"github.com/hackshel/tracker-server/routes"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupDatabase() *gorm.DB {

	db, db_err := gorm.Open(mysql.Open(setting.DSN), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   setting.TABLE_PREFIX, // 表前缀
			SingularTable: true,                 // 表名是否使用单数，例如：User -> user（默认 false，即复数 users）
		},
	})
	if db_err != nil {
		log.Fatalf("db init error ... %v", db_err)
	}
	return db
}

func main() {

	router := gin.Default()

	db := setupDatabase()
	router.Use(middleware.DBMiddleware(db))
	router.POST("/api/v1/login", controllers.Login)
	router.GET("/api/v1/tracker/announce", controllers.Announce)
	router.GET("/api/v1/tracker/scrape", controllers.Scrape)

	api := router.Group("/api/v1")

	api.Use(middleware.JWTAuthMiddleware())

	routes.RegisterAPIRoutes(api)

	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        router,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	srv.ListenAndServe()

	//router.Run(":8002")
	// 启动服务器（异步）
	go func() {
		fmt.Println("Server is running at http://localhost:10082")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	// 给 5 秒时间处理剩余请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exiting")

}
