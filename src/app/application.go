package app

import (
	"github.com/gin-gonic/gin"
	"github.com/twinemarron/bookstore_oauth-api/src/clients/cassandra"
	"github.com/twinemarron/bookstore_oauth-api/src/domain/access_token"
	http "github.com/twinemarron/bookstore_oauth-api/src/http/access_token"
	"github.com/twinemarron/bookstore_oauth-api/src/repository/db"
)

var (
	router = gin.Default()
)

func StartApplication() {
	session, dbErr := cassandra.GetSession()
	if dbErr != nil {
		panic(dbErr)
	}
	session.Close()

	atHandler := http.NewHandler(access_token.NewService(db.NewRepository()))

	router.GET("/oauth/access_token/:access_token_id", atHandler.GetById)
	router.POST("/oauth/access_token", atHandler.Create)

	router.Run(":8081")
}
