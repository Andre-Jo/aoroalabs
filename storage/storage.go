package storage

import (
	"aoroalabs/models"
	"sync"
	"time"
)

type Storage struct {
	users       map[uint]*models.User
	issues      map[uint]*models.Issue
	nextIssueID uint
	mutex       sync.RWMutex
}

var storage *Storage
var once sync.Once

// 싱글톤 패턴으로 Storage 인스턴스 반환
func GetStorage() *Storage {
	once.Do(func() {
		storage = &Storage{
			users:       make(map[uint]*models.User),
			issues:      make(map[uint]*models.Issue),
			nextIssueID: 1,
		}
		storage.initUsers()
	})
	return storage
}

// 기본 사용자 초기화
func (s *Storage) initUsers() {
	s.users[1] = &models.User{ID: 1, Name: "김개발"}
	s.users[2] = &models.User{ID: 2, Name: "이디자인"}
	s.users[3] = &models.User{ID: 3, Name: "박기획"}
}

// 사용자 조회
func (s *Storage) GetUser(id uint) *models.User {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.users[id]
}

// 모든 이슈 조회
func (s *Storage) GetAllIssues() []models.Issue {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	issues := make([]models.Issue, 0, len(s.issues))
	for _, issue := range s.issues {
		issues = append(issues, *issue)
	}
	return issues
}

// 상태별 이슈 조회
func (s *Storage) GetIssuesByStatus(status string) []models.Issue {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	
	var issues []models.Issue
	for _, issue := range s.issues {
		if issue.Status == status {
			issues = append(issues, *issue)
		}
	}
	return issues
}

// 이슈 ID로 조회
func (s *Storage) GetIssue(id uint) *models.Issue {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.issues[id]
}

// 이슈 생성
func (s *Storage) CreateIssue(title, description string, userID *uint) *models.Issue {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	now := time.Now()
	issue := &models.Issue{
		ID:          s.nextIssueID,
		Title:       title,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	
	// 담당자가 있는 경우 IN_PROGRESS, 없는 경우 PENDING
	if userID != nil {
		issue.User = s.users[*userID]
		issue.Status = models.StatusInProgress
	} else {
		issue.Status = models.StatusPending
	}
	
	s.issues[s.nextIssueID] = issue
	s.nextIssueID++
	
	return issue
}

// 이슈 수정
func (s *Storage) UpdateIssue(id uint, req *models.UpdateIssueRequest) *models.Issue {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	issue := s.issues[id]
	if issue == nil {
		return nil
	}
	
	// 제목 수정
	if req.Title != nil {
		issue.Title = *req.Title
	}
	
	// 설명 수정
	if req.Description != nil {
		issue.Description = *req.Description
	}
	
	// 담당자 수정
	if req.UserID != nil {
		if *req.UserID == 0 {
			// 담당자 제거
			issue.User = nil
			issue.Status = models.StatusPending
		} else {
			// 담당자 할당
			issue.User = s.users[*req.UserID]
			// 상태가 PENDING이고 담당자를 할당하는 경우 IN_PROGRESS로 변경
			if issue.Status == models.StatusPending && req.Status == nil {
				issue.Status = models.StatusInProgress
			}
		}
	}
	
	// 상태 수정
	if req.Status != nil {
		issue.Status = *req.Status
	}
	
	issue.UpdatedAt = time.Now()
	return issue
}