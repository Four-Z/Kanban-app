package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

type successResponseTask struct {
	User_id int    `json:"user_id"`
	Task_id int    `json:"task_id"`
	Message string `json:"message"`
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	taskID := r.URL.Query().Get("task_id")

	if taskID == "" {
		listTask, err := t.taskService.GetTasks(r.Context(), idLogin)
		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(listTask)
		return
	}

	if taskID != "" {
		taskIDNumber, _ := strconv.Atoi(taskID)
		task, err := t.taskService.GetTaskByID(r.Context(), taskIDNumber)

		if err != nil {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
			return
		}

		w.WriteHeader(200)
		json.NewEncoder(w).Encode(task)
		return
	}

}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	// TODO: answer here

	if task.Title == "" || task.Description == "" || task.CategoryID == 0 {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	id := fmt.Sprintf("%s", r.Context().Value("id"))

	if id == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	idNumber, _ := strconv.Atoi(id)
	taskEntity := &entity.Task{
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
		UserID:      idNumber,
	}

	taskFromServ := entity.Task{}
	taskFromServ, errServ := t.taskService.StoreTask(r.Context(), taskEntity)

	if errServ != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(errServ)
		return
	}

	response := successResponseTask{
		User_id: idNumber,
		Task_id: taskFromServ.ID,
		Message: "success create new task",
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(response)

}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	taskID := r.URL.Query().Get("task_id")
	taskIDNumber, _ := strconv.Atoi(taskID)
	err := t.taskService.DeleteTask(r.Context(), taskIDNumber)

	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	id := fmt.Sprintf("%s", r.Context().Value("id"))
	idNumber, _ := strconv.Atoi(id)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": idNumber,
		"task_id": taskIDNumber,
		"message": "success delete task",
	})
}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	// TODO: answer here
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = &entity.Task{
		ID:          task.ID,
		Title:       task.Title,
		Description: task.Description,
		CategoryID:  task.CategoryID,
		UserID:      int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task",
	})

}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
