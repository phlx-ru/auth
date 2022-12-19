// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package biz

import (
	"auth/ent"
	"context"
	"sync"
	"time"
)

// Ensure, that sessionRepoMock does implement sessionRepo.
// If this is not the case, regenerate this file with moq.
var _ sessionRepo = &sessionRepoMock{}

// sessionRepoMock is a mock implementation of sessionRepo.
//
//	func TestSomethingThatUsessessionRepo(t *testing.T) {
//
//		// make and configure a mocked sessionRepo
//		mockedsessionRepo := &sessionRepoMock{
//			CreateFunc: func(contextMoqParam context.Context, session *ent.Session) (*ent.Session, error) {
//				panic("mock out the Create method")
//			},
//			FindByTokenFunc: func(ctx context.Context, token string) (*ent.Session, error) {
//				panic("mock out the FindByToken method")
//			},
//		}
//
//		// use mockedsessionRepo in code that requires sessionRepo
//		// and then make assertions.
//
//	}
type sessionRepoMock struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(contextMoqParam context.Context, session *ent.Session) (*ent.Session, error)

	// FindByTokenFunc mocks the FindByToken method.
	FindByTokenFunc func(ctx context.Context, token string) (*ent.Session, error)

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Session is the session argument value.
			Session *ent.Session
		}
		// FindByToken holds details about calls to the FindByToken method.
		FindByToken []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Token is the token argument value.
			Token string
		}
	}
	lockCreate      sync.RWMutex
	lockFindByToken sync.RWMutex
}

// Create calls CreateFunc.
func (mock *sessionRepoMock) Create(contextMoqParam context.Context, session *ent.Session) (*ent.Session, error) {
	if mock.CreateFunc == nil {
		panic("sessionRepoMock.CreateFunc: method is nil but sessionRepo.Create was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Session         *ent.Session
	}{
		ContextMoqParam: contextMoqParam,
		Session:         session,
	}
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(contextMoqParam, session)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//
//	len(mockedsessionRepo.CreateCalls())
func (mock *sessionRepoMock) CreateCalls() []struct {
	ContextMoqParam context.Context
	Session         *ent.Session
} {
	var calls []struct {
		ContextMoqParam context.Context
		Session         *ent.Session
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
	return calls
}

// FindByToken calls FindByTokenFunc.
func (mock *sessionRepoMock) FindByToken(ctx context.Context, token string) (*ent.Session, error) {
	if mock.FindByTokenFunc == nil {
		panic("sessionRepoMock.FindByTokenFunc: method is nil but sessionRepo.FindByToken was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Token string
	}{
		Ctx:   ctx,
		Token: token,
	}
	mock.lockFindByToken.Lock()
	mock.calls.FindByToken = append(mock.calls.FindByToken, callInfo)
	mock.lockFindByToken.Unlock()
	return mock.FindByTokenFunc(ctx, token)
}

// FindByTokenCalls gets all the calls that were made to FindByToken.
// Check the length with:
//
//	len(mockedsessionRepo.FindByTokenCalls())
func (mock *sessionRepoMock) FindByTokenCalls() []struct {
	Ctx   context.Context
	Token string
} {
	var calls []struct {
		Ctx   context.Context
		Token string
	}
	mock.lockFindByToken.RLock()
	calls = mock.calls.FindByToken
	mock.lockFindByToken.RUnlock()
	return calls
}

// Ensure, that codeRepoMock does implement codeRepo.
// If this is not the case, regenerate this file with moq.
var _ codeRepo = &codeRepoMock{}

// codeRepoMock is a mock implementation of codeRepo.
//
//	func TestSomethingThatUsescodeRepo(t *testing.T) {
//
//		// make and configure a mocked codeRepo
//		mockedcodeRepo := &codeRepoMock{
//			CreateFunc: func(contextMoqParam context.Context, code *ent.Code) (*ent.Code, error) {
//				panic("mock out the Create method")
//			},
//			FindForUserFunc: func(ctx context.Context, userID int) (*ent.Code, error) {
//				panic("mock out the FindForUser method")
//			},
//		}
//
//		// use mockedcodeRepo in code that requires codeRepo
//		// and then make assertions.
//
//	}
type codeRepoMock struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(contextMoqParam context.Context, code *ent.Code) (*ent.Code, error)

	// FindForUserFunc mocks the FindForUser method.
	FindForUserFunc func(ctx context.Context, userID int) (*ent.Code, error)

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// Code is the code argument value.
			Code *ent.Code
		}
		// FindForUser holds details about calls to the FindForUser method.
		FindForUser []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserID is the userID argument value.
			UserID int
		}
	}
	lockCreate      sync.RWMutex
	lockFindForUser sync.RWMutex
}

// Create calls CreateFunc.
func (mock *codeRepoMock) Create(contextMoqParam context.Context, code *ent.Code) (*ent.Code, error) {
	if mock.CreateFunc == nil {
		panic("codeRepoMock.CreateFunc: method is nil but codeRepo.Create was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		Code            *ent.Code
	}{
		ContextMoqParam: contextMoqParam,
		Code:            code,
	}
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(contextMoqParam, code)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//
//	len(mockedcodeRepo.CreateCalls())
func (mock *codeRepoMock) CreateCalls() []struct {
	ContextMoqParam context.Context
	Code            *ent.Code
} {
	var calls []struct {
		ContextMoqParam context.Context
		Code            *ent.Code
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
	return calls
}

// FindForUser calls FindForUserFunc.
func (mock *codeRepoMock) FindForUser(ctx context.Context, userID int) (*ent.Code, error) {
	if mock.FindForUserFunc == nil {
		panic("codeRepoMock.FindForUserFunc: method is nil but codeRepo.FindForUser was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		UserID int
	}{
		Ctx:    ctx,
		UserID: userID,
	}
	mock.lockFindForUser.Lock()
	mock.calls.FindForUser = append(mock.calls.FindForUser, callInfo)
	mock.lockFindForUser.Unlock()
	return mock.FindForUserFunc(ctx, userID)
}

// FindForUserCalls gets all the calls that were made to FindForUser.
// Check the length with:
//
//	len(mockedcodeRepo.FindForUserCalls())
func (mock *codeRepoMock) FindForUserCalls() []struct {
	Ctx    context.Context
	UserID int
} {
	var calls []struct {
		Ctx    context.Context
		UserID int
	}
	mock.lockFindForUser.RLock()
	calls = mock.calls.FindForUser
	mock.lockFindForUser.RUnlock()
	return calls
}

// Ensure, that historyRepoMock does implement historyRepo.
// If this is not the case, regenerate this file with moq.
var _ historyRepo = &historyRepoMock{}

// historyRepoMock is a mock implementation of historyRepo.
//
//	func TestSomethingThatUseshistoryRepo(t *testing.T) {
//
//		// make and configure a mocked historyRepo
//		mockedhistoryRepo := &historyRepoMock{
//			CreateFunc: func(contextMoqParam context.Context, history *ent.History) (*ent.History, error) {
//				panic("mock out the Create method")
//			},
//			FindLastUserEventsFunc: func(ctx context.Context, userID int, types []string, interval time.Duration) ([]*ent.History, error) {
//				panic("mock out the FindLastUserEvents method")
//			},
//			FindUserEventsFunc: func(ctx context.Context, userID int, limit int, offset int) ([]*ent.History, error) {
//				panic("mock out the FindUserEvents method")
//			},
//		}
//
//		// use mockedhistoryRepo in code that requires historyRepo
//		// and then make assertions.
//
//	}
type historyRepoMock struct {
	// CreateFunc mocks the Create method.
	CreateFunc func(contextMoqParam context.Context, history *ent.History) (*ent.History, error)

	// FindLastUserEventsFunc mocks the FindLastUserEvents method.
	FindLastUserEventsFunc func(ctx context.Context, userID int, types []string, interval time.Duration) ([]*ent.History, error)

	// FindUserEventsFunc mocks the FindUserEvents method.
	FindUserEventsFunc func(ctx context.Context, userID int, limit int, offset int) ([]*ent.History, error)

	// calls tracks calls to the methods.
	calls struct {
		// Create holds details about calls to the Create method.
		Create []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// History is the history argument value.
			History *ent.History
		}
		// FindLastUserEvents holds details about calls to the FindLastUserEvents method.
		FindLastUserEvents []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserID is the userID argument value.
			UserID int
			// Types is the types argument value.
			Types []string
			// Interval is the interval argument value.
			Interval time.Duration
		}
		// FindUserEvents holds details about calls to the FindUserEvents method.
		FindUserEvents []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserID is the userID argument value.
			UserID int
			// Limit is the limit argument value.
			Limit int
			// Offset is the offset argument value.
			Offset int
		}
	}
	lockCreate             sync.RWMutex
	lockFindLastUserEvents sync.RWMutex
	lockFindUserEvents     sync.RWMutex
}

// Create calls CreateFunc.
func (mock *historyRepoMock) Create(contextMoqParam context.Context, history *ent.History) (*ent.History, error) {
	if mock.CreateFunc == nil {
		panic("historyRepoMock.CreateFunc: method is nil but historyRepo.Create was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		History         *ent.History
	}{
		ContextMoqParam: contextMoqParam,
		History:         history,
	}
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(contextMoqParam, history)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//
//	len(mockedhistoryRepo.CreateCalls())
func (mock *historyRepoMock) CreateCalls() []struct {
	ContextMoqParam context.Context
	History         *ent.History
} {
	var calls []struct {
		ContextMoqParam context.Context
		History         *ent.History
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
	return calls
}

// FindLastUserEvents calls FindLastUserEventsFunc.
func (mock *historyRepoMock) FindLastUserEvents(ctx context.Context, userID int, types []string, interval time.Duration) ([]*ent.History, error) {
	if mock.FindLastUserEventsFunc == nil {
		panic("historyRepoMock.FindLastUserEventsFunc: method is nil but historyRepo.FindLastUserEvents was just called")
	}
	callInfo := struct {
		Ctx      context.Context
		UserID   int
		Types    []string
		Interval time.Duration
	}{
		Ctx:      ctx,
		UserID:   userID,
		Types:    types,
		Interval: interval,
	}
	mock.lockFindLastUserEvents.Lock()
	mock.calls.FindLastUserEvents = append(mock.calls.FindLastUserEvents, callInfo)
	mock.lockFindLastUserEvents.Unlock()
	return mock.FindLastUserEventsFunc(ctx, userID, types, interval)
}

// FindLastUserEventsCalls gets all the calls that were made to FindLastUserEvents.
// Check the length with:
//
//	len(mockedhistoryRepo.FindLastUserEventsCalls())
func (mock *historyRepoMock) FindLastUserEventsCalls() []struct {
	Ctx      context.Context
	UserID   int
	Types    []string
	Interval time.Duration
} {
	var calls []struct {
		Ctx      context.Context
		UserID   int
		Types    []string
		Interval time.Duration
	}
	mock.lockFindLastUserEvents.RLock()
	calls = mock.calls.FindLastUserEvents
	mock.lockFindLastUserEvents.RUnlock()
	return calls
}

// FindUserEvents calls FindUserEventsFunc.
func (mock *historyRepoMock) FindUserEvents(ctx context.Context, userID int, limit int, offset int) ([]*ent.History, error) {
	if mock.FindUserEventsFunc == nil {
		panic("historyRepoMock.FindUserEventsFunc: method is nil but historyRepo.FindUserEvents was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		UserID int
		Limit  int
		Offset int
	}{
		Ctx:    ctx,
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	}
	mock.lockFindUserEvents.Lock()
	mock.calls.FindUserEvents = append(mock.calls.FindUserEvents, callInfo)
	mock.lockFindUserEvents.Unlock()
	return mock.FindUserEventsFunc(ctx, userID, limit, offset)
}

// FindUserEventsCalls gets all the calls that were made to FindUserEvents.
// Check the length with:
//
//	len(mockedhistoryRepo.FindUserEventsCalls())
func (mock *historyRepoMock) FindUserEventsCalls() []struct {
	Ctx    context.Context
	UserID int
	Limit  int
	Offset int
} {
	var calls []struct {
		Ctx    context.Context
		UserID int
		Limit  int
		Offset int
	}
	mock.lockFindUserEvents.RLock()
	calls = mock.calls.FindUserEvents
	mock.lockFindUserEvents.RUnlock()
	return calls
}

// Ensure, that userRepoMock does implement userRepo.
// If this is not the case, regenerate this file with moq.
var _ userRepo = &userRepoMock{}

// userRepoMock is a mock implementation of userRepo.
//
//	func TestSomethingThatUsesuserRepo(t *testing.T) {
//
//		// make and configure a mocked userRepo
//		mockeduserRepo := &userRepoMock{
//			ActivateFunc: func(ctx context.Context, userID int) (*ent.User, error) {
//				panic("mock out the Activate method")
//			},
//			CreateFunc: func(contextMoqParam context.Context, user *ent.User) (*ent.User, error) {
//				panic("mock out the Create method")
//			},
//			DeactivateFunc: func(ctx context.Context, userID int) (*ent.User, error) {
//				panic("mock out the Deactivate method")
//			},
//			FindByEmailFunc: func(ctx context.Context, email string) (*ent.User, error) {
//				panic("mock out the FindByEmail method")
//			},
//			FindByIDFunc: func(ctx context.Context, id int) (*ent.User, error) {
//				panic("mock out the FindByID method")
//			},
//			FindByPhoneFunc: func(ctx context.Context, phone string) (*ent.User, error) {
//				panic("mock out the FindByPhone method")
//			},
//			ListFunc: func(ctx context.Context, limit int64, offset int64, orderFields []string, orderDirection string) ([]*ent.User, error) {
//				panic("mock out the List method")
//			},
//			UpdateFunc: func(contextMoqParam context.Context, user *ent.User) (*ent.User, error) {
//				panic("mock out the Update method")
//			},
//		}
//
//		// use mockeduserRepo in code that requires userRepo
//		// and then make assertions.
//
//	}
type userRepoMock struct {
	// ActivateFunc mocks the Activate method.
	ActivateFunc func(ctx context.Context, userID int) (*ent.User, error)

	// CreateFunc mocks the Create method.
	CreateFunc func(contextMoqParam context.Context, user *ent.User) (*ent.User, error)

	// DeactivateFunc mocks the Deactivate method.
	DeactivateFunc func(ctx context.Context, userID int) (*ent.User, error)

	// FindByEmailFunc mocks the FindByEmail method.
	FindByEmailFunc func(ctx context.Context, email string) (*ent.User, error)

	// FindByIDFunc mocks the FindByID method.
	FindByIDFunc func(ctx context.Context, id int) (*ent.User, error)

	// FindByPhoneFunc mocks the FindByPhone method.
	FindByPhoneFunc func(ctx context.Context, phone string) (*ent.User, error)

	// ListFunc mocks the List method.
	ListFunc func(ctx context.Context, limit int64, offset int64, orderFields []string, orderDirection string) ([]*ent.User, error)

	// UpdateFunc mocks the Update method.
	UpdateFunc func(contextMoqParam context.Context, user *ent.User) (*ent.User, error)

	// calls tracks calls to the methods.
	calls struct {
		// Activate holds details about calls to the Activate method.
		Activate []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserID is the userID argument value.
			UserID int
		}
		// Create holds details about calls to the Create method.
		Create []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// User is the user argument value.
			User *ent.User
		}
		// Deactivate holds details about calls to the Deactivate method.
		Deactivate []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// UserID is the userID argument value.
			UserID int
		}
		// FindByEmail holds details about calls to the FindByEmail method.
		FindByEmail []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Email is the email argument value.
			Email string
		}
		// FindByID holds details about calls to the FindByID method.
		FindByID []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// ID is the id argument value.
			ID int
		}
		// FindByPhone holds details about calls to the FindByPhone method.
		FindByPhone []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Phone is the phone argument value.
			Phone string
		}
		// List holds details about calls to the List method.
		List []struct {
			// Ctx is the ctx argument value.
			Ctx context.Context
			// Limit is the limit argument value.
			Limit int64
			// Offset is the offset argument value.
			Offset int64
			// OrderFields is the orderFields argument value.
			OrderFields []string
			// OrderDirection is the orderDirection argument value.
			OrderDirection string
		}
		// Update holds details about calls to the Update method.
		Update []struct {
			// ContextMoqParam is the contextMoqParam argument value.
			ContextMoqParam context.Context
			// User is the user argument value.
			User *ent.User
		}
	}
	lockActivate    sync.RWMutex
	lockCreate      sync.RWMutex
	lockDeactivate  sync.RWMutex
	lockFindByEmail sync.RWMutex
	lockFindByID    sync.RWMutex
	lockFindByPhone sync.RWMutex
	lockList        sync.RWMutex
	lockUpdate      sync.RWMutex
}

// Activate calls ActivateFunc.
func (mock *userRepoMock) Activate(ctx context.Context, userID int) (*ent.User, error) {
	if mock.ActivateFunc == nil {
		panic("userRepoMock.ActivateFunc: method is nil but userRepo.Activate was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		UserID int
	}{
		Ctx:    ctx,
		UserID: userID,
	}
	mock.lockActivate.Lock()
	mock.calls.Activate = append(mock.calls.Activate, callInfo)
	mock.lockActivate.Unlock()
	return mock.ActivateFunc(ctx, userID)
}

// ActivateCalls gets all the calls that were made to Activate.
// Check the length with:
//
//	len(mockeduserRepo.ActivateCalls())
func (mock *userRepoMock) ActivateCalls() []struct {
	Ctx    context.Context
	UserID int
} {
	var calls []struct {
		Ctx    context.Context
		UserID int
	}
	mock.lockActivate.RLock()
	calls = mock.calls.Activate
	mock.lockActivate.RUnlock()
	return calls
}

// Create calls CreateFunc.
func (mock *userRepoMock) Create(contextMoqParam context.Context, user *ent.User) (*ent.User, error) {
	if mock.CreateFunc == nil {
		panic("userRepoMock.CreateFunc: method is nil but userRepo.Create was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		User            *ent.User
	}{
		ContextMoqParam: contextMoqParam,
		User:            user,
	}
	mock.lockCreate.Lock()
	mock.calls.Create = append(mock.calls.Create, callInfo)
	mock.lockCreate.Unlock()
	return mock.CreateFunc(contextMoqParam, user)
}

// CreateCalls gets all the calls that were made to Create.
// Check the length with:
//
//	len(mockeduserRepo.CreateCalls())
func (mock *userRepoMock) CreateCalls() []struct {
	ContextMoqParam context.Context
	User            *ent.User
} {
	var calls []struct {
		ContextMoqParam context.Context
		User            *ent.User
	}
	mock.lockCreate.RLock()
	calls = mock.calls.Create
	mock.lockCreate.RUnlock()
	return calls
}

// Deactivate calls DeactivateFunc.
func (mock *userRepoMock) Deactivate(ctx context.Context, userID int) (*ent.User, error) {
	if mock.DeactivateFunc == nil {
		panic("userRepoMock.DeactivateFunc: method is nil but userRepo.Deactivate was just called")
	}
	callInfo := struct {
		Ctx    context.Context
		UserID int
	}{
		Ctx:    ctx,
		UserID: userID,
	}
	mock.lockDeactivate.Lock()
	mock.calls.Deactivate = append(mock.calls.Deactivate, callInfo)
	mock.lockDeactivate.Unlock()
	return mock.DeactivateFunc(ctx, userID)
}

// DeactivateCalls gets all the calls that were made to Deactivate.
// Check the length with:
//
//	len(mockeduserRepo.DeactivateCalls())
func (mock *userRepoMock) DeactivateCalls() []struct {
	Ctx    context.Context
	UserID int
} {
	var calls []struct {
		Ctx    context.Context
		UserID int
	}
	mock.lockDeactivate.RLock()
	calls = mock.calls.Deactivate
	mock.lockDeactivate.RUnlock()
	return calls
}

// FindByEmail calls FindByEmailFunc.
func (mock *userRepoMock) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	if mock.FindByEmailFunc == nil {
		panic("userRepoMock.FindByEmailFunc: method is nil but userRepo.FindByEmail was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Email string
	}{
		Ctx:   ctx,
		Email: email,
	}
	mock.lockFindByEmail.Lock()
	mock.calls.FindByEmail = append(mock.calls.FindByEmail, callInfo)
	mock.lockFindByEmail.Unlock()
	return mock.FindByEmailFunc(ctx, email)
}

// FindByEmailCalls gets all the calls that were made to FindByEmail.
// Check the length with:
//
//	len(mockeduserRepo.FindByEmailCalls())
func (mock *userRepoMock) FindByEmailCalls() []struct {
	Ctx   context.Context
	Email string
} {
	var calls []struct {
		Ctx   context.Context
		Email string
	}
	mock.lockFindByEmail.RLock()
	calls = mock.calls.FindByEmail
	mock.lockFindByEmail.RUnlock()
	return calls
}

// FindByID calls FindByIDFunc.
func (mock *userRepoMock) FindByID(ctx context.Context, id int) (*ent.User, error) {
	if mock.FindByIDFunc == nil {
		panic("userRepoMock.FindByIDFunc: method is nil but userRepo.FindByID was just called")
	}
	callInfo := struct {
		Ctx context.Context
		ID  int
	}{
		Ctx: ctx,
		ID:  id,
	}
	mock.lockFindByID.Lock()
	mock.calls.FindByID = append(mock.calls.FindByID, callInfo)
	mock.lockFindByID.Unlock()
	return mock.FindByIDFunc(ctx, id)
}

// FindByIDCalls gets all the calls that were made to FindByID.
// Check the length with:
//
//	len(mockeduserRepo.FindByIDCalls())
func (mock *userRepoMock) FindByIDCalls() []struct {
	Ctx context.Context
	ID  int
} {
	var calls []struct {
		Ctx context.Context
		ID  int
	}
	mock.lockFindByID.RLock()
	calls = mock.calls.FindByID
	mock.lockFindByID.RUnlock()
	return calls
}

// FindByPhone calls FindByPhoneFunc.
func (mock *userRepoMock) FindByPhone(ctx context.Context, phone string) (*ent.User, error) {
	if mock.FindByPhoneFunc == nil {
		panic("userRepoMock.FindByPhoneFunc: method is nil but userRepo.FindByPhone was just called")
	}
	callInfo := struct {
		Ctx   context.Context
		Phone string
	}{
		Ctx:   ctx,
		Phone: phone,
	}
	mock.lockFindByPhone.Lock()
	mock.calls.FindByPhone = append(mock.calls.FindByPhone, callInfo)
	mock.lockFindByPhone.Unlock()
	return mock.FindByPhoneFunc(ctx, phone)
}

// FindByPhoneCalls gets all the calls that were made to FindByPhone.
// Check the length with:
//
//	len(mockeduserRepo.FindByPhoneCalls())
func (mock *userRepoMock) FindByPhoneCalls() []struct {
	Ctx   context.Context
	Phone string
} {
	var calls []struct {
		Ctx   context.Context
		Phone string
	}
	mock.lockFindByPhone.RLock()
	calls = mock.calls.FindByPhone
	mock.lockFindByPhone.RUnlock()
	return calls
}

// List calls ListFunc.
func (mock *userRepoMock) List(ctx context.Context, limit int64, offset int64, orderFields []string, orderDirection string) ([]*ent.User, error) {
	if mock.ListFunc == nil {
		panic("userRepoMock.ListFunc: method is nil but userRepo.List was just called")
	}
	callInfo := struct {
		Ctx            context.Context
		Limit          int64
		Offset         int64
		OrderFields    []string
		OrderDirection string
	}{
		Ctx:            ctx,
		Limit:          limit,
		Offset:         offset,
		OrderFields:    orderFields,
		OrderDirection: orderDirection,
	}
	mock.lockList.Lock()
	mock.calls.List = append(mock.calls.List, callInfo)
	mock.lockList.Unlock()
	return mock.ListFunc(ctx, limit, offset, orderFields, orderDirection)
}

// ListCalls gets all the calls that were made to List.
// Check the length with:
//
//	len(mockeduserRepo.ListCalls())
func (mock *userRepoMock) ListCalls() []struct {
	Ctx            context.Context
	Limit          int64
	Offset         int64
	OrderFields    []string
	OrderDirection string
} {
	var calls []struct {
		Ctx            context.Context
		Limit          int64
		Offset         int64
		OrderFields    []string
		OrderDirection string
	}
	mock.lockList.RLock()
	calls = mock.calls.List
	mock.lockList.RUnlock()
	return calls
}

// Update calls UpdateFunc.
func (mock *userRepoMock) Update(contextMoqParam context.Context, user *ent.User) (*ent.User, error) {
	if mock.UpdateFunc == nil {
		panic("userRepoMock.UpdateFunc: method is nil but userRepo.Update was just called")
	}
	callInfo := struct {
		ContextMoqParam context.Context
		User            *ent.User
	}{
		ContextMoqParam: contextMoqParam,
		User:            user,
	}
	mock.lockUpdate.Lock()
	mock.calls.Update = append(mock.calls.Update, callInfo)
	mock.lockUpdate.Unlock()
	return mock.UpdateFunc(contextMoqParam, user)
}

// UpdateCalls gets all the calls that were made to Update.
// Check the length with:
//
//	len(mockeduserRepo.UpdateCalls())
func (mock *userRepoMock) UpdateCalls() []struct {
	ContextMoqParam context.Context
	User            *ent.User
} {
	var calls []struct {
		ContextMoqParam context.Context
		User            *ent.User
	}
	mock.lockUpdate.RLock()
	calls = mock.calls.Update
	mock.lockUpdate.RUnlock()
	return calls
}