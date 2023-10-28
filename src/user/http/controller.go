package http

import (
	"net/http"

	"github.com/apm-dev/oha/src/domain"
	"github.com/apm-dev/oha/src/httputils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type userController struct {
	userSvc   domain.UserService
	validator *validator.Validate
}

func RegisterUserHandlers(
	us domain.UserService,
	r *mux.Router,
) {
	c := &userController{
		userSvc:   us,
		validator: validator.New(),
	}
	r.Methods(http.MethodPost).Path("/users").HandlerFunc(c.createUser)
	r.Methods(http.MethodGet).Path("/users/{id}").HandlerFunc(c.getUserByID)
}

func (c *userController) getUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	req := &GetUserByIDRequest{
		ID: params["id"],
	}
	err := c.validator.Struct(params)
	if err != nil {
		err = errors.Wrap(domain.ErrInvalidArgument, err.Error())
		httputils.RespondWithErrJson(w, err, "getUserByID", "", nil)
		return
	}
	user, err := c.userSvc.GetUserByID(r.Context(), req.ID)
	if err != nil {
		httputils.RespondWithErrJson(w, err, "getUserByID", "", nil)
		return
	}
	httputils.RespondWithStructJSON(w, http.StatusOK, user)
}

func (c *userController) createUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	err := httputils.ReadRequestBody(r, c.validator, &req, false)
	if err != nil {
		err = errors.Wrap(domain.ErrInvalidArgument, err.Error())
		httputils.RespondWithErrJson(w, err, "createUser", "", nil)
		return
	}
	user, err := c.userSvc.GetUserByID(r.Context(), req.Name)
	if err != nil {
		httputils.RespondWithErrJson(w, err, "createUser", "", nil)
		return
	}
	httputils.RespondWithStructJSON(w, http.StatusOK, user)
}
