package attendanceops

import (
	"github.com/xuri/excelize/v2"
)

var (
	HEADER_FONT  = excelize.Font{Bold: true, Size: 14}
	TEXT_FONT    = excelize.Font{Size: 14}
	ALIGN_CENTER = excelize.Alignment{Horizontal: "center", Vertical: "center"}
	ALIGN_LEFT   = excelize.Alignment{Horizontal: "left", Vertical: "center"}
	ALIGN_RIGHT  = excelize.Alignment{Horizontal: "right", Vertical: "center"}

	LEFT_BORDER = excelize.Border{Type: "left", Style: 1, Color: "000000"}
	TOP_BORDER = excelize.Border{Type: "top", Style: 1, Color: "000000"}
	BOTTOM_BORDER = excelize.Border{Type: "bottom", Style: 1, Color: "000000"}
	RIGHT_BORDER = excelize.Border{Type: "right", Style: 1, Color: "000000"}
	THICK_BORDER = []excelize.Border{
		LEFT_BORDER,
		TOP_BORDER,
		BOTTOM_BORDER,
		RIGHT_BORDER,
	}
	NUMBER_FORMAT = "#"
)

func TitleCellStyle() *excelize.Style {
	style := excelize.Style{
		Font:      &HEADER_FONT,
		Alignment: &ALIGN_RIGHT,
	}

	return &style
}


func DefaultCellStyle() *excelize.Style {
	style := excelize.Style{
		Font:      &TEXT_FONT,
		Alignment: &ALIGN_RIGHT,
		Border:   THICK_BORDER,
	}

	return &style
}

func HeaderCellStyle() *excelize.Style {
	style := excelize.Style{
		Font:      &HEADER_FONT,
		Alignment: &ALIGN_RIGHT,
		Border:   THICK_BORDER,
	}

	return &style
}

func NumericCellStyle() *excelize.Style {
	style := excelize.Style{
		Font:      &TEXT_FONT,
		Alignment: &ALIGN_RIGHT,
		CustomNumFmt: &NUMBER_FORMAT,
		Border:   THICK_BORDER,
	}

	return &style
}

