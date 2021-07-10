package endpoints

import (
	"context"
	"github.com/datatug/datatug/packages/api"
	"github.com/datatug/datatug/packages/dto"
	"net/http"
)

//var getQueries = api.GetQueries
var getQuery = api.GetQuery

//// GetQueries returns list of project queries
//func GetQueries(w http.ResponseWriter, r *http.Request) {
//	q := r.URL.Query()
//	folder := q.Get(urlQueryParamFolder)
//	ref := newProjectRef(r.URL.Query())
//	ctx, err := getContextFromRequest(r)
//	if err != nil {
//		handleError(err, w, r)
//	}
//	v, err := getQueries(ctx, ref, folder)
//	returnJSON(w, r, http.StatusOK, err, v)
//}

// getQueryHandler returns query definition
func getQueryHandler(w http.ResponseWriter, r *http.Request) {
	var ref dto.ProjectItemRef
	getProjectItem(w, r, &ref, nil, func(ctx context.Context) (responseDTO ResponseDTO, err error) {
		return getQuery(ctx, ref)
	})
}

// createQuery handles create query endpoint
var createQuery = func(w http.ResponseWriter, r *http.Request) {
	var ref dto.ProjectRef
	var request dto.CreateQuery
	saveFunc := func(ctx context.Context) (ResponseDTO, error) {
		return api.CreateQuery(ctx, request)
	}
	createProjectItem(w, r, &ref, &request, saveFunc)
}

// updateQuery handles update query endpoint
func updateQuery(w http.ResponseWriter, r *http.Request) {
	var ref dto.ProjectItemRef
	var request dto.UpdateQuery
	saveFunc := func(ctx context.Context) (ResponseDTO, error) {
		return api.UpdateQuery(ctx, request)
	}
	saveProjectItem(w, r, &ref, &request, saveFunc)
}

// deleteQuery handles delete query endpoint
var deleteQuery = deleteProjItem(api.DeleteQuery)
