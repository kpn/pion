package iam

import (
	"github.com/kpn/pion/test/e2e/framework"
	. "github.com/onsi/ginkgo"
)

var _ = framework.DescribeOSS("Bucket policy test-cases: ", func() {
	framework.NewFramework()

	Context("when a customer account has a bucket", func() {
		It("Anyone outside the customer account cannot access bucket", func() {

		})
	})

	Context("when a bucket does not have attached ACLs", func() {

		It("A person having customer's user-role can read and list bucket", func() {

		})

		It("A person having customer's editor-role can read, list, write, update and publish bucket", func() {

		})

		It("A person having customer's admin-role can read, list, write and update bucket", func() {

		})
	})

	Context("when a bucket has attached ACLs", func() {
		It("a person does not have updating role but denied ACLs can update bucket", func() {

		})

		It("a person has updating role but denied ACLs still can update bucket", func() {

		})
	})
})
