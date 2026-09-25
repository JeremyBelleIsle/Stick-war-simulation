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
	dir        float64
	sx, sy     float64
	img        *ebiten.Image
}

func Init(camps *[]Camp) {
	x := 500.0
	dir := 0.0
	sx := 1.5
	for i := 0; i < 2; i++ {
		if i == 1 {
			x *= 14
			sx = -sx
		}

		*camps = append(*camps, Camp{
			X:   x,
			Y:   screen.Height - 350,
			dir: dir,
			sx:  sx,
			sy:  1.5,
			img: gameutil.LoadImage("camp/campImg.png"),
		})
	}
}

func (c Camp) Draw(screenI *ebiten.Image, mplusSource *text.GoTextFaceSource) {
	op := &ebiten.DrawImageOptions{}

	bounds := c.img.Bounds()
	w := float64(bounds.Dx())
	h := float64(bounds.Dy())

	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Scale(c.sx, c.sy)
	op.GeoM.Rotate(c.dir)
	op.GeoM.Translate(c.X+cam.X, c.Y)
	screenI.DrawImage(c.img, op)

	if c.X > 3000 {
		gameutil.DrawText(fmt.Sprintf("%d$", ordi.Money), 125, screen.Widht, c.X+cam.X-200, c.Y-200, 0, screenI, color.RGBA{255, 255, 255, 255}, mplusSource)
	}
}
