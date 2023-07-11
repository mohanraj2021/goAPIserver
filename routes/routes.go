package routes

import (
	"context"
	"controller/controller"
	"controller/helpers"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, PortalUser")
		c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func StartGin() {
	var srv *http.Server
	quit := make(chan bool)
	gin.DisableConsoleColor()

	router := gin.New()
	router.Use(CORSMiddleware())
	res := helpers.Wakeup()
	if !res {
		// logger.Fatal("Wake up failed, check logs ")
		os.Exit(1)
	}

	if mariastart := helpers.InitMariaDB(); !mariastart {
		// logger.Fatal("Unable to start MariaDB, check logs")
		os.Exit(1)
	}
	srv = &http.Server{Addr: ":2207", Handler: router}

	router.NoRoute(func(c *gin.Context) {
		c.AbortWithStatus(http.StatusNotFound)
	})
	router.GET("/dimdim", func(c *gin.Context) {
		quit <- true
		c.IndentedJSON(http.StatusOK, gin.H{"message": "Starting Shutdown"})
	})
	router.GET("/ping", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, gin.H{"agentid": "Mohan", "epoch": time.Now().Unix(), "message": "status ok"})
	})

	web := router.Group("/web")
	{
		web.POST("create/user", controller.Createuser)
		web.POST("login/user", controller.Userlogin)
		web.POST("/product", controller.Getproduct)
		secured := web.Group("secured").Use(controller.AuthUser())
		{
			secured.GET("/product", controller.Getproduct)
			secured.GET("/getcart", controller.Getcart)
			secured.POST("/addcart/:productid", controller.Addcart)
		}
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			// logger.Info("sgt_portal_controller", zap.String("message", fmt.Sprintf("unable to start service because %v", err.Error())), zap.String("sendto", string(sgttypes.Local)))
			fmt.Println("a ...any")
		}
	}()
	q := <-quit
	if q {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			// logger.Info("sgt_portal_controller", zap.String("message", fmt.Sprintf("server forced to shutdown because %v", err.Error())), zap.String("sendto", string(sgttypes.Local)))
			//log to redis
			fmt.Println("Unable to Shutdown")
		}
	}
}
