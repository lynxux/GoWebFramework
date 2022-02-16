package framework

import (
	"errors"
	"fmt"
	"sync"
)

// Container 是一个服务容器，提供绑定服务和获取服务的功能   绑定服务和注册服务好像是一个东西？
type Container interface {
	// Bind 绑定一个服务提供者，如果关键字凭证已经存在，会进行替换操作，返回error
	Bind(provider ServiceProvider) error
	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务 如果这个关键字凭证未绑定服务提供者，那么会error,不会panic
	Make(key string) (interface{}, error)
	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会panic。
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者。
	MustMake(key string) interface{}
	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的params参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// MyWebContainer 是服务容器的具体实现
type MyWebContainer struct {
	Container
	providers map[string]ServiceProvider // providers 存储注册的服务提供者，key为字符串凭证
	instances map[string]interface{}     // instance 存储具体的实例，key为字符串凭证
	lock      sync.RWMutex               // lock 用于锁住对容器的变更操作
}

// NewMywebContainer 创建一个服务容器
func NewMywebContainer() *MyWebContainer {
	return &MyWebContainer{
		providers: map[string]ServiceProvider{},
		instances: map[string]interface{}{},
		lock:      sync.RWMutex{},
	}
}

// PrintProviders 输出服务容器中注册的关键字
func (mwf *MyWebContainer) PrintProviders() []string {
	ret := []string{}
	for _, provider := range mwf.providers {
		name := provider.Name()

		line := fmt.Sprint(name)
		ret = append(ret, line)
	}
	return ret
}

// Bind 将服务容器和关键字做了绑定
func (mwf *MyWebContainer) Bind(provider ServiceProvider) error {
	// mwf->myWebFramework
	mwf.lock.Lock()
	defer mwf.lock.Unlock()

	key := provider.Name()
	mwf.providers[key] = provider

	// 如果 IsDefer 方法标记这个服务实例要延迟实例化，即等到第一次 make 的时候再实例化，那么在 Bind 操作的时候，就什么都不需要做；
	// 而如果 IsDefer 方法为 false，即注册时就要实例化，那么我们就要在 Bind 函数中增加实例化的方法。
	// if provider is not defer
	if provider.IsDefer() == false {
		if err := provider.Boot(mwf); err != nil {
			return err
		}
		params := provider.Params(mwf)
		method := provider.Register(mwf)
		instance, err := method(params...)
		if err != nil {
			return errors.New(err.Error())
		}
		mwf.instances[key] = instance
	}
	return nil
}

func (mwf *MyWebContainer) IsBind(key string) bool {
	return mwf.findServiceProvider(key) != nil
}

func (mwf *MyWebContainer) findServiceProvider(key string) ServiceProvider {
	mwf.lock.RLock() //读取的时候禁止写
	defer mwf.lock.RUnlock()
	if sp, ok := mwf.providers[key]; ok {
		return sp
	}
	return nil
}

func (mwf *MyWebContainer) Make(key string) (interface{}, error) {
	return mwf.make(key, nil, false)
}

func (mwf *MyWebContainer) MustMake(key string) interface{} {
	serv, err := mwf.make(key, nil, false)
	if err != nil {
		panic(err)
	}
	return serv
}

func (mwf *MyWebContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return mwf.make(key, params, true)
}

func (mwf *MyWebContainer) newInstance(sp ServiceProvider, params []interface{}) (interface{}, error) {
	if err := sp.Boot(mwf); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(mwf)
	}
	method := sp.Register(mwf)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, err
}

// 真正的实例化一个服务
func (mwf *MyWebContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	mwf.lock.RLock()
	defer mwf.lock.RUnlock()
	// 查询是否已经注册了这个服务提供者，如果没有注册，则返回错误
	sp := mwf.findServiceProvider(key)
	if sp == nil { //没有注册
		return nil, errors.New("contract " + key + "have not register")
	}
	// 注册了，但是要强制重新实例化（在有参数的时候使用-MakeNew）
	if forceNew {
		mwf.newInstance(sp, params)
	}

	// 注册了，不需要强制重新实例化，如果容器中已经实例化了，那么就直接使用容器中的实例（用于Make，MustMake方法）
	if ins, ok := mwf.instances[key]; ok {
		return ins, nil
	}

	// 注册了，不需要强制重新实例化但是容器中还未实例化，则进行一次实例化 （用于Make，MustMake方法）
	inst, err := mwf.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}

	mwf.instances[key] = inst
	return inst, nil
}
