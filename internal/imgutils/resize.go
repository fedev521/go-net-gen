package imgutils

import (
	"image"
	"image/png"
	"os"

	"golang.org/x/image/draw"
)

func Resize() {
	// aprire l'immagine PNG
	file, err := os.Open("out.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	img, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	// calcolare la dimensione del padding
	width := img.Bounds().Size().X
	height := img.Bounds().Size().Y

	// creare un'immagine quadrata con il padding aggiunto
	squareSize := width
	if height > width {
		squareSize = height
	}
	squareImg := image.NewRGBA(image.Rect(0, 0, squareSize, squareSize))
	draw.Draw(squareImg, squareImg.Bounds(), image.White, image.Point{}, draw.Src)
	draw.Draw(squareImg, image.Rect(0, 0, width, height), img, img.Bounds().Min, draw.Src)

	// ridimensionare l'immagine quadrata a 1024x1024 pixel
	resizedImg := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	draw.NearestNeighbor.Scale(resizedImg, resizedImg.Bounds(), squareImg, squareImg.Bounds(), draw.Over, nil)

	// salvare l'immagine ridimensionata
	outFile, err := os.Create("out1024.png")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()
	err = png.Encode(outFile, resizedImg)
	if err != nil {
		panic(err)
	}
}
