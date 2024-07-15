package main

import (
	"fmt"
	"image"
	"os"
	"runtime"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const scaleFactor = 1.2

type ImageView struct {
	drawingArea            *gtk.DrawingArea
	imagePath              string
	pixbuf                 *gdk.Pixbuf
	surface                *cairo.Surface
	image                  image.Image
	newImagePath           string
	scale                  float64
	dragStartX, dragStartY int
	imgStartX, imgStartY   int
	imgWidth, imgHeight    int
}

func NewImageView() *ImageView {
	iv := &ImageView{
		scale: 1.0,
	}

	iv.drawingArea, err = gtk.DrawingAreaNew()
	check_error("Unable to create drawing area", err)

	iv.drawingArea.SetSizeRequest(1920, 1080)
	iv.drawingArea.AddEvents(int(gdk.BUTTON_PRESS_MASK | gdk.BUTTON1_MOTION_MASK | gdk.SCROLL_MASK))

	iv.drawingArea.Connect("draw", iv.draw)
	iv.drawingArea.Connect("button-press-event", iv.onButtonPress)
	iv.drawingArea.Connect("motion-notify-event", iv.onDrag)
	iv.drawingArea.Connect("scroll-event", iv.onScroll)

	return iv
}

func (iv *ImageView) draw(da *gtk.DrawingArea, cr *cairo.Context) {
	// , cr *cairo.Context
	if da == nil {
		panic("Drawing area is nil")
	}
	if cr == nil {
		panic("Cairo context is nil")
	}
	if iv == nil {
		panic("ImageView is nil")
	}

	if iv.imagePath != iv.newImagePath {
		// load image
		iv.pixbuf, err = gdk.PixbufNewFromFile(imageView.newImagePath)
		check_error("Unable to load image", err)
		iv.imgWidth = iv.pixbuf.GetWidth()
		iv.imgHeight = iv.pixbuf.GetHeight()

		drawAreaWidth := iv.drawingArea.GetAllocatedWidth()
		drawAreaHeight := iv.drawingArea.GetAllocatedHeight()

		// Reset values
		iv.imagePath = iv.newImagePath
		iv.scale = min(float64(drawAreaWidth)/float64(iv.imgWidth), float64(drawAreaHeight)/float64(iv.imgHeight))
		iv.dragStartX = 0
		iv.dragStartY = 0
		iv.imgStartX = 0
		iv.imgStartY = 0
	}

	cr.Scale(iv.scale, iv.scale)
	// cr.Translate(float64(iv.imgStartX-iv.imgStartXCurr), float64(iv.imgStartY-iv.imgStartYCurr))

	// draw image
	window, err := da.GetWindow()
	check_error("Unable to get window from drawing area", err)
	// iv.surface = nil
	// iv.surface.Close()

	iv.pixbuf, err = gdk.PixbufNewFromFile(imageView.newImagePath)
	check_error("Unable to load image", err)

	iv.surface, err = gdk.CairoSurfaceCreateFromPixbuf(iv.pixbuf, 1, window)
	check_error("Unable to create surface from pixbuf", err)

	cr.SetSourceSurface(iv.surface, float64(iv.imgStartX), float64(iv.imgStartY))
	cr.Paint()
}

func (iv *ImageView) onScroll(da *gtk.DrawingArea, event *gdk.Event) {
	scrollEvent := gdk.EventScrollNewFromEvent(event)
	direction := scrollEvent.Direction()
	x, y := scrollEvent.X(), scrollEvent.Y()
	deltaX, deltaY := 0, 0

	if direction == gdk.SCROLL_UP {
		iv.scale *= scaleFactor
		deltaX = int(x * (1 - scaleFactor))
		deltaY = int(y * (1 - scaleFactor))

	} else if direction == gdk.SCROLL_DOWN {
		iv.scale /= scaleFactor
		deltaX = int(x * (1 - 1/scaleFactor))
		deltaY = int(y * (1 - 1/scaleFactor))
	}
	iv.imgStartX = iv.imgStartX + deltaX
	iv.imgStartY = iv.imgStartY + deltaY

	// fmt.Println("Scrolling to scale: ", iv.scale)
	iv.drawingArea.QueueDraw()
}

func (iv *ImageView) onButtonPress(da *gtk.DrawingArea, event *gdk.Event) {
	// fmt.Println("Button Pressed")
	buttonEvent := gdk.EventButtonNewFromEvent(event)
	if buttonEvent.Type() == gdk.EVENT_2BUTTON_PRESS {
		// reset scaling and position
		drawAreaWidth := iv.drawingArea.GetAllocatedWidth()
		drawAreaHeight := iv.drawingArea.GetAllocatedHeight()
		iv.scale = min(float64(drawAreaWidth)/float64(iv.imgWidth), float64(drawAreaHeight)/float64(iv.imgHeight))
		iv.imgStartX = 0.0
		iv.imgStartY = 0.0
		iv.drawingArea.QueueDraw()

	} else if buttonEvent.Button() == gdk.BUTTON_PRIMARY {
		iv.dragStartX = int(buttonEvent.X())
		iv.dragStartY = int(buttonEvent.Y())
	}
}

func (iv *ImageView) onDrag(da *gtk.DrawingArea, event *gdk.Event) {
	motionEvent := gdk.EventMotionNewFromEvent(event)
	x, y := motionEvent.MotionVal()
	// fmt.Println("Dragging with motion values", x, y, "and start values", iv.dragStartX, iv.dragStartY)
	deltaX := int(x) - iv.dragStartX
	deltaY := int(y) - iv.dragStartY
	iv.dragStartX = int(x)
	iv.dragStartY = int(y)

	iv.imgStartX += deltaX
	iv.imgStartY += deltaY

	iv.drawingArea.QueueDraw()
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
	imageView.drawingArea.QueueDraw()
}
