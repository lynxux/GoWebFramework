package demo

import (
	"fmt"
	"github.com/lynxux/goWebFramework/framework"
)

// 实现具体的服务实例

// 具体的接口实例
type DemoService struct {
	// 实现接口
	Service

	// 参数
	c framework.Container
}

// 实现接口
func (s *DemoService) GetFoo() Foo {
	return Foo{
		Name: "i am foo",
	}
}

// 初始化实例的方法
func NewDemoService(params ...interface{}) (interface{}, error) {
	// 这里需要将参数展开
	c := params[0].(framework.Container)

	fmt.Println("new demo service")
	// 返回实例
	return &DemoService{c: c}, nil
}
