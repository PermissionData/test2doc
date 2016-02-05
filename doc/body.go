package doc

import (
	"bytes"
	"compress/gzip"
	"io"
	"text/template"
)

var (
	bodyTmpl *template.Template
	bodyFmt  = `    + Body

            {{.FormattedStr}}
`
)

func init() {
	bodyTmpl = template.Must(template.New("body").Parse(bodyFmt))
}

type Body struct {
	Content         []byte
	ContentType     string
	ContentEncoding string
}

func NewBody(content []byte, contentType string, contentEncoding string) (b *Body) {
	if len(content) > 0 {
		b = &Body{
			Content:         content,
			ContentType:     contentType,
			ContentEncoding: contentEncoding,
		}

		b.gzip()
	}
	return b
}

func (b *Body) Render() string {
	return render(bodyTmpl, b)
}

func (b *Body) gzip() {

	if b.ContentEncoding == "gzip" {
		var buf bytes.Buffer

		reader, err := gzip.NewReader(bytes.NewReader(b.Content))

		if err != nil {
			panic(err.Error())
		}

		_, err = io.Copy(&buf, reader)

		reader.Close()

		b.Content = buf.Bytes()

		if err != nil {
			panic(err.Error())
		}
	}
}

func (b *Body) FormattedStr() string {
	if b.ContentType == "application/json" {
		return b.FormattedJSON()
	}
	return string(b.Content)
}

func (b *Body) FormattedJSON() string {
	fbody, err := indentJSONBody(string(b.Content))
	if err != nil {
		panic(err.Error())
	}

	return fbody
}
