package stickman

import (
	rock "Sticks_War/Rock"
	"Sticks_War/cam"
	"Sticks_War/camp"
	"Sticks_War/ordi"
	"Sticks_War/screen"
	"Sticks_War/store"
	"Sticks_War/unit"
	"fmt"
	"image/color"
	"math"

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
	Attack         bool
	attackCooldown int
	targeted       int
	hp             float64
	damage         float64
	team           string
	class          unit.Type
	clr            color.RGBA
}

func Spawn(StickMans *[]StickMan, offers []store.Offer, campX, campY float64, class unit.Type, team string, money *int) {
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
	case unit.Miner:
		maxLimit = 4
	case unit.Attacker:
		maxLimit = 12
	case unit.Archer:
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
	clr := color.RGBA{255, 255, 255, 255} // Blanc par défaut (miner)

	switch class {
	case unit.Attacker:
		capacity = 0
		if team == "player" {
			clr = color.RGBA{255, 0, 0, 255} // Rouge
		} else {
			clr = color.RGBA{0, 255, 0, 255} // vert
		}
	case unit.Archer:
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
		targeted:      -1,
		damage:        10,
		hp:            20,
		clr:           clr,
	})
}

func (sm *StickMan) findCloserEnemy(stickmans []StickMan, targetCampY float32) *StickMan {
	// --- GESTION DE LA CIBLE ENNEMIE LA PLUS PROCHE ---
	minDist := float64(screen.Widht - 500) // Seuil ou distance maximale de détection (ici basé sur screen.Widht)

	for i := range stickmans {
		enemy := &stickmans[i]

		// si un autre attaquant l'a déjà, alors en choisir un autre
		if enemy.targeted != sm.id && enemy.targeted != -1 {
			if i == len(stickmans)+1 {
				// si on a trouvé aucune cible, on retourne à la position defensive
				enemy.moveToDefendPos(targetCampY)
			}

			continue
		}

		// On ne cible que l'équipe adverse
		if enemy.team != sm.team {
			// Calcul de la distance euclidienne (ou distance horizontale selon votre jeu)
			dx := enemy.x - sm.x
			dy := enemy.y - sm.y
			dist := dx*dx + dy*dy // Utiliser le carré de la distance évite un sqrt coûteux

			// On cherche le plus proche dans le rayon de détection (ex: screen.Widht au carré ou une autre limite)
			if dist < (minDist * minDist) {
				minDist = math.Sqrt(dist) // ou garder le carré directement si vous préférez
				enemy.targeted = sm.id
				return enemy
			}
		}
	}

	return nil
}

func (sm *StickMan) moveToDefendPos(targetCampY float32) {
	// Comportement par défaut si aucun ennemi n'est à portée
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
	y := float32(targetCampY-470) + float32(row*space)

	sm.targeted = -1
	sm.moveTo(x, y)
}

func (sm *StickMan) rushEnemy(closestEnemy *StickMan) {
	if closestEnemy.targeted != sm.id {
		return
	}

	sm.moveTo(float32(sm.x), float32(closestEnemy.y))
	if sm.team == "player" {
		sm.moveTo(float32(closestEnemy.x-sm.w-30), float32(sm.y))
	} else {
		sm.moveTo(float32(closestEnemy.x+sm.w+30), float32(sm.y))
	}

	closestEnemy.targeted = sm.id
}

func (sm *StickMan) Move(rocks []rock.Rock, camps []camp.Camp, stickmans []StickMan, gameMode string) {
	if len(rocks) == 0 || len(camps) == 0 {
		return
	}

	if sm.Attack {
		return
	}
	targetRock := rocks[0]
	targetCamp := camps[0]
	if sm.team != "player" {
		targetRock, targetCamp = rocks[1], camps[1]
	}

	switch sm.class {
	case unit.Miner:
		x := float32(targetRock.X + float64(sm.id)*150)
		y := float32(targetRock.Y + targetRock.H + 30.0)

		if sm.gold >= sm.capacity {
			x = float32(targetCamp.X)
			y = float32(targetCamp.Y)
		}

		sm.moveTo(x, y)

	case unit.Attacker:

		if gameMode == "attack" {
			closestEnemy := sm.findCloserEnemy(stickmans, float32(targetCamp.Y))

			// Si un ennemi est trouvé à portée, on fonce sur lui
			if closestEnemy != nil {
				sm.rushEnemy(closestEnemy)
			} else {
				prevY := sm.y
				sm.moveTo(float32(stickmans[0].x), float32(sm.y))
				if math.Abs(sm.y-prevY) <= float64(sm.mobilitySpeed)/4 {
					if sm.team == "player" {
						sm.x += float64(sm.mobilitySpeed)
					} else {
						sm.x -= float64(sm.mobilitySpeed)
					}
				}
			}
			return
		}
		sm.moveToDefendPos(float32(targetCamp.Y))

	case unit.Archer:
		if gameMode == "attack" && sm.team == "player" {
			sm.x += float64(sm.mobilitySpeed)
			return
		}

		x := float32(0)
		if sm.team == "player" {
			x = float32(targetCamp.X + 570)
		} else {
			x = float32(targetCamp.X - 570)
		}
		y := float32(targetCamp.Y - 470 + (float64(sm.id) * 150))

		sm.moveTo(x, y)
	}
}

// Petite fonction utilitaire interne pour éviter de répéter l'appel à DirigePointToPoint
func (sm *StickMan) moveTo(targetX, targetY float32) {
	xf32, yf32 := gameutil.DirigePointToPoint(float32(sm.mobilitySpeed), float32(sm.x), float32(sm.y), targetX, targetY)
	sm.x, sm.y = float64(xf32), float64(yf32)
}

func CheckColl(stickmans *[]StickMan) {
	for i := range *stickmans {

		smOrdi := &(*stickmans)[i]
		if smOrdi.team != "ordi" {
			continue
		}

		for j := range *stickmans {
			smPlayer := &(*stickmans)[j]
			if smPlayer.team != "player" {
				continue
			}

			if gameutil.RectColl(smPlayer.x+smPlayer.w, smPlayer.y, smPlayer.w, smPlayer.h, smOrdi.x-smOrdi.w, smOrdi.y, smOrdi.w, smOrdi.h) {
				smOrdi.Attack = true
				smPlayer.Attack = true
				smOrdi.attackCooldown = 90
				smPlayer.attackCooldown = 90
			}
		}
	}
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
	if sm.class == unit.Miner {
		// draw fluo rectangle (for the capacity)
		vector.StrokeRect(screenI, float32(sm.x+cam.X), float32(sm.y-70), 150, 70, 7, color.RGBA{255, 123, 25, 255}, true)
		// draw capacity (miner)
		gameutil.DrawText(fmt.Sprintf("%d/%d", sm.gold, sm.capacity), 50, screen.Widht, sm.x+cam.X, sm.y-50, 0, screenI, sm.clr, mplusSource)
	}
	// draw stickman
	vector.DrawFilledRect(screenI, float32(sm.x+cam.X), float32(sm.y), float32(sm.w), float32(sm.h), sm.clr, true)
}
