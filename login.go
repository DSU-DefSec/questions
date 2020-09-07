package main

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getUUID returns a randomly generated UUID from Google's UUID library.
func getUUID() string {
	return uuid.New().String()
}

// initCookies use gin-contrib/sessions{/cookie} to initalize a cookie store.
// It generates a random secret for the cookie store -- not ideal for continuity or invalidating previous cookies, but it's secure and it works
func initCookies(r *gin.Engine) {
	r.Use(sessions.Sessions("sarpedon", cookie.NewStore([]byte(getUUID()))))
}

var userkey = "user"

// login is a handler that parses a form and checks for specific data
func login(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")

	// Validate form input
	if strings.Trim(username, " ") == "" || strings.Trim(password, " ") == "" {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Username or password can't be empty 🙄"})
		return
	}

	if username != adminUsername || password != adminPassword {
		c.HTML(http.StatusBadRequest, "login.html", gin.H{"error": "Incorrect username or password."})
		return
	}

	// Save the username in the session
	session.Set(userkey, username)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/")
}

func getUser(c *gin.Context) string {
	session := sessions.Default(c)
	userName := session.Get("user")
	if userName == nil {
		return ""
	}
	return userName.(string)
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get(userkey)
	if user == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}
	session.Delete(userkey)
	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}
	c.Redirect(http.StatusSeeOther, "/login")
}
