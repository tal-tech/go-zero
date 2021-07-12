package router

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tal-tech/go-zero/rest/httpx"
	"github.com/tal-tech/go-zero/rest/internal/context"
)

const (
	applicationJsonWithUtf8 = "application/json; charset=utf-8"
	contentLength           = "Content-Length"
)

type mockedResponseWriter struct {
	code int
}

func (m *mockedResponseWriter) Header() http.Header {
	return http.Header{}
}

func (m *mockedResponseWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

func (m *mockedResponseWriter) WriteHeader(code int) {
	m.code = code
}

func TestPatRouterHandleErrors(t *testing.T) {
	tests := []struct {
		method string
		path   string
		err    error
	}{
		{"FAKE", "", ErrInvalidMethod},
		{"GET", "", ErrInvalidPath},
	}

	for _, test := range tests {
		t.Run(test.method, func(t *testing.T) {
			router := NewRouter()
			err := router.Handle(test.method, test.path, nil)
			assert.Error(t, ErrInvalidMethod, err)
		})
	}
}

func TestPatRouterNotFound(t *testing.T) {
	var notFound bool
	router := NewRouter()
	router.SetNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notFound = true
	}))
	err := router.Handle(http.MethodGet, "/a/b",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	assert.Nil(t, err)
	r, _ := http.NewRequest(http.MethodGet, "/b/c", nil)
	w := new(mockedResponseWriter)
	router.ServeHTTP(w, r)
	assert.True(t, notFound)
}

func TestPatRouterNotAllowed(t *testing.T) {
	var notAllowed bool
	router := NewRouter()
	router.SetNotAllowedHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAllowed = true
	}))
	err := router.Handle(http.MethodGet, "/a/b",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	assert.Nil(t, err)
	r, _ := http.NewRequest(http.MethodPost, "/a/b", nil)
	w := new(mockedResponseWriter)
	router.ServeHTTP(w, r)
	assert.True(t, notAllowed)
}

func TestPatRouterNotFoundAndNotAllowed(t *testing.T) {
	var notFound, notAllowed bool
	router := GetDefaultRouter()
	router.SetNotFoundHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notFound = true
	}))
	router.SetNotAllowedHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAllowed = true
	}))
	err := router.Handle(http.MethodGet, "/a/b",
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	assert.Nil(t, err)

	// not found
	r, _ := http.NewRequest(http.MethodGet, "/a/c", nil)
	w := new(mockedResponseWriter)
	router.ServeHTTP(w, r)
	assert.True(t, notFound)

	// not allowed
	r, _ = http.NewRequest(http.MethodPost, "/a/b", nil)
	w = new(mockedResponseWriter)
	router.ServeHTTP(w, r)
	assert.True(t, notAllowed)
}

func TestPatRouter(t *testing.T) {
	tests := []struct {
		method string
		path   string
		expect bool
		code   int
		err    error
	}{
		// we don't explicitly set status code, framework will do it.
		{http.MethodGet, "/a/b", true, 0, nil},
		{http.MethodGet, "/a/b/", true, 0, nil},
		{http.MethodGet, "/a/b?a=b", true, 0, nil},
		{http.MethodGet, "/a/b/?a=b", true, 0, nil},
		{http.MethodGet, "/a/b/c?a=b", true, 0, nil},
		{http.MethodGet, "/b/d", false, http.StatusNotFound, nil},
	}

	for _, test := range tests {
		t.Run(test.method+":"+test.path, func(t *testing.T) {
			routed := false
			router := NewRouter()
			err := router.Handle(test.method, "/a/:b", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				routed = true
				assert.Equal(t, 1, len(context.Vars(r)))
			}))
			assert.Nil(t, err)
			err = router.Handle(test.method, "/a/b/c", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				routed = true
				assert.Nil(t, context.Vars(r))
			}))
			assert.Nil(t, err)
			err = router.Handle(test.method, "/b/c", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				routed = true
			}))
			assert.Nil(t, err)

			w := new(mockedResponseWriter)
			r, _ := http.NewRequest(test.method, test.path, nil)
			router.ServeHTTP(w, r)
			assert.Equal(t, test.expect, routed)
			assert.Equal(t, test.code, w.code)

			if test.code == 0 {
				r, _ = http.NewRequest(http.MethodPut, test.path, nil)
				router.ServeHTTP(w, r)
				assert.Equal(t, http.StatusMethodNotAllowed, w.code)
			}
		})
	}
}

func TestParseSlice(t *testing.T) {
	body := `names=%5B%22first%22%2C%22second%22%5D`
	reader := strings.NewReader(body)
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/", reader)
	assert.Nil(t, err)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rt := NewRouter()
	err = rt.Handle(http.MethodPost, "/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := struct {
			Names []string `form:"names"`
		}{}

		err = httpx.Parse(r, &v)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(v.Names))
		assert.Equal(t, "first", v.Names[0])
		assert.Equal(t, "second", v.Names[1])
	}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	rt.ServeHTTP(rr, r)
}

func TestParseJsonPost(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017?nickname=whatever&zipcode=200000",
		bytes.NewBufferString(`{"location": "shanghai", "time": 20170912}`))
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, httpx.ApplicationJson)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		v := struct {
			Name     string `path:"name"`
			Year     int    `path:"year"`
			Nickname string `form:"nickname"`
			Zipcode  int64  `form:"zipcode"`
			Location string `json:"location"`
			Time     int64  `json:"time"`
		}{}

		err = httpx.Parse(r, &v)
		assert.Nil(t, err)
		_, err = io.WriteString(w, fmt.Sprintf("%s:%d:%s:%d:%s:%d", v.Name, v.Year,
			v.Nickname, v.Zipcode, v.Location, v.Time))
		assert.Nil(t, err)
	}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "kevin:2017:whatever:200000:shanghai:20170912", rr.Body.String())
}

func TestParseJsonPostWithIntSlice(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017",
		bytes.NewBufferString(`{"ages": [1, 2], "years": [3, 4]}`))
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, httpx.ApplicationJson)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(func(
		w http.ResponseWriter, r *http.Request) {
		v := struct {
			Name  string  `path:"name"`
			Year  int     `path:"year"`
			Ages  []int   `json:"ages"`
			Years []int64 `json:"years"`
		}{}

		err = httpx.Parse(r, &v)
		assert.Nil(t, err)
		assert.ElementsMatch(t, []int{1, 2}, v.Ages)
		assert.ElementsMatch(t, []int64{3, 4}, v.Years)
	}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseJsonPostError(t *testing.T) {
	payload := `[{"abcd": "cdef"}]`
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017?nickname=whatever&zipcode=200000",
		bytes.NewBufferString(payload))
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, httpx.ApplicationJson)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name     string `path:"name"`
				Year     int    `path:"year"`
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode"`
				Location string `json:"location"`
				Time     int64  `json:"time"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotNil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseJsonPostInvalidRequest(t *testing.T) {
	payload := `{"ages": ["cdef"]}`
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/",
		bytes.NewBufferString(payload))
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, httpx.ApplicationJson)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Ages []int `json:"ages"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotNil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseJsonPostRequired(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017",
		bytes.NewBufferString(`{"location": "shanghai"`))
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, httpx.ApplicationJson)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name     string `path:"name"`
				Year     int    `path:"year"`
				Location string `json:"location"`
				Time     int64  `json:"time"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotNil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParsePath(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/kevin/2017", nil)
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodGet, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name string `path:"name"`
				Year int    `path:"year"`
			}{}

			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
			_, err = io.WriteString(w, fmt.Sprintf("%s in %d", v.Name, v.Year))
			assert.Nil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "kevin in 2017", rr.Body.String())
}

func TestParsePathRequired(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/kevin", nil)
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodGet, "/:name/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name string `path:"name"`
				Year int    `path:"year"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotNil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseQuery(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/kevin/2017?nickname=whatever&zipcode=200000", nil)
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodGet, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode"`
			}{}

			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
			_, err = io.WriteString(w, fmt.Sprintf("%s:%d", v.Nickname, v.Zipcode))
			assert.Nil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "whatever:200000", rr.Body.String())
}

func TestParseQueryRequired(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017?nickname=whatever", nil)
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := struct {
			Nickname string `form:"nickname"`
			Zipcode  int64  `form:"zipcode"`
		}{}

		err = httpx.Parse(r, &v)
		assert.NotNil(t, err)
	}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseOptional(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/kevin/2017?nickname=whatever&zipcode=", nil)
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodGet, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode,optional"`
			}{}

			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
			_, err = io.WriteString(w, fmt.Sprintf("%s:%d", v.Nickname, v.Zipcode))
			assert.Nil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "whatever:0", rr.Body.String())
}

func TestParseNestedInRequestEmpty(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017", bytes.NewBufferString("{}"))
	assert.Nil(t, err)

	type (
		Request struct {
			Name string `path:"name"`
			Year int    `path:"year"`
		}

		Audio struct {
			Volume int `json:"volume"`
		}

		WrappedRequest struct {
			Request
			Audio Audio `json:"audio,optional"`
		}
	)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var v WrappedRequest
			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
			_, err = io.WriteString(w, fmt.Sprintf("%s:%d", v.Name, v.Year))
			assert.Nil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "kevin:2017", rr.Body.String())
}

func TestParsePtrInRequest(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017",
		bytes.NewBufferString(`{"audio": {"volume": 100}}`))
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, httpx.ApplicationJson)

	type (
		Request struct {
			Name string `path:"name"`
			Year int    `path:"year"`
		}

		Audio struct {
			Volume int `json:"volume"`
		}

		WrappedRequest struct {
			Request
			Audio *Audio `json:"audio,optional"`
		}
	)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var v WrappedRequest
			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
			_, err = io.WriteString(w, fmt.Sprintf("%s:%d:%d", v.Name, v.Year, v.Audio.Volume))
			assert.Nil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "kevin:2017:100", rr.Body.String())
}

func TestParsePtrInRequestEmpty(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin", bytes.NewBufferString("{}"))
	assert.Nil(t, err)

	type (
		Audio struct {
			Volume int `json:"volume"`
		}

		WrappedRequest struct {
			Audio *Audio `json:"audio,optional"`
		}
	)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/kevin", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var v WrappedRequest
			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseQueryOptional(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/kevin/2017?nickname=whatever&zipcode=", nil)
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodGet, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode,optional"`
			}{}

			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
			_, err = io.WriteString(w, fmt.Sprintf("%s:%d", v.Nickname, v.Zipcode))
			assert.Nil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "whatever:0", rr.Body.String())
}

func TestParse(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/kevin/2017?nickname=whatever&zipcode=200000", nil)
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodGet, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name     string `path:"name"`
				Year     int    `path:"year"`
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode"`
			}{}

			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
			_, err = io.WriteString(w, fmt.Sprintf("%s:%d:%s:%d", v.Name, v.Year, v.Nickname, v.Zipcode))
			assert.Nil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "kevin:2017:whatever:200000", rr.Body.String())
}

func TestParseWrappedRequest(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/kevin/2017", nil)
	assert.Nil(t, err)

	type (
		Request struct {
			Name string `path:"name"`
			Year int    `path:"year"`
		}

		WrappedRequest struct {
			Request
		}
	)

	router := NewRouter()
	err = router.Handle(http.MethodGet, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var v WrappedRequest
			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
			_, err = io.WriteString(w, fmt.Sprintf("%s:%d", v.Name, v.Year))
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "kevin:2017", rr.Body.String())
}

func TestParseWrappedGetRequestWithJsonHeader(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/kevin/2017", nil)
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, applicationJsonWithUtf8)

	type (
		Request struct {
			Name string `path:"name"`
			Year int    `path:"year"`
		}

		WrappedRequest struct {
			Request
		}
	)

	router := NewRouter()
	err = router.Handle(http.MethodGet, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var v WrappedRequest
			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
			_, err = io.WriteString(w, fmt.Sprintf("%s:%d", v.Name, v.Year))
			assert.Nil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "kevin:2017", rr.Body.String())
}

func TestParseWrappedHeadRequestWithJsonHeader(t *testing.T) {
	r, err := http.NewRequest(http.MethodHead, "http://hello.com/kevin/2017", nil)
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, applicationJsonWithUtf8)

	type (
		Request struct {
			Name string `path:"name"`
			Year int    `path:"year"`
		}

		WrappedRequest struct {
			Request
		}
	)

	router := NewRouter()
	err = router.Handle(http.MethodHead, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var v WrappedRequest
			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
			_, err = io.WriteString(w, fmt.Sprintf("%s:%d", v.Name, v.Year))
			assert.Nil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "kevin:2017", rr.Body.String())
}

func TestParseWrappedRequestPtr(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/kevin/2017", nil)
	assert.Nil(t, err)

	type (
		Request struct {
			Name string `path:"name"`
			Year int    `path:"year"`
		}

		WrappedRequest struct {
			*Request
		}
	)

	router := NewRouter()
	err = router.Handle(http.MethodGet, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var v WrappedRequest
			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
			_, err = io.WriteString(w, fmt.Sprintf("%s:%d", v.Name, v.Year))
			assert.Nil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "kevin:2017", rr.Body.String())
}

func TestParseWithAll(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017?nickname=whatever&zipcode=200000",
		bytes.NewBufferString(`{"location": "shanghai", "time": 20170912}`))
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, httpx.ApplicationJson)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v := struct {
			Name     string `path:"name"`
			Year     int    `path:"year"`
			Nickname string `form:"nickname"`
			Zipcode  int64  `form:"zipcode"`
			Location string `json:"location"`
			Time     int64  `json:"time"`
		}{}

		err = httpx.Parse(r, &v)
		assert.Nil(t, err)
		_, err = io.WriteString(w, fmt.Sprintf("%s:%d:%s:%d:%s:%d", v.Name, v.Year,
			v.Nickname, v.Zipcode, v.Location, v.Time))
		assert.Nil(t, err)
	}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "kevin:2017:whatever:200000:shanghai:20170912", rr.Body.String())
}

func TestParseWithAllUtf8(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017?nickname=whatever&zipcode=200000",
		bytes.NewBufferString(`{"location": "shanghai", "time": 20170912}`))
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, applicationJsonWithUtf8)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name     string `path:"name"`
				Year     int    `path:"year"`
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode"`
				Location string `json:"location"`
				Time     int64  `json:"time"`
			}{}

			err = httpx.Parse(r, &v)
			assert.Nil(t, err)
			_, err = io.WriteString(w, fmt.Sprintf("%s:%d:%s:%d:%s:%d", v.Name, v.Year,
				v.Nickname, v.Zipcode, v.Location, v.Time))
			assert.Nil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)

	assert.Equal(t, "kevin:2017:whatever:200000:shanghai:20170912", rr.Body.String())
}

func TestParseWithMissingForm(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017?nickname=whatever",
		bytes.NewBufferString(`{"location": "shanghai", "time": 20170912}`))
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name     string `path:"name"`
				Year     int    `path:"year"`
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode"`
				Location string `json:"location"`
				Time     int64  `json:"time"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotNil(t, err)
			assert.Equal(t, "field zipcode is not set", err.Error())
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseWithMissingAllForms(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017",
		bytes.NewBufferString(`{"location": "shanghai", "time": 20170912}`))
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name     string `path:"name"`
				Year     int    `path:"year"`
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode"`
				Location string `json:"location"`
				Time     int64  `json:"time"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotNil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseWithMissingJson(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017?nickname=whatever&zipcode=200000",
		bytes.NewBufferString(`{"location": "shanghai"}`))
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name     string `path:"name"`
				Year     int    `path:"year"`
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode"`
				Location string `json:"location"`
				Time     int64  `json:"time"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotEqual(t, io.EOF, err)
			assert.NotNil(t, httpx.Parse(r, &v))
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseWithMissingAllJsons(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/kevin/2017?nickname=whatever&zipcode=200000", nil)
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodGet, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name     string `path:"name"`
				Year     int    `path:"year"`
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode"`
				Location string `json:"location"`
				Time     int64  `json:"time"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotEqual(t, io.EOF, err)
			assert.NotNil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseWithMissingPath(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/2017?nickname=whatever&zipcode=200000",
		bytes.NewBufferString(`{"location": "shanghai", "time": 20170912}`))
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name     string `path:"name"`
				Year     int    `path:"year"`
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode"`
				Location string `json:"location"`
				Time     int64  `json:"time"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotNil(t, err)
			assert.Equal(t, "field name is not set", err.Error())
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseWithMissingAllPaths(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/?nickname=whatever&zipcode=200000",
		bytes.NewBufferString(`{"location": "shanghai", "time": 20170912}`))
	assert.Nil(t, err)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name     string `path:"name"`
				Year     int    `path:"year"`
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode"`
				Location string `json:"location"`
				Time     int64  `json:"time"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotNil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseGetWithContentLengthHeader(t *testing.T) {
	r, err := http.NewRequest(http.MethodGet, "http://hello.com/kevin/2017?nickname=whatever&zipcode=200000", nil)
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, httpx.ApplicationJson)
	r.Header.Set(contentLength, "1024")

	router := NewRouter()
	err = router.Handle(http.MethodGet, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name     string `path:"name"`
				Year     int    `path:"year"`
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode"`
				Location string `json:"location"`
				Time     int64  `json:"time"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotNil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseJsonPostWithTypeMismatch(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017?nickname=whatever&zipcode=200000",
		bytes.NewBufferString(`{"time": "20170912"}`))
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, applicationJsonWithUtf8)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name     string `path:"name"`
				Year     int    `path:"year"`
				Nickname string `form:"nickname"`
				Zipcode  int64  `form:"zipcode"`
				Time     int64  `json:"time"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotNil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func TestParseJsonPostWithInt2String(t *testing.T) {
	r, err := http.NewRequest(http.MethodPost, "http://hello.com/kevin/2017",
		bytes.NewBufferString(`{"time": 20170912}`))
	assert.Nil(t, err)
	r.Header.Set(httpx.ContentType, applicationJsonWithUtf8)

	router := NewRouter()
	err = router.Handle(http.MethodPost, "/:name/:year", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v := struct {
				Name string `path:"name"`
				Year int    `path:"year"`
				Time string `json:"time"`
			}{}

			err = httpx.Parse(r, &v)
			assert.NotNil(t, err)
		}))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, r)
}

func BenchmarkPatRouter(b *testing.B) {
	b.ReportAllocs()

	router := NewRouter()
	router.Handle(http.MethodGet, "/api/:user/:name", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	}))
	w := &mockedResponseWriter{}
	r, _ := http.NewRequest(http.MethodGet, "/api/a/b", nil)
	for i := 0; i < b.N; i++ {
		router.ServeHTTP(w, r)
	}
}
