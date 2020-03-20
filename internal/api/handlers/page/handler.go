package pagehandler

import (
	"context"
	"net/http"

	"github.com/worlve/sp-service/internal/models/pagetemplate"
	"github.com/worlve/sp-service/internal/models/version"

	"github.com/worlve/sp-service/internal/api/handlers/nextbatch"

	"github.com/worlve/sp-service/internal/api"
	"github.com/worlve/sp-service/internal/models/page"
	pageservice "github.com/worlve/sp-service/internal/services/page"
	"github.com/worlve/sp-service/internal/stores/storeerror"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// PageService see Service for more details
type PageService interface {
	CreatePage(ctx context.Context, params pageservice.CreatePageParams) (page.Page, error)
	UpdatePage(ctx context.Context, params pageservice.UpdatePageParams) error
	RemovePage(ctx context.Context, params pageservice.RemovePageParams) error
	GetPages(ctx context.Context, params pageservice.GetPagesParams) ([]page.Page, int, string, error)
	GetPage(ctx context.Context, params pageservice.GetPageParams) (page.Page, error)
	GetEntirePage(ctx context.Context, params pageservice.GetEntirePageParams) (page.Page, error)
}

// PageHandler is the handler for the associated API
type PageHandler struct {
	PageService PageService
}

// CreatePage see Service for more details
func (h PageHandler) CreatePage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	request, err := NewCreatePageRequest(r, p)
	if err != nil {
		api.RespondWith(r, w, http.StatusBadRequest, err, err)
		return
	}
	ctx := r.Context()
	authData, err := api.GetDataFromContext(ctx)
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, errors.Wrap(err, "failed to get auth data"))
		return
	}
	record, err := h.PageService.CreatePage(ctx, pageservice.CreatePageParams{
		Page: page.Page{
			Title:   request.Title,
			Summary: request.Summary,
			Version: version.Version{
				GUID: request.VersionID,
			},
			PermissionType: request.PermissionType,
			PageTemplate: pagetemplate.PageTemplate{
				GUID: request.PageTemplateID,
			},
		},
		OwnerID: authData.UserID,
	})
	if castErr, ok := err.(*storeerror.DupEntry); ok {
		api.RespondWith(r, w, http.StatusBadRequest, castErr, err)
		return
	}
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, err)
		return
	}
	api.RespondWith(r, w, http.StatusOK, map[string]string{"id": record.GUID}, nil)
}

// UpdatePage see Service for more details
func (h PageHandler) UpdatePage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	request, err := NewUpdatePageRequest(r, p)
	if err != nil {
		api.RespondWith(r, w, http.StatusBadRequest, err, err)
		return
	}
	ctx := r.Context()
	authData, err := api.GetDataFromContext(ctx)
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, errors.Wrap(err, "failed to get auth data"))
		return
	}
	err = h.PageService.UpdatePage(ctx, pageservice.UpdatePageParams{
		Page: page.Page{
			GUID:    request.GUID,
			Title:   request.Title,
			Summary: request.Summary,
			Version: version.Version{
				GUID: request.VersionID,
			},
			PermissionType: request.PermissionType,
			PageTemplate: pagetemplate.PageTemplate{
				GUID: request.PageTemplateID,
			},
		},
		UserID: authData.UserID,
	})
	if _, ok := err.(*storeerror.NotAuthorized); ok {
		api.RespondWith(r, w, http.StatusUnauthorized, &api.FailedAuthorization{}, err)
		return
	}
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, err)
		return
	}
	api.RespondWith(r, w, http.StatusOK, nil, nil)
}

// GetEntirePage see Service for more details
func (h PageHandler) GetEntirePage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	request, err := NewGetEntirePageRequest(r, p)
	if err != nil {
		api.RespondWith(r, w, http.StatusBadRequest, err, err)
		return
	}
	ctx := r.Context()
	authData, err := api.GetDataFromContext(ctx)
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, errors.Wrap(err, "failed to get auth data"))
		return
	}
	record, err := h.PageService.GetEntirePage(ctx, pageservice.GetEntirePageParams{
		Page: page.Page{
			GUID: request.GUID,
		},
		UserID: authData.UserID,
	})
	if _, ok := err.(*storeerror.NotAuthorized); ok {
		api.RespondWith(r, w, http.StatusUnauthorized, &api.FailedAuthorization{}, err)
		return
	}
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, err)
		return
	}
	conformedRecord := record.GetJSONConformed()
	api.RespondWith(r, w, http.StatusOK, conformedRecord, nil)
}

// GetPage see Service for more details
func (h PageHandler) GetPage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	request, err := NewGetPageRequest(r, p)
	if err != nil {
		api.RespondWith(r, w, http.StatusBadRequest, err, err)
		return
	}
	ctx := r.Context()
	authData, err := api.GetDataFromContext(ctx)
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, errors.Wrap(err, "failed to get auth data"))
		return
	}
	record, err := h.PageService.GetPage(ctx, pageservice.GetPageParams{
		Page: page.Page{
			GUID: request.GUID,
		},
		UserID: authData.UserID,
	})
	reducedPage := record.Reduce()
	if _, ok := err.(*storeerror.NotAuthorized); ok {
		api.RespondWith(r, w, http.StatusUnauthorized, &api.FailedAuthorization{}, err)
		return
	}
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, err)
		return
	}
	conformedRecord := reducedPage.GetJSONConformed()
	api.RespondWith(r, w, http.StatusOK, conformedRecord, nil)
}

// GetPages see Service for more details
func (h PageHandler) GetPages(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	request, err := NewGetPagesRequest(r, p)
	if err != nil {
		api.RespondWith(r, w, http.StatusBadRequest, err, err)
		return
	}
	ctx := r.Context()
	authData, err := api.GetDataFromContext(ctx)
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, errors.Wrap(err, "failed to get auth data"))
		return
	}
	records, total, nextBatchID, err := h.PageService.GetPages(ctx, pageservice.GetPagesParams{
		NextBatchID: request.NextBatchID,
		UserID:      authData.UserID,
	})
	if _, ok := err.(*storeerror.NotAuthorized); ok {
		api.RespondWith(r, w, http.StatusUnauthorized, &api.FailedAuthorization{}, err)
		return
	}
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, err)
		return
	}
	conformedRecords := make([]interface{}, 0)
	for _, record := range records {
		reducedPage := record.Reduce()
		conformedRecords = append(conformedRecords, reducedPage.GetJSONConformed())
	}
	if nextBatchID == "" {
		responseBody := struct {
			Batch []interface{} `json:"batch"`
			Total int           `json:"total"`
		}{
			Batch: conformedRecords,
			Total: total,
		}
		api.RespondWith(r, w, http.StatusOK, responseBody, nil)
		return
	}
	responseBody := struct {
		Batch     []interface{}       `json:"batch"`
		Total     int                 `json:"total"`
		NextBatch nextbatch.NextBatch `json:"nextBatch"`
	}{
		Batch: conformedRecords,
		Total: total,
		NextBatch: nextbatch.NextBatch{
			ParamKey:   "nextBatchId",
			ParamValue: nextBatchID,
		},
	}
	api.RespondWith(r, w, http.StatusOK, responseBody, nil)
}

// DeletePage see Service for more details
func (h PageHandler) DeletePage(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	request, err := NewDeletePageRequest(r, p)
	if err != nil {
		api.RespondWith(r, w, http.StatusBadRequest, err, err)
		return
	}
	ctx := r.Context()
	authData, err := api.GetDataFromContext(ctx)
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, errors.Wrap(err, "failed to get auth data"))
		return
	}
	err = h.PageService.RemovePage(ctx, pageservice.RemovePageParams{
		Page: page.Page{
			GUID: request.GUID,
		},
		UserID: authData.UserID,
	})
	if _, ok := err.(*storeerror.NotAuthorized); ok {
		api.RespondWith(r, w, http.StatusUnauthorized, &api.FailedAuthorization{}, err)
		return
	}
	if err != nil {
		api.RespondWith(r, w, http.StatusInternalServerError, &api.InternalErr{}, err)
		return
	}
	api.RespondWith(r, w, http.StatusOK, nil, nil)
}
