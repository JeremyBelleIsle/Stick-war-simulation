package ordi

import (
	"Sticks_War/store"
	"math/rand"
)

var Money = 5
var nextPurchase = "minner"

func ManageMoney(offers []store.Offer) string {
	// vérifier si on peut acheter ce stickMan
	purchase := nextPurchase
	for _, o := range offers {
		if o.TypeS == nextPurchase {
			if o.Cost > Money {
				return "not enough money"
			} else {
				break
			}
		}
	}
	switch rand.Intn(2) {
	case 0:
		nextPurchase = "minner"
	case 1:
		nextPurchase = "attacker"
	}

	if rand.Intn(7) == 0 {
		nextPurchase = "archer"
	}

	return purchase
}
