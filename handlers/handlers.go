package handlers

import (
	"aoroalabs/models"
	"aoroalabs/storage"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// 이슈 생성
func CreateIssue(c *gin.Context) {
	var req models.CreateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "필수 파라미터가 누락되었거나 유효하지 않습니다",
			Code:  http.StatusBadRequest,
		})
		return
	}

	store := storage.GetStorage()

	// 담당자가 지정된 경우 사용자 존재 여부 확인
	if req.UserID != nil {
		user := store.GetUser(*req.UserID)
		if user == nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "존재하지 않는 사용자입니다",
				Code:  http.StatusBadRequest,
			})
			return
		}
	}

	issue := store.CreateIssue(req.Title, req.Description, req.UserID)
	c.JSON(http.StatusCreated, issue)
}

// 이슈 목록 조회
func GetIssues(c *gin.Context) {
	store := storage.GetStorage()
	status := c.Query("status")

	var issues []models.Issue
	if status != "" {
		// 상태 유효성 검사
		if !models.IsValidStatus(status) {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "유효하지 않은 상태값입니다",
				Code:  http.StatusBadRequest,
			})
			return
		}
		issues = store.GetIssuesByStatus(status)
	} else {
		issues = store.GetAllIssues()
	}

	c.JSON(http.StatusOK, models.IssuesResponse{Issues: issues})
}

// 이슈 상세 조회
func GetIssue(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "유효하지 않은 이슈 ID입니다",
			Code:  http.StatusBadRequest,
		})
		return
	}

	store := storage.GetStorage()
	issue := store.GetIssue(uint(id))
	if issue == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "이슈를 찾을 수 없습니다",
			Code:  http.StatusNotFound,
		})
		return
	}

	c.JSON(http.StatusOK, issue)
}

// 이슈 수정
func UpdateIssue(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "유효하지 않은 이슈 ID입니다",
			Code:  http.StatusBadRequest,
		})
		return
	}

	store := storage.GetStorage()
	issue := store.GetIssue(uint(id))
	if issue == nil {
		c.JSON(http.StatusNotFound, models.ErrorResponse{
			Error: "이슈를 찾을 수 없습니다",
			Code:  http.StatusNotFound,
		})
		return
	}

	// COMPLETED 또는 CANCELLED 상태에서는 수정 불가
	if issue.Status == models.StatusCompleted || issue.Status == models.StatusCancelled {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "완료되었거나 취소된 이슈는 수정할 수 없습니다",
			Code:  http.StatusBadRequest,
		})
		return
	}

	var req models.UpdateIssueRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "유효하지 않은 요청 데이터입니다",
			Code:  http.StatusBadRequest,
		})
		return
	}

	// 상태 유효성 검사
	if req.Status != nil && !models.IsValidStatus(*req.Status) {
		c.JSON(http.StatusBadRequest, models.ErrorResponse{
			Error: "유효하지 않은 상태값입니다",
			Code:  http.StatusBadRequest,
		})
		return
	}

	// 담당자 존재 여부 확인
	if req.UserID != nil && *req.UserID != 0 {
		user := store.GetUser(*req.UserID)
		if user == nil {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "존재하지 않는 사용자입니다",
				Code:  http.StatusBadRequest,
			})
			return
		}
	}

	// 담당자 없이 PENDING, CANCELLED 이외의 상태로 변경 불가
	if req.Status != nil && *req.Status != models.StatusPending && *req.Status != models.StatusCancelled {
		willHaveUser := issue.User != nil
		if req.UserID != nil {
			willHaveUser = *req.UserID != 0
		}
		if !willHaveUser {
			c.JSON(http.StatusBadRequest, models.ErrorResponse{
				Error: "담당자 없이는 PENDING, CANCELLED 이외의 상태로 변경할 수 없습니다",
				Code:  http.StatusBadRequest,
			})
			return
		}
	}

	updatedIssue := store.UpdateIssue(uint(id), &req)
	c.JSON(http.StatusOK, updatedIssue)
}