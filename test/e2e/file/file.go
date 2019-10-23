package file

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/golang/glog"
	"github.com/kpn/pion/pkg/sts/secure_rand"
	"github.com/kpn/pion/test/e2e/framework"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = framework.DescribeOSS("File operations on a bucket: ", func() {
	f := framework.NewFramework()

	var targetBucket string
	var (
		accessKey = os.Getenv("U1_ACCESS_KEY")
		secretKey = os.Getenv("U1_SECRET_KEY")
	)
	BeforeEach(func() {
		var err error
		rndStr, err := secure_rand.SecureRandomString(8)
		Expect(err).To(BeNil())
		targetBucket = "bucket-" + rndStr

		// create a target bucket
		err = f.CreateBucket(accessKey, secretKey, targetBucket)
		Expect(err).To(BeNil())
	})

	Context("when having an existing bucket and a large file", func() {
		var tmpFilePath string
		BeforeEach(func() {
			var err error
			// create a large file
			tmpFilePath, err = createTmpFile(8 * 1024 * 1024)
			Expect(err).To(BeNil())
			glog.Infof("Created temp file %s of 8Gi", tmpFilePath)
		})

		It("user can upload and delete a large file to that bucket", func() {
			glog.Infof("Uploading '%s' to '%s'", tmpFilePath, targetBucket)
			location, err := f.UploadObject(accessKey, secretKey, tmpFilePath, targetBucket, "")
			Expect(err).To(BeNil())
			glog.Infof("Uploaded temp file %s to %s", tmpFilePath, location)

			Expect(f.DeleteObject(accessKey, secretKey, targetBucket, path.Base(tmpFilePath))).To(BeNil())
			glog.Infof("Deleted file %s from %s", tmpFilePath, location)
		})

		AfterEach(func() {
			Expect(os.Remove(tmpFilePath)).To(BeNil())
		})
	})

	Context("when having an existing bucket and a 0-length file", func() {
		var tmpFilePath string
		BeforeEach(func() {
			var err error
			tmpFilePath, err = createTmpFile(0)
			Expect(err).To(BeNil())
			glog.Infof("Created empty-content temp file '%s'", tmpFilePath)
		})

		It("user can upload and delete zero-length file to that bucket", func() {
			glog.Infof("Uploading '%s' to '%s'", tmpFilePath, targetBucket)
			location, err := f.UploadObject(accessKey, secretKey, tmpFilePath, targetBucket, "")
			Expect(err).To(BeNil())
			glog.Infof("Uploaded temp file %s to %s", tmpFilePath, location)

			Expect(f.DeleteObject(accessKey, secretKey, targetBucket, path.Base(tmpFilePath))).To(BeNil())
			glog.Infof("Deleted file %s from %s", tmpFilePath, location)
		})

		AfterEach(func() {
			Expect(os.Remove(tmpFilePath)).To(BeNil())
		})
	})

	AfterEach(func() {
		Expect(f.DeleteBucket(accessKey, secretKey, targetBucket)).To(BeNil())
	})
})

func createTmpFile(fileSizeInBytes int64) (filePath string, err error) {
	f, err := ioutil.TempFile("", "pion-e2e")
	if err != nil {
		return
	}

	err = f.Truncate(fileSizeInBytes)
	if err != nil {
		return
	}
	filePath = f.Name()
	return
}
