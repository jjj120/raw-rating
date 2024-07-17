package main

import (
	"fmt"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type InfoView struct {
	filenameLabel     *gtk.Label
	camTypeLabel      *gtk.Label
	lensTypeLabel     *gtk.Label
	dateTimeLabel     *gtk.Label
	focalLengthLabel  *gtk.Label
	exposureTimeLabel *gtk.Label
	apertureLabel     *gtk.Label
	isoLabel          *gtk.Label
	ratingLabel       *gtk.Label
	filenameVar       string
	camTypeVar        string
	lensTypeVar       string
	dateTimeVar       string
	focalLengthVar    string
	exposureTimeVar   string
	apertureVar       string
	isoVar            string
	ratingVar         string
	grid              *gtk.Grid
}

func (infoView *InfoView) createWidgets() {
	infoView.grid, _ = gtk.GridNew()
	infoView.grid.SetBorderWidth(0)
	infoView.grid.SetRowSpacing(2)
	infoView.grid.SetColumnSpacing(10)

	lineHeight := 1
	const SPACE_TO_BORDER = 0

	// filename
	infoView.filenameLabel, err = gtk.LabelNew(infoView.filenameVar)
	infoView.filenameLabel.SetProperty("height-request", lineHeight)
	// infoView.filenameLabel.SetXAlign(SPACE_TO_BORDER)

	infoView.gridAttach(infoView.filenameLabel, 1, 1, 2, 1)

	// camera type
	infoView.camTypeLabel, err = gtk.LabelNew(infoView.camTypeVar)
	infoView.camTypeLabel.SetProperty("height-request", lineHeight)
	infoView.camTypeLabel.SetXAlign(SPACE_TO_BORDER)

	infoView.gridAttach(infoView.camTypeLabel, 1, 2, 2, 1)

	// lens type
	infoView.lensTypeLabel, err = gtk.LabelNew(infoView.lensTypeVar)
	infoView.lensTypeLabel.SetProperty("height-request", lineHeight)
	infoView.lensTypeLabel.SetXAlign(SPACE_TO_BORDER)

	infoView.gridAttach(infoView.lensTypeLabel, 1, 3, 2, 1)

	// date and time
	infoView.dateTimeLabel, err = gtk.LabelNew(infoView.dateTimeVar)
	infoView.dateTimeLabel.SetProperty("height-request", lineHeight)
	infoView.dateTimeLabel.SetXAlign(SPACE_TO_BORDER)

	infoView.gridAttach(infoView.dateTimeLabel, 1, 4, 2, 1)

	// focal length
	focalLengthLabel, err := gtk.LabelNew("Focal length:")
	check_error("Unable to create Label for focal length", err)
	focalLengthLabel.SetProperty("height-request", lineHeight)
	focalLengthLabel.SetXAlign(SPACE_TO_BORDER)

	infoView.gridAttach(focalLengthLabel, 1, 5, 1, 1)

	infoView.focalLengthLabel, err = gtk.LabelNew(infoView.focalLengthVar)
	infoView.focalLengthLabel.SetProperty("height-request", lineHeight)
	infoView.focalLengthLabel.SetXAlign(1 - SPACE_TO_BORDER)

	infoView.gridAttach(infoView.focalLengthLabel, 2, 5, 1, 1)

	// exposure time
	exposureTimeLabel, err := gtk.LabelNew("Exposure time:")
	check_error("Unable to create Label for Exposure time", err)
	exposureTimeLabel.SetProperty("height-request", lineHeight)
	exposureTimeLabel.SetXAlign(SPACE_TO_BORDER)

	infoView.gridAttach(exposureTimeLabel, 1, 6, 1, 1)

	infoView.exposureTimeLabel, err = gtk.LabelNew(infoView.exposureTimeVar)
	infoView.exposureTimeLabel.SetProperty("height-request", lineHeight)
	infoView.exposureTimeLabel.SetXAlign(1 - SPACE_TO_BORDER)

	infoView.gridAttach(infoView.exposureTimeLabel, 2, 6, 1, 1)

	// aperture value
	apertureLabel, err := gtk.LabelNew("Aperture:")
	check_error("Unable to create Label for Aperture", err)
	apertureLabel.SetProperty("height-request", lineHeight)
	apertureLabel.SetXAlign(SPACE_TO_BORDER)

	infoView.gridAttach(apertureLabel, 1, 7, 1, 1)

	infoView.apertureLabel, err = gtk.LabelNew(infoView.apertureVar)
	infoView.apertureLabel.SetProperty("height-request", lineHeight)
	infoView.apertureLabel.SetXAlign(1 - SPACE_TO_BORDER)

	infoView.gridAttach(infoView.apertureLabel, 2, 7, 1, 1)

	// ISO
	isoLabel, err := gtk.LabelNew("ISO:")
	check_error("Unable to create Label for ISO", err)
	isoLabel.SetProperty("height-request", lineHeight)
	isoLabel.SetXAlign(SPACE_TO_BORDER)

	infoView.gridAttach(isoLabel, 1, 8, 1, 1)

	infoView.isoLabel, err = gtk.LabelNew(infoView.isoVar)
	infoView.isoLabel.SetProperty("height-request", lineHeight)
	infoView.isoLabel.SetXAlign(1 - SPACE_TO_BORDER)

	infoView.gridAttach(infoView.isoLabel, 2, 8, 1, 1)

	// rating
	infoView.ratingLabel, err = gtk.LabelNew(infoView.ratingVar)
	infoView.ratingLabel.SetProperty("height-request", lineHeight)
	infoView.ratingLabel.SetName("ratingLabel")
	infoView.gridAttach(infoView.ratingLabel, 1, 9, 2, 1)

	// formatting
	// infoView.filenameLabel
	// infoView.camTypeLabel
	// infoView.lensTypeLabel
	// infoView.dateTimeLabel
	// focalLengthLabel
	// infoView.focalLengthLabel
	// exposureTimeLabel
	// infoView.exposureTimeLabel
	// apertureLabel
	// infoView.apertureLabel
	// isoLabel
	// infoView.isoLabel
	// infoView.ratingLabel

	infoView.filenameLabel.SetMarginTop(LABEL_MARGIN)
	infoView.ratingLabel.SetMarginBottom(LABEL_MARGIN)

	infoView.filenameLabel.SetMarginStart(LABEL_MARGIN)
	infoView.camTypeLabel.SetMarginStart(LABEL_MARGIN)
	infoView.lensTypeLabel.SetMarginStart(LABEL_MARGIN)
	infoView.dateTimeLabel.SetMarginStart(LABEL_MARGIN)
	focalLengthLabel.SetMarginStart(LABEL_MARGIN)
	infoView.focalLengthLabel.SetMarginStart(LABEL_MARGIN)
	exposureTimeLabel.SetMarginStart(LABEL_MARGIN)
	infoView.exposureTimeLabel.SetMarginStart(LABEL_MARGIN)
	apertureLabel.SetMarginStart(LABEL_MARGIN)
	infoView.apertureLabel.SetMarginStart(LABEL_MARGIN)
	isoLabel.SetMarginStart(LABEL_MARGIN)
	infoView.isoLabel.SetMarginStart(LABEL_MARGIN)
	infoView.ratingLabel.SetMarginStart(LABEL_MARGIN)

	infoView.filenameLabel.SetMarginEnd(LABEL_MARGIN)
	infoView.camTypeLabel.SetMarginEnd(LABEL_MARGIN)
	infoView.lensTypeLabel.SetMarginEnd(LABEL_MARGIN)
	infoView.dateTimeLabel.SetMarginEnd(LABEL_MARGIN)
	focalLengthLabel.SetMarginEnd(LABEL_MARGIN)
	infoView.focalLengthLabel.SetMarginEnd(LABEL_MARGIN)
	exposureTimeLabel.SetMarginEnd(LABEL_MARGIN)
	infoView.exposureTimeLabel.SetMarginEnd(LABEL_MARGIN)
	apertureLabel.SetMarginEnd(LABEL_MARGIN)
	infoView.apertureLabel.SetMarginEnd(LABEL_MARGIN)
	isoLabel.SetMarginEnd(LABEL_MARGIN)
	infoView.isoLabel.SetMarginEnd(LABEL_MARGIN)
	infoView.ratingLabel.SetMarginEnd(LABEL_MARGIN)
}

func (infoView *InfoView) gridAttach(widget gtk.IWidget, left, top, width, height int) {
	infoView.grid.Attach(widget, left, top, width, height)
}

func (infoView *InfoView) resetVars(currImagePath string) {
	if currImagePath != "" {
		infoView.filenameVar = stripFilepath(currImagePath)
	} else {
		infoView.filenameVar = ""
	}
	infoView.camTypeVar = ""
	infoView.lensTypeVar = ""
	infoView.dateTimeVar = ""
	infoView.focalLengthVar = ""
	infoView.exposureTimeVar = ""
	infoView.apertureVar = ""
	infoView.isoVar = ""
	infoView.ratingVar = ""
}

func (infoView *InfoView) updateExifData(currImagePath string) {
	// currImagePath := infoView.images[infoView.currIndex][1]

	fileInfos := et.ExtractMetadata(currImagePath)
	if len(fileInfos) < 1 {
		infoView.resetVars(currImagePath)
		return
	}
	exifInfo := fileInfos[0].Fields

	if currImagePath != "" {
		infoView.filenameVar = stripFilepath(currImagePath)
	} else {
		infoView.filenameVar = ""
	}

	if val, ok := exifInfo["Model"].(string); ok {
		infoView.camTypeVar = val
	} else {
		infoView.camTypeVar = "Unknown"
	}

	if val, ok := exifInfo["LensModel"].(string); ok {
		infoView.lensTypeVar = val
	} else {
		infoView.lensTypeVar = "Unknown"
	}

	if val, ok := exifInfo["DateTimeOriginal"].(string); ok {
		infoView.dateTimeVar = formatDateTime(val)
	} else {
		infoView.dateTimeVar = "Unknown"
	}

	if val, ok := exifInfo["FocalLength"].(string); ok {
		if val[len(val)-1] == 'm' {
			infoView.focalLengthVar = val
		} else {
			infoView.focalLengthVar = fmt.Sprintf("%s mm", val)
		}
	} else {
		infoView.focalLengthVar = "Unknown"
	}

	if val, ok := exifInfo["ExposureTime"].(string); ok {
		infoView.exposureTimeVar = val
	} else {
		infoView.exposureTimeVar = "Unknown"
	}

	if val, ok := exifInfo["FNumber"].(float64); ok {
		infoView.apertureVar = fmt.Sprintf("F %.1f", val)
	} else {
		infoView.apertureVar = "Unknown"
	}

	if val, ok := exifInfo["ISO"].(float64); ok {
		infoView.isoVar = fmt.Sprintf("%.0f", val)
	} else {
		infoView.isoVar = "Unknown"
	}

	if val, ok := exifInfo["Rating"].(float64); ok {
		ratingStr := "                              " // 30 spaces for constant width
		for i := 1.; i <= 5; i++ {
			if i <= val {
				ratingStr += "★"
			} else {
				ratingStr += "☆"
			}
		}
		infoView.ratingVar = ratingStr + "                              " // 30 spaces for constant width
	} else {
		infoView.ratingVar = "Rating Unknown"
	}
}

func (infoView *InfoView) refreshLabels() {
	infoView.filenameLabel.SetText(infoView.filenameVar)
	infoView.camTypeLabel.SetText(infoView.camTypeVar)
	infoView.lensTypeLabel.SetText(infoView.lensTypeVar)
	infoView.dateTimeLabel.SetText(infoView.dateTimeVar)
	infoView.focalLengthLabel.SetText(infoView.focalLengthVar)
	infoView.exposureTimeLabel.SetText(infoView.exposureTimeVar)
	infoView.apertureLabel.SetText(infoView.apertureVar)
	infoView.isoLabel.SetText(infoView.isoVar)
	infoView.ratingLabel.SetText(infoView.ratingVar)
}

func setupInfoView() (*gtk.Frame, *InfoView) {
	frame, err := gtk.FrameNew("")
	frame.SetBorderWidth(0)
	check_error("Unable to create frame for info view", err)
	infoView := &InfoView{}

	infoView.createWidgets()
	infoView.updateExifData(currRAWImagePath)
	infoView.refreshLabels()

	frame.Add(infoView.grid)
	return frame, infoView
}

func refreshInfoView() {
	infoView.updateExifData(currRAWImagePath)
	infoView.refreshLabels()
}

func formatDateTime(dateTime string) string {
	return strings.Replace(dateTime, ":", ".", 2)
}

func stripFilepath(filepath string) string {
	return strings.Split(filepath, "/")[len(strings.Split(filepath, "/"))-1]
}
