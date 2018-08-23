// Code generated by counterfeiter. DO NOT EDIT.
package ldap_fake

import (
	"sync"
	"time"

	"code.cloudfoundry.org/goshims/ldapshim"
	ldap "gopkg.in/ldap.v2"
)

type FakeLdapConnection struct {
	SetTimeoutStub        func(timeout time.Duration)
	setTimeoutMutex       sync.RWMutex
	setTimeoutArgsForCall []struct {
		timeout time.Duration
	}
	CloseStub        func()
	closeMutex       sync.RWMutex
	closeArgsForCall []struct{}
	BindStub         func(string, string) error
	bindMutex        sync.RWMutex
	bindArgsForCall  []struct {
		arg1 string
		arg2 string
	}
	bindReturns struct {
		result1 error
	}
	bindReturnsOnCall map[int]struct {
		result1 error
	}
	SearchStub        func(*ldap.SearchRequest) (*ldap.SearchResult, error)
	searchMutex       sync.RWMutex
	searchArgsForCall []struct {
		arg1 *ldap.SearchRequest
	}
	searchReturns struct {
		result1 *ldap.SearchResult
		result2 error
	}
	searchReturnsOnCall map[int]struct {
		result1 *ldap.SearchResult
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeLdapConnection) SetTimeout(timeout time.Duration) {
	fake.setTimeoutMutex.Lock()
	fake.setTimeoutArgsForCall = append(fake.setTimeoutArgsForCall, struct {
		timeout time.Duration
	}{timeout})
	fake.recordInvocation("SetTimeout", []interface{}{timeout})
	fake.setTimeoutMutex.Unlock()
	if fake.SetTimeoutStub != nil {
		fake.SetTimeoutStub(timeout)
	}
}

func (fake *FakeLdapConnection) SetTimeoutCallCount() int {
	fake.setTimeoutMutex.RLock()
	defer fake.setTimeoutMutex.RUnlock()
	return len(fake.setTimeoutArgsForCall)
}

func (fake *FakeLdapConnection) SetTimeoutArgsForCall(i int) time.Duration {
	fake.setTimeoutMutex.RLock()
	defer fake.setTimeoutMutex.RUnlock()
	return fake.setTimeoutArgsForCall[i].timeout
}

func (fake *FakeLdapConnection) Close() {
	fake.closeMutex.Lock()
	fake.closeArgsForCall = append(fake.closeArgsForCall, struct{}{})
	fake.recordInvocation("Close", []interface{}{})
	fake.closeMutex.Unlock()
	if fake.CloseStub != nil {
		fake.CloseStub()
	}
}

func (fake *FakeLdapConnection) CloseCallCount() int {
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	return len(fake.closeArgsForCall)
}

func (fake *FakeLdapConnection) Bind(arg1 string, arg2 string) error {
	fake.bindMutex.Lock()
	ret, specificReturn := fake.bindReturnsOnCall[len(fake.bindArgsForCall)]
	fake.bindArgsForCall = append(fake.bindArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	fake.recordInvocation("Bind", []interface{}{arg1, arg2})
	fake.bindMutex.Unlock()
	if fake.BindStub != nil {
		return fake.BindStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.bindReturns.result1
}

func (fake *FakeLdapConnection) BindCallCount() int {
	fake.bindMutex.RLock()
	defer fake.bindMutex.RUnlock()
	return len(fake.bindArgsForCall)
}

func (fake *FakeLdapConnection) BindArgsForCall(i int) (string, string) {
	fake.bindMutex.RLock()
	defer fake.bindMutex.RUnlock()
	return fake.bindArgsForCall[i].arg1, fake.bindArgsForCall[i].arg2
}

func (fake *FakeLdapConnection) BindReturns(result1 error) {
	fake.BindStub = nil
	fake.bindReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeLdapConnection) BindReturnsOnCall(i int, result1 error) {
	fake.BindStub = nil
	if fake.bindReturnsOnCall == nil {
		fake.bindReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.bindReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeLdapConnection) Search(arg1 *ldap.SearchRequest) (*ldap.SearchResult, error) {
	fake.searchMutex.Lock()
	ret, specificReturn := fake.searchReturnsOnCall[len(fake.searchArgsForCall)]
	fake.searchArgsForCall = append(fake.searchArgsForCall, struct {
		arg1 *ldap.SearchRequest
	}{arg1})
	fake.recordInvocation("Search", []interface{}{arg1})
	fake.searchMutex.Unlock()
	if fake.SearchStub != nil {
		return fake.SearchStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.searchReturns.result1, fake.searchReturns.result2
}

func (fake *FakeLdapConnection) SearchCallCount() int {
	fake.searchMutex.RLock()
	defer fake.searchMutex.RUnlock()
	return len(fake.searchArgsForCall)
}

func (fake *FakeLdapConnection) SearchArgsForCall(i int) *ldap.SearchRequest {
	fake.searchMutex.RLock()
	defer fake.searchMutex.RUnlock()
	return fake.searchArgsForCall[i].arg1
}

func (fake *FakeLdapConnection) SearchReturns(result1 *ldap.SearchResult, result2 error) {
	fake.SearchStub = nil
	fake.searchReturns = struct {
		result1 *ldap.SearchResult
		result2 error
	}{result1, result2}
}

func (fake *FakeLdapConnection) SearchReturnsOnCall(i int, result1 *ldap.SearchResult, result2 error) {
	fake.SearchStub = nil
	if fake.searchReturnsOnCall == nil {
		fake.searchReturnsOnCall = make(map[int]struct {
			result1 *ldap.SearchResult
			result2 error
		})
	}
	fake.searchReturnsOnCall[i] = struct {
		result1 *ldap.SearchResult
		result2 error
	}{result1, result2}
}

func (fake *FakeLdapConnection) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.setTimeoutMutex.RLock()
	defer fake.setTimeoutMutex.RUnlock()
	fake.closeMutex.RLock()
	defer fake.closeMutex.RUnlock()
	fake.bindMutex.RLock()
	defer fake.bindMutex.RUnlock()
	fake.searchMutex.RLock()
	defer fake.searchMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeLdapConnection) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ ldapshim.LdapConnection = new(FakeLdapConnection)