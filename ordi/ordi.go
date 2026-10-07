package ordi

import (
	"Sticks_War/store"
	"Sticks_War/unit"
	"errors"
	"math/rand"
)

var Money = 10
var nextPurchase = unit.Miner

var ErrNotEnoughMoney = errors.New("not enough money")

func ManageMoney(offers []store.Offer) (unit.Type, error) {
	// vérifier si on peut acheter ce stickMan
	purchase := nextPurchase
	for _, o := range offers {
		if o.TypeS == nextPurchase {
			if o.Cost > Money {
				return unit.None, ErrNotEnoughMoney
			} else {
				break
			}
		}
	}
	switch rand.Intn(2) {
	case 0:
		nextPurchase = unit.Miner
	case 1:
		nextPurchase = unit.Attacker
	}

	if rand.Intn(7) == 0 {
		nextPurchase = unit.Archer
	}

	return purchase, nil
}
