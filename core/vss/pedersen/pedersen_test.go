package pedersen_test

import (
	"crypto/rand"
	"math/big"

	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/smpc-go/core/vss/pedersen"
)

var _ = Describe("Pedersen commitments", func() {

	const Trials = 50

	// perturbInt perturbs an element by a random (non-zero) number.
	perturbInt := func(ped Pedersen, n *big.Int) {
		var r *big.Int
		for {
			r, _ = rand.Int(rand.Reader, ped.SubgroupOrder())
			if r.Sign() != 0 {
				// Make sure that the perturbed element is different by
				// ensureing that adding the random number will change the
				// element.
				break
			}
		}
		n.Add(n, r)
	}

	// For each entry, q is chosen to be the largest prime less than 2^b for
	// various bit lengths b, and p is chosen to be the least prime such that q
	// divides  p - 1.
	table := []struct {
		p, q, g, h *big.Int
	}{
		{ // q ~ 8 bits
			big.NewInt(503),
			big.NewInt(251),
			big.NewInt(351), // q^176
			big.NewInt(8),   // q^248
		},
		{ // q ~ 16 bits
			big.NewInt(655211),
			big.NewInt(65521),
			big.NewInt(259323), // q^5387
			big.NewInt(617158), // q^26664
		},
		{ // q ~ 32 bits
			big.NewInt(8589934583),
			big.NewInt(4294967291),
			big.NewInt(592772542),  // q^3527860178
			big.NewInt(4799487786), // q^3349731522
		},
		{ // q ~ 64 bits
			big.NewInt(0).SetBytes([]byte{5, 255, 255, 255, 255, 255, 255, 254, 159}), // 110680464442257309343
			big.NewInt(0).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 197}),    // 18446744073709551557
			big.NewInt(0).SetBytes([]byte{1, 146, 9, 35, 48, 210, 219, 176, 237}),     // q^14829842343553222478
			big.NewInt(0).SetBytes([]byte{3, 131, 17, 181, 241, 96, 122, 74, 19}),     // q^11942310935152117490
		},
		{ // q ~ 128 bits
			big.NewInt(0).SetBytes([]byte{59, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 218, 189}), // 20416942015256307807802476445906092677821
			big.NewInt(0).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 97}),      // 340282366920938463463374607431768211297
			big.NewInt(0).SetBytes([]byte{48, 24, 73, 41, 213, 112, 168, 237, 147, 202, 198, 75, 175, 56, 240, 94, 176}),       // q^306940875728532667791655917128746165896
			big.NewInt(0).SetBytes([]byte{58, 11, 247, 139, 227, 253, 157, 242, 147, 134, 173, 116, 14, 142, 217, 73, 167}),    // q^30919156086154660785203002477014638086
		},
		{ // q ~ 256 bits
			big.NewInt(0).SetBytes([]byte{33, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 230, 231}), // 3936931034068750644401413490295388867011179478631779177341557856269046407751399
			big.NewInt(0).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 67}),      // 115792089237316195423570985008687907853269984665640564039457584007913129639747
			big.NewInt(0).SetBytes([]byte{11, 133, 57, 29, 12, 113, 44, 209, 192, 65, 252, 89, 247, 125, 23, 101, 33, 186, 56, 141, 173, 1, 232, 105, 250, 12, 160, 134, 138, 77, 229, 200, 143}),              // q^56152962986760337639031656357225957209699526761357461190857640909372757911263
			big.NewInt(0).SetBytes([]byte{29, 161, 146, 159, 245, 75, 165, 149, 2, 242, 7, 52, 228, 151, 187, 45, 149, 5, 17, 40, 201, 165, 189, 205, 45, 3, 189, 125, 94, 26, 56, 141, 34}),                   // q^55387777697025239469499746920689473451678181849429135457039521698992456692593
		},
		{ // q ~ 512 bits
			big.NewInt(0).SetBytes([]byte{4, 201, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 245, 91, 7}), // 16437972522109624044077754647800367352289702496046274281089086330002882700870168593559097889552623602347979058136631898346702260327446494754327653681458404103
			big.NewInt(0).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 253, 199}),      // 13407807929942597099574024998205846127479365820592393377723561443721764030073546976801874298166903427690031858186486050853753882811946569946433649006083527
			big.NewInt(0).SetBytes([]byte{34, 227, 113, 201, 108, 143, 123, 156, 129, 164, 87, 22, 70, 239, 97, 234, 180, 202, 246, 30, 194, 163, 113, 133, 37, 232, 77, 170, 159, 118, 95, 17, 116, 199, 55, 132, 156, 142, 126, 34, 61, 120, 229, 141, 197, 53, 201, 100, 33, 90, 229, 246, 218, 26, 18, 196, 62, 139, 136, 229, 124, 240, 86, 68, 83}),                       // q^2313223469486319304351557127107981179981811375871648778609488138158284033970434045960828878291659489151863079107926238730180855575203958252698604172393282
			big.NewInt(0).SetBytes([]byte{2, 213, 181, 190, 15, 59, 177, 180, 170, 72, 53, 6, 64, 134, 52, 149, 98, 2, 234, 118, 158, 179, 183, 190, 235, 244, 234, 42, 20, 137, 123, 53, 244, 223, 120, 41, 216, 5, 113, 216, 2, 181, 125, 40, 28, 118, 219, 199, 210, 79, 15, 194, 29, 3, 197, 237, 99, 208, 68, 99, 161, 107, 145, 12, 18, 43}),                              // q^2123660058119268294371654254614050015587085064079298406210266109396002109966413704619313534762901338642906339464900166903187025575969459600540403945834170
		},
		{ // q ~ 1024 bits
			big.NewInt(0).SetBytes([]byte{5, 169, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 253, 173, 71}), // 260665504555035806620749252664408586374606661946634453046473617678712379917976396542427292117490927230624165125813520368605245164780904102614628774427237480347763445466054953650638218421806836473206970280523258364815318489396571178137174010284199493886858240492650132447923506610795588642011266677955124998810951
			big.NewInt(0).SetBytes([]byte{255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 255, 151}),        // 179769313486231590772930519078902473361797697894230657273430081157732675805500963132708477322407536021120113879871393357658789768814416622492847430639474124377767893424865485276302219601246094119453082952085005768838150682342462881473913110540827237163350510684586298239947245938479716304835356329624224137111
			big.NewInt(0).SetBytes([]byte{1, 52, 192, 43, 8, 203, 47, 254, 188, 187, 95, 25, 110, 232, 113, 207, 85, 27, 187, 2, 51, 152, 144, 27, 220, 230, 132, 247, 246, 245, 113, 234, 48, 17, 197, 118, 202, 48, 119, 131, 147, 107, 104, 79, 221, 254, 74, 186, 207, 136, 20, 139, 165, 114, 219, 225, 210, 179, 162, 105, 56, 162, 193, 219, 209, 231, 102, 228, 5, 67, 156, 15, 22, 128, 24, 132, 189, 74, 20, 201, 93, 26, 43, 103, 185, 34, 170, 101, 175, 7, 7, 78, 129, 42, 133, 160, 215, 229, 243, 71, 202, 227, 189, 121, 232, 129, 59, 126, 176, 186, 177, 154, 145, 200, 148, 26, 17, 233, 149, 176, 153, 4, 3, 204, 242, 84, 192, 163, 132, 17}),                                                // q^76504997747089379379851214625213978887984955428323745694280408664704840586711756606456348617755031780834327963171139835068120850880276931679391188100919215516693562754769057716029964066946375140630991343845171906939443377157731505605857702213057849645661838880736934458114117087401472419720324377592216887124
			big.NewInt(0).SetBytes([]byte{5, 119, 199, 172, 247, 201, 206, 235, 191, 228, 185, 189, 106, 133, 159, 158, 182, 80, 209, 51, 164, 142, 159, 145, 21, 61, 93, 221, 53, 43, 242, 65, 192, 123, 26, 238, 231, 107, 240, 248, 122, 72, 23, 20, 239, 133, 184, 54, 222, 238, 38, 226, 30, 253, 120, 213, 47, 194, 240, 189, 147, 66, 147, 125, 24, 89, 156, 212, 54, 46, 159, 29, 201, 72, 140, 25, 44, 191, 235, 179, 132, 113, 227, 215, 168, 33, 48, 169, 244, 107, 24, 104, 97, 11, 242, 59, 0, 104, 71, 223, 221, 81, 11, 0, 25, 170, 237, 115, 86, 153, 112, 110, 97, 250, 71, 207, 253, 126, 9, 173, 120, 221, 179, 10, 94, 51, 167, 223, 3, 227}),                                                 // q^99698080805435051641410242084132915991510185641405271392786625720977706310966087784947122710719784813870715536932270384857937275577549106254692474198139325988180756021199621433268612322519320815873111449089532247920355054017684242152527838139624714558569578820724520518131956096960949189367403990735434862433
		},
	}

	for _, entry := range table {
		entry := entry

		Context("when using correctly constructed pedersen schemes", func() {
			ped, _ := New(entry.p, entry.q, entry.g, entry.h)
			// It("should construct without error", func(doneT Done) {
			// 	defer close(doneT)

			// 	Expect(err).To(BeNil())
			// })

			Context("when passing nil arguments to the verify function", func() {
				DescribeTable("an error is expected", func(s, t, commitment *big.Int, err error) {
					Expect(ped.Verify(s, t, commitment)).To(Equal(err))
				},
					Entry("when s is nil", nil, big.NewInt(1), big.NewInt(1), ErrNilArguments),
					Entry("when t is nil", big.NewInt(1), nil, big.NewInt(1), ErrNilArguments),
					Entry("when commitment is nil", big.NewInt(1), big.NewInt(1), nil, ErrNilArguments),
					Entry("when s, t are nil", nil, nil, big.NewInt(1), ErrNilArguments),
					Entry("when s, commitment are nil", nil, big.NewInt(1), nil, ErrNilArguments),
					Entry("when t, commitment are nil", big.NewInt(1), nil, nil, ErrNilArguments),
					Entry("when all arguments are nil", nil, nil, nil, ErrNilArguments),
				)
			})

			Context("when passing nil arguments to the commit function", func() {
				DescribeTable("it should return nil", func(s, t *big.Int) {
					Expect(ped.Commit(s, t)).To(BeNil())
				},
					Entry("when s is nil", nil, big.NewInt(1)),
					Entry("when t is nil", big.NewInt(1), nil),
					Entry("when all arguments are nil", nil, nil),
				)
			})

			Context("when verifying an incorrect commitment", func() {
				It("should return an error", func() {
					for i := 0; i < Trials; i++ {
						s, _ := rand.Int(rand.Reader, ped.SubgroupOrder())
						t, _ := rand.Int(rand.Reader, ped.SubgroupOrder())
						commitment := ped.Commit(s, t)

						// Make the commitment incorrect by perturbing it by a
						// random number.
						perturbInt(ped, commitment)

						Expect(ped.Verify(s, t, commitment)).To(Equal(ErrUnacceptableCommitment))
					}
				})
			})

			Context("when verifying a correct commitment", func() {
				It("should return a nil error", func() {
					for i := 0; i < Trials; i++ {
						s, _ := rand.Int(rand.Reader, ped.SubgroupOrder())
						t, _ := rand.Int(rand.Reader, ped.SubgroupOrder())
						commitment := ped.Commit(s, t)

						Expect(ped.Verify(s, t, commitment)).To(BeNil())
					}
				})
			})
		})

		Context("when using incorrectly constructed pedersen schemes", func() {
			DescribeTable("an error is expected", func(p, q, g, h *big.Int) {
				_, err := New(p, q, g, h)
				Expect(err).ToNot(BeNil())
			},
				Entry("when h is nil", entry.p, entry.q, entry.g, nil),
				Entry("when g is nil", entry.p, entry.q, nil, entry.h),
				Entry("when h, g are nil", entry.p, entry.q, nil, nil),
				Entry("when q is nil", entry.p, nil, entry.g, entry.h),
				Entry("when q, h are nil", entry.p, nil, entry.g, nil),
				Entry("when q, g are nil", entry.p, nil, nil, entry.h),
				Entry("when q, g, h are nil", entry.p, nil, nil, nil),
				Entry("when p is nil", nil, entry.q, entry.g, entry.h),
				Entry("when p, h are nil", nil, entry.q, entry.g, nil),
				Entry("when p, g are nil", nil, entry.q, nil, entry.h),
				Entry("when p, g, h are nil", nil, entry.q, nil, nil),
				Entry("when p, q are nil", nil, nil, entry.g, entry.h),
				Entry("when p, q, h are nil", nil, nil, entry.g, nil),
				Entry("when p, q, g are nil", nil, nil, nil, entry.h),
				Entry("when all arguments are nil", nil, nil, nil, nil),
			)

			Context("when picking p and q such that q does not divide p - 1", func() {
				ped, _ := New(entry.p, entry.q, entry.g, entry.h)

				It("should return an error", func() {
					for i := 0; i < Trials; i++ {
						perturbed := new(big.Int).Set(entry.p)
						perturbInt(ped, perturbed)
						_, err := New(perturbed, entry.q, entry.g, entry.h)

						Expect(err).ToNot(BeNil())
					}
				})
			})
		})
	}
})
