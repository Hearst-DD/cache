package mocks

import "github.com/Hearst-DD/htv-server/cache"
import "github.com/stretchr/testify/mock"

import "time"

type Cache struct {
	mock.Mock
}

func (m *Cache) Put(key string, value interface{}, ttl time.Duration) {
	m.Called(key, value, ttl)
}
func (m *Cache) Get(key string) (interface{}, cache.Result) {
	ret := m.Called(key)

	var r0 interface{}
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(interface{})
	}
	r1 := ret.Get(1).(cache.Result)

	return r0, r1
}
func (m *Cache) Size() int {
	ret := m.Called()

	r0 := ret.Get(0).(int)

	return r0
}
