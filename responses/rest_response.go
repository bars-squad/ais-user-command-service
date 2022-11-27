package responses

import (
	"encoding/json"
	"net/http"
)

func REST(w http.ResponseWriter, r Responses) {
	res := &ResponsesImpl{
		Data:       r.DataProperty(),
		Message:    r.MessageProperty(),
		Status:     r.StatusProperty(),
		Code:       r.CodeProperty(),
		Pagination: r.PaginationProperty(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	json.NewEncoder(w).Encode(res)
}
