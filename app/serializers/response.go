package serializers

import (
	"encoding/json"

	"github.com/revel/revel"
)

type ResponseSerializer struct {
	Serializer Serializer
	Status     int
}

func (r ResponseSerializer) Apply(req *revel.Request, resp *revel.Response) {
	respJSON, err := json.Marshal(r.Serializer)
	if err != nil {
		resp.Out.Write(respJSON)
	}

	resp.WriteHeader(r.Status, "application/json")
	resp.ContentType = "application/json"
	resp.Out.Write(respJSON)
}
