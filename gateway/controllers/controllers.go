package controllers

import (
	res "gateway/dto/response"
	"gateway/tracing"

	"github.com/gin-gonic/gin"
	ot "github.com/opentracing/opentracing-go"
)

// SendStandardGatewayResponse takes a gin context, a span, error code and an error message.
// It adds the error code and error message to the span, and sends a standard gateway response to the client.
func SendStandardGatewayResponse(c *gin.Context, span ot.Span, errorCode int32, errorMsg string) {
	// add the resulting error code to the span
	AddErrorTagsToSpan(span, errorCode, errorMsg)
	c.JSON(200, res.GatewayResponse{ErrorCode: errorCode})
}

// AddErrorTagsToSpan adds the custom service.errorCode and service.errorMsg tags to the given span
func AddErrorTagsToSpan(span ot.Span, errorCode int32, errorMsg string) {
	span.SetTag(tracing.ServiceErrorCode, errorCode)
	span.SetTag(tracing.ServiceErrorMsg, errorMsg)
}
