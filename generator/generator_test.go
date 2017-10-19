package generator_test

import (
	. "github.com/trusch/aliasd/generator"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generator", func() {
	It("should generate names", func() {
		n1 := Generate()
		n2 := Generate()
		Expect(n1).NotTo(BeEmpty())
		Expect(n2).NotTo(BeEmpty())
	})

	It("should generate the same names when inited with the same seed", func() {
		Seed(0)
		n1 := Generate()
		Seed(0)
		n2 := Generate()
		Expect(n2).To(Equal(n1))
	})
})
