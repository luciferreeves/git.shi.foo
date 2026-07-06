package meta

import (
	"git.shi.foo/utils/logger"

	"github.com/gofiber/fiber/v2"
)

type Param struct {
	Key   string
	Value string
}

type RequestInfo struct {
	Path        string
	Method      string
	Query       []Param
	Params      []Param
	Headers     []Param
	QueryString string
	IP          string
	URL         string
}

type RequestData struct {
	RequestInfo
	Context *fiber.Ctx
}

func BuildRequest(context *fiber.Ctx) RequestInfo {
	return RequestInfo{
		Path:        context.Path(),
		Method:      context.Method(),
		Query:       buildQueryParams(context),
		Params:      buildRouteParams(context),
		Headers:     buildHeaders(context),
		QueryString: string(context.Request().URI().QueryString()),
		IP:          context.IP(),
		URL:         context.OriginalURL(),
	}
}

func Request(context *fiber.Ctx) *RequestData {
	data, ok := context.Locals(RequestKey).(RequestInfo)
	if !ok {
		logger.Errorf(LogPrefix, RequestContextMissing)
		return nil
	}

	return &RequestData{
		RequestInfo: data,
		Context:     context,
	}
}

func (self *RequestData) Param(key string) string {
	if self == nil || self.Context == nil {
		return ""
	}

	return self.Context.Params(key)
}

func (self *RequestData) Query(key string) string {
	if self == nil {
		return ""
	}

	return findParam(self.RequestInfo.Query, key)
}

func (self *RequestData) Header(key string) string {
	if self == nil {
		return ""
	}

	return findParam(self.RequestInfo.Headers, key)
}

func buildQueryParams(context *fiber.Ctx) []Param {
	params := make([]Param, 0)
	context.Request().URI().QueryArgs().VisitAll(func(name []byte, paramValue []byte) {
		params = append(params, Param{
			Key:   string(name),
			Value: string(paramValue),
		})
	})

	return params
}

func buildRouteParams(context *fiber.Ctx) []Param {
	params := make([]Param, 0)
	for name, routeValue := range context.AllParams() {
		params = append(params, Param{
			Key:   name,
			Value: routeValue,
		})
	}

	return params
}

func buildHeaders(context *fiber.Ctx) []Param {
	params := make([]Param, 0)
	context.Request().Header.VisitAll(func(name []byte, headerValue []byte) {
		params = append(params, Param{
			Key:   string(name),
			Value: string(headerValue),
		})
	})

	return params
}

func findParam(params []Param, key string) string {
	for _, param := range params {
		if param.Key == key {
			return param.Value
		}
	}

	return ""
}
