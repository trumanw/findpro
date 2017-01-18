package logrus

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/urfave/negroni"
)

// Timer is the interface for time util
type Timer interface {
	Now() time.Time
	Since(time.Time) time.Duration
}

// Clock is the real instance for Timer
type Clock struct{}

// Now retrieves the current timestamp
func (c *Clock) Now() time.Time {
	return time.Now()
}

// Since returns the duration from the specific time to now.
func (c *Clock) Since(t time.Time) time.Duration {
	return time.Since(t)
}

// Middleware is a middleware handler that logs the request as it goes in and the response as it goes out.
type Middleware struct {
	// Logger is the log.Logger instance used to log messages with the Logger middleware
	Logger *log.Logger
	// Name is the name of the application as recorded in latency metrics
	Name   string
	Before func(*log.Entry, *http.Request, string) *log.Entry
	After  func(*log.Entry, negroni.ResponseWriter, time.Duration, string) *log.Entry

	logStarting bool
	clock       Timer
	// Exclude URLs from logging
	excludeURLs []string
}

// NewDefultMiddleware returns a new instance of logrus handler
func NewDefultMiddleware() *Middleware {
	// return NewMiddleware(log.InfoLevel, &log.TextFormatter{}, "web")
	return NewMiddleware(log.InfoLevel, &log.JSONFormatter{}, "web")
}

// NewMiddlewareWithLogrus returns a new *Middleware which writes to a given logrus logger.
func NewMiddlewareWithLogrus(logger *log.Logger, name string) *Middleware {
	return &Middleware{
		Logger: logger,
		Name:   name,
		Before: DefaultBefore,
		After:  DefaultAfter,

		logStarting: true,
		clock:       &Clock{},
	}
}

// NewMiddleware builds a LogrusHandler with the given level and formatter
func NewMiddleware(level log.Level, formatter log.Formatter, name string) *Middleware {
	logger := log.New()
	logger.Level = level
	logger.Formatter = formatter

	return &Middleware{
		Logger: logger,
		Name:   name,
		Before: DefaultBefore,
		After:  DefaultAfter,

		logStarting: true,
		clock:       &Clock{},
	}
}

// SetLogStarting accepts a bool val to control the logging of the "started handling
// request" prior to passing to the next middleware
func (m *Middleware) SetLogStarting(v bool) {
	m.logStarting = v
}

// ExcludeURL adds a new URL u to be ignored during logging.
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

// ServeHTTP is the main func which can be used by negroni
func (m *Middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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

	// Record the starting time
	start := m.clock.Now()

	// Try to get the real IP
	remoteAddr := r.RemoteAddr
	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		remoteAddr = realIP
	}

	entry := log.NewEntry(m.Logger)

	if reqID := r.Header.Get("X-Request-Id"); reqID != "" {
		entry = entry.WithField("request_id", reqID)
	}

	// Execute Before func
	entry = m.Before(entry, r, remoteAddr)

	if m.logStarting {
		entry.Info("started handling request.")
	}

	next(rw, r)

	// Record the latency of the request request
	latency := m.clock.Since(start)
	res := rw.(negroni.ResponseWriter)

	m.After(entry, res, latency, m.Name).Info("completed handling request.")
}

// BeforeFunc is the func type used to modify or replace the *log.Entry prior
// to calling the next func in the middleware chain
type BeforeFunc func(*log.Entry, *http.Request, string) *log.Entry

// AfterFunc is the func type used to modify or replace the *log.Entry after
// calling the next func in the middleware chain
type AfterFunc func(*log.Entry, negroni.ResponseWriter, time.Duration, string) *log.Entry

// DefaultBefore is the default func assigned to *Middleware.Before
func DefaultBefore(entry *log.Entry, req *http.Request, remoteAddr string) *log.Entry {
	return entry.WithFields(log.Fields{
		"request": req.RequestURI,
		"method":  req.Method,
		"remote":  remoteAddr,
	})
}

// DefaultAfter is the default func assigned to *Middleware.After
func DefaultAfter(entry *log.Entry, res negroni.ResponseWriter, latency time.Duration, name string) *log.Entry {
	return entry.WithFields(log.Fields{
		"status":      res.Status(),
		"text_status": http.StatusText(res.Status()),
		"took":        latency,
		fmt.Sprintf("measure#%s.latency", name): latency.Nanoseconds(),
	})
}
