package HttpFC

type HTTPRequest struct {
	Path                  string              `json:"path"`
	HTTPMethod            string              `json:"httpMethod"`
	Headers               map[string]string   `json:"headers"`
	MultiValueHeaders     map[string][]string `json:"multiValueHeaders"`
	QueryStringParameters map[string]string   `json:"queryStringParameters"`
	RequestContext        map[string]string   `json:"requestContext"`
	Body                  string              `json:"body"`
	IsBase64Encoded       bool                `json:"isBase64Encoded"`
}

// i cant find the real model so i fake it emmm no!
