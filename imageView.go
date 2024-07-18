package main

import (
	"github.com/barasher/go-exiftool"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const SCALE_FACTOR = 1.5

type ImageView struct {
	box, boxLeft, boxRight *gtk.Box
	scrolledWindow         *gtk.ScrolledWindow
	zoomView               *ZoomView
	image                  *gtk.Image
	pixbuf_orig            *gdk.Pixbuf
	imagePath              string
	newImagePath           string
	currHeight, currWidth  int
	xOffset, yOffset       float64
}

func NewImageView() *ImageView {
	iv := &ImageView{}

	iv.box, err = gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 10)
	check_error("Unable to create scrolled window", err)

	iv.zoomView = NewZoomView()

	iv.boxLeft, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	check_error("Unable to create box", err)
	iv.boxLeft.SetHExpand(true)

	iv.boxRight, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	check_error("Unable to create box", err)

	iv.scrolledWindow, err = gtk.ScrolledWindowNew(nil, nil)
	check_error("Unable to create ScrolledWindow", err)
	iv.scrolledWindow.SetHExpand(true)
	iv.scrolledWindow.SetVExpand(true)

	iv.image, err = gtk.ImageNew()
	check_error("Unable to create image", err)
	iv.image.SetCanFocus(false)

	iv.box.SetSizeRequest(1920, 1080)
	iv.scrolledWindow.AddEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK))

	iv.image.Connect("draw", iv.draw)
	iv.scrolledWindow.Connect("button-press-event", iv.onButtonPress)
	iv.scrolledWindow.Connect("motion-notify-event", iv.onDrag)

	iv.scrolledWindow.Add(iv.image)
	iv.boxLeft.Add(iv.scrolledWindow)
	// iv.boxLeft.Add(iv.image)
	iv.boxRight.Add(iv.zoomView.box)

	iv.box.PackStart(iv.boxLeft, true, true, 10)
	iv.box.PackStart(iv.boxRight, true, true, 10)

	return iv
}

func (iv *ImageView) draw(img *gtk.Image) {
	boxHeight := iv.boxLeft.GetAllocatedHeight()
	boxWidth := iv.boxLeft.GetAllocatedWidth()

	if (iv.imagePath != iv.newImagePath) || (iv.currHeight != boxHeight) || (iv.currWidth != boxWidth) {
		// load image
		img.SetFromFile(iv.newImagePath)

		iv.pixbuf_orig = nil
		iv.pixbuf_orig = img.GetPixbuf()

		// Rotate image if needed (based on EXIF data)
		iv.pixbuf_orig = rotatePixbuf(iv.newImagePath, iv.pixbuf_orig)

		iv.zoomView.pixbuf_orig = iv.pixbuf_orig

		// Reset values
		iv.imagePath = iv.newImagePath

		imageRatio := float64(iv.pixbuf_orig.GetWidth()) / float64(iv.pixbuf_orig.GetHeight())

		newHeight := boxHeight
		newWidthFromHeight := int(float64(newHeight) * imageRatio)

		newWidth := boxWidth
		newHeightFromWidth := int(float64(newWidth) / imageRatio)

		if newHeight <= boxHeight && newWidthFromHeight <= boxWidth {
			newWidth = newWidthFromHeight
		} else {
			newHeight = newHeightFromWidth
		}

		pixbuf, err := iv.pixbuf_orig.ScaleSimple(newWidth, newHeight, gdk.INTERP_BILINEAR)
		check_error("Unable to scale pixbuf", err)

		img.SetFromPixbuf(pixbuf)

		iv.xOffset = float64(boxWidth-newWidth) / 2
		iv.yOffset = float64(boxHeight-newHeight) / 2

		iv.zoomView.zoomVal = float64(iv.pixbuf_orig.GetWidth()) / float64(newWidth)
		iv.zoomView.zoomVal = float64(iv.pixbuf_orig.GetHeight()) / float64(newHeight)

		iv.zoomView.QueueDraw()
	}
	iv.currHeight = boxHeight
	iv.currWidth = boxWidth
}

func (iv *ImageView) onButtonPress(sw *gtk.ScrolledWindow, event *gdk.Event) {
	buttonEvent := gdk.EventButtonNewFromEvent(event)

	if buttonEvent.Button() == gdk.BUTTON_PRIMARY {
		log.Debug("Button Pressed")

		x, y := buttonEvent.MotionVal()
		iv.zoomView.valPosX, iv.zoomView.valPosY = (x-iv.xOffset)/(float64(iv.boxLeft.GetAllocatedWidth())-2*iv.xOffset), (y-iv.yOffset)/(float64(iv.boxLeft.GetAllocatedHeight())-2*iv.yOffset)
		iv.zoomView.QueueDraw()
	}
}

func (iv *ImageView) onDrag(sw *gtk.ScrolledWindow, event *gdk.Event) {
	dragEvent := gdk.EventMotionNewFromEvent(event)
	if dragEvent.State()&gdk.BUTTON1_MASK != 0 {
		log.Debug("Dragging")

		x, y := dragEvent.MotionVal()
		iv.zoomView.valPosX, iv.zoomView.valPosY = (x-iv.xOffset)/(float64(iv.boxLeft.GetAllocatedWidth())-2*iv.xOffset), (y-iv.yOffset)/(float64(iv.boxLeft.GetAllocatedHeight())-2*iv.yOffset)
		iv.zoomView.QueueDraw()
	}
}

func rotatePixbuf(filepath string, pixbuf *gdk.Pixbuf) *gdk.Pixbuf {
	// Implement rotation logic here if needed
	metadata := et.ExtractMetadata(filepath)

	rot, err := metadata[0].GetString("Orientation")
	if err == exiftool.ErrKeyNotFound {
		log.Debug("No orientation found. Defaulting to Horizontal (normal)")
		rot = "Horizontal (normal)"
		err = nil
	}

	// 1 = Horizontal (normal)
	// 2 = Mirror horizontal
	// 3 = Rotate 180
	// 4 = Mirror vertical
	// 5 = Mirror horizontal and rotate 270 CW
	// 6 = Rotate 90 CW
	// 7 = Mirror horizontal and rotate 90 CW
	// 8 = Rotate 270 CW

	switch rot {
	case "Horizontal (normal)":
		// Do nothing

	case "Mirror horizontal":
		pixbuf, err = pixbuf.Flip(true)
		check_error("Unable to flip pixbuf", err)

	case "Rotate 180":
		pixbuf, err = pixbuf.RotateSimple(gdk.PIXBUF_ROTATE_UPSIDEDOWN)
		check_error("Unable to rotate pixbuf", err)

	case "Mirror vertical":
		pixbuf, err = pixbuf.Flip(false)
		check_error("Unable to flip pixbuf", err)

	case "Mirror horizontal and rotate 270 CW":
		pixbuf, err = pixbuf.RotateSimple(gdk.PIXBUF_ROTATE_CLOCKWISE)
		check_error("Unable to rotate pixbuf", err)
		pixbuf, err = pixbuf.Flip(true)
		check_error("Unable to flip pixbuf", err)

	case "Rotate 90 CW":
		pixbuf, err = pixbuf.RotateSimple(gdk.PIXBUF_ROTATE_CLOCKWISE)
		check_error("Unable to rotate pixbuf", err)

	case "Mirror horizontal and rotate 90 CW":
		pixbuf, err = pixbuf.RotateSimple(gdk.PIXBUF_ROTATE_COUNTERCLOCKWISE)
		check_error("Unable to rotate pixbuf", err)
		pixbuf, err = pixbuf.Flip(true)
		check_error("Unable to flip pixbuf", err)

	case "Rotate 270 CW":
		pixbuf, err = pixbuf.RotateSimple(gdk.PIXBUF_ROTATE_COUNTERCLOCKWISE)
		check_error("Unable to rotate pixbuf", err)
	}

	return pixbuf
}

func setupImageView() *ImageView {
	return NewImageView()
}

func refreshImageView() {
	imageView.newImagePath = imageRawToDispMap[currRAWImagePath]
	log.Debug("Refreshing image view with new image path: ", imageView.newImagePath, "from raw path: ", currRAWImagePath)
	imageView.image.QueueDraw()
}
