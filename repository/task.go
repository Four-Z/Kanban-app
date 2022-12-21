package repository

import (
	"a21hc3NpZ25tZW50/entity"
	"context"

	"gorm.io/gorm"
)

type TaskRepository interface {
	GetTasks(ctx context.Context, id int) ([]entity.Task, error)
	StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error)
	GetTaskByID(ctx context.Context, id int) (entity.Task, error)
	GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error)
	UpdateTask(ctx context.Context, task *entity.Task) error
	DeleteTask(ctx context.Context, id int) error
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db}
}

func (r *taskRepository) GetTasks(ctx context.Context, id int) ([]entity.Task, error) {
	tasks := []entity.Task{}
	err := r.db.Raw("SELECT * FROM tasks WHERE user_id = ?", id).Scan(&tasks).Error

	if err != nil {
		return []entity.Task{}, err
	}

	return tasks, nil
}

func (r *taskRepository) StoreTask(ctx context.Context, task *entity.Task) (taskId int, err error) {
	// taskStored := entity.Task{}
	err = r.db.Create(task).Error

	if err != nil {
		return 0, err
	}

	return task.ID, nil
}

func (r *taskRepository) GetTaskByID(ctx context.Context, id int) (entity.Task, error) {
	task := entity.Task{}
	err := r.db.First(&entity.Task{}, id).Scan(&task).Error

	if err != nil {
		return entity.Task{}, err
	}

	return task, nil
}

func (r *taskRepository) GetTasksByCategoryID(ctx context.Context, catId int) ([]entity.Task, error) {
	task := []entity.Task{}
	err := r.db.Raw("SELECT * FROM tasks WHERE category_id = ?", catId).Scan(&task).Error

	if err != nil {
		return []entity.Task{}, err
	}

	return task, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, task *entity.Task) error {
	err := r.db.Table("tasks").Where("id = ?", task.ID).Updates(&task).Error

	return err

}

func (r *taskRepository) DeleteTask(ctx context.Context, id int) error {
	// err := r.db.Where("id = ?", id).Delete(&entity.Task{}).Error
	err := r.db.Delete(&entity.Task{}, id).Error

	return err
}
