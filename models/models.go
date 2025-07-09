package models

import (
	"time"
)

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Issue struct {
	ID          uint      `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	User        *User     `json:"user,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// 유효한 상태값들
const (
	StatusPending    = "PENDING"
	StatusInProgress = "IN_PROGRESS"
	StatusCompleted  = "COMPLETED"
	StatusCancelled  = "CANCELLED"
)

// 상태 유효성 검사
func IsValidStatus(status string) bool {
	switch status {
	case StatusPending, StatusInProgress, StatusCompleted, StatusCancelled:
		return true
	default:
		return false
	}
}

// 에러 응답 구조체
type ErrorResponse struct {
	Error string `json:"error"`
	Code  int    `json:"code"`
}

// 이슈 생성 요청 구조체
type CreateIssueRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	UserID      *uint  `json:"userId"`
}

// 이슈 수정 요청 구조체
type UpdateIssueRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	UserID      *uint   `json:"userId"`
}

// 이슈 목록 응답 구조체
type IssuesResponse struct {
	Issues []Issue `json:"issues"`
}