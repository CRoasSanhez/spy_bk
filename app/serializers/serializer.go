package serializers

import (
	"bytes"
	"compress/gzip"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"reflect"
	"spyc_backend/app/core"

	"github.com/revel/revel"
)

// ENCRYPT flag to encrypt data on serializer
const ENCRYPT = false

// Serializer is an interface for make Alias JSON response
type Serializer interface {
	Cast(data interface{}) Serializer
}

// SerializerModel struct for JSON response
type SerializerModel struct {
	HTTPStatus int               `json:"-"`
	Code       int               `json:"code"`
	Success    bool              `json:"success"`
	Message    string            `json:"message"`
	Errors     interface{}       `json:"errors"`
	Data       interface{}       `json:"data"`
	Headers    map[string]string `json:"-"`
	Compressed bool              `json:"-"`
}

// Apply make JSON response
func (s SerializerModel) Apply(req *revel.Request, resp *revel.Response) {
	respJSON, _ := json.Marshal(s)

	if s.Compressed {
		resp.Out.Header().Set("Content-Encoding", "gzip")

		var b bytes.Buffer
		w := gzip.NewWriter(&b)
		w.Write(respJSON)
		w.Close()

		resp.GetWriter().Write([]byte(b.Bytes()))
		return
	}

	if len(s.Headers) > 0 {
		for k, v := range s.Headers {
			resp.Out.Header().Set(k, v)
		}
	}

	resp.WriteHeader(s.HTTPStatus, core.ApplicationJSON)
	resp.ContentType = core.ApplicationJSON

	if ENCRYPT {
		pubkey, err := core.GetPubKey()
		if err != nil {
			resp.GetWriter().Write(respJSON)
		} else {
			encryptedPKCS1v15, _ := rsa.EncryptPKCS1v15(rand.Reader, pubkey, respJSON)
			resp.GetWriter().Write([]byte(base64.URLEncoding.EncodeToString(encryptedPKCS1v15)))
		}
	} else {
		resp.GetWriter().Write(respJSON)
	}
}

// Serialize ...
func Serialize(data interface{}, s Serializer) interface{} {
	var serializer interface{}

	if IsArray(data) {
		arr := make([]interface{}, 0)
		slice := reflect.ValueOf(data)
		for i := 0; i < slice.Len(); i++ {
			if s != nil {
				arr = append(arr, s.Cast(slice.Index(i).Interface()))
			} else {
				arr = append(arr, slice.Index(i).Interface())
			}
		}

		serializer = arr
	} else {
		if s != nil {
			serializer = s.Cast(data)
		} else {
			serializer = data
		}
	}

	return serializer
}

// SuccessResponse for serializer. HTTP code 200
func SuccessResponse(data interface{}, message string, code int, s Serializer) SerializerModel {
	var model SerializerModel

	model.Code = code
	model.Success = true
	model.Message = message
	model.HTTPStatus = http.StatusOK

	if s == nil {
		model.Data = data
		return model
	}

	model.Data = Serialize(data, s)

	return model
}

// ErrorResponse for serializer. HTTP code 400
func ErrorResponse(errors interface{}, message string, code int) SerializerModel {
	var model SerializerModel

	model.Code = code
	model.Success = false
	model.Message = message
	model.Errors = errors
	model.HTTPStatus = http.StatusBadRequest

	return model
}

// UnauthorizedResponse for serializer, HTTP code 401
func UnauthorizedResponse() SerializerModel {
	var model SerializerModel

	model.Code = 401
	model.Success = false
	model.Message = "Unauthorized"
	model.Errors = nil
	model.HTTPStatus = http.StatusUnauthorized

	return model
}

// ForbiddenResponse for serializer, HTTP code 403
func ForbiddenResponse() SerializerModel {
	var model SerializerModel

	model.Code = 403
	model.Success = false
	model.Message = "Not found"
	model.Errors = nil
	model.HTTPStatus = http.StatusForbidden

	return model
}

// ServerErrorResponse ...
func ServerErrorResponse() SerializerModel {
	var model SerializerModel

	model.Code = 500
	model.Success = false
	model.Message = "Server error"
	model.Errors = nil
	model.HTTPStatus = http.StatusInternalServerError

	return model
}

// Compress response
func (s SerializerModel) Compress() SerializerModel {
	s.Compressed = true

	return s
}

// IsArray verify data type if array
func IsArray(data interface{}) bool {
	r := reflect.TypeOf(data)
	for r.Kind() == reflect.Ptr {
		r = r.Elem()
	}

	if r.Kind() == reflect.Slice {
		return true
	}

	return false
}
