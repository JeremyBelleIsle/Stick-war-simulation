package stickman

import (
	rock "Sticks_War/Rock"
	"Sticks_War/cam"
	"Sticks_War/camp"
	"Sticks_War/ordi"
	"Sticks_War/screen"
	"Sticks_War/store"
	"fmt"
	"image/color"

	"github.com/JeremyBelleIsle/gameutil"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type StickMan struct {
	x, y, w, h     float64
	capacity       int
	mobilitySpeed  int
	miningSpeed    int
	miningCooldown int
	gold           int
	id             int
	team           string
	class          string
	clr            color.RGBA
}

func Spawn(StickMans *[]StickMan, offers []store.Offer, campX, campY float64, class string, team string, money *int) {
	// 1. Trouver le coût de l'offre pour la classe demandée
	var cost int
	found := false
	for _, o := range offers {
		if o.TypeS == class {
			cost = o.Cost
			found = true
			break
		}
	}
	if !found {
		return
	}

	// 2. Vérifier et déduire l'argent selon l'équipe
	if team == "player" {
		if *money >= cost {
			*money -= cost
		} else {
			return
		}
	} else {
		if ordi.Money >= cost {
			ordi.Money -= cost
		} else {
			return
		}
	}

	// 3. Compter combien de StickMen de cette classe existent déjà pour cette équipe
	count := 0
	for _, sm := range *StickMans {
		if sm.team == team && sm.class == class {
			count++
		}
	}

	// 4. Définir les limites maximales par classe
	maxLimit := 0
	switch class {
	case "minner":
		maxLimit = 4
	case "attacker":
		maxLimit = 12
	case "archer":
		maxLimit = 4
	}

	// 5. Vérifier si la limite est atteinte (avec remboursement si besoin)
	if count >= maxLimit {
		if team == "player" {
			*money += cost
		} else {
			ordi.Money += cost
		}
		return
	}

	// 6. Définir la couleur et la capacité selon la classe
	capacity := 5
	clr := color.RGBA{255, 255, 255, 255} // Blanc par défaut (minner)

	switch class {
	case "attacker":
		capacity = 0
		clr = color.RGBA{255, 0, 0, 255} // Rouge
	case "archer":
		capacity = 0
		clr = color.RGBA{255, 165, 0, 255} // Orange
	}

	// 7. Ajouter le nouveau StickMan
	*StickMans = append(*StickMans, StickMan{
		x:             campX + 100,
		y:             campY,
		w:             130,
		h:             130,
		capacity:      capacity,
		mobilitySpeed: 7,
		id:            count, // L'ID correspond directement au nombre actuel (0, 1, 2...)
		team:          team,
		class:         class,
		clr:           clr,
	})
}

func (sm *StickMan) Move(rocks []rock.Rock, camps []camp.Camp, gameMode string) {
	if len(rocks) == 0 || len(camps) == 0 {
		return
	}

	// Sélection sécurisée des ressources selon l'équipe
	// Si l'index 2 n'existe pas, on se replie sur le premier élément pour éviter les crashs
	rIdx, cIdx := 0, 0
	if sm.team != "player" && len(rocks) > 1 && len(camps) > 1 {
		rIdx, cIdx = 1, 1
	}
	targetRock := rocks[rIdx]
	targetCamp := camps[cIdx]

	switch sm.class {
	case "minner":
		x := float32(targetRock.X + float64(sm.id)*150)
		y := float32(targetRock.Y + targetRock.H + 30.0)

		// Si le mineur a son sac plein, il retourne au camp
		if sm.gold >= sm.capacity {
			x = float32(targetCamp.X)
			y = float32(targetCamp.Y)
		}

		sm.moveTo(x, y)

	case "attacker":
		if gameMode == "create" {
			space := 150
			maxPerColumn := 4
			col := sm.id / maxPerColumn
			row := sm.id % maxPerColumn
			x := float32(0)

			if sm.team == "ordi" {
				x = float32(screen.Widht-800) + float32(col*space) + 4000
			} else {
				x = float32(screen.Widht-800) - float32(col*space)
			}
			y := float32(targetCamp.Y-470) + float32(row*space)

			sm.moveTo(x, y)
		} else {
			sm.x += float64(sm.mobilitySpeed)
		}

	case "archer":
		if gameMode == "create" {
			x := float32(0)
			if sm.team == "player" {
				x = float32(targetCamp.X + 570)
			} else {
				x = float32(targetCamp.X - 570)
			}
			y := float32(targetCamp.Y - 470 + (float64(sm.id) * 150))

			sm.moveTo(x, y)
		} else {
			sm.x += float64(sm.mobilitySpeed)
		}
	}
}

// Petite fonction utilitaire interne pour éviter de répéter l'appel à DirigePointToPoint
func (sm *StickMan) moveTo(targetX, targetY float32) {
	xf32, yf32 := gameutil.DirigePointToPoint(float32(sm.mobilitySpeed), float32(sm.x), float32(sm.y), targetX, targetY)
	sm.x, sm.y = float64(xf32), float64(yf32)
}

func (sm *StickMan) Mine(rocks []rock.Rock) {
	for _, rock := range rocks {
		if sm.y > rock.Y+rock.H+30.0 || sm.x < rock.X-30 || sm.x > rock.X+rock.W+30 {
			continue
		}

		if sm.miningCooldown <= 0 {
			sm.gold++
			sm.miningCooldown = 60
		} else {
			sm.miningCooldown--
		}
		break
	}
}

func (sm *StickMan) DropMoney(money *int, camps []camp.Camp) {
	for _, camp := range camps {
		targetX := camp.X
		targetY := camp.Y

		dx := sm.x - targetX
		dy := sm.y - targetY
		if (dx*dx + dy*dy) < 25 { // 5² = 25
			if sm.team == "player" {
				*money += sm.gold
			} else {
				ordi.Money += sm.gold
			}
			sm.gold = 0
			break
		}
	}
}

func (sm *StickMan) Draw(screenI *ebiten.Image, mplusSource *text.GoTextFaceSource) {
	if sm.class == "minner" {
		// draw fluo rectangle (for the capacity)
		vector.StrokeRect(screenI, float32(sm.x+cam.X), float32(sm.y-70), 150, 70, 7, color.RGBA{255, 123, 25, 255}, true)
		// draw capacity (minner)
		gameutil.DrawText(fmt.Sprintf("%d/%d", sm.gold, sm.capacity), 50, screen.Widht, sm.x+cam.X, sm.y-50, 0, screenI, sm.clr, mplusSource)
	}
	// draw stickman
	vector.DrawFilledRect(screenI, float32(sm.x+cam.X), float32(sm.y), float32(sm.w), float32(sm.h), sm.clr, true)
}
