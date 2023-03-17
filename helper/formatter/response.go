package responseFormatter

type QueryData struct {
	Code    int         `json:"code"`
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type HttpData struct {
	Code    int         `json:"code"`
	Message interface{} `json:"message"`
	Data    interface{} `json:"data"`
}

func QueryResponse(code int, status bool, message string, data interface{}) QueryData {
	return QueryData{
		Code:    code,
		Status:  status,
		Message: message,
		Data:    data,
	}
}

func HttpResponse(code int, message interface{}, data interface{}) HttpData {
	if code < 100 && code > 599 {
		code = 500
	}

	dataHttp := HttpData{
		Code:    code,
		Message: message,
		Data:    data,
	}

	return dataHttp

}
