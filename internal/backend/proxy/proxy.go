package proxy

import (
	"crypto/sha256"
	"fmt"
	"net"
	"sync"

	"github.com/go-logr/logr"
	"github.com/kernel164/go389/internal/cfg"
	"github.com/kernel164/go389/internal/model"

	"github.com/nmcclain/ldap"
)

type session struct {
	id   string
	c    net.Conn
	ldap *ldap.Conn
}

type handler struct {
	model.BackendHandler
	log      logr.Logger
	sessions map[string]session
	lock     sync.Mutex
	ldapHost string
}

// BackendSettings - backend settings
type BackendSettings struct {
	model.BaseBackend
	Host string
}

// New - new proxy backend handler
func New(name string, log logr.Logger) (model.BackendHandler, error) {
	settings := BackendSettings{}
	cfg.GetBackendCfg(name, &settings)
	return &handler{sessions: make(map[string]session), log: log, ldapHost: settings.Host}, nil
}

func (h *handler) Bind(bindDN, bindSimplePw string, conn net.Conn) (ldap.LDAPResultCode, error) {
	s, err := h.getSession(conn)
	if err != nil {
		return ldap.LDAPResultOperationsError, err
	}
	if err := s.ldap.Bind(bindDN, bindSimplePw); err != nil {
		return ldap.LDAPResultOperationsError, err
	}
	return ldap.LDAPResultSuccess, nil
}

func (h *handler) Search(boundDN string, searchReq ldap.SearchRequest, conn net.Conn) (ldap.ServerSearchResult, error) {
	s, err := h.getSession(conn)
	if err != nil {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, nil
	}
	search := ldap.NewSearchRequest(
		searchReq.BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		searchReq.Filter,
		searchReq.Attributes,
		nil)
	sr, err := s.ldap.Search(search)
	if err != nil {
		return ldap.ServerSearchResult{}, err
	}
	//log.Printf("P: Search OK: %s -> num of entries = %d\n", search.Filter, len(sr.Entries))
	return ldap.ServerSearchResult{sr.Entries, []string{}, []ldap.Control{}, ldap.LDAPResultSuccess}, nil
}

func (h *handler) Close(boundDN string, conn net.Conn) error {
	conn.Close() // close connection to the server when then client is closed
	h.lock.Lock()
	defer h.lock.Unlock()
	delete(h.sessions, connID(conn))
	return nil
}

func connID(conn net.Conn) string {
	h := sha256.New()
	h.Write([]byte(conn.LocalAddr().String() + conn.RemoteAddr().String()))
	sha := fmt.Sprintf("% x", h.Sum(nil))
	return string(sha)
}

func (h *handler) getSession(conn net.Conn) (session, error) {
	id := connID(conn)
	h.lock.Lock()
	s, ok := h.sessions[id] // use server connection if it exists
	h.lock.Unlock()
	if !ok { // open a new server connection if not
		l, err := ldap.Dial("tcp", h.ldapHost)
		if err != nil {
			return session{}, err
		}
		s = session{id: id, c: conn, ldap: l}
		h.lock.Lock()
		h.sessions[s.id] = s
		h.lock.Unlock()
	}
	return s, nil
}
