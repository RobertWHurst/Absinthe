package abshttp

const ABSINTHE_SUBJECT_NAMESPACE = "ABS"

func NewSubjects(method string, path string, requestId string) Subjects {
	base := ABSINTHE_SUBJECT_NAMESPACE + ".HTTP"
	tail := method + "." + path + "." + requestId
	return Subjects{
		Request:         NewRequestSubject(method, path),
		GetRequestBody:  base + ".GET_REQ_BODY." + tail,
		RequestBody:     base + ".REQ_BODY." + tail,
		Response:        base + ".RES." + tail,
		GetResponseBody: base + ".GET_RES_BODY." + tail,
		ResponseBody:    base + ".RES_BODY." + tail,
	}
}

func NewRequestSubject(method string, path string) string {
	return ABSINTHE_SUBJECT_NAMESPACE + ".HTTP.REQ." + method + "." + path
}

type Subjects struct {
	Request         string
	GetRequestBody  string
	RequestBody     string
	Response        string
	GetResponseBody string
	ResponseBody    string
}
