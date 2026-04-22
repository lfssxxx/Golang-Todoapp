package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/lfssxxx/Golang-Todoapp/internal/core/domain"
	core_logger "github.com/lfssxxx/Golang-Todoapp/internal/core/logger"
	core_http_request "github.com/lfssxxx/Golang-Todoapp/internal/core/transport/http/request"
	core_http_response "github.com/lfssxxx/Golang-Todoapp/internal/core/transport/http/response"
	core_http_types "github.com/lfssxxx/Golang-Todoapp/internal/core/transport/http/types"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title" swaggertype:"string" example:"Prepare for Camping"`
	Description core_http_types.Nullable[string] `json:"description" swaggertype:"string" example:"null"`
	Completed   core_http_types.Nullable[bool]   `json:"completed" swaggertype:"boolean"`
}

type PatchTaskResponse TaskDTOResponse

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("`Title` can't be NULL")
		}
		titleLen := len([]rune(*r.Title.Value))
		if titleLen < 1 || titleLen > 100 {
			return fmt.Errorf("`Title` must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLen := len([]rune(*r.Description.Value))
			if descriptionLen < 1 || descriptionLen > 100 {
				return fmt.Errorf("`Description` must be between 1 and 1000 symbols")

			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can't be NULL")

		}
	}

	return nil
}

// PatchTask 		godoc
// @Summary 		Patch User
// @Description 	Patch info of already existing Task
// @Description 	### Field Update Logic (Three-state logic):
// @Description 	1. **Filed not set**: `description` is ignored, value does not change
// @Description 	2. **Value is set**: `"description": "Prepare for Camping"` - Setting new description in db
// @Description 	3. **Value is null**: `"description": null` - Deleting field in db (set to NULL)
// @Description 	Restrictions: `title` and `completed` can't be set to null
// @Tags 			tasks
// @Accept 			json
// @Produce 		json
// @Param			id path int true "ID of Task you want to PATCH"
// @Param 			request body PatchTaskRequest true "PatchTask body request"
// @Success 		200 	{object} PatchTaskResponse "Succesfully patched Task"
// @Failure 		400 	{object} core_http_response.ErrorResponse "Bad Request"
// @Failure 		404		{object} core_http_response.ErrorResponse "Task Not Found"
// @Failure 		409		{object} core_http_response.ErrorResponse "Conflict"
// @Failure 		500 	{object} core_http_response.ErrorResponse "Internal Server Error"
// @Router 			/tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get taskID path value",
		)

		return
	}

	var request PatchTaskRequest

	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request",
		)

		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to patch task",
		)
	}

	response := PatchTaskResponse(taskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
