package client

import (
	"errors"
	"fmt"
	"github.com/no-src/gofs/server"
	"github.com/no-src/gofs/util/httputil"
	"github.com/no-src/log"
	"net/http"
	"net/url"
)

// ErrSignIn the current user sign in failed
var ErrSignIn = errors.New("file server sign in failed")

// SignIn sign in the file server
func SignIn(scheme, host, userName, password string) ([]*http.Cookie, error) {
	loginUrl := fmt.Sprintf("%s://%s%s", scheme, host, server.LoginSignInFullRoute)
	form := url.Values{}
	form.Set(server.ParamUserName, userName)
	form.Set(server.ParamPassword, password)
	log.Debug("try to auto login file server %s=%s %s=%s", server.ParamUserName, userName, server.ParamPassword, password)
	loginResp, err := httputil.HttpPostWithoutRedirect(loginUrl, form)
	if err != nil {
		return nil, err
	}
	if loginResp.StatusCode == http.StatusFound {
		cookies := loginResp.Cookies()
		if len(cookies) > 0 {
			return cookies, nil
		}
	}
	return nil, ErrSignIn
}
