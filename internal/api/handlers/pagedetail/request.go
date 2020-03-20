package pagedetailhandler

import (
	"encoding/json"
	"net/http"

	"github.com/worlve/sp-service/internal/models/pagedetail"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// UpdatePageDetailRequest parameters from the UpdatePageDetail call
type UpdatePageDetailRequest struct {
	PageGUID       string
	PageDetailGUID string
	Title          string                 `json:"title"`
	Summary        string                 `json:"summary"`
	Partitions     []pagedetail.Partition `json:"partitions"`
}

// NewUpdatePageDetailRequest extracts the UpdatePageDetailRequest
func NewUpdatePageDetailRequest(r *http.Request, p httprouter.Params) (UpdatePageDetailRequest, error) {
	var request UpdatePageDetailRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return request, errors.New("invalid request")
	}
	err = pagedetail.UnmarshalPartitions(request.Partitions)
	if err != nil {
		return request, errors.New("not valid page partitions")
	}
	request.PageGUID = p.ByName(PageIDRouteKey)
	request.PageDetailGUID = p.ByName(PageDetailIDRouteKey)
	return request.validate()
}

func (request UpdatePageDetailRequest) validate() (UpdatePageDetailRequest, error) {
	if request.Title == "" {
		return request, errors.New("a page detail must retain a title")
	}
	return request, nil
}
