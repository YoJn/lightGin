package lightGin

import (
	"net/http"

	"github.com/YoJn/lightGin/internal/common"
)

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
}
