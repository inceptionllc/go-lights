package lights_test

import (
	. "github.com/inceptionllc/go-lights"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {

	Describe("Schedule Parser", func() {

		It("should parse color codes", func() {
			r, err := ParseHex("AB")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(r).Should(Equal(0xab))
		})
	})
})
