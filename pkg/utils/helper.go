package utils

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/shakinm/xlsReader/xls"
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

// Function to clean up email addresses
func cleanEmailAddress(email string) string {
	// Define a regular expression to match valid email characters
	validEmailRegex := regexp.MustCompile(`[[:alnum:]!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[[:alnum:]!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[[:alnum:]-]+\.)+[[:alpha:]]{2,7}`)

	// Find valid email addresses in the string
	validEmail := validEmailRegex.FindString(email)

	// Remove any leading or trailing whitespace
	cleanEmail := strings.TrimSpace(validEmail)

	return cleanEmail
}

func GetRowsFromSheet(xlsFile xls.Workbook, name string, colTotal int) ([]RowData, error) {
	var rows []RowData

	sheet, err := getSheetByName(xlsFile, name)
	if err != nil {
		return nil, err
	}

	numRows := sheet.GetNumberRows()
	for i := 0; i < numRows; i++ {
		rowData, err := extractRowData(sheet, i, colTotal)
		if err != nil {
			return nil, err
		}

		if rowData != nil {
			rows = append(rows, *rowData)
		}
	}

	return rows, nil
}

func extractRowData(sheet *xls.Sheet, rowIndex, colTotal int) (*RowData, error) {
	rowData := RowData{}

	row, err := sheet.GetRow(rowIndex)
	if err != nil {
		return nil, err
	}

	// Check if the first cell contains a numeric index
	cellFirst, err := row.GetCol(0)
	if err != nil {
		return nil, err
	}
	cellIndex := cellFirst.GetString()
	if !isNumeric(cellIndex) || strings.TrimSpace(cellIndex) == "" {
		// Skip non-numeric index rows
		return nil, nil
	}

	// Extract data from each column
	for j := 0; j <= colTotal; j++ {
		cell, err := row.GetCol(j)
		if err != nil {
			return nil, err
		}
		cellValue := cell.GetString()
		rowData.Cells = append(rowData.Cells, cellValue)
	}

	return &rowData, nil
}

func getSheetByName(workbook xls.Workbook, name string) (*xls.Sheet, error) {
	numSheets := workbook.GetNumberSheets()
	for i := 0; i < numSheets; i++ {
		sheet, err := workbook.GetSheet(i)
		if err != nil {
			return nil, err
		}
		if sheet.GetName() == name {
			return sheet, nil
		}
	}
	return nil, fmt.Errorf("sheet not found: %s", name)
}
