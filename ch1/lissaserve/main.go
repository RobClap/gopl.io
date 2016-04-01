package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

//!+main

var palette = []color.Color{color.Black}

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	fillpalette()

	//!+http
	//handler := func(w http.ResponseWriter, r *http.Request) {
	//	lissajous(w)
	//}
	http.HandleFunc("/", handler)
	//!-http
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
	return
}

func handler(w http.ResponseWriter, r *http.Request) {
	phase := parse("phase", r, 0.0)
	freq := parse("freq", r, rand.Float64()*3.0)
	size := int(parse("size", r, 500.0))
	lissajous(w, phase, freq, size)
}
func parse(field string, r *http.Request, _default float64) float64 {
	tmpstring := r.FormValue(field)
	value, err := strconv.ParseFloat(tmpstring, 64)
	if err != nil {
		log.Println(err)
		return _default
	} else {
		return value
	}
}

//!-handler
func fillpalette() {
	frequency := 0.3
	for i := 0; i < 254; i++ {
		red := uint8(math.Sin(frequency*float64(i)+0.0)*127.5 + 128.0)
		green := uint8(math.Sin(frequency*float64(i)+2.0)*127.5 + 128.0)
		blue := uint8(math.Sin(frequency*float64(i)+4.0)*127.5 + 128.0)

		palette = append(palette, color.RGBA{red, green, blue, 0xff})
	}
}

func lissajous(out io.Writer, phase float64, freq float64, size int) {
	const (
		cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		nframes = 64    // number of animation frames
		delay   = 8     // delay between frames in 10ms units
	)
	anim := gif.GIF{LoopCount: nframes}
	index := 0.0
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5), uint8(index/12.0+1))
			index++
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	err := gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
	}
}
