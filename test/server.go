package test

import (
	"log"
	"net/http"
	"net/http/httptest"

	"github.com/adams-sarah/test2doc/doc"
	"github.com/adams-sarah/test2doc/doc/parse"
)

// resources = map[uri]Resource
//var resources = map[string]*doc.Resource{}

type BlueprintGenerator struct {
	Resources map[string]*doc.Resource
	doc       *doc.Doc
}

var generator BlueprintGenerator

type Server struct {
	*httptest.Server
	gen *BlueprintGenerator
}


// TODO: filter out 404 responses
func NewServer(handler http.Handler, fn parse.URLVarExtractor) (s *Server) {
	// check if url var extractor func is set
	if fn == nil {
		panic("please set a URLVarExtractor.")
	}
	
	httptestServer := httptest.NewServer(handleAndRecord(handler, generator.doc, fn))

	return &Server{
		httptestServer,
		&generator,
	}
}

func BeginGenerating(){
	generator = BlueprintGenerator{}
	generator.Resources = make(map[string]*doc.Resource)
	generator.doc, _ = doc.NewDoc(".")
}

func FinishGenerating(){
	err := generator.doc.Write()
	if err != nil {
		panic(err.Error())
	}
}

func (s *Server) Finish() {
	s.Close()

	for _, r := range s.gen.Resources {
		s.gen.doc.AddResource(r)
	}

}

func handleAndRecord(handler http.Handler, outDoc *doc.Doc, fn parse.URLVarExtractor) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// copy request body into Request object
		docReq, err := doc.NewRequest(req)
		if err != nil {
			log.Println("Error:", err.Error())
			return
		}

		// record response
		rw := httptest.NewRecorder()
		resp := NewResponseWriter(rw)

		handler.ServeHTTP(resp, req)
						
		if (resp.W.Code < 400 || resp.W.Code >= 500) {
								
			// setup resource
			u := doc.NewURL(req, fn)
			path := u.ParameterizedPath
	
			if generator.Resources[path] == nil {
				generator.Resources[path] = doc.NewResource(u)
			}
	
			// store response body in Response object
			docResp := doc.NewResponse(resp.W)
	
			// find action
			action := generator.Resources[path].FindAction(req.Method)
			if action == nil {
				// make new action
				action, err = doc.NewAction(req.Method, resp.HandlerInfo.FuncName)
				if err != nil {
					log.Println("Error:", err.Error())
					return
				}
	
				// add Action to Resource's list of Actions
				generator.Resources[path].AddAction(action)
			}
	
			// add request, response to action
			action.AddRequest(docReq, docResp)
	
			// copy response over to w
			doc.CopyHeader(w.Header(), resp.Header())
		}
		
		w.WriteHeader(resp.W.Code)
		w.Write(resp.W.Body.Bytes())
	}
}
