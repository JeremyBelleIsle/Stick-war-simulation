package gamemode

import (
	"Sticks_War/screen"
	"image/color"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var (
	x, y, w, h            = screen.Widht - 700, screen.Height - 250, 670, 220
	Mode                  = "create"
	defenseIconUnselected = gameutil.LoadImage("gameMode/defense icon.png")
	attackIconUnselected  = gameutil.LoadImage("gameMode/attack icon.png")
	defenseIconSelected   = gameutil.LoadImage("gameMode/defense icon (selected).png")
	attackIconSelected    = gameutil.LoadImage("gameMode/attack icon (selected).png")
)

func DetectClickToChangeGameMode(gameMode *string) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}

	cx, cy := ebiten.CursorPosition()

	if gameutil.ClickInARect(float64(cx), float64(cy), float64(x), float64(y), float64(w)/2, float64(h)) {
		*gameMode = "create"
	}
	if gameutil.ClickInARect(float64(cx), float64(cy), float64(x+(w/2)), float64(y), float64(w)/2, float64(h)) {
		*gameMode = "attack"
	}
}

func Draw(screenI *ebiten.Image) {
	vector.DrawFilledRect(screenI, float32(x), float32(y), float32(w), float32(h), color.RGBA{255, 255, 255, 255}, true)
	vector.StrokeRect(screenI, float32(x), float32(y), float32(w), float32(h), 20, color.RGBA{255, 255, 0, 255}, true)

	op := &ebiten.DrawImageOptions{}

	bounds := defenseIconUnselected.Bounds()
	w := float64(bounds.Dx())
	h := float64(bounds.Dy())

	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Scale(.5, .5)
	op.GeoM.Translate(float64(x+100), float64(y+110))
	screenI.DrawImage(defenseIconUnselected, op)

	op = &ebiten.DrawImageOptions{}

	bounds = attackIconUnselected.Bounds()
	w = float64(bounds.Dx())
	h = float64(bounds.Dy())

	op.GeoM.Translate(-w/2, -h/2)
	op.GeoM.Scale(.5, .5)
	op.GeoM.Translate(float64(x+int(w)+100), float64(y+110))
	screenI.DrawImage(attackIconUnselected, op)
}
