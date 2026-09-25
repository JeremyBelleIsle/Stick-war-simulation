package cam

import (
	"Sticks_War/screen"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
)

var (
	X          float64
	isDragging bool
	lastMouseX int
)

func ScrollCamX() {
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		if !gameutil.ClickInARect(float64(x), float64(y), 0, 0, screen.Widht-600, screen.Height-400) {
			return
		}

		if !isDragging {
			// Début du glissement : on mémorise la position initiale
			isDragging = true
			lastMouseX = x
		} else {
			// Pendant le glissement : on calcule la distance parcourue
			dx := float64(x - lastMouseX)

			// On met à jour la position de la caméra
			X += dx

			if X > 0 {
				X = 0
			}

			// On actualise la dernière position connue
			lastMouseX = x
		}
	} else {
		// Relâchement du bouton
		isDragging = false
	}
}
