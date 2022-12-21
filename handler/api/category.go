package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
)

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

type successResponseCategory struct {
	User_id     int    `json:"user_id"`
	Category_id int    `json:"category_id"`
	Message     string `json:"message"`
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	id := fmt.Sprintf("%s", r.Context().Value("id"))
	idNumber, _ := strconv.Atoi(id)

	if id == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
	}

	category, err := c.categoryService.GetCategories(r.Context(), idNumber)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(err)
		return
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(category)

}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	// TODO: answer here

	if category.Type == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	id := fmt.Sprintf("%s", r.Context().Value("id"))

	if id == "" {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	idNumber, _ := strconv.Atoi(id)

	categoryEntity := &entity.Category{
		Type:   category.Type,
		UserID: idNumber,
	}

	catID, errServ := c.categoryService.StoreCategory(r.Context(), categoryEntity)

	if errServ != nil {
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(errServ)
		return
	}

	succResp := successResponseCategory{
		User_id:     catID.UserID,
		Category_id: catID.ID,
		Message:     "success create new category",
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(succResp)

}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	// TODO: answer here
	categoryID := r.URL.Query().Get("category_id")
	categoryIDNumber, _ := strconv.Atoi(categoryID)

	err := c.categoryService.DeleteCategory(r.Context(), categoryIDNumber)

	if err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	id := fmt.Sprintf("%s", r.Context().Value("id"))
	idNumber, _ := strconv.Atoi(id)

	succResp := successResponseCategory{
		User_id:     idNumber,
		Category_id: categoryIDNumber,
		Message:     "success delete category",
	}

	w.WriteHeader(200)
	json.NewEncoder(w).Encode(succResp)

}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)

}
