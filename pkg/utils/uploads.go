package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func ProcessFile(file *multipart.FileHeader) error {
	// Save the uploaded file (optional)
	savedFile, _, err := SaveFile(file)
	if err != nil {
		return err
	}

	// dataSalary, err := ExtractSheetSalary(pathFile)
	// if err != nil {
	// 	return nil, err
	// }

	_, err = ExtractFile("./uploads/2024-03/เงินเดือน รพ.สต. มี.ค. 67(1).xlsx", "Detal")
	if err != nil {
		return err
	}
	defer savedFile.Close()
	return nil
}

func GetFolderName() (folderPath string, err error) {
	// Get current date in yyyy-mm-dd format
	date := time.Now().Format("2006-01")

	// Create folder with date name if not exists
	folder := fmt.Sprintf("./uploads/%s", date)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0755)
		if err != nil {
			return "", err
		}
	}

	return folder, nil
}

func SaveFile(file *multipart.FileHeader) (savedFile *os.File, pathFile string, err error) {
	// Open uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer uploadedFile.Close()

	// Create folder with date name if not exists
	folder, err := GetFolderName()
	if err != nil {
		return nil, "", err
	}

	// Save file to folder with date name
	fileName := filepath.Join(folder, file.Filename)
	savedFile, err = os.Create(fileName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create file: %w", err)
	}

	// Copy the uploaded file contents to the saved file
	_, err = io.Copy(savedFile, uploadedFile)
	if err != nil {
		return nil, "", fmt.Errorf("failed to save file: %w", err)
	}
	defer savedFile.Close()
	return savedFile, fileName, nil
}

func ExtractFile(pathFile string, sheet string) ([]map[string]string, error) {
	f, err := excelize.OpenFile(pathFile)
	if err != nil {
		return nil, err
	}

	// Get all the rows in the vegan section.
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, err
	}

	// Initialize a slice to hold the key names
	var keys []string

	// Extract data from the fourth row and use it as key names
	for _, colCell := range rows[2] {
		colCell = strings.ReplaceAll(colCell, "\n", "")
		colCell = strings.ReplaceAll(colCell, " ", "")
		// keys = append(keys, strconv.Itoa(idx))
		keys = append(keys, colCell)
	}

	var data []map[string]string

	// Process the data from subsequent rows starting from the tenth row
	for rowIndex := 3; rowIndex < len(rows)-11; rowIndex++ {
		row := rows[rowIndex]

		rowData := make(map[string]string)

		for colIndex, colCell := range row {
			// Use the keys extracted from the ninth row to create a map
			if colIndex < len(keys) {
				rowData[keys[colIndex]] = colCell
			}
		}

		// Append the row data to the slice
		data = append(data, rowData)
	}

	// Marshal the data into JSON format
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return nil, err
	}

	// Write the JSON data to a file
	jsonFile, err := os.Create("./uploads/2024-03/data.json")
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	_, err = jsonFile.Write(jsonData)
	if err != nil {
		return nil, err
	}

	fmt.Println("JSON file created successfully")
	return data, nil
}

// check col
func CheckFormatCol(rows []string, formatKeys []string, skip int) error {
	// check format file
	var keys []string

	// Extract data from the fourth row and use it as key names
	for _, colCell := range rows {
		colCell = strings.ReplaceAll(colCell, "\n", "")
		colCell = strings.ReplaceAll(colCell, " ", "")
		keys = append(keys, colCell)
	}

	// Check if each key exists in the expected format keys
	for idx, key := range keys {
		found := false
		for _, formatKey := range formatKeys {
			if key == formatKey || idx >= skip {
				found = true
				break
			}
		}
		if !found {
			// If key does not match any expected format key, return an error
			return fmt.Errorf("key '%s' does not match expected format", key)
		}
	}

	return nil
}
