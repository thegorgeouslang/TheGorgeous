// Author: James Mallon <jamesmallondev@gmail.com>
// session package -
package session

import (
	conf "GoAuthentication/configs"
	. "GoAuthentication/models"
	. "GoAuthentication/models/dao"
	"github.com/satori/go.uuid"
	"net/http"
	"strconv"
)

// Struct type session -
type sessionHelper struct {
	sessExp int
}

// init function - data and process initialization
func SessionHelper() *sessionHelper {
	se, _ := strconv.Atoi(conf.Env["session_exp"])
	return &sessionHelper{se}
}

// GetSession method -
func (this *sessionHelper) Start(w http.ResponseWriter, r *http.Request) string {
	c, e := r.Cookie("session") // create the cookie
	if e != nil {
		sID := uuid.NewV4() // create the universal unique id
		c = &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
	}
	c.MaxAge = this.sessExp
	http.SetCookie(w, c)
	return c.Value
}

// Get method -
func (this *sessionHelper) GetSession(w http.ResponseWriter, r *http.Request) (e error) {
	c, e := r.Cookie("session") // create the cookie
	if e != nil {
		return
	}
	session := SessionDAO.GetSession(c.Value)
	if session.Email != "" {
		SessionDAO.Renew(c.Value)
	}
	// refresh session
	c.MaxAge = this.sessExp
	http.SetCookie(w, c)
	return
}

//
func (this *sessionHelper) User(w http.ResponseWriter, r *http.Request) (user User) {
	c, e := r.Cookie("session") // create the cookie
	if e == nil {
		c.MaxAge = this.sessExp
		http.SetCookie(w, c)

		// if the user exists already, get user
		session := SessionDAO.GetSession(c.Value) // retrieve the session
		if session.UserId > 0 {                   // check for the user id
			SessionDAO.Renew(c.Value)              // update LastActivity
			user = UserDAO.GetUser(session.UserId) // retrieve user
			return
		}
	}
	return
}

// CloseSession method -
func (this *sessionHelper) Close(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session") // create the cookie
	if err == nil {
		c.MaxAge = -1
		http.SetCookie(w, c)
	}
}