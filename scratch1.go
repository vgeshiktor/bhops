package main

import (
	"fmt"

	"github.com/xuri/excelize/v2"
)

func main() {
	f := excelize.NewFile()

	// Define a style with specific font and alignment
	style, err := f.NewStyle(
		&excelize.Style{
			Font: &excelize.Font{
				Bold:  true,
				Italic : true,
				Size:  12,
				// Color: "#FF0000",
			},
			Alignment: &excelize.Alignment{
				Horizontal: "center",
				Vertical:  "center",
			},
			// Fill: excelize.Fill{
			// 	Type:    "pattern",
			// 	Color:   []string{"#FFFF00"},
			// 	Pattern: 1,
			// },
		},
	)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Set value and apply the style to a specific cell
	cell := "A1"
	f.SetCellValue("Sheet1", cell, "Formatted Text")
	if err := f.SetCellStyle("Sheet1", cell, cell, style); err != nil {
		fmt.Println(err)
		return
	}

	// Save the file
	if err := f.SaveAs("formatted.xlsx"); err != nil {
		fmt.Println(err)
	}
}
