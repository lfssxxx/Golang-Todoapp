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
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
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
