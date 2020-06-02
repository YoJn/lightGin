package lightGin

// HandlerFunc middleware 包括处理逻辑
type HandlerFunc func(*Context)

// HandlersChain middleware array.
type HandlersChain []HandlerFunc

// Last 最后一个中间件
func (c HandlersChain) Last() HandlerFunc {
	if length := len(c); length > 0 {
		return c[length-1]
	}
	return nil
}

// RouteInfo represents a request route's specification which contains method and path and its handler.
type RouteInfo struct {
	Method      string
	Path        string
	Handler     string
	HandlerFunc HandlerFunc
}

// RoutesInfo defines a RouteInfo array.
type RoutesInfo []RouteInfo
