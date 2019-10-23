package mock

import (
	"github.com/kpn/pion/pkg/pion/shared/multi-tenancy/model"
	"github.com/stretchr/testify/mock"
)

type mockBucketsWatcher struct {
	mock.Mock
}

func NewBucketsWatcher() *mockBucketsWatcher {
	return &mockBucketsWatcher{}
}

func (m *mockBucketsWatcher) GetBucket(name string) *model.Bucket {
	args := m.Called(name)
	return args.Get(0).(*model.Bucket)
}

func (m *mockBucketsWatcher) Watch() {
	m.Called()
}
