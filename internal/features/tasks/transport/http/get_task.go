package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/lfssxxx/Golang-Todoapp/internal/core/logger"
	core_http_request "github.com/lfssxxx/Golang-Todoapp/internal/core/transport/http/request"
	core_http_response "github.com/lfssxxx/Golang-Todoapp/internal/core/transport/http/response"
)

type GetTaskResponse TaskDTOResponse

// GetTask 		godoc
// @Summary 	Get Task
// @Description Get existing Task
// @Tags 		tasks
// @Param 		id path int true "Id of task you want to get"
// @Produce 	json
// @Success 	200 {object} GetTaskResponse "Task was found"
// @Failure 	400 {object} core_http_response.ErrorResponse "Bad Request"
// @Failure		404 {object} core_http_response.ErrorResponse "Task Not Found"
// @Failure 	500 {object} core_http_response.ErrorResponse "Internal Server Error"
// @Router		/tasks/{id} [get]
func (h *TasksHTTPHandler) GetTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get tasklID path value",
		)

		return
	}

	task, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get task",
		)
	}

	response := GetTaskResponse(taskDTOFromDomain(task))

	responseHandler.JSONResponse(response, http.StatusOK)
}
