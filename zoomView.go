package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

type ZoomView struct {
	box              *gtk.Box
	scrolledWindow   *gtk.ScrolledWindow
	image            *gtk.Image
	pixbuf_orig      *gdk.Pixbuf
	valPosX, valPosY float64
	zoomVal          float64
}

func NewZoomView() *ZoomView {
	zv := &ZoomView{}

	zv.valPosX = 0
	zv.valPosY = 0

	zv.box, err = gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	check_error("Unable to create scrolled window", err)

	zv.scrolledWindow, err = gtk.ScrolledWindowNew(nil, nil)
	check_error("Unable to create ScrolledWindow", err)
	zv.scrolledWindow.SetVExpand(true)

	zv.image, err = gtk.ImageNew()
	check_error("Unable to create image", err)
	zv.image.Connect("draw", zv.draw)

	zv.scrolledWindow.SetSizeRequest(400, 1080)
	zv.scrolledWindow.Add(zv.image)
	zv.box.Add(zv.scrolledWindow)

	return zv
}

func (zv *ZoomView) draw(img *gtk.Image) {
	img.SetFromPixbuf(zv.pixbuf_orig)

	// calculate max values of the adjustments of the scrollbars, since the value is between 0 and max - page_size
	hMaxVal := zv.scrolledWindow.GetHAdjustment().GetUpper() - zv.scrolledWindow.GetHAdjustment().GetPageSize()
	vMaxVal := zv.scrolledWindow.GetVAdjustment().GetUpper() - zv.scrolledWindow.GetVAdjustment().GetPageSize()

	zv.scrolledWindow.GetHScrollbar().SetValue(zv.valPosX * hMaxVal)
	zv.scrolledWindow.GetVScrollbar().SetValue(zv.valPosY * vMaxVal)
}

func (zv *ZoomView) QueueDraw() {
	zv.draw(zv.image)
}

func (zv *ZoomView) changeImage(pixbuf *gdk.Pixbuf) {
	zv.pixbuf_orig = pixbuf
	zv.QueueDraw()
}
