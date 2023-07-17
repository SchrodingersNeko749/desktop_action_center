package main

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func DownloadImage(url string) *gtk.Image {
	response, err := http.Get(url)

	if err != nil || response.StatusCode != 200 {
		return GetDefaultImage()
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, response.Body)

	if err != nil {
		return GetDefaultImage()
	}

	defer response.Body.Close()
	decoded := decode(buf.Bytes())
	if decoded == nil {
		return GetDefaultImage()
	}
	gtkImg := toGtkImage(decoded)
	resizedImg := resize(gtkImg, 64, 64)
	return resizedImg
}

func decode(buf []byte) image.Image {
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
		return GetDefaultImage()
	}
	gtkimg, err := gtk.ImageNewFromPixbuf(pixbuf)
	if err != nil {
		return GetDefaultImage()
	}
	return gtkimg
}

func GetDefaultImage() *gtk.Image {
	img, err := gtk.ImageNewFromIconName("radio", gtk.ICON_SIZE_LARGE_TOOLBAR)
	img.SetPixelSize(64)
	if err != nil {
		fmt.Println(err)
	}
	return img
}

func resize(img *gtk.Image, width int, height int) *gtk.Image {
	pixbuf := img.GetPixbuf()
	pixbuf, err := pixbuf.ScaleSimple(width, height, gdk.INTERP_HYPER)
	if err != nil {
		fmt.Println(err)
	}
	img.SetFromPixbuf(pixbuf)
	return img
}
