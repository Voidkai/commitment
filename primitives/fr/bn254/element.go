package bn254

type Element [4]uint64

const (
	Words = 4   // number of Words for a field element
	Bits  = 254 // number of Bits for a field element
	Bytes = 32  // number of Bytes for a field element
)

// Field modulus q
const (
	q0 uint64 = 4891460686036598785
	q1 uint64 = 2896914383306846353
	q2 uint64 = 13281191951274694749
	q3 uint64 = 3486998266802970665
)

// One returns 1
func One() Element {
	var one Element
	one[0] = 12436184717236109307
	one[1] = 3962172157175319849
	one[2] = 7381016538464732718
	one[3] = 1011752739694698287
	return one
}
