package main

import (
	"log"
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/danielchwi/gbim/backend/controllers"
	"github.com/danielchwi/gbim/backend/database"
	"github.com/danielchwi/gbim/backend/middleware"

	"github.com/gin-gonic/gin"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// r.POST("/login", controllers.Login)
	// r.POST("/register", controllers.Register)

	// r.GET("/users", controllers.UserIndex)
	// r.POST("/user", controllers.UserStore)
	// r.GET("/user/:id", controllers.UserShow)
	// r.PUT("/user/:id", controllers.UserUpdate)
	// r.DELETE("/user/:id", controllers.UserDestroy)

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
		"foo":  "bar", // user:foo password:bar
		"manu": "123", // user:manu password:123
	}))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	authorized.POST("admin", func(c *gin.Context) {
		user := c.MustGet(gin.AuthUserKey).(string)

		// Parse JSON
		var json struct {
			Value string `json:"value" binding:"required"`
		}

		if c.Bind(&json) == nil {
			db[user] = json.Value
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		}
	})

	return r
}

func main() {
	database.Connect()
	authMiddleware, err := middleware.JwtMiddleware()
	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	if errInit := authMiddleware.MiddlewareInit(); errInit != nil {
		log.Fatal("authMiddleware.MiddlewareInit() Error:" + errInit.Error())
	}

	r := setupRouter()

	r.POST("/login", authMiddleware.LoginHandler)

	r.NoRoute(authMiddleware.MiddlewareFunc(), func(c *gin.Context) {
		claims := jwt.ExtractClaims(c)
		log.Printf("NoRoute claims: %#v\n", claims)
		c.JSON(404, gin.H{"code": "PAGE_NOT_FOUND", "message": "Page not found"})
	})

	// // Refresh time can be longer than token timeout
	// auth.GET("/refresh_token", authMiddleware.RefreshHandler)
	// r.Use(authMiddleware.MiddlewareFunc())
	// {
	r.GET("/users", authMiddleware.MiddlewareFunc(), controllers.UserIndex)
	r.POST("/user", authMiddleware.MiddlewareFunc(), controllers.UserStore)
	r.GET("/user/:id", authMiddleware.MiddlewareFunc(), controllers.UserShow)
	r.PUT("/user/:id", authMiddleware.MiddlewareFunc(), controllers.UserUpdate)
	r.DELETE("/user/:id", authMiddleware.MiddlewareFunc(), controllers.UserDestroy)

	r.GET("/persons", authMiddleware.MiddlewareFunc(), controllers.PersonIndex)
	r.POST("/person", authMiddleware.MiddlewareFunc(), controllers.PersonStore)
	r.GET("/person/:id", authMiddleware.MiddlewareFunc(), controllers.PersonShow)
	r.PUT("/person/:id", authMiddleware.MiddlewareFunc(), controllers.PersonUpdate)
	r.DELETE("/person/:id", authMiddleware.MiddlewareFunc(), controllers.PersonDestroy)
	//}

	//r.POST("/login", controllers.Login)
	r.POST("/register", controllers.Register)

	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
