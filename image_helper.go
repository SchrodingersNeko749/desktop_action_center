package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/kbinani/screenshot"
)

func ImgDefault() *gtk.Image {
	img, err := gtk.ImageNewFromIconName("radio", gtk.ICON_SIZE_LARGE_TOOLBAR)
	if err != nil {
		fmt.Println(err)
	}
	img.SetPixelSize(64)
	return img
}
func ImgDefaultSized(size int) *gtk.Image {
	img, err := gtk.ImageNewFromIconName("radio", gtk.ICON_SIZE_LARGE_TOOLBAR)
	if err != nil {
		fmt.Println(err)
	}
	img.SetPixelSize(size)
	return img
}

func ImgFromTheme(name string, size int) *gtk.Image {
	img, err := gtk.ImageNewFromIconName(name, gtk.ICON_SIZE_LARGE_TOOLBAR)
	if err != nil {
		fmt.Println(err)
	}
	img.SetPixelSize(size)
	return img
}

func ImgGetFromPath(path string, size int) *gtk.Image {
	img, err := gtk.ImageNewFromFile(path)
	if err != nil {
		fmt.Println(err)
	}
	img = ImgResize(img, size, size)
	return img
}

func ImgDownload(url string, size int) *gtk.Image {
	response, err := http.Get(url)

	if err != nil || response.StatusCode != 200 {
		return ImgDefault()
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, response.Body)

	if err != nil {
		return ImgDefaultSized(size)
	}

	defer response.Body.Close()
	decoded := ImgDecode(buf.Bytes())
	if decoded == nil {
		return ImgDefaultSized(size)
	}
	gtkImg := toGtkImage(decoded)
	resizedImg := ImgResize(gtkImg, size, size)
	return resizedImg
}

func ImgGetScreenshot(path string) *gtk.Image {
	bounds := screenshot.GetDisplayBounds(0)
	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		panic(err)
	}
	file, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
	return ImgGetFromPath(path, 0)
}

func ImgGetFromRaw(width int32, height int32, stride int32, hasAlpha bool, bits int32, bytes []byte) *gtk.Image {
	pixbuf, err := gdk.PixbufNewFromData(bytes, gdk.COLORSPACE_RGB, hasAlpha, int(bits), int(width), int(height), int(stride))
	img, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		fmt.Println(err)
	}
	return img
}

func ImgDecode(buf []byte) image.Image {
	rdr := bytes.NewReader(buf)
	img, err := png.Decode(rdr)
	if err != nil {
		img, err = jpeg.Decode(rdr)
		if err != nil {
			return nil
		}
	}
	return img
}

func toGtkImage(image image.Image) *gtk.Image {
	width := image.Bounds().Size().X
	height := image.Bounds().Size().Y
	stride := width * 4
	pixels := make([]byte, height*stride)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, a := image.At(x, y).RGBA()
			p := (y*width + x) * 4
			pixels[p] = byte(r)
			pixels[p+1] = byte(g)
			pixels[p+2] = byte(b)
			pixels[p+3] = byte(a)
		}
	}

	pixbuf, err := gdk.PixbufNewFromData(pixels, gdk.COLORSPACE_RGB, true, 8, width, height, stride)
	if err != nil {
		return ImgDefault()
	}
	gtkimg, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		return ImgDefault()
	}
	return gtkimg
}

func ImgResize(img *gtk.Image, width int, height int) *gtk.Image {
	if width <= 0 {
		return img
	}
	pixbuf := img.GetPixbuf()
	w := img.GetPixbuf().GetWidth()
	h := img.GetPixbuf().GetHeight()

	aspectRatio := float64(w) / float64(h)

	// if width is set and height is not set, calculate height
	if width > 0 && height <= 0 {
		height = int(float64(width) / aspectRatio)
	}

	pixbuf, err := pixbuf.ScaleSimple(width, height, gdk.INTERP_HYPER)
	if err != nil {
		fmt.Println(err)
	}
	img.SetFromPixbuf(pixbuf)
	return img
}
