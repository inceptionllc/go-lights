package schedule_test

import (
	. "github.com/inceptionllc/go-schedule"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {

	Describe("Schedule", func() {

		It("should retrieve patterns by name", func() {
			c := NewConfig()
			p := &Pattern{ID: "pid", Name: "pattern1"}
			c.AddPattern(p)
			Ω(len(c.Patterns)).Should(Equal(1))
			_, ok := c.PatternByName("pattern0")
			Ω(ok).Should(BeFalse())
			found, ok := c.PatternByName("pattern1")
			Ω(ok).Should(BeTrue())
			Ω(found).Should(Equal(p))
		})
	})
})
