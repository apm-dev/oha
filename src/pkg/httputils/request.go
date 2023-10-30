package httputils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

func ReadRequestBody(r *http.Request, v *validator.Validate, result interface{}, allowEmptyBody bool) error {
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return errors.New("can't read request body data: " + err.Error())
	}
	if len(data) == 0 || result == nil {
		if allowEmptyBody {
			return nil
		} else {
			return errors.New("request body is missing")
		}
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		return errors.New("can't unmarshal request body: " + err.Error())
	}
	if v != nil {
		err = v.Struct(result)
		if err != nil {
			return err
		}
	}
	return nil
}
