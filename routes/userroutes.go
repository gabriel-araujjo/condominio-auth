package routes

import (
	"encoding/json"
	"net/http"

	"github.com/gabriel-araujjo/condominio-auth/errors"
)

type userContext struct {
	*context
}

func (c *userContext) login(w http.ResponseWriter, req *http.Request) {
	var params struct {
		Credential string `json:"cred" form:"cred" binding:"required"`
		Password   string `json:"passwd" binding:"required"`
	}

	decoder := json.NewDecoder(req.Body)

	if err := decoder.Decode(&params); err != nil {
		errors.WriteErrorWithCode(w, http.StatusBadRequest, "cannot decode json")
		return
	}
	userID, err := c.dao.User.Authenticate(params.Credential, params.Password)
	if err != nil {
		errors.WriteErrorWithCode(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	c.context.SetCurrentUserID(req, userID)
	c.context.PersistSession(req, w)
	w.WriteHeader(http.StatusNoContent)
}

func (c *userContext) get(w http.ResponseWriter, req *http.Request) {

}

func (c *userContext) create(w http.ResponseWriter, req *http.Request) {

}

func (c *userContext) delete(w http.ResponseWriter, req *http.Request) {

}
