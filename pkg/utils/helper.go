package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type RowData struct {
	Cells []string
}

// var isInt = regexp.MustCompile(`/^-?\d+\.?\d*$/`)
func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

func takesFloat(arr []string, index int) float64 {
	// Check if the index is out of range
	if index < 0 || index >= len(arr) {
		return 0.0
	}

	return convertFloat(arr[index])
}

func convertFloat(str string) float64 {
	// Remove commas from the string
	valueStr := strings.ReplaceAll(str, ",", "")

	// Trim whitespace from the string
	valueStr = strings.TrimSpace(valueStr)

	// Convert the modified string value to a float64
	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		// Handle the case where the value cannot be converted to a float64
		//fmt.Printf("Error converting %q to float64: %v\n", arr[index], err)
		return 0.0
	}
	return value
}

func takes(arr []string, index int) string {

	if index < 0 || index >= len(arr) {
		return ""
	}
	return arr[index]
}

// remove space ` f_name    l_name ` to `f_name l_name`
func trimAllSpace(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

// Function to generate a random email address template
func randomEmailTemplate() string {
	// Get the current time
	currentTime := time.Now()

	// Generate a random number between 1000 and 9999
	randomNumber := rand.Intn(9000) + 1000

	// Construct the email template using the current time and random number
	emailTemplate := fmt.Sprintf("user%d_%s@payroll.com", randomNumber, currentTime.Format("20060102_150405"))

	return emailTemplate
}

// Function to generate a random string of numbers with a specified length
func randomNumericString(length int) string {
	// Create a new source using the current time as a seed
	source := rand.NewSource(time.Now().UnixNano())

	// Create a new random generator using the source
	random := rand.New(source)

	// Define the characters to be used in the random string
	const charset = "0123456789"

	// Create a byte slice with the specified length
	randomString := make([]byte, length)

	// Fill the byte slice with random characters from the charset
	for i := range randomString {
		randomString[i] = charset[random.Intn(len(charset))]
	}

	return string(randomString)
}
