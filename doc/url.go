package doc

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/adams-sarah/test2doc/doc/parse"
)

type URL struct {
	rawURL            *url.URL
	ParameterizedPath string
	Parameters        []Parameter
}

func NewURL(req *http.Request, fn parse.URLVarExtractor) *URL {
	u := &URL{
		rawURL: req.URL,
	}
	u.ParameterizedPath, u.Parameters = paramPath(req, fn)
	return u
}

func paramPath(req *http.Request, fn parse.URLVarExtractor) (string, []Parameter) {
	uri, err := url.QueryUnescape(req.URL.Path)
	if err != nil {
		// fall back to unescaped uri
		uri = req.URL.Path
	}

	vars := fn(req)

	params := []Parameter{}

	for k, v := range vars {

		uri = strings.Replace(uri, "/"+v, "/{"+k+"}", 1)

		params = append(params, MakeParameter(k, v))
	}

	var queryKeys []string
	queryParams := req.URL.Query()

	log.Println("Query Params: ", queryParams)

	for k, vs := range queryParams {

		queryKeys = append(queryKeys, k)
		// just take first value
		params = append(params, MakeParameter(k, vs[0]))

	}

	var queryKeysStr string
	if len(queryKeys) > 0 {
		queryKeysStr = "{?" + strings.Join(queryKeys, ",") + "}"
	}

	uri = uri + queryKeysStr

	return uri, params
}
