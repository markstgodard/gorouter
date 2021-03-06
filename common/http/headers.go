package http

import (
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"

	"code.cloudfoundry.org/gorouter/common/secure"
	"code.cloudfoundry.org/gorouter/common/uuid"
	"code.cloudfoundry.org/lager"
)

const (
	VcapBackendHeader     = "X-Vcap-Backend"
	CfRouteEndpointHeader = "X-Cf-RouteEndpoint"
	VcapRouterHeader      = "X-Vcap-Router"
	VcapRequestIdHeader   = "X-Vcap-Request-Id"
	VcapTraceHeader       = "X-Vcap-Trace"
	CfInstanceIdHeader    = "X-CF-InstanceID"
	B3TraceIdHeader       = "X-B3-TraceId"
	B3SpanIdHeader        = "X-B3-SpanId"
	CfAppInstance         = "X-CF-APP-INSTANCE"
)

func SetVcapRequestIdHeader(request *http.Request, logger lager.Logger) {
	guid, err := uuid.GenerateUUID()
	if err == nil {
		request.Header.Set(VcapRequestIdHeader, guid)
		if logger != nil {
			logger.Debug("vcap-request-id-header-set", lager.Data{VcapRequestIdHeader: guid})
		}
	}
}

func SetTraceHeaders(responseWriter http.ResponseWriter, routerIp, addr string) {
	responseWriter.Header().Set(VcapRouterHeader, routerIp)
	responseWriter.Header().Set(VcapBackendHeader, addr)
	responseWriter.Header().Set(CfRouteEndpointHeader, addr)
}
func SetB3SpanIdHeader(request *http.Request, logger lager.Logger) {
	randBytes, err := secure.RandomBytes(8)
	if err != nil {
		logger.Debug("failed-to-create-b3-span-id")
		return
	}
	id := hex.EncodeToString(randBytes)
	request.Header.Set(B3SpanIdHeader, id)
	if logger != nil {
		logger.Debug("b3-span-id-header-set", lager.Data{B3SpanIdHeader: id})
	}
}
func SetB3Headers(request *http.Request, logger lager.Logger) {
	SetB3TraceIdHeader(request, logger)
	SetB3SpanIdHeader(request, logger)
}
func SetB3TraceIdHeader(request *http.Request, logger lager.Logger) {
	existingTraceId := request.Header.Get(B3TraceIdHeader)
	if existingTraceId != "" {
		if logger != nil {
			logger.Debug("b3-trace-id-header-exists", lager.Data{B3TraceIdHeader: existingTraceId})
		}
		return
	}

	randBytes, err := secure.RandomBytes(8)
	if err != nil {
		logger.Debug("failed-to-create-b3-trace-id")
		return
	}
	id := hex.EncodeToString(randBytes)
	request.Header.Set(B3TraceIdHeader, id)
	if logger != nil {
		logger.Debug("b3-trace-id-header-set", lager.Data{B3TraceIdHeader: id})
	}
}

func ValidateCfAppInstance(appInstanceHeader string) (string, string, error) {
	appDetails := strings.Split(appInstanceHeader, ":")
	if len(appDetails) != 2 {
		return "", "", fmt.Errorf("Incorrect %s header : %s", CfAppInstance, appInstanceHeader)
	}

	if appDetails[0] == "" || appDetails[1] == "" {
		return "", "", fmt.Errorf("Incorrect %s header : %s", CfAppInstance, appInstanceHeader)
	}

	return appDetails[0], appDetails[1], nil
}
