package algebra_test

import (
	"crypto/rand"
	"math/big"

	. "github.com/onsi/ginkgo/extensions/table"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	. "github.com/republicprotocol/smpc-go/core/vss/algebra"
)

var _ = Describe("Finite field Fp", func() {
	const Trials = 100

	Context("when constructing a field", func() {
		Context("with a prime number", func() {
			DescribeTable("no panic is expected", func(prime *big.Int) {
				Expect(func() { NewField(prime) }).ToNot(Panic())
			},
				PrimeEntries...,
			)
		})

		Context("with a composite number", func() {
			DescribeTable("a panic is expected", func(composite *big.Int) {
				Expect(func() { NewField(composite) }).To(Panic())
			},
				CompositeEntries...,
			)
		})

		Context("with a negative number", func() {
			It("should panic", func() {
				for i := 0; i < Trials; i++ {
					negative, err := rand.Int(rand.Reader, big.NewInt(0).SetUint64(^uint64(0)))
					Expect(err).To(BeNil())

					negative.Neg(negative)
					Expect(func() { NewField(negative) }).To(Panic())
				}
			})
		})
	})

	Context("when creating a field element from a field and a value", func() {
		DescribeTable("it should panic when the value is not in the field", func(prime *big.Int) {
			field := NewField(prime)
			for i := 0; i < Trials; i++ {
				value := RandomNotInField(prime)
				Expect(func() { field.NewInField(value) }).To(Panic())
			}
		},
			PrimeEntries...,
		)

		DescribeTable("it should succeed when the value is in the field", func(prime *big.Int) {
			field := NewField(prime)
			for i := 0; i < Trials; i++ {
				value, _ := rand.Int(rand.Reader, prime)
				Expect(func() { field.NewInField(value) }).ToNot(Panic())
			}
		},
			PrimeEntries...,
		)
	})

	Context("when comparing two fields", func() {
		DescribeTable("it should return false when the fields are defined by different primes", func(prime *big.Int) {
			field := NewField(prime)
			otherField := NewField(big.NewInt(7))
			Expect(field.Eq(otherField)).To(BeFalse())
		},
			PrimeEntries...,
		)

		DescribeTable("it should return true when the fields are defined by the same prime", func(prime *big.Int) {
			field := NewField(prime)
			otherField := NewField(prime)
			Expect(field.Eq(otherField)).To(BeTrue())
		},
			PrimeEntries...,
		)
	})

	Context("when checking if an integer is an element of the field", func() {
		prime, _ := big.NewInt(0).SetString("11415648579556416673", 10)
		field := NewField(prime)

		Context("when the integer is too big", func() {
			It("should return false", func() {
				for i := 0; i < Trials; i++ {
					toobig, err := rand.Int(rand.Reader, big.NewInt(0).SetUint64(^uint64(0)))
					Expect(err).To(BeNil())

					toobig.Add(toobig, prime)
					Expect(field.Contains(toobig)).To(BeFalse())
				}
			})
		})

		Context("when the integer is negative", func() {
			It("should return false", func() {
				for i := 0; i < Trials; i++ {
					negative, err := rand.Int(rand.Reader, big.NewInt(0).SetUint64(^uint64(0)))
					Expect(err).To(BeNil())

					negative.Neg(negative)
					Expect(field.Contains(negative)).To(BeFalse())
				}
			})
		})

		Context("when the integer is in the field", func() {
			It("should return false", func() {
				for i := 0; i < Trials; i++ {
					correct, err := rand.Int(rand.Reader, prime)
					Expect(err).To(BeNil())

					Expect(field.Contains(correct)).To(BeTrue())
				}
			})
		})
	})

	Context("when creating a random field element", func() {
		DescribeTable("no panic is expected", func(prime *big.Int) {
			field := NewField(prime)

			for i := 0; i < Trials; i++ {
				Expect(func() { field.Random() }).ToNot(Panic())
			}
		},
			PrimeEntries...,
		)
	})
})
