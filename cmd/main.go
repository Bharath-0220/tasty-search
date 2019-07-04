package main

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"

	"tasty-search/pkg/review"
)

func init() {
	review.ParseReviews()

	review.SortWordIndexes()
}

func main() {

	mainRouter := gin.Default()

	mainRouter.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	mainRouter.GET("/search", func(c *gin.Context) {
		f := review.SearchTopReviews(c)

		b, _ := json.Marshal(f)
		c.Data(http.StatusOK, "application/json", b)
	})

	_ = mainRouter.Run("0.0.0.0:8080")
}
