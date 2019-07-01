// Code generated by counterfeiter. DO NOT EDIT.
package waiterfakes

import (
	"sync"

	"code.cloudfoundry.org/eirini/waiter"
	v1 "k8s.io/api/apps/v1"
	v1a "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type FakeDeploymentLister struct {
	ListStub        func(v1a.ListOptions) (*v1.DeploymentList, error)
	listMutex       sync.RWMutex
	listArgsForCall []struct {
		arg1 v1a.ListOptions
	}
	listReturns struct {
		result1 *v1.DeploymentList
		result2 error
	}
	listReturnsOnCall map[int]struct {
		result1 *v1.DeploymentList
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeDeploymentLister) List(arg1 v1a.ListOptions) (*v1.DeploymentList, error) {
	fake.listMutex.Lock()
	ret, specificReturn := fake.listReturnsOnCall[len(fake.listArgsForCall)]
	fake.listArgsForCall = append(fake.listArgsForCall, struct {
		arg1 v1a.ListOptions
	}{arg1})
	fake.recordInvocation("List", []interface{}{arg1})
	fake.listMutex.Unlock()
	if fake.ListStub != nil {
		return fake.ListStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.listReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeDeploymentLister) ListCallCount() int {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	return len(fake.listArgsForCall)
}

func (fake *FakeDeploymentLister) ListCalls(stub func(v1a.ListOptions) (*v1.DeploymentList, error)) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = stub
}

func (fake *FakeDeploymentLister) ListArgsForCall(i int) v1a.ListOptions {
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	argsForCall := fake.listArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeDeploymentLister) ListReturns(result1 *v1.DeploymentList, result2 error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = nil
	fake.listReturns = struct {
		result1 *v1.DeploymentList
		result2 error
	}{result1, result2}
}

func (fake *FakeDeploymentLister) ListReturnsOnCall(i int, result1 *v1.DeploymentList, result2 error) {
	fake.listMutex.Lock()
	defer fake.listMutex.Unlock()
	fake.ListStub = nil
	if fake.listReturnsOnCall == nil {
		fake.listReturnsOnCall = make(map[int]struct {
			result1 *v1.DeploymentList
			result2 error
		})
	}
	fake.listReturnsOnCall[i] = struct {
		result1 *v1.DeploymentList
		result2 error
	}{result1, result2}
}

func (fake *FakeDeploymentLister) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.listMutex.RLock()
	defer fake.listMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeDeploymentLister) recordInvocation(key string, args []interface{}) {
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

var _ waiter.DeploymentLister = new(FakeDeploymentLister)
