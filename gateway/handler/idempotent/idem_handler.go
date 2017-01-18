package idempotent

import (
	"fmt"
	"net/http"
	"net/url"

	"gopkg.in/unrolled/render.v1"
)

// Middleware defines the typical class of negroni handler
type Middleware struct {
    Before func(*http.Request, bool) *http.Request
    After  func(*http.Request, bool) *http.Request
    // Exclude URLs from idempotent
    excludeURLs []string
}

// NewDefaultMiddeleware returns an instance of idempotent handler with default settings
func NewDefaultMiddeleware() *Middleware {
    return NewMiddleware()
}

// NewMiddleware returns an instace of idempotentn handler
func NewMiddleware() *Middleware {
    return &Middleware {
        Before: DefaultBefore,
        After: DefaultAfter,
    }
}

// ExcludeURL add a new URL u to be ignored during idempotent checking
// The URL u is parsed, hence the returned error
func (m *Middleware) ExcludeURL(u string) error {
    if _, err := url.Parse(u); err != nil {
        return err
	}
	m.excludeURLs = append(m.excludeURLs, u)
	return nil
}

// ExcludeURLs returns the list of excluded URLs
func (m *Middleware) ExcludeURLs() []string {
	return m.excludeURLs
}

// ServeHTTP is the main function of controlling the idempotent
func (m *Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
    // Checking the before and after functions
    if nil == m.Before {
        m.Before = DefaultBefore
    }
    if nil == m.After {
        m.After = DefaultAfter
    }

    for _, u := range m.excludeURLs {
        if r.URL.Path == u {
            return
        }
    }

    // Execute the Before func
    r = m.Before(r, true)

    // Validate idempotent through request unique id
    isvalid, err := m.ValidateRequestIdFromHeaders(r)
    if err != nil {
        jsonr := render.New(render.Options{})
        jsonr.JSON(rw, http.StatusInternalServerError, "Idempotent validating internal errors.")
        r = m.After(r, false)
        return
    }

    if isvalid {
        next(rw, r)
        r = m.After(r, true)
    } else {
        jsonr := render.New(render.Options{})
        jsonr.JSON(rw, http.StatusConflict, "Duplicated requests conflict.")
        r = m.After(r, false)
        return
    }
}

// ValidateRequestIdFromHeaders extracts the X-Request-Id from headers and validates it.
func (m *Middleware) ValidateRequestIdFromHeaders(r *http.Request) (bool, error) {
    reqid := r.Header.Get("X-Request-Id")
    fmt.Println("X-Request-Id: " + reqid)
    return m.ValidateRequestId(reqid)
}

// ValidateRequestId defines the logics of checking whether the request is idempotent
func (m *Middleware) ValidateRequestId(reqid string) (bool, error) {
    return true, nil
}

// BeforeFunc is the customized type of the Before function
type BeforeFunc func(*http.Request, bool) *http.Request
// AfterFunc is the customized type of the After function
type AfterFunc func(*http.Request, bool) *http.Request

// DefaultBefore will do nothing but returns the original request
func DefaultBefore(r *http.Request, isValid bool) *http.Request {
    return r
}

// DefaultAfter will do nothing but return the original request
func DefaultAfter(r *http.Request, isValid bool) *http.Request {
    return r
}
