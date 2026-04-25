package handler

import (
	"comp-math-5/internal/algo"
	"comp-math-5/internal/numeric"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Point struct {
	X *float64 `json:"x" binding:"required"`
	Y *float64 `json:"y" binding:"required"`
}
type InterpolateRequest struct {
	Points []Point `json:"points" binding:"required,dive"`
	X      float64 `json:"x" binding:"required"`
}

func Solve() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req InterpolateRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
			return
		}

		if len(req.Points) < 2 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Need at least 2 points for interpolation"})
			return
		}

		pointsToInterpolate := make([]numeric.Point, len(req.Points))
		for i, p := range req.Points {
			pointsToInterpolate[i] = numeric.Point{
				X: *p.X,
				Y: *p.Y,
			}
		}

		results, err := algo.Interpolate(pointsToInterpolate, req.X)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Errorf("failed to compute interpolation: %w", err).Error()})
			return
		}

		c.JSON(http.StatusOK, results)
	}
}
