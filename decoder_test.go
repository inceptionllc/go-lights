package schedule_test

import (
	"image/color"

	. "github.com/inceptionllc/go-schedule"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {

	Describe("Schedule Parser", func() {

		It("should parse color codes", func() {
			c, err := ParseColorCode("#ABCDEF")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(c).Should(Equal(color.RGBA{0xab, 0xcd, 0xef, 0x00}))
			c, err = ParseColorCode("#abcdef")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(c).Should(Equal(color.RGBA{0xab, 0xcd, 0xef, 0x00}))
			c, err = ParseColorCode("#ABC")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(c).Should(Equal(color.RGBA{0xaa, 0xbb, 0xcc, 0x00}))
		})
	})
})
