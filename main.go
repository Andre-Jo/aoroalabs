package main

import (
	"aoroalabs/handlers"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Gin 라우터 초기화
	r := gin.Default()

	// 라우트 설정
	r.POST("/issue", handlers.CreateIssue)
	r.GET("/issues", handlers.GetIssues)
	r.GET("/issue/:id", handlers.GetIssue)
	r.PATCH("/issue/:id", handlers.UpdateIssue)

	// 서버 시작
	log.Println("서버가 포트 8080에서 시작됩니다...")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("서버 시작 실패:", err)
	}
}