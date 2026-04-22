package users_transport_http

import (
	"fmt"
	"net/http"
	"regexp"

	"github.com/lfssxxx/Golang-Todoapp/internal/core/domain"
	core_logger "github.com/lfssxxx/Golang-Todoapp/internal/core/logger"
	core_http_request "github.com/lfssxxx/Golang-Todoapp/internal/core/transport/http/request"
	core_http_response "github.com/lfssxxx/Golang-Todoapp/internal/core/transport/http/response"
	core_http_types "github.com/lfssxxx/Golang-Todoapp/internal/core/transport/http/types"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"John Doe"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+3801233456789"`
}

func (r *PatchUserRequest) Validate() error {
	var phoneRegex = regexp.MustCompile(`^\+?[1-9]\d{7,14}$`)

	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("`FullName` can't be NULL")
		}

		fullNameLen := len([]rune(*r.FullName.Value))
		if fullNameLen < 3 || fullNameLen > 100 {
			return fmt.Errorf("`FullName` must be between 3 and 100")
		}

	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phone := *r.PhoneNumber.Value
			if !phoneRegex.MatchString(phone) {
				return fmt.Errorf("`PhoneNumber was not validated. Phone Number must be edited with e164 standart.")
			}
		}
	}
	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser 		godoc
// @Summary 		Patch User
// @Description 	Patch info of already existing User
// @Description 	### Field Update Logic (Three-state logic):
// @Description 	1. **Filed not set**: `phone_number` is ignored, value does not change
// @Description 	2. **Value is set**: `"phone_number": "+380123456789"` - Setting new phone_number in db
// @Description 	3. **Value is null**: `"phone_number": null` - Deleting field in db (set to NULL)
// @Description 	Restrictions: `full_name` can't be set to null
// @Tags 			users
// @Accept 			json
// @Produce 		json
// @Param			id path int true "ID of user you want to PATCH"
// @Param 			request body PatchUserRequest true "PatchUser body request"
// @Success 		200 	{object} PatchUserResponse "Succesfully patched user"
// @Failure 		400 	{object} core_http_response.ErrorResponse "Bad Request"
// @Failure 		404		{object} core_http_response.ErrorResponse "User Not Found"
// @Failure 		409		{object} core_http_response.ErrorResponse "Conflict"
// @Failure 		500 	{object} core_http_response.ErrorResponse "Internal Server Error"
// @Router 			/users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validaete http request")

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")

		return
	}

	response := PatchUserResponse(userDTOFromDomain(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(request.FullName.ToDomain(), request.PhoneNumber.ToDomain())
}
