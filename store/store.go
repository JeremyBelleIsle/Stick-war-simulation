package store

import (
	"Sticks_War/screen"
	"fmt"
	"image/color"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Offer struct {
	x, y, w, h float32
	bodyClr    color.RGBA
	outLineClr color.RGBA
	Cost       int
	TypeS      string
}

func Init(offers *[]Offer) {
	*offers = []Offer{
		{x: screen.Widht - 500, y: 50, w: 470, h: 70, outLineClr: color.RGBA{255, 223, 0, 255}, bodyClr: color.RGBA{0, 32, 255, 255}, Cost: 5, TypeS: "minner"},
		{x: screen.Widht - 500, y: 150, w: 470, h: 70, outLineClr: color.RGBA{255, 223, 0, 255}, bodyClr: color.RGBA{0, 32, 255, 255}, Cost: 5, TypeS: "attacker"},
		{x: screen.Widht - 500, y: 250, w: 470, h: 70, outLineClr: color.RGBA{255, 223, 0, 255}, bodyClr: color.RGBA{0, 32, 255, 255}, Cost: 7, TypeS: "archer"},
	}
}

func (o *Offer) DetectClickOnOffer() bool {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return false
	}

	cx, cy := ebiten.CursorPosition()
	if gameutil.ClickInARect(float64(cx), float64(cy), float64(o.x), float64(o.y), float64(o.w), float64(o.h)) {
		return true
	}

	return false
}

func (o *Offer) Draw(screenI *ebiten.Image, mplusSource *text.GoTextFaceSource) {
	// draw box
	vector.DrawFilledRect(screenI, o.x, o.y, o.w, o.h, o.bodyClr, true)
	// draw the outline of the box
	vector.StrokeRect(screenI, o.x, o.y, o.w, o.h, 10, o.outLineClr, true)
	// draw description
	gameutil.DrawText(o.TypeS, 40, screen.Widht, float64(o.x)+20, float64(o.y)+20, 0, screenI, color.RGBA{255, 255, 255, 255}, mplusSource)
	// draw cost
	gameutil.DrawText(fmt.Sprintf("%d$", o.Cost), 40, screen.Widht, float64(o.x)-150, float64(o.y)+20, 0, screenI, color.RGBA{255, 255, 255, 255}, mplusSource)
}
