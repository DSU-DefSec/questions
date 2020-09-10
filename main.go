package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	questions     = []question{}
	eventTitle    string
	adminUsername string
	adminPassword string
)

type question struct {
	Time time.Time
	Ip   string
	Text string
}

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage:", os.Args[0], "event username password")
		os.Exit(1)
	}

	eventTitle = os.Args[1]
	adminUsername = os.Args[2]
	adminPassword = os.Args[3]

	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets", "./assets")

	initCookies(r)

	// Routes
	routes := r.Group("/")
	{
		routes.GET("/login", func(c *gin.Context) {
			c.HTML(http.StatusOK, "login.html", pageData(c, "login", nil))
		})
		routes.POST("/login", login)
		routes.GET("/logout", logout)
		routes.GET("/", viewQuestions)
		routes.POST("/", submitQuestion)
	}

	r.Run(":1157")
}

func viewQuestions(c *gin.Context) {
	c.HTML(http.StatusOK, "questions.html", pageData(c, "Questions", gin.H{"questions": questions}))
}

func submitQuestion(c *gin.Context) {
	message := "Question successfully submitted!"
	c.Request.ParseForm()
	questionText := c.Request.Form.Get("text")
	if questionText == "" || len(questionText) > 150 {
		message = "Question was empty or too long!"
	} else {
		questions = append([]question{{time.Now(), c.ClientIP(), questionText}}, questions...)
	}
	c.HTML(http.StatusOK, "questions.html", pageData(c, "Questions", gin.H{"questions": questions, "message": message}))
}

func pageData(c *gin.Context, title string, ginMap gin.H) gin.H {
	newGinMap := gin.H{}
	if len(getUser(c)) > 0 {
		newGinMap["user"] = getUser(c)
	}
	newGinMap["event"] = eventTitle
	for key, value := range ginMap {
		newGinMap[key] = value
	}
	return newGinMap
}
