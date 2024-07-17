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

var RAW_SUFFIXES = []string{"360", "3fr", "3g2", "3gp", "3gp2", "3gpp", "7z", "a", "aa", "aae", "aax", "acfm", "acr", "afm", "ai", "aif", "aifc", "aiff", "ait", "amfm", "ape", "apng", "arq", "arw", "asf", "avi", "avif", "azw", "azw3", "bmp", "bpg", "btf", "chm", "ciff", "cos", "cr2", "cr3", "crm", "crw", "cs1", "csv", "cur", "czi", "dc3", "dcm", "dcp", "dcr", "dfont", "dib", "dic", "dicm", "divx", "djv", "djvu", "dll", "dng", "doc", "docm", "docx", "dot", "dotm", "dotx", "dpx", "dr4", "ds2", "dss", "dv", "dvb", "dvr-ms", "dylib", "eip", "eps", "eps2", "eps3", "epsf", "epub", "erf", "exe", "exif", "exr", "exv", "f4a", "f4b", "f4p", "f4v", "fff", "fit", "fits", "fla", "flac", "flif", "flir", "flv", "fpf", "fpx", "gif", "glv", "gpr", "gz", "gzip", "hdp", "hdr", "heic", "heif", "hif", "htm", "html", "ical", "icc", "icm", "ico", "ics", "idml", "iiq", "ind", "indd", "indt", "insp", "insv", "inx", "iso", "itc", "j2c", "j2k", "jng", "jp2", "jpc", "jpe", "jpf", "jpm", "jps", "jpx", "json", "jxl", "jxr", "k25", "kdc", "key", "kth", "la", "lfp", "lfr", "lif", "lnk", "lrv", "m2t", "m2ts", "m2v", "m4a", "m4b", "m4p", "m4v", "macos", "max", "mef", "mie", "mif", "miff", "mka", "mks", "mkv", "mng", "mobi", "modd", "moi", "mos", "mov", "mp3", "mp4", "mpc", "mpeg", "mpg", "mpo", "mqv", "mrc", "mrw", "mts", "mxf", "nef", "newer", "nksc", "nmbtemplate", "nrw", "numbers", "o", "odb", "odc", "odf", "odg", "odi", "odp", "ods", "odt", "ofr", "ogg", "ogv", "onp", "opus", "orf", "ori", "otf", "pac", "pages", "pbm", "pcd", "pct", "pcx", "pdb", "pdf", "pef", "pfa", "pfb", "pfm", "pgf", "pgm", "pict", "plist", "pmp", "pot", "potm", "potx", "ppam", "ppax", "ppm", "pps", "ppsm", "ppsx", "ppt", "pptm", "pptx", "prc", "ps", "ps2", "ps3", "psb", "psd", "psdt", "psp", "pspframe", "pspimage", "pspshape", "psptube", "qif", "qt", "qti", "qtif", "r3d", "ra", "raf", "ram", "rar", "raw", "rif", "riff", "rm", "rmvb", "rpm", "rsrc", "rtf", "rv", "rw2", "rwl", "rwz", "seq", "sketch", "so", "sr2", "srf", "srw", "svg", "swf", "thm", "thmx", "tif", "tiff", "torrent", "ts", "ttc", "ttf", "tub", "txt", "vcard", "vcf", "vnt", "vob", "vrd", "vsd", "wav", "wdp", "webm", "webp", "wma", "wmv", "woff", "woff2", "wpg", "wtv", "wv", "x3f", "xcf", "xhtml", "xla", "xlam", "xls", "xlsb", "xlsm", "xlsx", "xlt", "xltm", "xltx", "xmp", "zip"}

var DISPLAY_SUFFIXES = []string{"jpg", "jpeg", "png"}

const DEFAULT_PATH = "~/Pictures"

var (
	et                *exiftool.Exiftool
	imageDirectory    string
	imageRawToDispMap map[string]string
	currRAWImagePath  string
	application       *gtk.Application
	win               *gtk.ApplicationWindow
	box               *gtk.Box
	listView          *ListView
	infoViewFrame     *gtk.Frame
	infoView          *InfoView
	imageView         *ImageView
	header            *gtk.HeaderBar
	err               error
)

func main() {
	const appID = "com.github.jjj120.raw-rating"
	application, err = gtk.ApplicationNew(appID, glib.APPLICATION_FLAGS_NONE)
	check_error("Could not create application:", err)

	err = godotenv.Load()
	check_error("Could not load dotenv", err)

	imageDirectory = os.Getenv("START_FOLDER")
	if imageDirectory == "" {
		imageDirectory = DEFAULT_PATH
	}

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

	win.SetTitle("RAW Rating")

	infoViewFrame, infoView = setupInfoView()
	refreshInfoView()
	listView = setupListView()
	refreshListView()
	imageView = setupImageView()
	refreshImageView()

	listView.listBox.Connect("row-selected", func(box *gtk.ListBox, row *gtk.ListBoxRow) {
		if row != nil {
			// Call the function to get the label text from the selected row
			text := getLabelTextFromRow(row)
			currRAWImagePath = imageDirectory + "/" + text

			refreshImageChanged()
		} else {
			// clear all data of the application
			// println("resetup application because of none selection")
			// setupWindow(application)
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
	vBox.PackStart(listView.scrolledWindow, true, true, 10)

	// Add the vertical box and the right label to the horizontal box
	box.PackStart(vBox, false, false, 10)
	box.PackStart(imageView.box, true, true, 10)

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
	header.SetTitle("RAW Rating")
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
	case gdk.KEY_Left:
		prevImage()
	case gdk.KEY_Right:
		nextImage()
	case gdk.KEY_Up:
		prevImage()
	case gdk.KEY_Down:
		nextImage()
	default:
		fmt.Println("Key pressed:", keyVal)
	}
	win.ShowAll()
}

func nextImage() {
	// change to the next image
	incrementImageBy(1)
}

func prevImage() {
	// change to the previous image
	incrementImageBy(-1)
}

func incrementImageBy(increment int) {
	selectedRow := listView.listBox.GetSelectedRow()
	if selectedRow != nil {
		selectedIndex := selectedRow.GetIndex()
		nextIndex := selectedIndex + increment
		if nextIndex >= 0 && nextIndex < int(listView.listBox.GetChildren().Length()) {
			listView.listBox.SelectRow(listView.listBox.GetRowAtIndex(nextIndex))
		}
	} else if listView.listBox.GetChildren().Length() > 0 {
		// if there is no selected row, select the first row
		firstRow := listView.listBox.GetRowAtIndex(0)
		listView.listBox.SelectRow(firstRow)
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
