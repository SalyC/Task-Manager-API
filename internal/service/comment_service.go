package service

import (
	"errors"
	"taskmanager/internal/models"
)

type CommentService struct {
	commentRepo CommentRepositoryInterface
	taskSvc     *TaskService
}

func NewCommentService(CommentRepo CommentRepositoryInterface, taskSvc *TaskService) *CommentService {
	return &CommentService{
		commentRepo: CommentRepo,
		taskSvc:     taskSvc,
	}
}

var (
	ErrCommentNotFound        = errors.New("comment not found")
	ErrCommentContentRequired = errors.New("comment content is required")
)

func (c *CommentService) Create(comment *models.Comment) error {
	if comment.Content == "" {
		return ErrCommentContentRequired
	}
	if _, err := c.taskSvc.GetByID(comment.TaskID); err != nil {
		return err
	}
	return c.commentRepo.Create(comment)
}
func (c *CommentService) GetByTaskID(taskID int) ([]models.Comment, error) {
	return c.commentRepo.GetByTaskID(taskID)
}
func (c *CommentService) Delete(id int) error {
	comment, err := c.commentRepo.GetByID(id)
	if err != nil {
		return err
	}
	if comment == nil {
		return ErrCommentNotFound
	}
	return c.commentRepo.Delete(id)
}
