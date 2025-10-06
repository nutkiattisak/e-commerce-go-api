package entity

type ResponseData struct {
	Data interface{} `json:"data"`
}

type ResponseSuccess struct {
	Message string      `json:"message"`
	Status  string      `json:"status"`
	Data    interface{} `json:"data,omitempty"`
}

type ResponseError struct {
	Message string `json:"message"`
	Status  string `json:"status"`
}



type ResponsePagination struct {
	Data  interface{} `json:"data"`
	Total int         `json:"total"`
}


