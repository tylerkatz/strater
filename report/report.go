package report

import (
	"fmt"
	"os"
	"strconv"

	"encoding/csv"
	"encoding/json"

	"github.com/tylerkatz/strater/strategy"
	"github.com/xuri/excelize/v2"
)

func Generate(plans []*strategy.Plan, format string, outputPath string) error {
	switch format {
	case "csv":
		return generateCSV(plans, outputPath)
	case "xlsx":
		return generateXLSX(plans, outputPath)
	case "json":
		return generateJSON(plans, outputPath)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func generateCSV(plans []*strategy.Plan, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Month", "Starting Balance", "Ending Balance", "Profit Target", "Risk Per Trade"}
	if err := writer.Write(header); err != nil {
		return err
	}

	// Write data for each plan
	for _, plan := range plans {
		for _, result := range plan.MonthlyResults {
			row := []string{
				strconv.Itoa(result.Month),
				fmt.Sprintf("%.2f", result.StartingBalance),
				fmt.Sprintf("%.2f", result.EndingBalance),
				fmt.Sprintf("%.2f", result.ProfitTarget),
				fmt.Sprintf("%.2f", result.RiskPerTrade),
			}
			if err := writer.Write(row); err != nil {
				return err
			}
		}
	}

	return nil
}

func generateXLSX(plans []*strategy.Plan, outputPath string) error {
	f := excelize.NewFile()
	defer f.Close()

	// Set headers
	headers := []string{"Month", "Starting Balance", "Ending Balance", "Profit Target", "Risk Per Trade"}
	for i, header := range headers {
		cell := fmt.Sprintf("%c1", 'A'+i)
		f.SetCellValue("Sheet1", cell, header)
	}

	// Write data
	row := 2
	for _, plan := range plans {
		for _, result := range plan.MonthlyResults {
			f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), result.Month)
			f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), result.StartingBalance)
			f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), result.EndingBalance)
			f.SetCellValue("Sheet1", fmt.Sprintf("D%d", row), result.ProfitTarget)
			f.SetCellValue("Sheet1", fmt.Sprintf("E%d", row), result.RiskPerTrade)
			row++
		}
	}

	return f.SaveAs(outputPath)
}

func generateJSON(plans []*strategy.Plan, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	return encoder.Encode(plans)
}

// Implementation of generateCSV, generateXLSX, and generateJSON methods...
