package commitment

type Commitment interface {
	Setup()
	Commit()
	Reveal()
}
