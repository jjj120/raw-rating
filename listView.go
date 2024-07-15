package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

func setupListView() *gtk.ListBox {
	listBox, err := gtk.ListBoxNew()
	check_error("Unable to create list box", err)

	if imageDirectory == "" {
		return listBox
	}

	listBox.SetActivateOnSingleClick(true)

	return listBox
}

func addListEntry(rawFilename string, dispFilename string, listBox *gtk.ListBox) {
	rawFilenameStripped := strings.Split(rawFilename, "/")[len(strings.Split(rawFilename, "/"))-1]
	label, err := gtk.LabelNew(rawFilenameStripped)
	label.SetXAlign(0)
	label.SetMarginStart(LABEL_MARGIN)
	check_error("Unable to create label for listbox row", err)

	row, err := gtk.ListBoxRowNew()
	check_error("Unable to create listbox row", err)

	row.Add(label)
	listBox.Add(row)
}

func getImageMapping(dirName string) map[string]string {
	rawFileList := []string{}
	displayFileList := []string{}

	// Open the directory
	d, err := os.Open(dirName)
	if err != nil {
		fmt.Printf("Error opening directory: %v\n", err)
		panic(1)
	}
	defer d.Close()

	// Read the directory entries
	files, err := d.Readdir(-1)
	if err != nil {
		fmt.Printf("Error reading directory entries: %v\n", err)
		panic(1)
	}

	// Iterate over the directory entries
	for _, file := range files {
		// Check if it is a file (not a directory)
		if file.IsDir() {
			continue
		}

		// Check if the filename ends with the specified suffix
		for _, suffix := range RAW_SUFFIXES {
			if strings.HasSuffix(strings.ToLower(file.Name()), strings.ToLower(suffix)) {
				// Print or process the file
				rawFileList = append(rawFileList, (filepath.Join(dirName, file.Name())))
			}
		}

		for _, suffix := range DISPLAY_SUFFIXES {
			if strings.HasSuffix(strings.ToLower(file.Name()), strings.ToLower(suffix)) {
				// Print or process the file
				displayFileList = append(displayFileList, (filepath.Join(dirName, file.Name())))
			}
		}
	}

	fmt.Println(rawFileList)
	fmt.Println(displayFileList)

	// TODO: make this more performant, this is ugly haha
	fileList := map[string]string{}
	for _, filenameRAW := range rawFileList {
		filenameSplitRAW := strings.Split(filenameRAW, ".")
		filenameStrippedRAW := strings.Join(filenameSplitRAW[:len(filenameSplitRAW)-1], ".")
		for _, filenameDisp := range displayFileList {
			filenameSplitDisp := strings.Split(filenameDisp, ".")
			filenameStrippedDisp := strings.Join(filenameSplitDisp[:len(filenameSplitDisp)-1], ".")
			if filenameStrippedRAW == filenameStrippedDisp {
				fileList[filenameRAW] = filenameDisp
			}
		}
	}

	fmt.Println(fileList)
	return fileList
}

func refreshListView() {
	clearListBox(listView)
	fmt.Println("refresh list view")

	imageRawToDispMap = getImageMapping(imageDirectory)
	fmt.Println("refresh list view")

	fmt.Println(imageRawToDispMap)

	for rawFilename, dispFilename := range imageRawToDispMap {
		addListEntry(rawFilename, dispFilename, listView)
		fmt.Println("Added ", rawFilename, "to list with display filename", dispFilename)
	}

	listView.SetSortFunc(func(row1, row2 *gtk.ListBoxRow) int {
		return strings.Compare(getLabelTextFromRow(row1), getLabelTextFromRow(row2))
	})

}

func clearListBox(listBox *gtk.ListBox) {
	listView = setupListView()
}

func getLabelTextFromRow(row *gtk.ListBoxRow) string {
	child, err := row.GetChild()
	check_error("unable to get child of row", err)

	label, ok := child.(*gtk.Label)
	if !ok {
		log.Println("The child widget is not a label")
		return ""
	}

	text, err := label.GetText()
	if err != nil {
		log.Fatalf("Unable to get text from label: %v", err)
	}

	return text
}
