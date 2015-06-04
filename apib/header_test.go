package apib

import "net/http"

func (t *suite) TestRenderHeader() {

}

func (t *suite) TestRenderHeader_HeaderIsNil() {

}

func (t *suite) TestContentType_OneContentType() {
	ct := "text/plain"

	h := http.Header{}
	h.Add("Content-Type", ct)
	t.Equal(Header(h).ContentType(), ct)
}

func (t *suite) TestContentType_MultipleContentTypes_Mistakenly() {
	ct := "text/plain"

	h := http.Header{}
	h.Add("Content-Type", ct)
	h.Add("Content-Type", "application/json")
	t.Equal(Header(h).ContentType(), ct)
}

func (t *suite) TestContentType_NoContentType() {
	h := http.Header{}
	t.Equal(Header(h).ContentType(), "")
}
