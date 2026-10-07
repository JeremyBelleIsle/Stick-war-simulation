package rock

import (
	"Sticks_War/cam"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Rock struct {
	X, Y, W, H float64
	dir        float64
	sx, sy     float64
	clr        color.RGBA
	img        *ebiten.Image
}

func New(x, y float64) Rock {
	dir := 0.0
	sx := 1.5
	W, H := 600, 300
	// 1. Créer l'image avec la taille de la roche (w et h)
	img := ebiten.NewImage(600, 300)

	// 2. Dessiner le rectangle de base à partir de (0, 0) sur l'image r.img
	vector.DrawFilledRect(img, 0, 0, float32(W), float32(H), color.RGBA{128, 128, 128, 255}, true)

	return Rock{
		X:   x,
		Y:   y,
		W:   float64(W),
		H:   float64(H),
		dir: dir,
		sx:  sx,
		sy:  1.5,
		img: img,
	}
}
func (r *Rock) Draw(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	// Positionner la roche aux coordonnées mondiales (r.x, r.y) sur l'écran
	op.GeoM.Translate(r.X+cam.X, r.Y)

	// Dessiner l'image de la roche SUR l'écran
	screen.DrawImage(r.img, op)
}
