package main

import (
	"fmt"
	"image"
	"os"
	"runtime"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const SCALE_FACTOR = 1.5

type ImageView struct {
	scrolledWindow               *gtk.ScrolledWindow
	image                        *gtk.Image
	pixbuf_orig                  *gdk.Pixbuf
	imagePath                    string
	newImagePath                 string
	scale                        float64
	dragStartX, dragStartY       float64
	prevScrollX, prevScrollY     float64
	hScrollbarVal, vScrollbarVal float64 // needed to get around the issue of scrolling to the top left corner on double click
	imgOffsetX, imgOffsetY       int
}

func NewImageView() *ImageView {
	iv := &ImageView{
		scale: 1.0,
	}

	iv.scrolledWindow, err = gtk.ScrolledWindowNew(nil, nil)
	check_error("Unable to create scrolled window", err)

	iv.image, err = gtk.ImageNew()
	check_error("Unable to create image", err)
	iv.image.SetCanFocus(false)

	iv.scrolledWindow.SetSizeRequest(1920, 1080)
	iv.scrolledWindow.AddEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK | gdk.SCROLL_MASK))

	// iv.scrolledWindow.SetOverlayScrolling(false)
	iv.scrolledWindow.SetCaptureButtonPress(true)

	iv.image.Connect("draw", iv.draw)
	iv.scrolledWindow.Connect("button-press-event", iv.onButtonPress)
	iv.scrolledWindow.Connect("motion-notify-event", iv.onDrag)
	iv.scrolledWindow.Connect("scroll-event", iv.onScroll)

	iv.scrolledWindow.Add(iv.image)

	return iv
}

func (iv *ImageView) draw(img *gtk.Image) {
	scrollWidth := iv.scrolledWindow.GetAllocatedWidth()
	scrollHeight := iv.scrolledWindow.GetAllocatedHeight()

	if iv.imagePath != iv.newImagePath {
		// load image
		img.SetFromFile(iv.newImagePath)

		iv.pixbuf_orig = nil
		iv.pixbuf_orig = img.GetPixbuf()

		// Reset values
		iv.imagePath = iv.newImagePath
		iv.scale = 1
		iv.dragStartX = 0
		iv.dragStartY = 0
	}

	// fmt.Println("Drawing image with scale: ", iv.scale)

	imageRatio := float64(iv.pixbuf_orig.GetWidth()) / float64(iv.pixbuf_orig.GetHeight())

	newHeight := int(float64(scrollHeight) * iv.scale)
	newWidthFromHeight := int(float64(newHeight) * imageRatio)

	newWidth := int(float64(scrollWidth) * iv.scale)
	newHeightFromWidth := int(float64(newWidth) / imageRatio)

	if newHeight <= scrollHeight && newWidthFromHeight <= scrollWidth {
		newWidth = newWidthFromHeight
	} else {
		newHeight = newHeightFromWidth
	}

	iv.imgOffsetX = (newWidth - scrollWidth) / 2
	iv.imgOffsetY = (newHeight - scrollHeight) / 2

	pixbuf, err := iv.pixbuf_orig.ScaleSimple(newWidth, newHeight, gdk.INTERP_BILINEAR)
	check_error("Unable to scale pixbuf", err)

	img.SetFromPixbuf(pixbuf)

	fmt.Println("Setting scrollbar values to: ", iv.hScrollbarVal, iv.vScrollbarVal)

	iv.scrolledWindow.GetHAdjustment().SetValue(iv.hScrollbarVal)
	iv.scrolledWindow.GetVAdjustment().SetValue(iv.vScrollbarVal)

	// fmt.Println(iv.image.GetAllocatedWidth(), iv.image.GetAllocatedHeight(), int(float64(scrollWidth)*iv.scale), int(float64(scrollHeight)*iv.scale), pixbuf.GetWidth(), pixbuf.GetHeight())
}

func (iv *ImageView) onScroll(sw *gtk.ScrolledWindow, event *gdk.Event) {
	scrollEvent := gdk.EventScrollNewFromEvent(event)
	direction := scrollEvent.Direction()

	// zoom if ctrl is pressed
	if scrollEvent.State()&gdk.CONTROL_MASK == 0 {
		return
	}

	scaleBefore := iv.scale
	x, y := scrollEvent.X(), scrollEvent.Y()

	fmt.Println(x, y)

	hScrollbarValueBefore := sw.GetHScrollbar().GetValue()
	vScrollbarValueBefore := sw.GetVScrollbar().GetValue()

	switch direction {
	case gdk.SCROLL_UP:
		iv.scale *= SCALE_FACTOR
	case gdk.SCROLL_DOWN:
		iv.scale /= SCALE_FACTOR
	case gdk.SCROLL_SMOOTH:
		iv.scale -= scrollEvent.DeltaY() / 100 * SCALE_FACTOR
	}

	fmt.Println("Scrolling to scale: ", iv.scale, scrollEvent.DeltaY())

	iv.hScrollbarVal = (float64(x)+hScrollbarValueBefore)*(iv.scale/scaleBefore) - float64(x)
	iv.vScrollbarVal = (float64(y)+vScrollbarValueBefore)*(iv.scale/scaleBefore) - float64(y) - scrollEvent.DeltaY()*sw.GetVScrollbar().GetAdjustment().GetMinimumIncrement()/5

	iv.image.QueueDraw()
}

func (iv *ImageView) onButtonPress(sw *gtk.ScrolledWindow, event *gdk.Event) {
	// fmt.Println("Button Pressed")
	buttonEvent := gdk.EventButtonNewFromEvent(event)
	// if buttonEvent.Type() == gdk.EVENT_2BUTTON_PRESS {
	if buttonEvent.Button() == gdk.BUTTON_SECONDARY {
		// reset scaling
		if iv.scale != 1.0 {
			iv.scale = 1.0
		} else {
			iv.scale = 3.0

			iv.hScrollbarVal = (buttonEvent.X()/float64(sw.GetAllocatedWidth())*(sw.GetHAdjustment().GetUpper()-sw.GetHAdjustment().GetLower()) + sw.GetHAdjustment().GetLower()) * 2
			iv.vScrollbarVal = (buttonEvent.Y()/float64(sw.GetAllocatedHeight())*(sw.GetVAdjustment().GetUpper()-sw.GetVAdjustment().GetLower()) + sw.GetVAdjustment().GetLower()) * 2
		}
		iv.image.QueueDraw()

	} else if buttonEvent.Button() == gdk.BUTTON_PRIMARY {
		iv.dragStartX = buttonEvent.X()
		iv.dragStartY = buttonEvent.Y()
		iv.prevScrollX = sw.GetHScrollbar().GetValue()
		iv.prevScrollY = sw.GetVScrollbar().GetValue()
	}
}

func (iv *ImageView) onDrag(sw *gtk.ScrolledWindow, event *gdk.Event) {
	motionEvent := gdk.EventMotionNewFromEvent(event)
	if motionEvent.State()&gdk.BUTTON1_MASK != 0 {
		x, y := motionEvent.MotionVal()
		deltaX := x - iv.dragStartX
		deltaY := y - iv.dragStartY

		iv.hScrollbarVal = iv.prevScrollX - float64(deltaX)
		iv.vScrollbarVal = iv.prevScrollY - float64(deltaY)
	}
}

func rotateImage(img image.Image) image.Image {
	// Implement rotation logic here if needed
	return img
}

func resizeImage(img image.Image, width, height int) image.Image {
	// Implement resize logic here
	return img
}

func main_old() {
	runtime.LockOSThread()
	gtk.Init(nil)

	window, _ := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Image Viewer")
	window.SetDefaultSize(1920, 1080)
	window.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Load images here
	var images []image.Image
	imgFile, _ := os.Open("image.png") // Use actual file paths
	img, _, _ := image.Decode(imgFile)
	images = append(images, img)

	NewImageView()

	gtk.Main()
}

func setupImageView() *ImageView {
	return NewImageView()
}

func refreshImageView() {
	imageView.newImagePath = imageRawToDispMap[currRAWImagePath]
	fmt.Println("Refreshing image view with new image path: ", imageView.newImagePath, "from raw path: ", currRAWImagePath)
	imageView.image.QueueDraw()
}
