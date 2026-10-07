package camp

import (
	"Sticks_War/cam"
	"Sticks_War/ordi"
	"Sticks_War/screen"
	"fmt"
	"image/color"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Camp struct {
	X, Y, w, h float64
	sx, sy     float64
	img        *ebiten.Image
}

func New(x, y, sizeX float64) Camp {
	return Camp{
		X:   x,
		Y:   y,
		sx:  sizeX,
		sy:  1.5,
		img: gameutil.LoadImage("camp/campImg.png"),
	}
}

func (c Camp) Draw(screenI *ebiten.Image, mplusSource *text.GoTextFaceSource) {
	op := &ebiten.DrawImageOptions{}

	bounds := c.img.Bounds()
	w := float64(bounds.Dx())
	h := float64(bounds.Dy())

	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Scale(c.sx, c.sy)
	op.GeoM.Translate(c.X+cam.X, c.Y)
	screenI.DrawImage(c.img, op)

	if c.X > 3000 {
		gameutil.DrawText(fmt.Sprintf("%d$", ordi.Money), 125, screen.Widht, c.X+cam.X-200, c.Y-200, 0, screenI, color.RGBA{255, 255, 255, 255}, mplusSource)
	}
}
