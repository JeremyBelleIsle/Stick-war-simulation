package unit

type Type int

const (
	None Type = iota
	Miner
	Attacker
	Archer
)

func (u Type) String() string {
	switch u {
	case Miner:
		return "miner"
	case Attacker:
		return "attacker"
	case Archer:
		return "archer"
	}

	return "unknown"
}
