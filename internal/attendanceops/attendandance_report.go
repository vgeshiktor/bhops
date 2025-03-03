package attendanceops

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/xuri/excelize/v2"
)

type WorkerDetails struct {
	WorkerID           string  `json:"worker_id"`
	Name               string  `json:"name"`
	Type               string  `json:"worker_type"`
	DailyHours         float64 `json:"daily_hours"`
	PerHour            float64 `json:"per_hour"`
	PerHour125         float64 `json:"per_hour_125"`
	MonthlySal         float64 `json:"monthly_sal"`
	TransExpanses      float64 `json:"trans_expanses"`
	Holidays           float64 `json:"holidays"`
	HolidayPresent     float64 `json:"holiday_present"`
	HoursAdjustment    float64 `json:"hours_adjustment"`
	Hours125Adjustment float64 `json:"hours_125_adjustment"`
	VacDaysAdjustment float64 `json:"vac_days_adjustment"`
}

type Worker struct {
	WorkerID        string  `json:"id"`
	Name            string  `json:"name"`
	WorkerType      string  `json:"worker_type"`
	DailyHours      float64 `json:"daily_hours"`
	Hours           float64 `json:"hours"`
	PerHour         float64 `json:"per_hour"`
	RegularHoursSal float64 `json:"reg_hours_sal"`
	Hours125        float64 `json:"hours_125"`
	PerHour125      float64 `json:"per_hour_125"`
	ExtraHoursSal   float64 `json:"extra_hours_sal"`
	MonthlySal      float64 `json:"monthly_sal"`
	TransExpanses   float64 `json:"trans_expanses"`
	TotalHours      float64 `json:"total_hours"`
	WorkDays        float64 `json:"work_days"`
	Holidays        float64 `json:"holidays"`
	HolidayPresent  float64 `json:"holiday_present"`
	SickDays        float64 `json:"sick_days"`
	VacDays         float64 `json:"vac_days"`
	AbsenseHours    float64 `json:"absense_hours"`
	TotalSal        float64 `json:"total_sal"`
}

type AttendanceReport struct {
	AttendanceReportPath    string
	NonAttendanceReportPath string
	WorkerDetailsPath       string
	workerDetails           map[string]WorkerDetails
	workers                 []Worker
	// file                    *excelize.File
}

func NewAttendanceReport(
	attendanceReportPath,
	nonAttendanceReportPath,
	workerDetailsPath string) *AttendanceReport {
	return &AttendanceReport{
		AttendanceReportPath:    attendanceReportPath,
		NonAttendanceReportPath: nonAttendanceReportPath,
		WorkerDetailsPath:       workerDetailsPath,
	}
}

func CreateAttendanceReport(
	attendanceReportPath,
	nonAttendanceReportPath,
	workerDetailsPath string,
) (*AttendanceReport, error) {
	// create workers attendance monthly report
	attendanceReport := NewAttendanceReport(
		attendanceReportPath,
		nonAttendanceReportPath,
		workerDetailsPath,
	)

	// load worker details
	err := attendanceReport.loadWorkerDetails(workerDetailsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load worker details, error: %w", err)
	}

	// add attendance workers to monthly report
	err = attendanceReport.addAttendanceWorkers(attendanceReportPath)
	if err != nil {
		return nil, fmt.Errorf("failed to add attendance workers to monthly report, error: %w", err)
	}

	// add non-attendance workers to monthly report
	if err = attendanceReport.addNonAttendanceWorkers(nonAttendanceReportPath); err != nil {
		return nil, fmt.Errorf("failed to add non-attendance workers from: %s to monthly report, error: %w", nonAttendanceReportPath, err)
	}

	return attendanceReport, nil
}

func SaveAttendanceReport(
	attendanceReport *AttendanceReport,
	attendanceReportPath string,
) error {
	// create excel file from attendance report
	f, err := attendanceReport.createExcelSheet()
	if err != nil {
		return fmt.Errorf("failed to create excel sheet from attendance report, error: %w", err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Error().Err(err)
		}
	}()

	// save excel file
	if err := f.SaveAs(attendanceReportPath); err != nil {
		return fmt.Errorf("failed to save excel file: %s, error: %w", attendanceReportPath, err)
	}

	return nil
}

func (a *AttendanceReport) loadWorkerDetails(workerDetailsPath string) error {
	// open worker details file
	file, err := os.Open(workerDetailsPath)
	if err != nil {
		return fmt.Errorf("failed to open worker details file: %s, error: %v",
			workerDetailsPath, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error().Msgf("failed to close worker details file: %s, error: %v",
				workerDetailsPath, err)
		}
	}()

	// read worker details file
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("failed to read worker details file: %s, error: %v",
			workerDetailsPath, err)
	}

	// unmarshal worker details
	if err := json.Unmarshal(content, &a.workerDetails); err != nil {
		return fmt.Errorf("failed to unmarshal worker details, error: %v", err)
	}

	// fix worker ids
	for workerID, details := range a.workerDetails {
		details.WorkerID = workerID
		a.workerDetails[workerID] = details
	}

	return nil
}

func (a *AttendanceReport) addAttendanceWorkers(attendanceReportPath string,
) error {
	// open attendance report file
	f, err := excelize.OpenFile(attendanceReportPath)
	if err != nil {
		return fmt.Errorf("failed to open attendance report file: %s, error: %v",
			attendanceReportPath, err)
	}

	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			log.Error().Msgf("failed to close attendance report file: %s  error: %v", attendanceReportPath, err)
		}
	}()

	// get sheet to use
	sheet := f.GetSheetName(0)
	log.Info().Msgf("Processing sheet: %s\n", sheet)

	rows, err := f.GetRows(sheet)
	if err != nil {
		return fmt.Errorf("failed to get rows from sheet: %s, error: %w", sheet, err)
	}

	// get start and end row for workers
	startRow := 2
	endRow := startRow + 5

	// add workers to attendance report
	for i := startRow; i < endRow; i++ {
		// create worker report
		worker, err := a.createWorkerReport(rows[i])
		if err != nil {
			return fmt.Errorf("failed to create worker report for worker id: %s, error: %w", rows[i][4], err)
		}

		// add worker report to attendance report
		a.workers = append(a.workers, worker)
	}

	return nil
}

func (a *AttendanceReport) createWorkerReport(row []string) (Worker, error) {
	workerID := row[3]
	if _, ok := a.workerDetails[workerID]; !ok {
		return Worker{}, fmt.Errorf("worker details not found for workerID: %s", workerID)
	}
	worker := Worker{
		WorkerID:        workerID,
		Name:            a.workerDetails[workerID].Name,
		WorkerType:      a.workerDetails[workerID].Type,
		DailyHours:      a.workerDetails[workerID].DailyHours,
		Hours:           TimeStrToFloat64(row[10]) + a.workerDetails[workerID].HoursAdjustment,
		PerHour:         a.workerDetails[workerID].PerHour,
		RegularHoursSal: 0,
		Hours125:        TimeStrToFloat64(row[12]) + a.workerDetails[workerID].Hours125Adjustment,
		PerHour125:      a.workerDetails[workerID].PerHour125,
		ExtraHoursSal:   0,
		MonthlySal:      a.workerDetails[workerID].MonthlySal,
		TransExpanses:   a.workerDetails[workerID].TransExpanses,
		TotalHours:      0,
		WorkDays:        StrToFloat64(row[7]),
		Holidays:        a.workerDetails[workerID].Holidays,
		HolidayPresent:  a.workerDetails[workerID].HolidayPresent,
		SickDays:        StrToFloat64(row[20]),
		VacDays:         StrToFloat64(row[21]) + 
			a.workerDetails[workerID].VacDaysAdjustment,
		AbsenseHours: func() float64 {
			if a.workerDetails[workerID].Type == "monthly" {
				return TimeStrToFloat64(row[15])
			} else {
				return 0
			}
		}(),
		TotalSal: 0,
	}

	return worker, nil
}

func (a *AttendanceReport) addNonAttendanceWorkers(
	nonAttendanceReportPath string) error {
	f, err := os.Open(nonAttendanceReportPath)
	if err != nil {
		return fmt.Errorf("failed to open non-attendance report file: %s, error: %w",
			nonAttendanceReportPath, err)
	}
	defer func() {
		if err := f.Close(); err != nil {
			log.Error().Msgf("failed to close non-attendance report file: %s, error: %v",
				nonAttendanceReportPath, err)
		}
	}()

	content, err := io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("failed to read non-attendance report file: %s, error: %w",
			nonAttendanceReportPath, err)
	}

	var nonAttendanceWorkers []Worker
	if err := json.Unmarshal(content, &nonAttendanceWorkers); err != nil {
		return fmt.Errorf("failed to unmarshal non-attendance workers, error: %w", err)
	}

	// update worker reports
	for i := range nonAttendanceWorkers {
		nonAttendanceWorkers[i].MonthlySal =
			a.workerDetails[nonAttendanceWorkers[i].WorkerID].MonthlySal
		nonAttendanceWorkers[i].TransExpanses =
			a.workerDetails[nonAttendanceWorkers[i].WorkerID].TransExpanses
		nonAttendanceWorkers[i].Hours = a.workerDetails[
			nonAttendanceWorkers[i].WorkerID].DailyHours * 
			nonAttendanceWorkers[i].WorkDays
	}

	a.workers = append(a.workers, nonAttendanceWorkers...)

	return nil

}

func (a *AttendanceReport) createExcelSheet() (*excelize.File, error) {
	f := excelize.NewFile()

	// Create a new sheet.
	sheetName := "Sheet1"
	index, err := f.NewSheet("Sheet1")
	if err != nil {
		return nil, fmt.Errorf("failed to create new sheet, error: %w", err)
	}

	// Set the active sheet to the RTL sheet
	f.SetActiveSheet(index)

	// Set the sheet view to RTL
	RightToLeft := true
	err = f.SetSheetView(sheetName, 0, &excelize.ViewOptions{
		RightToLeft: &RightToLeft, // Set RightToLeft to true
	})
	if err != nil {
		log.Error().Msgf("failed to set sheet view to RTL, error: %w", err)
	}

	// Set the width of column A to 20
	if err := f.SetColWidth(
		sheetName,
		"A",
		"D",
		20); err != nil {
		log.Error().Err(err)
	}

	// iterate over workers and add them to the sheet
	start_row := 1
	for _, worker := range a.workers {
		start_row = writeWorkerToSheet(f, worker, start_row)
	}

	// set column width

	return f, nil
}

func writeWorkerToSheet(f *excelize.File, worker Worker, startRow int) int {
	writeCell :=
		func(row, col int, value string, style *excelize.Style) {
			styleID, err := f.NewStyle(style)
			if err != nil {
				log.Error().Msgf("failed to create cell style for cell: %s, error: %v", cellName(row, col), err)
			}

			cell := cellName(row, col)

			if err := f.SetCellStyle(
				"Sheet1",
				cell,
				cell,
				styleID); err != nil {
				log.Error().Msgf(
					"failed to set cell style for cell: %s, error: %v", cell, err)
			}

			if err := f.SetCellValue(
				"Sheet1", cell, value); err != nil {
				log.Error().Msgf(
					"failed to set cell value for cell: %s, error: %v", cell, err)
			}
		}

	writeCellFormula :=
		func(row, col int, value string, style *excelize.Style) {
			styleID, err := f.NewStyle(style)
			if err != nil {
				log.Error().Msgf("failed to create cell style for cell: %s, error: %v", cellName(row, col), err)
			}

			cell := cellName(row, col)

			if err := f.SetCellStyle(
				"Sheet1",
				cell,
				cell,
				styleID); err != nil {
				log.Error().Msgf(
					"failed to set cell style for cell: %s, error: %v", cell, err)
			}

			if err := f.SetCellFormula(
				"Sheet1", cell, value); err != nil {
				log.Error().Msgf(
					"failed to set cell value for cell: %s, error: %v", cell, err)
			}
		}

	// worker name
	writeCell(
		startRow, 1,
		worker.Name,
		TitleCellStyle())

	// writing title in the second row of the table
	writeCell(
		startRow+1, 1,
		"",
		HeaderCellStyle())

	writeCell(
		startRow+1, 2,
		"שעות",
		HeaderCellStyle())

	writeCell(
		startRow+1, 3,
		"לשעה ₪",
		HeaderCellStyle())

	writeCell(
		startRow+1, 4,
		"סה״כ ₪",
		HeaderCellStyle())

	// Writing regular hours row in the 3rd row of the table
	writeCell(
		startRow+2, 1,
		"ש.רגילות",
		HeaderCellStyle())

	writeCell(
		startRow+2, 2,
		strconv.FormatFloat(
			worker.Hours, 'f', 0, 64),
		NumericCellStyle())

	writeCell(
		startRow+2, 3,
		strconv.FormatFloat(
			worker.PerHour, 'f', 0, 64),
		NumericCellStyle())

	if worker.WorkerType == "hourly" || worker.WorkerType == "daily" {
		writeCellFormula(
			startRow+2, 4,
			fmt.Sprintf(
				"=%s*%s",
				cellName(startRow+2, 2),
				cellName(startRow+2, 3)),
			NumericCellStyle())
	} else {
		writeCell(
			startRow+2, 4,
			strconv.FormatFloat(
				worker.MonthlySal, 'f', 0, 64),
			NumericCellStyle())
	}

	// Empty row (start_row + 3)

	// Writing extra hours row in the 4th row of the table
	writeCell(
		startRow+4, 1,
		"ש.נ. 125%",
		HeaderCellStyle())

	writeCell(
		startRow+4, 2,
		strconv.FormatFloat(
			worker.Hours125, 'f', 0, 64),
		NumericCellStyle())

	writeCell(
		startRow+4, 3,
		strconv.FormatFloat(
			worker.PerHour125, 'f', 0, 64),
		NumericCellStyle())

	if worker.WorkerType == "hourly" {
		writeCellFormula(
			startRow+4, 4,
			fmt.Sprintf(
				"=%s*%s",
				cellName(startRow+4, 2),
				cellName(startRow+4, 3)),
			NumericCellStyle())
	}

	// Empty row (start_row + 5)

	// Writing transportation expenses row in the 6th row of the table
	writeCell(
		startRow+6, 1,
		"נסיעות",
		HeaderCellStyle())

	writeCell(
		startRow+6, 4,
		strconv.FormatFloat(
			worker.TransExpanses, 'f', 0, 64),
		NumericCellStyle())

	// Writing total salary row in the 7th row of the table
	writeCell(
		startRow+7, 1,
		"סה״כ ₪",
		HeaderCellStyle())

	writeCellFormula(
		startRow+7, 2,
		fmt.Sprintf(
			"=%s+%s",
			cellName(startRow+2, 2),
			cellName(startRow+4, 2)),
		NumericCellStyle())

	writeCellFormula(
		startRow+7, 4,
		fmt.Sprintf(
			"=%s+%s+%s",
			cellName(startRow+2, 4),
			cellName(startRow+4, 4),
			cellName(startRow+6, 4)),
		NumericCellStyle())

	// Empty row (start_row + 8)

	// Writing work days row in the 9th row of the table
	writeCell(
		startRow+9, 1,
		"ימי עבודה",
		HeaderCellStyle())

	writeCell(
		startRow+9, 2,
		strconv.FormatFloat(
			worker.WorkDays, 'f', 0, 64),
		NumericCellStyle())

	// Writing holiday row in the 10th row of the table
	writeCell(
		startRow+10, 1,
		"חג",
		HeaderCellStyle())

	writeCell(
		startRow+10, 2,
		strconv.FormatFloat(
			worker.Holidays, 'f', 1, 64),
		NumericCellStyle())

	// Writing gift row in the 11th row of the table
	writeCell(
		startRow+11, 1,
		"מתנה",
		HeaderCellStyle())

	writeCell(
		startRow+11, 2, 
		strconv.FormatFloat(
			worker.HolidayPresent, 'f', 1, 64),
		NumericCellStyle())

	// Writing sick days row in the 12th row of the table
    writeCell(
		startRow + 12, 1, 
		"ימי מחלה", 
		HeaderCellStyle())

    writeCell(
		startRow + 12, 2, 
		strconv.FormatFloat(
			worker.SickDays, 'f', 1, 64),
		HeaderCellStyle())


    // Writing vacation days row in the 13th row of the table
    writeCell(
		startRow + 13, 1, 
		"ימי חופש", 
		HeaderCellStyle())

    writeCell(
		startRow + 13, 2, 
		strconv.FormatFloat(
			worker.VacDays, 'f', 1, 64),
		NumericCellStyle())

    // Write absense hours row in the 14th row of the table
    writeCell(
		startRow + 14, 1, 
		"שעות להוריד", 
		HeaderCellStyle())

    writeCell(
		startRow + 14, 2, 
		strconv.FormatFloat(
			worker.AbsenseHours, 'f', 1, 64),
		NumericCellStyle())

	// Adding some space between tablesCRturn startRow + 18
	return startRow + 18
}

func TimeStrToFloat64(timeStr string) float64 {
	var hours, minutes, seconds int
	var err error

	parts := strings.Split(timeStr, ":")

	if len(parts) > 0 {
		hours, err = strconv.Atoi(parts[0])
		if err != nil {
			panic(fmt.Errorf("failed to convert hours to int, error: %w", err))
		}
	} else {
		panic(fmt.Errorf("time string is empty"))
	}

	if len(parts) > 1 {
		minutes, err = strconv.Atoi(parts[1])
		if err != nil {
			panic(fmt.Errorf("failed to convert minutes to int, error: %w", err))
		}
	} else {
		minutes = 0
	}

	if len(parts) > 2 {
		seconds, err = strconv.Atoi(parts[2])
		if err != nil {
			panic(fmt.Errorf("failed to convert seconds to int, error: %w", err))
		}
	} else {
		seconds = 0
	}

	return float64(hours) + float64(minutes)/60 + float64(seconds)/3600
}

func StrToFloat64(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(fmt.Errorf("failed to convert string %s to float64, error: %w", str, err))
	}
	return val
}

// Convert row and column to Excel cell name
func cellName(row, col int) string {
	colName, err := excelize.ColumnNumberToName(col)
	if err != nil {
		panic(fmt.Errorf("failed to convert column number to name, error: %w", err))
	}
	return colName + strconv.Itoa(row)
}
