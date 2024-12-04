package consts

import "github.com/mehakhanaa/complex-micro-blog/utils/serializers"

const (
	SUCCESS serializers.ResponseCode = 0

	SERVER_ERROR serializers.ResponseCode = 1

	PARAMETER_ERROR serializers.ResponseCode = 2

	AUTH_ERROR serializers.ResponseCode = 3

	NETWORK_ERROR serializers.ResponseCode = 4
)
