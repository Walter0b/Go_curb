package components

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PaginateWithEmbed performs pagination on a GORM query with embedding and returns the paginated result.
func PaginateWithEmbed(c *gin.Context, query *gorm.DB, results interface{}, resultsEmbed interface{} , embedType reflect.Type, embedField string) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page number"})
		return
	}

	// Retrieve total count for metadata
	var totalRowCount int64
	if err := query.Count(&totalRowCount).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Calculate pageSize
	pageSize := int(totalRowCount) // Set pageSize to totalRowCount by default

	// If pageSize is specified by the user, parse it
	if pageSizeStr := c.Query("pageSize"); pageSizeStr != "" {
		pageSizeInt, err := strconv.Atoi(pageSizeStr)
		if err != nil || pageSizeInt < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page size"})
			return
		}
		pageSize = pageSizeInt
	}

	// Now, use the pageSize in your pagination logic
	offset := (page - 1) * pageSize

	if embedField != "" {
		if _, found := embedType.FieldByName(embedField); found {
			// Check if the field is a slice or not

			// Retrieve records with association
			if err := query.Limit(pageSize).Offset(offset).Preload(embedField).Find(resultsEmbed).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			// Combine association and Invoice information in the response
			results = resultsEmbed
		}
	} else if err := query.Limit(pageSize).Offset(offset).Find(results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Create a response object with paginated data and metadata
	response := gin.H{
		"data":          results,                                        // Data for the current page
		"totalRowCount": totalRowCount,                                  // Total count of records
		"currentPage":   page,                                           // Current page
		"pageSize":      pageSize,                                       // Page size
		"totalPages":    (int(totalRowCount) + pageSize - 1) , // Total pages
	}

	c.JSON(http.StatusOK, response)
}
