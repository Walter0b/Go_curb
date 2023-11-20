package components

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func GenerateUniqueInvoiceNumber() string {
	return fmt.Sprintf("INV-%d", time.Now().Unix())
}
func ConvertStringToFloat64(amountStr string) float64 {
	// Remove non-numeric characters from the string
	cleanedStr := strings.ReplaceAll(amountStr, "$", "")
	cleanedStr = strings.ReplaceAll(cleanedStr, ",", "")

	// Parse the cleaned string to a float64
	amountFloat, err := strconv.ParseFloat(cleanedStr, 64)
	if err != nil {
		return 0
	}

	return amountFloat
}

func GenerateRandomSlug() int64 {
	rand.Seed(time.Now().UnixNano())
	min := int64(100)
	max := int64(999)
	return rand.Int63n(max-min+1) + min
}
