package golum

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Server : TODO
type Server struct {
	SocketName   string
	Handler      Handler
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	inShutdown int32 // accessed atomically (non-zero means we're in Shutdown)
	mu         sync.Mutex
	listeners  map[*net.Listener]struct{}
	activeConn map[*conn]struct{}
}

func (s *Server) trackListener(l *net.Listener, add bool) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.listeners == nil {
		s.listeners = make(map[*net.Listener]struct{})
	}
	if add {
		if s.shuttingDown() {
			return false
		}
		s.listeners[l] = struct{}{}
	} else {
		delete(s.listeners, l)
	}
	return true
}

func (s *Server) trackConn(c *conn, add bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.activeConn == nil {
		s.activeConn = make(map[*conn]struct{})
	}
	if add {
		s.activeConn[c] = struct{}{}
	} else {
		delete(s.activeConn, c)
	}
}

// ErrServerClosed : TODO
var ErrServerClosed = errors.New("server closed")

func (s *Server) shuttingDown() bool {
	return atomic.LoadInt32(&s.inShutdown) != 0
}

// Shutdown : TODO
func (s *Server) Shutdown(ctx context.Context) error {
	atomic.StoreInt32(&s.inShutdown, 1)

	s.mu.Lock()
	lerr := s.closeListenersLocked()
	s.mu.Unlock()

	for {
		if s.closeIdleConns() {
			return lerr
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

func (s *Server) closeIdleConns() bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	quiescent := true
	for c := range s.activeConn {
		st, unixSec := c.getState()
		if st == StateNew && unixSec < time.Now().Unix()-5 {
			st = StateIdle
		}
		if st != StateIdle || unixSec == 0 {
			quiescent = false
			continue
		}
		c.rwc.Close()
		delete(s.activeConn, c)
	}
	return quiescent
}

func (s *Server) closeListenersLocked() error {
	var err error
	for l := range s.listeners {
		if cerr := (*l).Close(); cerr != nil && err == nil {
			err = cerr
		}
		delete(s.listeners, l)
	}
	return err
}

// ListenAndServe : TODO
func (s *Server) ListenAndServe() error {
	_ = os.Remove(s.SocketName)

	if s.shuttingDown() {
		return ErrServerClosed
	}

	socketName := s.SocketName
	if socketName == "" {
		socketName = "golum.sock"
	}
	l, err := net.Listen("unix", socketName)
	if err != nil {
		return err
	}
	return s.Serve(l)
}

// Serve : TODO
func (s *Server) Serve(l net.Listener) error {
	l = &onceCloseListener{Listener: l}
	defer l.Close()

	if !s.trackListener(&l, true) {
		return ErrServerClosed
	}
	defer s.trackListener(&l, false)

	ctx := context.Background()

	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}

		c := s.newConn(conn)
		c.setState(c.rwc, StateNew)
		go c.serve(ctx)
	}
}

// onceCloseListener wraps a net.Listener, protecting it from multiple Close calls.
type onceCloseListener struct {
	net.Listener
	once     sync.Once
	closeErr error
}

func (o *onceCloseListener) Close() error {
	o.once.Do(o.close)
	return o.closeErr
}

func (o *onceCloseListener) close() {
	o.closeErr = o.Listener.Close()
}

type conn struct {
	server    *Server
	cancelCtx context.CancelFunc
	rwc       net.Conn

	bufr *bufio.Reader
	bufw *bufio.Writer

	curState struct{ atomic uint64 }
}

func (s *Server) newConn(rwc net.Conn) *conn {
	c := &conn{
		server: s,
		rwc:    rwc,
	}
	return c
}

func (c *conn) serve(ctx context.Context) {
	ctx, cancelCtx := context.WithCancel(ctx)
	defer cancelCtx()

	c.bufr = bufio.NewReader(c.rwc)
	c.bufw = bufio.NewWriter(c.rwc)

	w, err := c.readRequest(ctx)
	c.setState(c.rwc, StateActive)
	if err != nil {
		fmt.Fprintf(c.rwc, "erroe: bad request")
		return
	}

	serverHandler{c.server}.ServeSocket(w, w.req)
	w.cancelCtx()
	w.finishRequest()
	c.setState(c.rwc, StateIdle)
}

func (c *conn) setState(nc net.Conn, state ConnState) {
	s := c.server
	switch state {
	case StateNew:
		s.trackConn(c, true)
	case StateClosed:
		s.trackConn(c, false)
	}
	if state > 0xff || state < 0 {
		return
	}
	packedState := uint64(time.Now().Unix()<<8) | uint64(state)
	atomic.StoreUint64(&c.curState.atomic, packedState)
}

func (c *conn) getState() (state ConnState, unixSec int64) {
	packedState := atomic.LoadUint64(&c.curState.atomic)
	return ConnState(packedState & 0xff), int64(packedState >> 8)
}

func (c *conn) finalFlush() {
	if c.bufr != nil {
		c.bufr = nil
	}
	if c.bufw != nil {
		c.bufw.Flush()
		c.bufw = nil
	}
}

func (c *conn) close() {
	c.finalFlush()
	c.rwc.Close()
}

// ConnState : TODO
type ConnState int

const (
	// StateNew : TODO
	StateNew ConnState = iota
	// StateActive : TODO
	StateActive
	// StateIdle : TODO
	StateIdle
	// StateClosed : TODO
	StateClosed
)

var stateName = map[ConnState]string{
	StateNew:    "new",
	StateActive: "active",
	StateIdle:   "idle",
	StateClosed: "closed",
}

func (c ConnState) String() string {
	return stateName[c]
}

type response struct {
	conn      *conn
	cancelCtx context.CancelFunc
	req       *Request
}

func (w *response) Write(data []byte) (n int, err error) {
	if len(data) == 0 {
		return 0, nil
	}
	return w.conn.bufw.Write(data)
}

func (w *response) finishRequest() {
	w.conn.bufw.Flush()
	w.conn.close()
}

func (c *conn) readRequest(ctx context.Context) (*response, error) {
	var deadline time.Time
	t0 := time.Now()
	if d := c.server.ReadTimeout; d != 0 {
		deadline = t0.Add(d)
	}
	c.rwc.SetReadDeadline(deadline)
	if d := c.server.WriteTimeout; d != 0 {
		defer func() {
			c.rwc.SetWriteDeadline(time.Now().Add(d))
		}()
	}

	req, err := readRequest(c.bufr)
	if err != nil {
		return nil, err
	}

	ctx, cancelCtx := context.WithCancel(ctx)

	w := &response{
		conn:      c,
		cancelCtx: cancelCtx,
		req:       req,
	}

	return w, nil
}

type serverHandler struct {
	srv *Server
}

func (s serverHandler) ServeSocket(w io.Writer, r *Request) {
	h := s.srv.Handler
	if h == nil {
		h = DefaultServeMux
	}
	h.ServeSocket(w, r)
}

// Handler : TODO
type Handler interface {
	ServeSocket(io.Writer, *Request)
}

// HandlerFunc : TODO
type HandlerFunc func(w io.Writer, r *Request)

// ServeSocket : TODO
func (f HandlerFunc) ServeSocket(w io.Writer, r *Request) {
	f(w, r)
}

// HandleFunc : TODO
func HandleFunc(typ string, handler func(w io.Writer, r *Request)) {
	DefaultServeMux.HandleFunc(typ, handler)
}

// ServeMux : TODO
type ServeMux struct {
	mu sync.RWMutex
	m  map[string]Handler
}

// ServeSocket : TODO
func (m *ServeMux) ServeSocket(w io.Writer, r *Request) {
	h, _ := m.Handler(r)
	h.ServeSocket(w, r)
}

// HandleFunc : TODO
func (m *ServeMux) HandleFunc(typ string, handler func(w io.Writer, r *Request)) {
	if handler == nil {
		return
	}
	m.Handle(typ, HandlerFunc(handler))
}

// Handle : TODO
func (m *ServeMux) Handle(typ string, handler Handler) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if typ == "" {
		return
	}
	if handler == nil {
		return
	}
	if _, exist := m.m[typ]; exist {
		return
	}
	if m.m == nil {
		m.m = make(map[string]Handler)
	}
	m.m[typ] = handler
}

// Handler : TODO
func (m *ServeMux) Handler(r *Request) (Handler, string) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	h, ok := m.m[r.Type]
	if !ok {
		return NotFoundHandler(), ""
	}
	return h, r.Type
}

// DefaultServeMux : TODO
var DefaultServeMux = &ServeMux{}

// NotFoundHandler : TODO
func NotFoundHandler() Handler { return HandlerFunc(NotFound) }

// NotFound : TODO
func NotFound(w io.Writer, r *Request) {
	RenderError(w, errors.New("not found"))
}
