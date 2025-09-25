package helper

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAuthenticatedUserID extracts the authenticated user's ID from the context.
func GetAuthenticatedUserID(c *gin.Context) (uint64, bool) {
	val, exists := c.Get("aud")
	if !exists {
		return 0, false
	}

	audStr, ok := val.(string)
	if !ok {
		return 0, false
	}

	aud, err := strconv.ParseUint(audStr, 10, 64)
	if err != nil || aud == 0 {
		return 0, false
	}
	return aud, true
}

// ParsePaginationParams reads limit and page query parameters and applies defaults.
func ParsePaginationParams(c *gin.Context) (limit, page int) {
	limit, _ = strconv.Atoi(c.Query("limit"))
	if limit <= 0 {
		limit = 10 // default page size is 10
	}

	page, err := strconv.Atoi(c.Query("page"))
	if err != nil || page < 1 {
		page = 1
	}

	return limit, page
}

// GetUint8QueryParam converts a query param string to uint8, defaulting to 0 if missing/invalid.
func GetUint8QueryParam(c *gin.Context, key string) uint8 {
	valStr := c.Query(key)
	val, _ := strconv.Atoi(valStr)
	if val < 0 {
		return 0
	}
	return uint8(val)
}
