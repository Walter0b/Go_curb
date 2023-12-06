package components

import (
	"net/http"
	"reflect"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Get performs pagination on a GORM query with embedding and returns the paginated result.
func Get(c *gin.Context, query *gorm.DB, results interface{}, embedField string, id string) {
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

	// Pagination logic
	offset := (page - 1) * pageSize

	// Use reflection to check if the field exists
	resultsValue := reflect.ValueOf(results)
	if resultsValue.Kind() != reflect.Ptr || resultsValue.Elem().Kind() != reflect.Slice {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Results must be a pointer to a slice"})
		return
	}

	sliceElemType := resultsValue.Elem().Type().Elem()
	if _, found := sliceElemType.FieldByName(embedField); found {
		// Retrieve records with association
		if err := query.Limit(pageSize).Offset(offset).Preload(embedField).Find(results, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		// If the field is not found, fetch records without association
		if err := query.Limit(pageSize).Offset(offset).Find(results, id).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Create a response object with paginated data and metadata
	response := gin.H{
		"data":          results,                             // Data for the current page
		"totalRowCount": totalRowCount,                       // Total count of records
		"currentPage":   page,                                // Current page
		"pageSize":      pageSize,                            // Page size
		"totalPages":    (int(totalRowCount) + pageSize - 1), // Total pages
	}

	c.JSON(http.StatusOK, response)
}
