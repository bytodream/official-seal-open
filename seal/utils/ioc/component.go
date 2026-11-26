package ioc

import "reflect"

// ComponentOption 组件注册可选配置项。
type ComponentOption func(component *Component)

// WithComponentName 组件名
func WithComponentName(name string) ComponentOption {
	return func(c *Component) {
		c.Name = name
	}
}

// Component 组件
type Component struct {
	Name  string // 组件名
	Value any    // 组件实例
}

// GetType 获取组件实例类型
func (c *Component) GetType() reflect.Type {
	return reflect.TypeOf(c.Value)
}

// NewComponent 创建组件实例。
func NewComponent(val any, opts ...ComponentOption) *Component {
	component := &Component{
		Value: val,
	}

	for _, o := range opts {
		o(component)
	}
	return component
}
