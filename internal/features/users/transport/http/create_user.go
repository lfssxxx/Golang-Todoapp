package users_transport_http

import (
	"net/http"

	"github.com/lfssxxx/Golang-Todoapp/internal/core/domain"
	core_logger "github.com/lfssxxx/Golang-Todoapp/internal/core/logger"
	core_http_request "github.com/lfssxxx/Golang-Todoapp/internal/core/transport/http/request"
	core_http_response "github.com/lfssxxx/Golang-Todoapp/internal/core/transport/http/response"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100" example:"John Doe"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,e164" example:"+380123456789"`
}

type CreateUserResponse UserDTOResponse

// CreateUser 		godoc
// @Summary 		Create user
// @Description 	Create new user in db
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param 			request body CreateUserRequest true "CreateUser body request"
// @Success 		201 	{object} CreateUserResponse "Successfully created user"
// @Failure 		400 	{object} core_http_response.ErrorResponse "Bad Request"
// @Failure 		500 	{object} core_http_response.ErrorResponse "Internal Server Error"
// @Router 			/users [post]
func (h *UsersHTTPHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "falied to decode and validate HTTP request")

		return
	}

	userDomain := domainFromDTO(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create User")

		return
	}

	response := CreateUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusCreated)

}

func domainFromDTO(dto CreateUserRequest) domain.User {
	return domain.NewUserUnitialized(dto.FullName, dto.PhoneNumber)
}
