package rest

import (
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twinemarron/golang-restclient/rest"
)

func TestMain(m *testing.M) {
	rest.StartMockupServer()
	os.Exit(m.Run())
}

func TestLoginUserTimeoutFromApi(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "localhost/users/login",
		ReqBody:      `{"email":"tz3@tz.tz","password":"test"}`,
		RespHTTPCode: -1,
		RespBody:     `{}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("tz3@tz.tz", "test")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid restclient response when trying login user", err.Message)
}

func TestLoginUserInvalidErrorInterface(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "localhost/users/login",
		ReqBody:      `{"email":"tz3@tz.tz","password":"test"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid user credencials", "status": "404", "error": "not_found"}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("tz3@tz.tz", "test")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "invalid error interface when trying to login user", err.Message)
}

func TestLoginUserInvalidLoginCredentials(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "localhost/users/login",
		ReqBody:      `{"email":"tz3@tz.tz","password":"test"}`,
		RespHTTPCode: http.StatusNotFound,
		RespBody:     `{"message": "invalid user credencials", "status": 404, "error": "not_found"}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("tz3@tz.tz", "test")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusNotFound, err.Status)
	assert.EqualValues(t, "invalid user credencials", err.Message)
}

func TestLoginUserInvalidUserJsonResponse(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "localhost/users/login",
		ReqBody:      `{"email":"tz3@tz.tz","password":"test"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": "4", "first_name": "test", "last_name": "test", "email": "tz3@tz.tz"}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("tz3@tz.tz", "test")

	assert.Nil(t, user)
	assert.NotNil(t, err)
	assert.EqualValues(t, http.StatusInternalServerError, err.Status)
	assert.EqualValues(t, "error when trying unmarshall users response", err.Message)
}

func TestLoginUserNoError(t *testing.T) {
	rest.FlushMockups()
	rest.AddMockups(&rest.Mock{
		HTTPMethod:   http.MethodPost,
		URL:          "localhost/users/login",
		ReqBody:      `{"email":"tz3@tz.tz","password":"test"}`,
		RespHTTPCode: http.StatusOK,
		RespBody:     `{"id": 4, "first_name": "test", "last_name": "test", "email": "tz3@tz.tz"}`,
	})

	repository := usersRepository{}

	user, err := repository.LoginUser("tz3@tz.tz", "test")

	assert.Nil(t, err)
	assert.NotNil(t, user)
	assert.EqualValues(t, 4, user.Id)
	assert.EqualValues(t, "test", user.FirstName)
	assert.EqualValues(t, "test", user.LastName)
	assert.EqualValues(t, "tz3@tz.tz", user.Email)
}
