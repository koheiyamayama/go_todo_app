package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/koheiyamayama/go_todo_app/entity"
	"github.com/koheiyamayama/go_todo_app/store"
	"gopkg.in/go-playground/validator.v9"
)

type AddTask struct {
	Store     *store.TaskStore
	Validator *validator.Validate
}

func (at *AddTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var b struct {
		Title string `json:"title" validate:"required"`
	}
	// リクエストボディをb構造体に書き込む
	// 書き込む時にエラーが有る場合はレスポンスにエラーを書き込み返す
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	// Structの引数の構造体のvalidateタグに応じてvalidationを行う
	err := at.Validator.Struct(b)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusBadRequest)
		return
	}

	t := &entity.Task{
		Title:   b.Title,
		Status:  "todo",
		Created: time.Now(),
	}
	id, err := store.Tasks.Add(t)
	if err != nil {
		RespondJSON(ctx, w, &ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := struct {
		ID int `json:"id"`
	}{ID: int(id)}
	RespondJSON(ctx, w, rsp, http.StatusOK)
}
