package main

import (
	"fmt"

	"github.com/vgeshiktor/bhops/internal/attendanceops"
)

const (
    MONTHLY_ATTENDANCE_REPORT_XLS_PATH = "/Users/vadimgeshiktor/repos/github.com/vgeshiktor/bhops/internal/attendanceops/input/02-2025.xlsx"
    WORKER_DETAILS_JSON_PATH = "/Users/vadimgeshiktor/repos/github.com/vgeshiktor/bhops/internal/attendanceops/config/id2worker.json"
    
    WORKER_HOURS_JSON_PATH = "/Users/vadimgeshiktor/repos/github.com/vgeshiktor/bhops/internal/attendanceops/input/workershours.json"
    SALARY_DETAILS_OUTPUT_PATH = "/Users/vadimgeshiktor/repos/github.com/vgeshiktor/bhops/internal/attendanceops/output/salary_details.xlsx"
)

func main() {
    // create workers attendance report sheet
    AttendanceReport, err := attendanceops.CreateAttendanceReport(
        MONTHLY_ATTENDANCE_REPORT_XLS_PATH,
        WORKER_HOURS_JSON_PATH,
        WORKER_DETAILS_JSON_PATH,
    )
    if err != nil {
        fmt.Println("Failed to create attendance report: ", err)
        return
    }

    // save workers attendance report sheet
    err = attendanceops.SaveAttendanceReport(
        AttendanceReport,    SALARY_DETAILS_OUTPUT_PATH)
    if err != nil {
        fmt.Printf(
            "Failed to save attendance report: %s, error: %v", SALARY_DETAILS_OUTPUT_PATH, err)
        }
}
