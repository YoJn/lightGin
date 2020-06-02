package lightGin

import (
	"github.com/YoJn/lightGin/render"
	"io/ioutil"
	"math"
	"net/http"
	"net/url"

	"github.com/YoJn/lightGin/internal/common"
)

const abortIndex int8 = math.MaxInt8 / 2

const (
	MIMEJSON              = common.MIMEJSON
	MIMEHTML              = common.MIMEHTML
	MIMEXML               = common.MIMEXML
	MIMEXML2              = common.MIMEXML2
	MIMEPlain             = common.MIMEPlain
	MIMEPOSTForm          = common.MIMEPOSTForm
	MIMEMultipartPOSTForm = common.MIMEMultipartPOSTForm
	MIMEYAML              = common.MIMEYAML
	BodyBytesKey          = "_yojn/lightGin/babibabu"
)

type Context struct {
	writerEntity common.ResponseWriterEntity
	Request      *http.Request
	Writer       common.ResponseWriter

	handlers HandlersChain
	fullPath string
	index    int8

	Keys map[string]interface{}

	Accepted []string

	queryCache url.Values
	formCache  url.Values
}

// reset 重置context
func (c *Context) reset() {
	c.Writer = &c.writerEntity
	//c.Params = c.Params[0:0]
	c.handlers = nil
	c.index = -1
	c.fullPath = ""
	c.Keys = nil
	//c.Errors = c.Errors[0:0]
	c.Accepted = nil
	c.queryCache = nil
	c.formCache = nil
}

// Copy deep copy 一个新的
func (c *Context) Copy() *Context {
	var cp = *c
	cp.writerEntity.ResponseWriter = nil
	cp.Writer = &cp.writerEntity
	//cp.index = abortIndex
	cp.handlers = nil
	cp.Keys = map[string]interface{}{}
	for k, v := range c.Keys {
		cp.Keys[k] = v
	}
	//paramCopy := make([]Param, len(cp.Params))
	//copy(paramCopy, cp.Params)
	//cp.Params = paramCopy
	return &cp
}

//HandlerName 获得最后一个Handler的name
func (c *Context) HandlerName() string {
	return common.NameOfFunction(c.handlers.Last())
}

// HandlerNames returns a list of all registered handlers for this context in descending order,
// following the semantics of HandlerName()
func (c *Context) HandlerNames() []string {
	hn := make([]string, 0, len(c.handlers))
	for _, val := range c.handlers {
		hn = append(hn, common.NameOfFunction(val))
	}
	return hn
}

// Handler returns the main handler.
func (c *Context) Handler() HandlerFunc {
	return c.handlers.Last()
}

// FullPath returns a matched route full path. For not found routes
// returns an empty string.
//     router.GET("/user/:id", func(c *gin.Context) {
//         c.FullPath() == "/user/:id" // true
//     })
func (c *Context) FullPath() string {
	return c.fullPath
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *Context) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

// IsAborted returns true if the current context was aborted.
func (c *Context) IsAborted() bool {
	return c.index >= abortIndex
}

// Abort prevents pending handlers from being called. Note that this will not stop the current handler.
// Let's say you have an authorization middleware that validates that the current request is authorized.
// If the authorization fails (ex: the password does not match), call Abort to ensure the remaining handlers
// for this request are not called.
func (c *Context) Abort() {
	c.index = abortIndex
}

// Status sets the HTTP response code.
func (c *Context) Status(code int) {
	c.writerEntity.WriteHeader(code)
}

// AbortWithStatus calls `Abort()` and writes the headers with the specified status code.
// For example, a failed attempt to authenticate a request could use: context.AbortWithStatus(401).
func (c *Context) AbortWithStatus(code int) {
	c.Status(code)
	c.Writer.WriteHeaderNow()
	c.Abort()
}

// Header is a intelligent shortcut for c.Writer.Header().Set(key, value).
// It writes a header in the response.
// If value == "", this method removes the header `c.Writer.Header().Del(key)`
func (c *Context) Header(key, value string) {
	if value == "" {
		c.Writer.Header().Del(key)
		return
	}
	c.Writer.Header().Set(key, value)
}

func (c *Context) requestHeader(key string) string {
	return c.Request.Header.Get(key)
}

// GetHeader returns value from request headers.
func (c *Context) GetHeader(key string) string {
	return c.requestHeader(key)
}

// GetRawData return stream data.
func (c *Context) GetRawData() ([]byte, error) {
	return ioutil.ReadAll(c.Request.Body)
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// Cookie returns the named cookie provided in the request or
// ErrNoCookie if not found. And return the named cookie is unescaped.
// If multiple cookies match the given name, only one cookie will
// be returned.
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (c *Context) AbortWithStatusJSON(code int, jsonObj interface{}) {
	c.Abort()
	c.JSON(code, jsonObj)
}

// bodyAllowedForStatus is a copy of http.bodyAllowedForStatus non-exported function.
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}

// Render writes the response headers and calls render.Render to render data.
func (c *Context) Render(code int, r render.Render) {
	c.Status(code)

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(c.Writer)
		c.Writer.WriteHeaderNow()
		return
	}

	if err := r.Render(c.Writer); err != nil {
		panic(err)
	}
}

// JSON serializes the given struct as JSON into the response body.
// It also sets the Content-Type as "application/json".
func (c *Context) JSON(code int, obj interface{}) {
	c.Render(code, render.JSON{Data: obj})
}
