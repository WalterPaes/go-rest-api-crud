// Code generated by mockery v2.41.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/WalterPaes/go-rest-api-crud/internal/domain"
	mock "github.com/stretchr/testify/mock"

	resterrors "github.com/WalterPaes/go-rest-api-crud/pkg/rest_errors"
)

// UserRepository is an autogenerated mock type for the UserRepository type
type UserRepository struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: _a0, _a1
func (_m *UserRepository) CreateUser(_a0 context.Context, _a1 *domain.User) (*domain.User, *resterrors.RestErr) {
	ret := _m.Called(_a0, _a1)

	if len(ret) == 0 {
		panic("no return value specified for CreateUser")
	}

	var r0 *domain.User
	var r1 *resterrors.RestErr
	if rf, ok := ret.Get(0).(func(context.Context, *domain.User) (*domain.User, *resterrors.RestErr)); ok {
		return rf(_a0, _a1)
	}
	if rf, ok := ret.Get(0).(func(context.Context, *domain.User) *domain.User); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, *domain.User) *resterrors.RestErr); ok {
		r1 = rf(_a0, _a1)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*resterrors.RestErr)
		}
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: ctx, userID
func (_m *UserRepository) DeleteUser(ctx context.Context, userID string) *resterrors.RestErr {
	ret := _m.Called(ctx, userID)

	if len(ret) == 0 {
		panic("no return value specified for DeleteUser")
	}

	var r0 *resterrors.RestErr
	if rf, ok := ret.Get(0).(func(context.Context, string) *resterrors.RestErr); ok {
		r0 = rf(ctx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*resterrors.RestErr)
		}
	}

	return r0
}

// FindAll provides a mock function with given fields: parentCtx, itemsPerPage, currentPage
func (_m *UserRepository) FindAll(parentCtx context.Context, itemsPerPage int, currentPage int) ([]*domain.User, *resterrors.RestErr) {
	ret := _m.Called(parentCtx, itemsPerPage, currentPage)

	if len(ret) == 0 {
		panic("no return value specified for FindAll")
	}

	var r0 []*domain.User
	var r1 *resterrors.RestErr
	if rf, ok := ret.Get(0).(func(context.Context, int, int) ([]*domain.User, *resterrors.RestErr)); ok {
		return rf(parentCtx, itemsPerPage, currentPage)
	}
	if rf, ok := ret.Get(0).(func(context.Context, int, int) []*domain.User); ok {
		r0 = rf(parentCtx, itemsPerPage, currentPage)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, int, int) *resterrors.RestErr); ok {
		r1 = rf(parentCtx, itemsPerPage, currentPage)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*resterrors.RestErr)
		}
	}

	return r0, r1
}

// FindUserByEmail provides a mock function with given fields: parentCtx, email
func (_m *UserRepository) FindUserByEmail(parentCtx context.Context, email string) (*domain.User, *resterrors.RestErr) {
	ret := _m.Called(parentCtx, email)

	if len(ret) == 0 {
		panic("no return value specified for FindUserByEmail")
	}

	var r0 *domain.User
	var r1 *resterrors.RestErr
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.User, *resterrors.RestErr)); ok {
		return rf(parentCtx, email)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.User); ok {
		r0 = rf(parentCtx, email)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *resterrors.RestErr); ok {
		r1 = rf(parentCtx, email)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*resterrors.RestErr)
		}
	}

	return r0, r1
}

// FindUserById provides a mock function with given fields: parentCtx, userID
func (_m *UserRepository) FindUserById(parentCtx context.Context, userID string) (*domain.User, *resterrors.RestErr) {
	ret := _m.Called(parentCtx, userID)

	if len(ret) == 0 {
		panic("no return value specified for FindUserById")
	}

	var r0 *domain.User
	var r1 *resterrors.RestErr
	if rf, ok := ret.Get(0).(func(context.Context, string) (*domain.User, *resterrors.RestErr)); ok {
		return rf(parentCtx, userID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) *domain.User); ok {
		r0 = rf(parentCtx, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) *resterrors.RestErr); ok {
		r1 = rf(parentCtx, userID)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*resterrors.RestErr)
		}
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: ctx, userID, user
func (_m *UserRepository) UpdateUser(ctx context.Context, userID string, user *domain.User) (*domain.User, *resterrors.RestErr) {
	ret := _m.Called(ctx, userID, user)

	if len(ret) == 0 {
		panic("no return value specified for UpdateUser")
	}

	var r0 *domain.User
	var r1 *resterrors.RestErr
	if rf, ok := ret.Get(0).(func(context.Context, string, *domain.User) (*domain.User, *resterrors.RestErr)); ok {
		return rf(ctx, userID, user)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *domain.User) *domain.User); ok {
		r0 = rf(ctx, userID, user)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.User)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *domain.User) *resterrors.RestErr); ok {
		r1 = rf(ctx, userID, user)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*resterrors.RestErr)
		}
	}

	return r0, r1
}

// NewUserRepository creates a new instance of UserRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewUserRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *UserRepository {
	mock := &UserRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}