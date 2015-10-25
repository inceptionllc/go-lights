package lights_test

import (
	"github.com/inceptionllc/go-lights"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {

	Describe("Message", func() {
		Describe("Execute Commands", func() {
			It("should parse color messages", func() {
				cmd, err := lights.NewCommand("!#F00")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd.Action).Should(Equal("execute"))
				Ω(cmd.Type).Should(Equal("color"))
				Ω(cmd.ID).Should(Equal("#F00"))
			})
			It("should parse pattern messages", func() {
				cmd, err := lights.NewCommand("!:ab:|#F00,2,1|#FFF,2,1|#00F,2,1")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd.Action).Should(Equal("execute"))
				Ω(cmd.Type).Should(Equal("pattern"))
				Ω(cmd.ID).Should(Equal("ab"))
				cmd, err = lights.NewCommand("!:1:3|#F00,1,2|#0F0,1,2|#00f,1,")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd.Action).Should(Equal("execute"))
				Ω(cmd.Type).Should(Equal("pattern"))
				Ω(cmd.ID).Should(Equal("1"))
				cmd, err = lights.NewCommand("!:ab|#F00,2,1|#FFF,2,1|#00F,2,1")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd.Action).Should(Equal("execute"))
				Ω(cmd.Type).Should(Equal("pattern"))
				Ω(cmd.ID).Should(Equal("ab"))
			})
			It("should parse schedule messages", func() {
				cmd, err := lights.NewCommand("!~4|||0 30 * * * *|#000|")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd.Action).Should(Equal("execute"))
				Ω(cmd.Type).Should(Equal("schedule"))
				Ω(cmd.ID).Should(Equal("4"))
				cmd, err = lights.NewCommand("!~8|2015-07-04|2015-07-05|0 0 20 * * *|:ab|1")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd.Action).Should(Equal("execute"))
				Ω(cmd.Type).Should(Equal("schedule"))
				Ω(cmd.ID).Should(Equal("8"))
				cmd, err = lights.NewCommand("!~9|2015-07-04|2015-07-05|0 0 23 * * *|#00|1")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd.Action).Should(Equal("execute"))
				Ω(cmd.Type).Should(Equal("schedule"))
				Ω(cmd.ID).Should(Equal("9"))
			})
			It("should parse scene messages", func() {
				cmd, err := lights.NewCommand("!^32|#F00,2|1|3|ab")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd.Action).Should(Equal("execute"))
				Ω(cmd.Type).Should(Equal("scene"))
				Ω(cmd.ID).Should(Equal("32"))
				cmd, err = lights.NewCommand("!^2|#00F|4|56")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd.Action).Should(Equal("execute"))
				Ω(cmd.Type).Should(Equal("scene"))
				Ω(cmd.ID).Should(Equal("2"))
			})
		})
		Describe("Query Commands", func() {
			It("should parse color messages", func() {
				cmd, err := lights.NewCommand("?-version")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd.Action).Should(Equal("query"))
				Ω(cmd.Type).Should(Equal("property"))
				Ω(cmd.ID).Should(Equal("version"))
			})
		})
		Describe("Add Commands", func() {
			It("should parse color messages", func() {
				cmd, err := lights.NewCommand("+:ab:|#F00,2,1|#FFF,2,1|#00F,2,1")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd.Action).Should(Equal("add"))
				Ω(cmd.Type).Should(Equal("pattern"))
				Ω(cmd.ID).Should(Equal("ab"))
			})
		})
		Describe("Remove Commands", func() {
			It("should parse color messages", func() {
				cmd, err := lights.NewCommand("-:ab:|#F00,2,1|#FFF,2,1|#00F,2,1")
				Ω(err).ShouldNot(HaveOccurred())
				Ω(cmd.Action).Should(Equal("remove"))
				Ω(cmd.Type).Should(Equal("pattern"))
				Ω(cmd.ID).Should(Equal("ab"))
			})
		})
	})
})
