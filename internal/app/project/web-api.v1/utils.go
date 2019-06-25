package api

import (
	"bytes"
	"context"
	"encoding/json"
	"math"
	"net/http"
	"strconv"
	"strings"

	"gitlab.com/standard-go/project/internal/app/project"
	"gitlab.com/standard-go/project/internal/app/responses"
	"gitlab.com/standard-go/project/internal/app/config/env"
)

/*
====================================
    HELPER FUNCTIONS
====================================
*/

func printError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	ret := responses.SingleResponse{}

	switch err.Error() {
	case "400":
		ret.Meta.StatusCode = http.StatusBadRequest
		ret.Meta.Message = M{
			"errors": "Bad Request",
		}
		break
	case "400-Token-Parse":
		ret.Meta.StatusCode = http.StatusBadRequest
		ret.Meta.Message = M{
			"errors": "Cannot Read JWT Payload",
		}
		break
	case "400-Signature":
		ret.Meta.StatusCode = http.StatusBadRequest
		ret.Meta.Message = M{
			"errors": "Invalid JWT Signature",
		}
		break
	case "401-Expired":
		ret.Meta.StatusCode = http.StatusBadRequest
		ret.Meta.Message = M{
			"errors": "Invalid JWT Signature",
		}
		break
	case "401":
		ret.Meta.StatusCode = http.StatusUnauthorized
		ret.Meta.Message = M{
			"errors": "JWT Token Is Required",
		}
		break
	case "404":
		ret.Meta.StatusCode = http.StatusNotFound
		ret.Meta.Message = M{
			"errors": "Page Not Found",
		}
		break
	case "405":
		ret.Meta.StatusCode = http.StatusMethodNotAllowed
		ret.Meta.Message = M{
			"errors": "Method Not Allowed",
		}
		break
	case "500":
		ret.Meta.StatusCode = http.StatusInternalServerError
		ret.Meta.Message = M{
			"errors": err.Error(),
		}
		break
	default:
		ret.Meta.StatusCode = http.StatusBadRequest
		ret.Meta.Message = M{
			"errors": err.Error(),
		}
		break
	}

	w.WriteHeader(ret.Meta.StatusCode)
	json.NewEncoder(w).Encode(ret)

	return
}

func decimalCheck(max int, value float64) bool {
	stringFloat := strconv.FormatFloat(value, 'f', -1, 64)
	data := strings.Split(stringFloat, ".")
	if len(data) > 2 {
		return false
	}

	if len(data) == 1 {
		return true
	}

	return len(data[1]) <= max
}

func inArray(val string, array []string) (ok bool) {
	for _, data := range array {
		if ok = data == val; ok {
			return
		}
	}
	return
}

func difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

//
// Response helpers
//

// M represents JSON response body.
type M map[string]interface{}

// GetInt returns int from the map for the given key.
func (m M) GetInt(key string) (int, error) {
	val, ok := m[key]
	if !ok {
		return 0, project.ErrMapKeyDoesNotExist
	}
	n, ok := val.(int)
	if !ok {
		return 0, project.ErrUnknownMapValueType
	}

	return n, nil
}

// GetString returns string from the map for the given key.
func (m M) GetString(key string) (string, error) {
	val, ok := m[key]
	if !ok {
		return "", project.ErrMapKeyDoesNotExist
	}
	s, ok := val.(string)
	if !ok {
		return "", project.ErrUnknownMapValueType
	}

	return s, nil
}

func setStatus(r *http.Request, status int) {
	*r = *r.WithContext(context.WithValue(r.Context(), statusCtxKey, status))
}

func setPaginate(r *http.Request, pageRequest *project.PageRequest, data interface{}, count int, pageUrl string) interface{} {
	var res interface{}

	buffer := new(bytes.Buffer)
	isFirstParam := true
	for k, v := range r.URL.Query() {
		if k != "page" {
			if !isFirstParam {
				buffer.WriteString("&")
			} else {
				isFirstParam = false
			}
			buffer.WriteString(k)
			buffer.WriteString("=")
			buffer.WriteString(v[0])
		}
	}

	if pageRequest.Paginate == 1 {
		res = setResponsePagination(pageRequest, data, count, pageUrl, buffer.String(), isFirstParam)
	} else {
		res = responses.NewResponse(data)
		return res
	}

	return res
}

func setResponsePagination(pageRequest *project.PageRequest, data interface{}, count int, pageUrl string, params string, isFirstParam bool) interface{} {
	response := &responses.VueTablePaginateResponse{
		Data: data,
		Meta: responses.SimpleMeta{
			StatusCode: 200,
			Message:    []string{"Success"},
		},
		CurrentPage: pageRequest.Page,
		From:        (pageRequest.Page - 1) * 10 + 1,
		LastPage:    int64(math.Ceil(float64(count) / float64(pageRequest.PerPage))),
		PerPage:     pageRequest.PerPage,
		To:          int64(int(count) % int(pageRequest.PerPage)),
		Total:       int64(count),
		Path: 		 env.Get("APP_HOST") + pageUrl + "?" + params,
	}

	if pageRequest.Page == response.LastPage {
		response.To = pageRequest.Page * 10
	}

	if pageRequest.Page >= 1 && pageRequest.Page < response.LastPage {
		nextPage := "page=" + strconv.Itoa(int(pageRequest.Page)+1)
		if !isFirstParam {
			nextPage = "&" + nextPage
		}
		response.NextPageUrl = env.Get("APP_HOST") + pageUrl + "?" + params + nextPage
	}

	if pageRequest.Page > 1 && pageRequest.Page <= response.LastPage {
		prevPage := "page=" + strconv.Itoa(int(pageRequest.Page)-1)
		if !isFirstParam {
			prevPage = "&" + prevPage
		}
		response.PrevPageUrl = env.Get("APP_HOST") + pageUrl + "?" + params + prevPage
	}

	return response;
}
