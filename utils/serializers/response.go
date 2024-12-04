package serializers

type ResponseCode uint64

type BasicResponse struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
}

type DataResponse struct {
	Code    ResponseCode `json:"code"`
	Message string       `json:"message"`
	Data    interface{}  `json:"data"`
}

func NewResponse(code ResponseCode, message string, data ...interface{}) interface{} {

	if len(data) == 0 {
		return BasicResponse{
			Code:    code,
			Message: message,
		}
	}

	if len(data) == 1 {
		return DataResponse{
			Code:    code,
			Message: message,
			Data:    data[0],
		}
	}

	return DataResponse{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
