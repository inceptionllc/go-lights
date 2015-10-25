package lights_test

import (
	"image/color"
	"time"

	"github.com/inceptionllc/go-lights"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {
	Describe("Slot", func() {
		It("should parse slot specs", func() {
			slot, err := lights.NewSlot("#F00,2s,1s")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(slot.Color).Should(Equal(color.RGBA{0xff, 0x00, 0x00, 0x00}))
			Ω(slot.Fade).Should(Equal(2 * time.Second))
			Ω(slot.Hold).Should(Equal(1 * time.Second))
			Ω(slot.Transition).Should(Equal("ease"))
		})
	})
	Describe("Pattern", func() {
		It("should parse pattern messages", func() {
			pattern, err := lights.NewPattern(":ab:|#F00,2s,1s|#FFF,2s,1s|#00F,2s,1s")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(pattern.ID).Should(Equal("ab"))
			Ω(pattern.Loops).Should(Equal(-1))
			Ω(pattern.Slots).Should(HaveLen(3))
			pattern, err = lights.NewPattern(":1:3|#F00,1s,2s|#0F0,1s,2s|#00f,1s,")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(pattern.ID).Should(Equal("1"))
			Ω(pattern.Loops).Should(Equal(3))
			Ω(pattern.Slots).Should(HaveLen(3))
			pattern, err = lights.NewPattern(":ab|#F00,2s,1s|#FFF,2s,1s|#00F,2s,1s")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(pattern.ID).Should(Equal("ab"))
			Ω(pattern.Loops).Should(Equal(-1))
			Ω(pattern.Slots).Should(HaveLen(3))
		})
	})
})
