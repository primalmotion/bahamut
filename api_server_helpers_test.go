package bahamut

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/aporeto-inc/elemental"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAPIServerHelper_corsHandler(t *testing.T) {

	Convey("Given I call the corsHandler", t, func() {

		h := http.Header{}
		h.Add("Origin", "toto")

		w := httptest.NewRecorder()
		corsHandler(w, &http.Request{Header: h})

		Convey("Then the response should be correct", func() {
			So(w.Code, ShouldEqual, http.StatusOK)
		})
	})
}

func TestAPIServerHelper_notFoundHandler(t *testing.T) {

	Convey("Given I call the notFoundHandler", t, func() {

		h := http.Header{}
		h.Add("Origin", "toto")

		w := httptest.NewRecorder()
		notFoundHandler(w, &http.Request{Header: h})

		Convey("Then the response should be correct", func() {
			So(w.Code, ShouldEqual, http.StatusNotFound)
		})
	})
}

func TestAPIServerHelper_writeHTTPResponse(t *testing.T) {

	type Entity struct {
		Name string `json:"name"`
	}

	e1 := &Entity{Name: "e1"}
	e2 := &Entity{Name: "e2"}

	Convey("Given I create Context from a request with pagination info", t, func() {

		u, _ := url.Parse("http://link.com/path?page=2&per_page=10")
		req := &http.Request{
			Host:   "link.com",
			URL:    u,
			Method: http.MethodGet,
		}
		request, _ := elemental.NewRequestFromHTTPRequest(req)

		c := NewContext(elemental.OperationRetrieveMany)
		c.ReadElementalRequest(request)

		c.Count.Total = 40

		Convey("When I write the response from a context with no error for a retrieve", func() {

			w := httptest.NewRecorder()
			c.Operation = elemental.OperationRetrieveMany
			c.OutputData = []*Entity{e1, e2}
			req.Method = http.MethodGet
			writeHTTPResponse(w, c)

			Convey("Then the status code should be default to 200", func() {
				So(w.Code, ShouldEqual, 200)
			})

			Convey("Then the pagination headers should be correct", func() {
				So(w.Header().Get("X-Page-First"), ShouldEqual, "1")
				So(w.Header().Get("X-Page-Prev"), ShouldEqual, "1")
				So(w.Header().Get("X-Page-Next"), ShouldEqual, "3")
				So(w.Header().Get("X-Page-Last"), ShouldEqual, "4")
			})

			Convey("Then the status should be 200", func() {
				So(string(w.Body.Bytes()), ShouldEqual, "[{\"name\":\"e1\"},{\"name\":\"e2\"}]\n")
			})
		})

		Convey("When I write the response from a context with no error for a info", func() {

			w := httptest.NewRecorder()
			c.Operation = elemental.OperationInfo
			req.Method = http.MethodHead
			writeHTTPResponse(w, c)

			Convey("Then the status code should be default to 204", func() {
				So(w.Code, ShouldEqual, 204)
			})

			Convey("Then the pagination headers should be correct", func() {
				So(w.Header().Get("X-Page-First"), ShouldEqual, "1")
				So(w.Header().Get("X-Page-Prev"), ShouldEqual, "1")
				So(w.Header().Get("X-Page-Next"), ShouldEqual, "3")
				So(w.Header().Get("X-Page-Last"), ShouldEqual, "4")
			})

			Convey("Then the body should be empty", func() {
				So(len(w.Body.Bytes()), ShouldEqual, 0)
			})
		})

		Convey("When I write the response from a context with no error for a create", func() {

			w := httptest.NewRecorder()
			c.Operation = elemental.OperationCreate
			req.Method = http.MethodPost
			writeHTTPResponse(w, c)

			Convey("Then the status code should be default to 201", func() {
				So(w.Code, ShouldEqual, 201)
			})
		})

		Convey("When I try write the response with an unmarshallable object", func() {

			w := httptest.NewRecorder()
			c.OutputData = NewUnmarshalableList()
			writeHTTPResponse(w, c)

			Convey("Then err should not be nil", func() {
				So(w.Code, ShouldEqual, http.StatusInternalServerError)
			})
		})
	})
}

func TestAPIServerHelpers_writeHTTPError(t *testing.T) {

	Convey("Given I create a http.ResponseWriter", t, func() {

		w := httptest.NewRecorder()

		Convey("When I use writeHTTPError with a simple elemental.Error", func() {

			writeHTTPError(w, "origin", elemental.NewError("title", "description", "subject", 42))

			Convey("Then the status should be 42", func() {
				So(w.Code, ShouldEqual, 42)
			})

			Convey("Then the body should be correct", func() {
				So(string(w.Body.Bytes()), ShouldEqual, "[{\"code\":42,\"description\":\"description\",\"subject\":\"subject\",\"title\":\"title\",\"data\":null}]\n")
			})
		})

		Convey("When I use writeHTTPError with an elemental.Errors", func() {

			errs := elemental.NewErrors(elemental.NewError("title", "description", "subject", 43))
			writeHTTPError(w, "origin", errs)

			Convey("Then the status should be 43", func() {
				So(w.Code, ShouldEqual, 43)
			})

			Convey("Then the body should be correct", func() {
				So(string(w.Body.Bytes()), ShouldEqual, "[{\"code\":43,\"description\":\"description\",\"subject\":\"subject\",\"title\":\"title\",\"data\":null}]\n")
			})
		})
	})
}

func TestAPIServerHelpers_commonHeaders(t *testing.T) {

	Convey("Given I create a http.ResponseWriter", t, func() {

		w := httptest.NewRecorder()

		Convey("When I use setCommonHeader with a referer", func() {

			setCommonHeader(w, "http://toto.com:8443")

			Convey("Then the common headers should be set", func() {
				So(w.Header().Get("Content-Type"), ShouldEqual, "application/json; charset=UTF-8")
				So(w.Header().Get("Access-Control-Allow-Origin"), ShouldEqual, "http://toto.com:8443")
				So(w.Header().Get("Access-Control-Expose-Headers"), ShouldEqual, "X-Requested-With, X-Count-Local, X-Count-Total, X-PageCurrent, X-Page-Size, X-Page-Prev, X-Page-Next, X-Page-First, X-Page-Last, X-Namespace")
				So(w.Header().Get("Access-Control-Allow-Methods"), ShouldEqual, "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS")
				So(w.Header().Get("Access-Control-Allow-Headers"), ShouldEqual, "Authorization, Content-Type, Cache-Control, If-Modified-Since, X-Requested-With, X-Count-Local, X-Count-Total, X-PageCurrent, X-Page-Size, X-Page-Prev, X-Page-Next, X-Page-First, X-Page-Last, X-Namespace")
				So(w.Header().Get("Access-Control-Allow-Credentials"), ShouldEqual, "true")
			})
		})

		Convey("When I use setCommonHeader without a referer", func() {

			setCommonHeader(w, "")

			Convey("Then the common headers should be set", func() {
				So(w.Header().Get("Content-Type"), ShouldEqual, "application/json; charset=UTF-8")
				So(w.Header().Get("Access-Control-Allow-Origin"), ShouldEqual, "*")
				So(w.Header().Get("Access-Control-Expose-Headers"), ShouldEqual, "X-Requested-With, X-Count-Local, X-Count-Total, X-PageCurrent, X-Page-Size, X-Page-Prev, X-Page-Next, X-Page-First, X-Page-Last, X-Namespace")
				So(w.Header().Get("Access-Control-Allow-Methods"), ShouldEqual, "GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS")
				So(w.Header().Get("Access-Control-Allow-Headers"), ShouldEqual, "Authorization, Content-Type, Cache-Control, If-Modified-Since, X-Requested-With, X-Count-Local, X-Count-Total, X-PageCurrent, X-Page-Size, X-Page-Prev, X-Page-Next, X-Page-First, X-Page-Last, X-Namespace")
				So(w.Header().Get("Access-Control-Allow-Credentials"), ShouldEqual, "true")
			})
		})
	})
}
