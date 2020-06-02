package lightGin

// HandlerFunc middleware
type HandlerFunc func(*Context)

// HandlersChain middleware array.
type HandlersChain []HandlerFunc
