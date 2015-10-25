package lights_test

import (
	"github.com/inceptionllc/go-lights"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Core", func() {

	Describe("Store tests", func() {
		It("should support mock stores", func() {
			s := &lights.MockStore{}
			found, err := s.Read("foo", "bar")
			Ω(err).Should(HaveOccurred())
			err = s.Write("foo", "bar", "baz")
			Ω(err).ShouldNot(HaveOccurred())
			found, err = s.Read("foo", "bar")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(found).Should(Equal("baz"))
			loaded, err := s.Load("foo")
			Ω(err).ShouldNot(HaveOccurred())
			Ω(loaded).Should(Equal([]string{"baz"}))
			err = s.Remove("foo", "bar")
			Ω(err).ShouldNot(HaveOccurred())
			found, err = s.Read("foo", "bar")
			Ω(err).Should(HaveOccurred())
		})
	})
})
