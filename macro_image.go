package main

import (
	"fmt"
	"os"
	"bytes"
	"path/filepath"
	"image"
	"image/png"
	_ "image/gif"
	_ "image/jpeg"
	"github.com/rogpeppe/misc/svg"
)

type MacroImage struct {
	Path string
	Title string
	Alt string
	Height int
	filename string
	Data []byte
	// border bool
}




func newMacroImage(path, title, alt string ) (*MacroImage, error) {
	m := new (MacroImage)
	m.Path = path

	filename := filepath.Base(path)
	extension := filepath.Ext(path)
 
	if extension == ".svg" {
		logger.Debug("Convertiong SVG %s to PNG", path)
		// convert SVG to PNG
		extension = ".png"
		filename = filename[0:len(filename)-len(extension)] + ".png"

		file, err := os.Open(path)
		if (err != nil){
			return nil, err;
		}
		defer file.Close()
		size := image.Point{1000, 1000}
		dest, err := svg.Render(file, size)
		if (err != nil){
			return nil, err;
		}


		b := new(bytes.Buffer)

		err = png.Encode(b, dest)
		if err != nil {
			return nil, err
		}
		
		m.Data =  b.Bytes()
	}

	m.filename = filename
	if title == "" {
		extension := filepath.Ext(path)
		title = m.filename[0:len(m.filename)-len(extension)]
	} else {
		m.Title = title
	}
	m.Alt = alt

	if (m.Data == nil) {
		reader, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		
		buf := new(bytes.Buffer)
		buf.ReadFrom(reader)
		m.Data = buf.Bytes()
		reader.Close()
	}

	img, _, err := image.Decode(bytes.NewReader(m.Data))
	  
  	if err != nil {
  		return nil, err
  	}
	bounds := img.Bounds()

	// wider than high?
	if (bounds.Max.X > bounds.Max.Y) {
		if bounds.Max.Y > 250 {
			m.Height = 250
		} else {
			m.Height = bounds.Max.Y
		}
	} else {
		if bounds.Max.X > 450 {
			m.Height = (bounds.Max.Y) / (bounds.Max.X / 450)
		} else {
			m.Height = bounds.Max.Y
		}
	}

	fmt.Printf("%#v -> %d", bounds.Max, m.Height)
	  
	return m, nil
}

func (image MacroImage) Render() string {
	return fmt.Sprintf(
		`<ac:image ac:alt="%s" ac:height="%d"  ac:title="%s" ><ri:attachment ri:filename="%s" /></ac:image>`,
		// ac:border="%s" 
		image.Alt,
		image.Height, 
		image.Title,
	//	image.border,
		filepath.Base(image.Path),
	)
}