package helper

import (
	"fmt"
	"net/http"
	"strings"
)

type PaginationResponse struct {
	Page        int `json:"page"`
	Limit       int `json:"limit"`
	Offset      int `json:"offset"`
	TotalRecord int `json:"total_record"`
	TotalPage   int `json:"total_page"`
}

type WithPagination struct {
	Pagination PaginationResponse `json:"pagination"`
	Data       interface{}        `json:"data"`
	Message    string             `json:"message"`
}

func SuccessResponse(code int, message string, data ...any) (int, map[string]interface{}) {
	response := make(map[string]interface{})

	response["message"] = message

	switch len(data) {
	case 1:
		response["data"] = data[0]
	case 2:
		response["data"] = data[0]
		response["token"] = data[1]
	}

	return code, response
}

func ErrorResponse(msg string) (int, interface{}) {
	resp := map[string]interface{}{}
	code := http.StatusInternalServerError

	if msg != "" {
		resp["message"] = msg
	}

	switch true {
	case strings.Contains(msg, "server"):
		code = http.StatusInternalServerError
	case strings.Contains(msg, "format"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "not found"):
		code = http.StatusNotFound
	case strings.Contains(msg, "bad request"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "please upload the"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "conflict"):
		code = http.StatusConflict
	case strings.Contains(msg, "duplicated"):
		code = http.StatusConflict
	case strings.Contains(msg, "syntax"):
		code = http.StatusNotFound
		resp["message"] = "not found"
	case strings.Contains(msg, "input invalid"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "input value"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "validation"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "unmarshal"):
		resp["message"] = "failed to unmarshal json"
		code = http.StatusBadRequest
	case strings.Contains(msg, "upload"):
		code = http.StatusInternalServerError
	case strings.Contains(msg, "denied"):
		code = http.StatusUnauthorized
	case strings.Contains(msg, "jwt"):
		msg = "access is denied due to invalid credential"
		code = http.StatusUnauthorized
	case strings.Contains(msg, "Unauthorized"):
		code = http.StatusUnauthorized
	case strings.Contains(msg, "empty"):
		code = http.StatusBadRequest
	}

	return code, resp
}

func ErrorPage(msg string) (int, string) {
	code := http.StatusInternalServerError

	switch true {
	case strings.Contains(msg, "server"):
		code = http.StatusInternalServerError
	case strings.Contains(msg, "format"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "not found"):
		code = http.StatusNotFound
	case strings.Contains(msg, "conflict"):
		code = http.StatusConflict
	case strings.Contains(msg, "duplicated"):
		code = http.StatusConflict
	case strings.Contains(msg, "syntax"):
		code = http.StatusNotFound
		msg = "not found"
	case strings.Contains(msg, "input invalid"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "input value"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "validation"):
		code = http.StatusBadRequest
	case strings.Contains(msg, "unmarshal"):
		msg = "failed to unmarshal json"
		code = http.StatusBadRequest
	case strings.Contains(msg, "upload"):
		code = http.StatusInternalServerError
	case strings.Contains(msg, "denied"):
		code = http.StatusUnauthorized
	case strings.Contains(msg, "jwt"):
		msg = "access is denied due to invalid credential"
		code = http.StatusUnauthorized
	case strings.Contains(msg, "Unauthorized"):
		code = http.StatusUnauthorized
	}

	baseURL := "https://campyuk.vercel.app"
	css := "margin:auto; position:relative; text-align:center ;top:20%; width:100%; color:#D9D9D9; font-size:40px;"
	html := fmt.Sprintf("<meta http-equiv='refresh' content='3;url=%s/camplist' /><body style='background:#1E3231'><div style ='%s'><h1>ERROR %d</h1><p>%s</p></div></body>", baseURL, css, code, msg)

	return code, html
}
