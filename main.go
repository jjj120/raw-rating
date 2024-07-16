package main

import (
	"fmt"
	"log"
	"os"

	"github.com/barasher/go-exiftool"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/joho/godotenv"
)

const BOX_SPACING int = 2

const LABEL_MARGIN = 10

var RAW_SUFFIXES = []string{"cr3", "cr2", "dng"}
var DISPLAY_SUFFIXES = []string{"jpg", "jpeg", "png"}

var (
	et                *exiftool.Exiftool
	imageDirectory    string
	imageRawToDispMap map[string]string
	currRAWImagePath  string
	application       *gtk.Application
	win               *gtk.ApplicationWindow
	box               *gtk.Box
	listView          *gtk.ListBox
	infoViewFrame     *gtk.Frame
	infoView          *InfoView
	imageView         *ImageView
	header            *gtk.HeaderBar
	err               error
)

func main() {
	const appID = "com.github.jjj120.cr3-rating"
	application, err = gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	check_error("Could not create application:", err)
	// currImagePath = "./test_images/IMG_2055.CR3"
	err = godotenv.Load()
	check_error("Could not load dotenv", err)

	imageDirectory = os.Getenv("START_FOLDER")

	currRAWImagePath = "test_images/IMG_0371.CR3"
	imageDirectory = "test_images"

	et, err = exiftool.NewExiftool()
	if err != nil {
		fmt.Printf("Error when intializing: %v\n", err)
		return
	}
	defer et.Close()

	application.Connect("activate", func() {
		win := setupWindow(application)

		aNew := glib.SimpleActionNew("new", nil)
		aNew.Connect("activate", func() {
			setupWindow(application).ShowAll()
		})
		application.AddAction(aNew)

		aQuit := glib.SimpleActionNew("quit", nil)
		aQuit.Connect("activate", func() {
			application.Quit()
		})
		application.AddAction(aQuit)

		win.ShowAll()
	})

	os.Exit(application.Run(os.Args))
}

func check_error(message string, err error) {
	if err != nil {
		log.Fatal(message, " - ", err)
	}
}

func setupWindow(application *gtk.Application) *gtk.ApplicationWindow {
	win, err = gtk.ApplicationWindowNew(application)
	check_error("Unable to create window", err)

	win.SetTitle("CR3 Rating")

	infoViewFrame, infoView = setupInfoView()
	refreshInfoView()
	listView = setupListView()
	refreshListView()
	imageView = setupImageView()
	refreshImageView()

	listView.Connect("row-selected", func(box *gtk.ListBox, row *gtk.ListBoxRow) {
		if row != nil {
			// Call the function to get the label text from the selected row
			text := getLabelTextFromRow(row)
			currRAWImagePath = imageDirectory + "/" + text

			refreshImageChanged()
		} else {
			// clear all data of the application
			println("resetup application because of none selection")
			setupWindow(application)
		}
	})

	// Create a horizontal box to hold the left and right elements
	box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	check_error("Unable to create horizontal box:", err)

	// Create a vertical box for the left elements
	vBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10)
	check_error("Unable to create vertical box:", err)

	// Add the top-left element
	vBox.PackStart(infoViewFrame, false, false, 10)

	// Add the bottom-left element
	vBox.PackStart(listView, true, true, 10)

	// Create the right element
	imageView.scrolledWindow.SetHExpand(true)

	// Add the vertical box and the right label to the horizontal box
	box.PackStart(vBox, false, false, 10)
	box.PackStart(imageView.scrolledWindow, true, true, 10)

	// Connect to key presses
	win.Connect("key-press-event", keyPress)

	// Add the horizontal box to the window
	win.Add(box)

	setupHeaderBar()
	win.SetTitlebar(header)
	win.SetPosition(gtk.WIN_POS_MOUSE)
	win.SetDefaultSize(1200, 800)

	return win
}

func setupHeaderBar() {
	// Create a header bar
	header, err = gtk.HeaderBarNew()
	check_error("Could not create header bar", err)

	header.SetShowCloseButton(true)
	header.SetTitle("CR3 Rating")
	// header.SetSubtitle("Actions Example")

	// Create a new menu button
	mbtn, err := gtk.MenuButtonNew()
	check_error("Could not create menu button", err)

	// Set up the menu model for the button
	menu := glib.MenuNew()
	if menu == nil {
		log.Fatal("Could not create menu (nil)")
	}

	// Actions with the prefix 'app' reference actions on the application
	// Actions with the prefix 'win' reference actions on the current window (specific to ApplicationWindow)
	// Other prefixes can be added to widgets via InsertActionGroup
	menu.Append("New Window", "app.new")
	menu.Append("Open Directory", "custom.open_dir")
	menu.Append("Close Window", "win.close")
	menu.Append("Quit", "app.quit")

	// Create the action "win.close"
	aClose := glib.SimpleActionNew("close", nil)
	aClose.Connect("activate", func() {
		win.Close()
	})
	win.AddAction(aClose)

	// Create and insert custom action group with prefix "custom"
	customActionGroup := glib.SimpleActionGroupNew()
	win.InsertActionGroup("custom", customActionGroup)

	// Create an action in the custom action group
	aOpenDir := glib.SimpleActionNew("open_dir", nil)
	aOpenDir.Connect("activate", func() {
		fileChooser, err := gtk.FileChooserDialogNewWith1Button(
			"Choose the image directory",
			win,
			gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
			"Select",
			gtk.RESPONSE_ACCEPT,
		)

		if err != nil {
			log.Fatal("Unable to create file chooser dialog:", err)
		}

		// Run the file chooser dialog and check for the response
		response := fileChooser.Run()

		if response == gtk.RESPONSE_ACCEPT {
			imageDirectory = fileChooser.GetFilename()
			fmt.Println("Selected image dir:", imageDirectory)
			refreshAll()
		}

		// Destroy the file chooser dialog after use
		fileChooser.Destroy()
	})
	customActionGroup.AddAction(aOpenDir)
	win.AddAction(aOpenDir)

	mbtn.SetMenuModel(&menu.MenuModel)

	// add the menu button to the header
	header.PackStart(mbtn)
}

// filesSelected: callback function for "file-set" signal
func refreshAll() {
	refreshImageView()
	refreshInfoView()
	refreshListView()
	win.ShowAll()
}

func refreshImageChanged() {
	//refresh views (but not file list)

	refreshImageView()
	refreshInfoView()

	win.ShowAll()
}

func keyPress(win *gtk.ApplicationWindow, event *gdk.Event) {
	// button press event handler

	keyEvent := gdk.EventKeyNewFromEvent(event)
	keyVal := keyEvent.KeyVal()

	switch keyVal {
	case gdk.KEY_Escape:
		win.Close()

	case gdk.KEY_q:
		win.Close()

	case gdk.KEY_0:
		changeRating(0)

	case gdk.KEY_1:
		changeRating(1)

	case gdk.KEY_2:
		changeRating(2)

	case gdk.KEY_3:
		changeRating(3)

	case gdk.KEY_4:
		changeRating(4)

	case gdk.KEY_5:
		changeRating(5)

	}

}

func changeRating(rating int) {
	// change rating of the current image
	// fmt.Println("Change rating of ", currRAWImagePath, "to", rating)

	rawArgs := et.ExtractMetadata(currRAWImagePath)
	rawArgs[0].SetString("Rating", fmt.Sprintf("%d", rating))
	et.WriteMetadata(rawArgs)

	dispArgs := et.ExtractMetadata(imageRawToDispMap[currRAWImagePath])
	dispArgs[0].SetString("Rating", fmt.Sprintf("%d", rating))
	et.WriteMetadata(dispArgs)

	refreshInfoView()

}
