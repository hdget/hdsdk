package ws

// HttpMethod http 方法
type HttpMethod int

const (
	HttpMethodGet HttpMethod = iota
	HttpMethodPost
)

var (
	name2HttpMethod = map[string]HttpMethod{
		"GET":  HttpMethodGet,
		"POST": HttpMethodPost,
	}

	httpMethod2name = map[HttpMethod]string{
		HttpMethodGet:  "GET",
		HttpMethodPost: "POST",
	}
)

func ToHttpMethodName(method HttpMethod) string {
	return httpMethod2name[method]
}

func ToHttpMethod(name string) HttpMethod {
	return name2HttpMethod[name]
}
