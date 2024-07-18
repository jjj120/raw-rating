package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/gotk3/gotk3/gtk"
)

type ListView struct {
	listBox        *gtk.ListBox
	scrolledWindow *gtk.ScrolledWindow
}

func NewListView() *ListView {
	listView := &ListView{}

	listView.listBox, err = gtk.ListBoxNew()
	check_error("Unable to create list box", err)

	listView.scrolledWindow, err = gtk.ScrolledWindowNew(nil, nil)
	check_error("Unable to create scrolled window", err)

	listView.scrolledWindow.Add(listView.listBox)

	return listView
}

func (lv *ListView) addListEntry(rawFilename string, dispFilename string) {
	rawFilenameStripped := strings.Split(rawFilename, "/")[len(strings.Split(rawFilename, "/"))-1]

	label, err := gtk.LabelNew(rawFilenameStripped)
	label.SetXAlign(0)
	label.SetMarginStart(LABEL_MARGIN)
	check_error("Unable to create label for listbox row", err)

	// pixbuf, err := gdk.PixbufNewFromFile(imageRawToDispMap[rawFilename])
	// check_error("Unable to create pixbuf from file", err)

	// pixbuf, err = pixbuf.ScaleSimple(100, 100, gdk.INTERP_BILINEAR)
	// check_error("Unable to scale pixbuf", err)

	// pixbuf = rotatePixbuf(imageRawToDispMap[rawFilename], pixbuf)

	// image, err := gtk.ImageNewFromPixbuf(pixbuf)
	// check_error("Unable to create image from pixbuf", err)
	// image.SetMarginStart(LABEL_MARGIN)

	row, err := gtk.ListBoxRowNew()
	check_error("Unable to create listbox row", err)

	// rowBox, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0)
	// check_error("Unable to create box for listbox row", err)

	// rowBox.Add(label)
	// rowBox.Add(image)
	row.Add(label)
	lv.listBox.Add(row)
}

func (lv *ListView) refreshListView() {
	lv.clearListBox()

	imageRawToDispMap = getImageMapping(imageDirectory)

	for rawFilename, dispFilename := range imageRawToDispMap {
		lv.addListEntry(rawFilename, dispFilename)
		log.Info("Added ", rawFilename, " to list with display filename ", dispFilename)

	}

	lv.listBox.SetSortFunc(func(row1, row2 *gtk.ListBoxRow) int {
		return strings.Compare(getLabelTextFromRow(row1), getLabelTextFromRow(row2))
	})
}

func (lv *ListView) clearListBox() {
	for {
		row := lv.listBox.GetRowAtIndex(0)
		if row == nil {
			break
		}
		row.Destroy()
	}
}

func refreshListView() {
	listView.refreshListView()
}

func setupListView() *ListView {
	return NewListView()
}

func getImageMapping(dirName string) map[string]string {
	rawFileList := []string{}
	displayFileList := []string{}

	// Open the directory
	d, err := os.Open(dirName)
	if err != nil {
		check_error("Error opening directory", err)
	}
	defer d.Close()

	// Read the directory entries
	files, err := d.Readdir(-1)
	if err != nil {
		check_error("Error reading directory entries", err)
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

	// fmt.Println(fileList)
	return fileList
}

func getLabelTextFromRow(row *gtk.ListBoxRow) string {
	child, err := row.GetChild()
	check_error("unable to get child of row", err)

	label, ok := child.(*gtk.Label)
	if !ok {
		log.Debug("The child widget is not a label")
		return ""
	}

	text, err := label.GetText()
	if err != nil {
		log.Fatalf("Unable to get text from label: %v", err)
	}

	return text
}
