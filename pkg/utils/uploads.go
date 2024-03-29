package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"maxl3oss/app/models"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
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

func SaveFile(file *multipart.FileHeader) (savedFile *os.File, pathFile string, err error) {
	// Open uploaded file
	uploadedFile, err := file.Open()
	if err != nil {
		return nil, "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer uploadedFile.Close()

	// Get current date in yyyy-mm-dd format
	date := time.Now().Format("2006-01")

	// Create folder with date name if not exists
	folder := fmt.Sprintf("./uploads/%s", date)
	if _, err := os.Stat(folder); os.IsNotExist(err) {
		err = os.MkdirAll(folder, 0755)
		if err != nil {
			return nil, "", err
		}
	}

	// Save file to folder with date name
	fileName := filepath.Join(folder, file.Filename)
	savedFile, err = os.Create(fileName)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create file: %w", err)
	}
	defer savedFile.Close()

	// Copy the uploaded file contents to the saved file
	_, err = io.Copy(savedFile, uploadedFile)
	if err != nil {
		return nil, "", fmt.Errorf("failed to save file: %w", err)
	}

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
	jsonFile, err := os.Create("./uploads/2024-03/data_datel.json")
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

func takesFloat(arr []string, index int) float64 {
	// Check if the index is out of range
	if index < 0 || index >= len(arr) {
		return 0.0
	}

	// Remove commas from the string
	valueStr := strings.ReplaceAll(arr[index], ",", "")

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

// for file salary
// for sheet เงินเดือน
func ExtractSheetSalary(pathFile string) ([]models.Salary, error) {
	f, err := excelize.OpenFile(pathFile)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows("เงินเดือน ") // Change "Sheet1" to your sheet name
	if err != nil {
		return nil, err
	}

	// add data form xlsx to models
	var employees []models.Salary

	for rowIndex := 7; rowIndex < len(rows)-1; rowIndex++ {
		row := rows[rowIndex]
		employee := models.Salary{
			FullName:            takes(row, 1),
			BankAccountNumber:   takes(row, 2),
			Salary:              takesFloat(row, 3),
			AdditionalBenefits:  takesFloat(row, 4),
			FixedIncome:         takesFloat(row, 5),
			MonthlyCompensation: takesFloat(row, 6),
			TotalIncome:         takesFloat(row, 7),

			Tax:                     takesFloat(row, 8),
			Cooperative:             takesFloat(row, 9),
			PublicHealthCooperative: takesFloat(row, 10),
			RevenueDepartment:       takesFloat(row, 11),
			DGS:                     takesFloat(row, 12),
			BangkokBank:             takesFloat(row, 13),
			ActualPay:               takesFloat(row, 14),
			Received:                takesFloat(row, 15),

			SocialSecurity: takes(row, 16),
			BankTransfer:   takes(row, 17),
		}
		employees = append(employees, employee)
	}
	f.Close()
	return employees, nil
}

func CheckFormatFileSalary(pathFile string) error {
	f, err := excelize.OpenFile(pathFile)
	if err != nil {
		return err
	}

	rows, err := f.GetRows("เงินเดือน ") // Change "Sheet1" to your sheet name
	if err != nil {
		return err
	}

	// check format file
	formatKeys := []string{
		"ลำดับที่",
		"ชื่อ-สกุล",
		"เลขบัญชีธนาคาร",
		"เงินเดือน",
		"เงินเพิ่มค่าครองชีพ",
		"เงินประจำตำแหน่ง",
		"ค่าตอบแทนรายเดือน",
		"รวมจ่ายจริง",
		"ภาษี",
		"กบข.",
		"สหกรณ์ฯ(สาธารณสุข)",
		"กรมสรรพากร(กยศ.)",
		"ฌกส.",
		"ธ.กรุงไทย",
		"รวมรับจริง",
		"รับจริง",
		"รวมประกันสังคม",
		"ส่งธนาคาร",
	}
	var keys []string

	// Extract data from the fourth row and use it as key names
	for _, colCell := range rows[3] {
		colCell = strings.ReplaceAll(colCell, "\n", "")
		colCell = strings.ReplaceAll(colCell, " ", "")
		keys = append(keys, colCell)
	}

	// Check if each key exists in the expected format keys
	for _, key := range keys {
		found := false
		for _, formatKey := range formatKeys {
			if key == formatKey {
				found = true
				break
			}
		}
		if !found {
			// If key does not match any expected format key, return an error
			return fmt.Errorf("key '%s' does not match expected format", key)
		}
	}

	f.Close()
	return nil
}

// for sheet Detal
func ExtractSheetDetail(pathFile string) ([]models.TransferInfo, error) {
	f, err := excelize.OpenFile(pathFile)
	if err != nil {
		return nil, err
	}

	rows, err := f.GetRows("Detal")
	if err != nil {
		return nil, err
	}

	// add data form xlsx to models
	var transferInfos []models.TransferInfo

	for rowIndex := 3; rowIndex < len(rows)-10; rowIndex++ {
		row := rows[rowIndex]
		transfer := models.TransferInfo{
			ReceivingBankCode: takes(row, 0),
			ReceivingACNo:     takes(row, 1),
			ReceiverName:      takes(row, 2),
			TransferAmount:    takesFloat(row, 3),
			CitizenIDTaxID:    takes(row, 4),
			DDARef:            takes(row, 5),
			ReferenceNoDDARef: takes(row, 6),
			Email:             takes(row, 7),
			MobileNo:          takes(row, 8),
		}
		transferInfos = append(transferInfos, transfer)
	}
	f.Close()
	return transferInfos, nil
}
