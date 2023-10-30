package userhttp

import (
	"net/http"

	"github.com/apm-dev/oha/src/domain"
	"github.com/apm-dev/oha/src/pkg/httputils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type UserController struct {
	userSvc   domain.UserService
	validator *validator.Validate
}

func NewUserController(
	us domain.UserService,
) *UserController {
	return &UserController{
		userSvc:   us,
		validator: validator.New(),
	}
}

func (c *UserController) GetUserByID(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	req := &GetUserByIDRequest{
		ID: params["id"],
	}
	err := c.validator.Struct(req)
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

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	err := httputils.ReadRequestBody(r, c.validator, &req, false)
	if err != nil {
		err = errors.Wrap(domain.ErrInvalidArgument, err.Error())
		httputils.RespondWithErrJson(w, err, "createUser", "", nil)
		return
	}
	user, err := c.userSvc.AddNewUser(r.Context(), req.Name)
	if err != nil {
		httputils.RespondWithErrJson(w, err, "createUser", "", nil)
		return
	}
	httputils.RespondWithStructJSON(w, http.StatusOK, user)
}
