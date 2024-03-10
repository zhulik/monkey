package object_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	obj "github.com/zhulik/monkey/evaluator/object"
)

var _ = Describe("Integer", func() {
	integer := obj.New[obj.Integer](123)

	Describe(".TypeName", func() {
		It("returns type name", func() {
			Expect(integer.TypeName()).To(Equal("Integer"))
		})
	})

	Describe(".Inspect", func() {
		It("returns string representation", func() {
			Expect(integer.Inspect()).To(Equal("123"))
		})
	})
})
