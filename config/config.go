package config

const (
	FromHeader = 1
	FromCtx    = 2
	FromMock   = 3
	FromCustom = 4
)

type Config struct {
	Input  Input
	Output Output
}

type Output struct {
	RequestIDFromMock   bool
	RequestIDFromCustom string
}

type Input struct {
	RequestIDOrigin     int
	RequestIDFromCustom string

	AccountIDOrigin     int
	AccountIDFromCustom string

	DomainNameOrigin     int
	DomainNameFromCustom string

	DomainPrefixOrigin     int
	DomainPrefixFromCustom string
}

//RequestContext.RequestId = StringPtr(c.Request.Header.Get("X-Fc-Request-Id"))
//RequestContext.AccountId = StringPtr(c.Request.Header.Get("X-Fc-Account-Id"))
//
//RequestContext.DomainName = StringPtr(c.Request.Header.Get("X-Fc-Domain-Name"))
//RequestContext.DomainPrefix = StringPtr(c.Request.Header.Get("X-Fc-Domain-Prefix"))
