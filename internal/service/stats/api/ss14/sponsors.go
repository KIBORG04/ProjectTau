package ss14

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"ssstatistics/internal/repository"
)

func GetSponsor(c *gin.Context) (int, any) {
	userId := c.Param("userid")
	_, err := repository.FindSponsor(userId)
	if err != nil {
		return 404, gin.H{
			"code":  "404",
			"error": fmt.Sprintf("%s not round", userId),
		}
	}
	return 200, gin.H{
		"tier":         1,
		"oocColor":     "#9b59b6",
		"priorityJoin": true,
	}
}
