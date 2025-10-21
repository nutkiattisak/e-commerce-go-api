package errmap

import "errors"

var (
	ErrProvinceIDRequired = errors.New("provinceId is required")
	ErrInvalidProvinceID  = errors.New("invalid province id")
	ErrDistrictIDRequired = errors.New("districtId is required")
	ErrInvalidDistrictID  = errors.New("invalid district id")
)
