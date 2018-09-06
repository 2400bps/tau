package vss_test

import (
	"math/big"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/smpc-go/core/vss"
	"github.com/republicprotocol/smpc-go/core/vss/algebra"
	"github.com/republicprotocol/smpc-go/core/vss/pedersen"
	"github.com/republicprotocol/smpc-go/core/vss/shamir"
)

var _ = Describe("Verifiable secret sharing", func() {

	const Trials = 10

	// For each entry, q is chosen to be the largest prime less than 2^b for
	// various bit lengths b, and p is chosen to be the least prime such that q
	// divides  p - 1.
	table := []struct {
		p, q, g, h *big.Int
	}{
		{ // q ~ 8 bits
			big.NewInt(503),
			big.NewInt(251),
			big.NewInt(351),
			big.NewInt(8),
		},
		{ // q ~ 16 bits
			big.NewInt(655211),
			big.NewInt(65521),
			big.NewInt(259323),
			big.NewInt(617158),
		},
		{ // q ~ 32 bits
			big.NewInt(8589934583),
			big.NewInt(4294967291),
			big.NewInt(592772542),
			big.NewInt(4799487786),
		},
		{ // q ~ 64 bits
			big.NewInt(0).SetBytes([]byte{5, 255, 255, 255, 255, 255, 255, 254, 159}), // 110680464442257309343
			big.NewInt(0).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 197}),    // 18446744073709551557
			big.NewInt(0).SetBytes([]byte{2, 143, 225, 91, 153, 86, 87, 209, 90}),     // 47261156678739415386
			big.NewInt(0).SetBytes([]byte{119, 242, 71, 138, 30, 234, 113, 106}),      // 8643049293427143018
		},
		{ // q ~ 128 bits
			big.NewInt(0).SetBytes([]byte{59, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 218, 189}), // 20416942015256307807802476445906092677821
			big.NewInt(0).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 97}),      // 340282366920938463463374607431768211297
			big.NewInt(0).SetBytes([]byte{46, 217, 224, 127, 231, 86, 139, 236, 205, 5, 30, 157, 110, 39, 97, 111, 203}),       // 15942597022139317939475358237557751115723
			big.NewInt(0).SetBytes([]byte{1, 165, 222, 96, 116, 29, 93, 210, 13, 91, 221, 196, 71, 130, 67, 247, 62}),          // 560759632438921603770492308279030708030
		},
		{ // q ~ 256 bits
			big.NewInt(0).SetBytes([]byte{33, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 230, 231}), // 3936931034068750644401413490295388867011179478631779177341557856269046407751399
			big.NewInt(0).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 67}),      // 115792089237316195423570985008687907853269984665640564039457584007913129639747
			big.NewInt(0).SetBytes([]byte{23, 207, 104, 159, 54, 244, 165, 239, 185, 226, 64, 252, 177, 68, 64, 86, 43, 223, 127, 48, 34, 19, 120, 201, 91, 190, 153, 241, 82, 8, 108, 93, 119}),               // 2757031663069016310286585750987615067052580546540373697271521881678374018768247
			big.NewInt(0).SetBytes([]byte{15, 227, 1, 107, 16, 126, 133, 122, 215, 11, 100, 192, 213, 253, 7, 35, 83, 5, 22, 245, 130, 103, 208, 224, 96, 167, 193, 252, 3, 87, 73, 234, 255}),                 // 1839558860966751692336441212815126320032775739119224251233662653792562960788223
		},
		{ // q ~ 512 bits
			big.NewInt(0).SetBytes([]byte{4, 201, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 245, 91, 7}), // 16437972522109624044077754647800367352289702496046274281089086330002882700870168593559097889552623602347979058136631898346702260327446494754327653681458404103
			big.NewInt(0).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 253, 199}),      // 13407807929942597099574024998205846127479365820592393377723561443721764030073546976801874298166903427690031858186486050853753882811946569946433649006083527
			big.NewInt(0).SetBytes([]byte{2, 126, 192, 45, 136, 9, 71, 62, 58, 2, 223, 0, 89, 208, 66, 75, 201, 19, 200, 155, 134, 80, 138, 46, 54, 138, 227, 163, 77, 251, 91, 134, 52, 26, 8, 140, 28, 146, 21, 176, 35, 119, 110, 4, 31, 128, 228, 174, 239, 118, 71, 162, 17, 45, 203, 7, 96, 67, 192, 220, 50, 179, 111, 154, 76, 47}),                                     // 8564246630377680631564400639972671137439956622643029828960025335891206374791129288542024756892236314055627529701878434918124899605710839776272675378392878127
			big.NewInt(0).SetBytes([]byte{1, 155, 85, 23, 121, 47, 206, 89, 217, 195, 239, 219, 208, 221, 155, 7, 77, 185, 23, 202, 190, 184, 207, 0, 182, 204, 33, 223, 172, 100, 135, 241, 117, 23, 131, 147, 107, 168, 97, 100, 61, 163, 61, 174, 189, 107, 110, 45, 40, 232, 94, 213, 185, 147, 94, 158, 32, 193, 144, 208, 20, 78, 214, 220, 127, 23}),                     // 5515065672780666600016583677548739776362501856726989816378508129242471762114311340533829327584004624008535874878447467980995509656465881362270331691809406743
		},
		{ // q ~ 1024 bits
			big.NewInt(0).SetBytes([]byte{5, 169, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 253, 173, 71}), // 260665504555035806620749252664408586374606661946634453046473617678712379917976396542427292117490927230624165125813520368605245164780904102614628774427237480347763445466054953650638218421806836473206970280523258364815318489396571178137174010284199493886858240492650132447923506610795588642011266677955124998810951
			big.NewInt(0).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 151}),        // 179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111
			big.NewInt(0).SetBytes([]byte{4, 12, 171, 80, 59, 112, 87, 191, 45, 62, 114, 36, 14, 113, 114, 216, 228, 169, 251, 61, 49, 179, 111, 253, 229, 214, 103, 26, 42, 255, 120, 112, 91, 242, 126, 103, 164, 80, 138, 186, 91, 217, 71, 191, 28, 125, 23, 247, 85, 136, 165, 191, 215, 224, 57, 4, 14, 28, 233, 32, 204, 235, 169, 49, 214, 94, 158, 180, 42, 71, 205, 201, 3, 201, 237, 155, 99, 123, 213, 98, 53, 252, 36, 3, 96, 43, 176, 97, 179, 0, 228, 249, 87, 250, 233, 73, 184, 180, 230, 3, 173, 196, 176, 250, 126, 16, 21, 219, 6, 119, 67, 29, 71, 131, 121, 183, 134, 182, 42, 160, 0, 154, 251, 87, 219, 68, 53, 238, 182, 187}),                                                           // 186361309137208710974894049867971390884596024321835358104739536820286472561461673760624224024865922331260979934707584105746546371103170986057625555541509748960323162179180607923216896019581463196390584432440066404363336674053782358131452589782688367820448015655355858056800340479792388452041371605165341934532283
			big.NewInt(0).SetBytes([]byte{5, 115, 143, 55, 30, 2, 255, 139, 83, 33, 93, 221, 74, 161, 22, 16, 149, 252, 11, 23, 243, 81, 72, 37, 230, 237, 175, 209, 97, 27, 23, 169, 190, 106, 12, 250, 100, 13, 69, 183, 85, 167, 71, 168, 240, 11, 144, 242, 129, 245, 154, 14, 116, 39, 23, 184, 156, 209, 142, 184, 145, 85, 104, 125, 233, 56, 81, 172, 162, 184, 52, 166, 148, 54, 27, 140, 106, 141, 53, 156, 224, 187, 132, 149, 79, 215, 86, 60, 123, 37, 126, 193, 187, 32, 149, 34, 219, 34, 17, 23, 15, 253, 60, 133, 150, 32, 78, 101, 215, 39, 41, 107, 12, 125, 219, 183, 135, 51, 229, 240, 8, 80, 216, 77, 36, 221, 102, 186, 61, 35}),                                                          // 250878761518238235686252428989631348547141621723716412508485472897013986492003239743842235204997618869447912693568041917266321874458142230126810698260150923539891354843218189254127822287513674484995332415052827826236946280584908050185003422831150586546376125046671376629115512822481674295764992544448758062660899
		},
	}

	for _, entry := range table {
		entry := entry

		Context("when creating verifiable shares", func() {
			g := algebra.NewFpElement(entry.g, entry.p)
			h := algebra.NewFpElement(entry.h, entry.p)
			ped := pedersen.New(g, h)

			It("should panic when there are no commitments", func() {
				field := algebra.NewField(entry.q)

				for i := 0; i < Trials; i++ {
					secret := field.Random()
					n := uint64(24)
					k := uint64(16)

					verifiableShares := Share(&ped, secret, n, k)
					for _, share := range verifiableShares {
						share.SetCommitments(make([]algebra.FpElement, 0))
						Expect(func() { Verify(&ped, share) }).To(Panic())
					}
				}
			})

			It("should verify correct shares", func() {
				field := algebra.NewField(entry.q)

				for i := 0; i < Trials; i++ {
					secret := field.Random()
					n := uint64(24)
					k := uint64(16)

					verifiableShares := Share(&ped, secret, n, k)
					for _, share := range verifiableShares {
						Expect(Verify(&ped, share)).To(BeTrue())
					}

					// Check that the secret can be correctly reconstructed
					shares := make(shamir.Shares, 24)
					for i := range shares {
						shares[i] = verifiableShares[i].Share()
					}

					for i := uint64(0); i < n-k; i++ {
						val, _ := shamir.Join(shares[i : i+k])
						Expect(val.Eq(secret)).To(BeTrue())
					}
				}
			})

			It("should catch incorrect shares", func() {
				field := algebra.NewField(entry.q)

				for i := 0; i < Trials; i++ {
					secret := field.Random()
					r := field.Random()
					if r.IsZero() {
						continue
					}
					n := uint64(24)
					k := uint64(16)
					verifiableShares := Share(&ped, secret, n, k)

					for _, share := range verifiableShares {
						innerShare := share.Share()
						rShare := shamir.New(innerShare.Index(), r)
						share.SetShare(innerShare.Add(rShare))

						Expect(Verify(&ped, share)).To(BeFalse())
					}
				}
			})

			Specify("addition should correspond to addition of the underlying secret", func() {
				field := algebra.NewField(entry.q)

				for i := 0; i < Trials; i++ {
					secretA := field.Random()
					secretB := field.Random()
					n := uint64(24)
					k := uint64(16)

					sharesA := Share(&ped, secretA, n, k)
					sharesB := Share(&ped, secretB, n, k)
					addedShares := make(VShares, n)
					for i := range addedShares {
						addedShares[i] = sharesA[i].Add(&sharesB[i])
						Expect(Verify(&ped, addedShares[i])).To(BeTrue())
					}

					// Check that the secret can be correctly reconstructed
					shares := make(shamir.Shares, 24)
					for i := range shares {
						shares[i] = addedShares[i].Share()

					}

					for i := uint64(0); i < n-k; i++ {
						val, _ := shamir.Join(shares[i : i+k])
						Expect(val.Eq(secretA.Add(secretB))).To(BeTrue())
					}
				}
			})
		})
	}
})
