package lightGin

import (
	"io"
	"net/http"
)

const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

// ResponseWriter 自己定义的http writer，至少要包含http.ResponseWriter
type ResponseWriter interface {
	http.ResponseWriter

	// Returns the HTTP response status code of the current request.
	Status() int

	// Returns the number of bytes already written into the response http body.
	// See Written()
	Size() int

	// Writes the string into the response body.
	WriteString(string) (int, error)

	// Returns true if the response body was already written.
	Written() bool

	// Forces to write the http header (status code + headers).
	WriteHeaderNow()
}

type ResponseWriterEntity struct {
	http.ResponseWriter
	size   int
	status int
}

// 检验ResponseWriterEntity 是否实现了ResponseWriter的全部接口
var _ ResponseWriter = &ResponseWriterEntity{}

// reset 初始化ResponseWriterEntity
func (w *ResponseWriterEntity) Reset(writer http.ResponseWriter) {
	w.ResponseWriter = writer
	w.size = noWritten
	w.status = defaultStatus
}

// Written 判断是否已经写入
func (w *ResponseWriterEntity) Written() bool {
	return w.size != noWritten
}

// Status http code
func (w *ResponseWriterEntity) Status() int {
	return w.status
}

// Size 多少字节
func (w *ResponseWriterEntity) Size() int {
	return w.size
}

//WriteHeader 写入状态码
func (w *ResponseWriterEntity) WriteHeader(code int) {
	if code > 0 && w.status != code {
		if w.Written() {
			//todo logger
			//debugPrint("[WARNING] Headers were already written. Wanted to override status code %d with %d", w.status, code)
		}
		w.status = code
	}
}

//WriteHeaderNow 立刻写入header
func (w *ResponseWriterEntity) WriteHeaderNow() {
	if !w.Written() {
		w.size = 0
		w.ResponseWriter.WriteHeader(w.status)
	}
}

// Write 写入内容
func (w *ResponseWriterEntity) Write(data []byte) (n int, err error) {
	w.WriteHeaderNow()
	n, err = w.ResponseWriter.Write(data)
	w.size += n
	return
}

// Write 写入内容
func (w *ResponseWriterEntity) WriteString(s string) (n int, err error) {
	w.WriteHeaderNow()
	n, err = io.WriteString(w.ResponseWriter, s)
	w.size += n
	return
}
