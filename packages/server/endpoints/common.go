package endpoints

import (
	"context"
	"encoding/json"
	"github.com/datatug/datatug/packages/dto"
	"net/http"
)

func deleteProjItem(del func(ctx context.Context, ref dto.ProjectItemRef) error) func(w http.ResponseWriter, request *http.Request) {
	return func(w http.ResponseWriter, request *http.Request) {
		query := request.URL.Query()
		ref := newProjectItemRef(query)
		err := del(request.Context(), ref)
		returnJSON(w, request, http.StatusOK, err, true)
	}
}

func saveItem(
	w http.ResponseWriter, r *http.Request,
	target interface{},
	saveFunc func(ctx context.Context, ref dto.ProjectItemRef) (result interface{}, err error),
) {
	projectIemRef := newProjectItemRef(r.URL.Query())

	decoder := json.NewDecoder(r.Body)

	var err error
	if err = decoder.Decode(target); err != nil {
		handleError(err, w, r)
	}
	var result interface{}
	result, err = saveFunc(r.Context(), projectIemRef)
	returnJSON(w, r, http.StatusOK, err, result)
}
