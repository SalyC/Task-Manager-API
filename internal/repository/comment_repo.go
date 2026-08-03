package repository

import (
	"database/sql"
	"taskmanager/internal/models"
)

type CommentRepository struct {
	db *sql.DB
}

func NewCommentRepository(db *sql.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (c *CommentRepository) Create(comment *models.Comment) error {
	query := "INSERT INTO comments (task_id, author, content) VALUES ($1, $2, $3) RETURNING id, created_at"
	return c.db.QueryRow(query, &comment.TaskID, &comment.Author, &comment.Content).Scan(&comment.ID, &comment.CreatedAt)
}
func (r *CommentRepository) GetByID(id int) (*models.Comment, error) {
	var comment models.Comment
	query := `SELECT id, task_id, author, content, created_at FROM comments WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&comment.ID, &comment.TaskID, &comment.Author, &comment.Content, &comment.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &comment, nil
}
func (c *CommentRepository) GetByTaskID(taskID int) ([]models.Comment, error) {
	query := "SELECT id, task_id, author, content, created_at FROM comments WHERE task_id = $1 ORDER BY created_at"
	rows, err := c.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var cmt models.Comment
		err := rows.Scan(&cmt.ID, &cmt.TaskID, &cmt.Author, &cmt.Content, &cmt.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, cmt)
	}

	return comments, rows.Err()
}
func (c *CommentRepository) Delete(id int) error {
	_, err := c.db.Exec("DELETE FROM comments WHERE id = $1", id)
	return err
}
