package bucket

import (
	"os"

	"github.com/golang/glog"
	"github.com/kpn/pion/test/e2e/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = framework.DescribeOSS("Bucket operations for customers: ", func() {
	f := framework.NewFramework()

	testCases := []struct {
		AccessKey string
		SecretKey string
		Error     error
		Buckets   []interface{}
	}{
		{
			AccessKey: os.Getenv("U1_ACCESS_KEY"),
			SecretKey: os.Getenv("U1_SECRET_KEY"),
			Error:     nil,
			Buckets:   []interface{}{"c1-b1", "c1-b2"},
		},
		{
			AccessKey: os.Getenv("U2_ACCESS_KEY"),
			SecretKey: os.Getenv("U2_SECRET_KEY"),
			Error:     nil,
			Buckets:   []interface{}{"c2-b1", "c2-b2"},
		},
	}

	Context("when creating different buckets for different customers", func() {
		BeforeEach(func() {
			// TODO setup environments:
			// - Two Customers 'c1', 'c2', each has buckets 'c1-b1', 'c1-b2', 'c2-b1', 'c2-b2'

			// - Two users: u1 of c1, u2 of c2 with their access/secret tokens

			// create buckets for customers
			for _, c := range testCases {
				Expect(c.AccessKey).NotTo(BeEmpty())
				Expect(c.SecretKey).NotTo(BeEmpty())
				glog.Infof("Creating buckets %v", c.Buckets)
				for _, bkt := range c.Buckets {
					err := f.CreateBucket(c.AccessKey, c.SecretKey, bkt.(string))
					Expect(err).To(BeNil())
				}
			}
		})

		It("each customer can only list his own buckets", func() {
			for _, c := range testCases {
				Expect(c.AccessKey).NotTo(BeEmpty())
				Expect(c.SecretKey).NotTo(BeEmpty())
				actualBuckets, err := f.ListBuckets(c.AccessKey, c.SecretKey)
				Expect(err).To(BeNil())
				Expect(actualBuckets).To(HaveLen(len(c.Buckets)))
				Expect(actualBuckets).To(ConsistOf(c.Buckets...))
			}
		})

		AfterEach(func() {
			// clean up
			for _, c := range testCases {
				Expect(c.AccessKey).NotTo(BeEmpty())
				Expect(c.SecretKey).NotTo(BeEmpty())
				for _, bkt := range c.Buckets {
					err := f.DeleteBucket(c.AccessKey, c.SecretKey, bkt.(string))
					Expect(err).To(BeNil())
				}
			}
		})
	})
})
