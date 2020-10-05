package yaml

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"path"
	"strings"

	"github.com/go-logr/logr"
	"github.com/kernel164/go389/internal/auth"
	"github.com/kernel164/go389/internal/cfg"
	"github.com/kernel164/go389/internal/model"

	ber "github.com/nmcclain/asn1-ber"
	"github.com/nmcclain/ldap"
	"gopkg.in/yaml.v3"
)

// YUserSettings - user settings
type YUserSettings struct {
	DN          string            `yaml:"dn"`
	ObjectClass string            `yaml:"objectclass"`
	BindAttr    string            `yaml:"bindAttr"`
	SearchAttr  string            `yaml:"searchAttr"`
	Alias       map[string]string `yaml:"alias"`
}

// YGroupSettings - group settings
type YGroupSettings struct {
	DN          string            `yaml:"dn"`
	ObjectClass string            `yaml:"objectclass"`
	BindAttr    string            `yaml:"bindAttr"`
	SearchAttr  string            `yaml:"searchAttr"`
	Alias       map[string]string `yaml:"alias"`
}

// YSASettings - service account settings
type YSASettings struct {
	DN       string `yaml:"dn"`
	BindAttr string `yaml:"bindAttr"`
}

// YSettings - settings
type YSettings struct {
	SA    YSASettings    `yaml:"sa"`
	User  YUserSettings  `yaml:"user"`
	Group YGroupSettings `yaml:"group"`
}

// YServiceAccount - service account config
type YServiceAccount struct {
	Auths []string `yaml:"auths"`
}

// YUser - user config
type YUser struct {
	Attrs map[string]interface{} `yaml:"attrs"`
	Alias map[string]string      `yaml:"alias"`
	Auths []string               `yaml:"auths"`
}

// YGroup - group config
type YGroup struct {
	Attrs map[string]interface{} `yaml:"attrs"`
	Alias map[string]string      `yaml:"alias"`
}

// YLdapDB - yaml ldap db
type YLdapDB struct {
	Settings        YSettings                  `yaml:"settings"`
	ServiceAccounts map[string]YServiceAccount `yaml:"serviceAccounts"`
	Groups          map[string]YGroup          `yaml:"groups"`
	Users           map[string]YUser           `yaml:"users"`
}

type handler struct {
	model.BackendHandler
	log logr.Logger
	db  YLdapDB
	cfg BackendSettings
}

// BackendSettings - backend settings
type BackendSettings struct {
	model.BaseBackend
	Path string
}

// New - new yaml backend handler
func New(name string, log logr.Logger, args *model.ServerArgs) (model.BackendHandler, error) {
	settings := BackendSettings{}
	cfg.GetBackendCfg(name, &settings)
	db := YLdapDB{}
	dbPath := settings.Path
	if !strings.HasPrefix(dbPath, "/") {
		dbPath = path.Join(path.Dir(args.Config), settings.Path)
	}
	content, err := ioutil.ReadFile(dbPath)
	if err == nil {
		if unmarshalerr := yaml.Unmarshal(content, &db); unmarshalerr != nil {
			return nil, unmarshalerr
		}
	} else {
		return nil, err
	}
	return &handler{log: log, db: db, cfg: settings}, nil
}

func getAuth(part string, bindAttr string, fn func(string) (string, []string, error)) (string, []string, error) {
	id := ""
	if strings.HasPrefix(part, bindAttr+"=") {
		id = strings.TrimPrefix(part, bindAttr+"=")
	}
	return fn(id)
}

// Bind - handler
func (h *handler) Bind(bindDN, bindSimplePw string, conn net.Conn) (resultCode ldap.LDAPResultCode, err error) {
	bindDN = strings.ToLower(bindDN)
	h.log.Info("Bind Request", "BindDN", bindDN, "BaseDN", h.cfg.BaseDN, "Remote", conn.RemoteAddr().String())
	//stats_frontend.Add("bind_reqs", 1)

	// parse the bindDN
	if !strings.HasSuffix(bindDN, h.cfg.BaseDN) {
		h.log.Info(fmt.Sprintf("Bind Error: BindDN %s not our BaseDN %s", bindDN, h.cfg.BaseDN))
		return ldap.LDAPResultInvalidCredentials, nil
	}
	parts := strings.Split(strings.TrimSuffix(bindDN, ","+h.cfg.BaseDN), ",")
	id, auths, err := getAuth(parts[0], h.db.Settings.SA.BindAttr, func(id string) (string, []string, error) {
		sa, ok := h.db.ServiceAccounts[id]
		if ok {
			return id, sa.Auths, nil
		}
		return id, nil, nil
	})
	if err != nil {
		h.log.Info(fmt.Sprintf("Bind Error: BindDN %s - %s", bindDN, err))
		return ldap.LDAPResultInvalidCredentials, nil
	}
	if auths == nil {
		id, auths, err = getAuth(parts[0], h.db.Settings.User.BindAttr, func(id string) (string, []string, error) {
			if id == "" {
				return "", nil, fmt.Errorf("cannot extract login id")
			}
			user, ok := h.db.Users[id]
			if ok {
				return id, user.Auths, nil
			}
			return id, nil, fmt.Errorf("user not found")
		})
	}
	if err != nil {
		h.log.Info(fmt.Sprintf("Bind Error: BindDN %s - %s", bindDN, err))
		return ldap.LDAPResultInvalidCredentials, nil
	}
	authx := false
	for _, authstr := range auths {
		aparts := strings.Split(authstr, ":")
		backendPassword := ""
		if len(aparts) > 1 {
			backendPassword = aparts[1]
		}
		a, _ := auth.GetAuthHandler(aparts[0])
		if a.Auth(id, backendPassword, bindSimplePw) == nil { // if auth successful
			authx = true
			break
		}
	}
	if !authx {
		h.log.Info(fmt.Sprintf("Bind Error: invalid credentials as %s from %s", bindDN, conn.RemoteAddr().String()))
		return ldap.LDAPResultInvalidCredentials, nil
	}

	//stats_frontend.Add("bind_successes", 1)
	h.log.Info("Bind Success")
	return ldap.LDAPResultSuccess, nil
}

func toStringValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		return v
	case *string:
		return *v
	}
	//h.log.Info(fmt.Sprintf("toStringValue.Type=%T", value))
	return fmt.Sprintf("%v", value)
}

func toStringArray(value interface{}) []string {
	switch v := value.(type) {
	case string:
		return []string{v}
	case []string:
		return v
	case []interface{}:
		x := make([]string, len(v))
		for i, it := range v {
			x[i] = toStringValue(it)
		}
		return x
	}
	//h.log.Info(fmt.Sprintf("toAttrValue.Type=%T", value))
	return []string{fmt.Sprintf("%v", value)}
}

// Search - search handler
func (h *handler) Search(bindDN string, searchReq ldap.SearchRequest, conn net.Conn) (result ldap.ServerSearchResult, err error) {
	bindDN = strings.ToLower(bindDN)
	//baseDN := strings.ToLower("," + h.baseDN)
	searchBaseDN := strings.ToLower(searchReq.BaseDN)
	h.log.Info("Search Request", "BindDN", bindDN, "From", conn.RemoteAddr().String(), "Req", searchReq)
	//stats_frontend.Add("search_reqs", 1)

	// validate the user is authenticated and has appropriate access
	if len(bindDN) < 1 {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultInsufficientAccessRights}, fmt.Errorf("Search Error: Anonymous BindDN not allowed %s", bindDN)
	}
	if !strings.HasSuffix(bindDN, h.cfg.BaseDN) {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultInsufficientAccessRights}, fmt.Errorf("Search Error: BindDN %s not in our BaseDN %s", bindDN, h.cfg.BaseDN)
	}
	if !strings.HasSuffix(searchBaseDN, h.cfg.BaseDN) {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultInsufficientAccessRights}, fmt.Errorf("Search Error: search BaseDN %s is not in our BaseDN %s", searchBaseDN, h.cfg.BaseDN)
	}
	// return all users in the config file - the LDAP library will filter results for us
	f, err := ldap.CompileFilter(searchReq.Filter)
	if err != nil {
		panic(err)
	}
	m := map[string]string{}
	err = extractAttrs(f, m)
	if err != nil {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, fmt.Errorf("Search Error: error parsing filter: %s", searchReq.Filter)
	}
	entries := []*ldap.Entry{}
	objectclass := m["objectclass"]
	switch objectclass {
	case h.db.Settings.User.ObjectClass:
		userID, ok := m[h.db.Settings.User.SearchAttr]
		if !ok {
			h.log.Info(fmt.Sprintf("Search Error: Missing user info in search. attr=%s", h.db.Settings.User.SearchAttr))
			return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, fmt.Errorf("Search Error: user not found %s", userID)
		}
		// find the user
		user, found := h.db.Users[userID]
		if !found {
			h.log.Info(fmt.Sprintf("Search Error: User %s not found.", user))
			return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, fmt.Errorf("Search Error: user not found %s", userID)
		}
		entries = append(entries, h.getUserEntry(userID))
	case h.db.Settings.Group.ObjectClass:
		userName := m[h.db.Settings.User.SearchAttr]
		// find the user
		user, found := h.db.Users[userName] // find user by member
		if !found {
			h.log.Info(fmt.Sprintf("Search Error: User %s not found.", user))
			return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, fmt.Errorf("Search Error: user not found %s", userName)
		}
		groups, ok := user.Attrs["groups"]
		if ok {
			for _, grp := range toStringArray(groups) {
				entries = append(entries, h.getGroupEntry(grp))
			}
		} else {
			h.log.Info(fmt.Sprintf("Search Error: User %s doesn't have groups attribute.", user))
			return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, fmt.Errorf("Search Error: User %s doesn't have groups attribute", userName)
		}
	default:
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, fmt.Errorf("Search Error: unhandled filter type: %s [%s]", objectclass, searchReq.Filter)
	}
	//stats_frontend.Add("search_successes", 1)
	h.log.Info("Search OK")
	/*
		for _, e := range entries {
			h.log.Info(e.DN)
			for _, a := range e.Attributes {
				h.log.Info(a.Name, a.Values)
			}
		}
	*/
	return ldap.ServerSearchResult{
		Entries:    entries,
		Referrals:  []string{},
		Controls:   []ldap.Control{},
		ResultCode: ldap.LDAPResultSuccess,
	}, nil
}

func (h *handler) getUserEntry(userID string) *ldap.Entry {
	userAttrs := map[string]interface{}{}
	userAlias := map[string]string{}
	user, ok := h.db.Users[userID]
	if ok {
		userAttrs = user.Attrs
		userAlias = user.Alias
	}
	if h.db.Settings.User.SearchAttr != "" {
		userAttrs[h.db.Settings.User.SearchAttr] = userID
	}
	return &ldap.Entry{
		DN:         buildDN(userAttrs, h.db.Settings.User.BindAttr, h.db.Settings.User.DN, h.cfg.BaseDN),
		Attributes: buildAttrs(userAttrs, h.db.Settings.User.Alias, userAlias),
	}
}

func (h *handler) getGroupEntry(groupID string) *ldap.Entry {
	groupAttrs := map[string]interface{}{}
	groupAlias := map[string]string{}
	group, ok := h.db.Groups[groupID]
	if ok {
		groupAttrs = group.Attrs
		groupAlias = group.Alias
	}
	if h.db.Settings.Group.SearchAttr != "" {
		groupAttrs[h.db.Settings.Group.SearchAttr] = groupID
	}
	return &ldap.Entry{
		DN:         buildDN(groupAttrs, h.db.Settings.Group.BindAttr, h.db.Settings.Group.DN, h.cfg.BaseDN),
		Attributes: buildAttrs(groupAttrs, h.db.Settings.Group.Alias, groupAlias),
	}
}

func buildAttrs(attrs map[string]interface{}, galias map[string]string, alias map[string]string) []*ldap.EntryAttribute {
	entries := []*ldap.EntryAttribute{}
	for attrKey, attrValue := range attrs {
		av := toStringArray(attrValue)
		entries = append(entries, &ldap.EntryAttribute{Name: attrKey, Values: av})
	}
	for attrKey, attrValue := range alias {
		att, ok := attrs[attrValue]
		if ok {
			av := toStringArray(att)
			entries = append(entries, &ldap.EntryAttribute{Name: attrKey, Values: av})
		}
	}
	for attrKey, attrValue := range galias {
		att, ok := attrs[attrValue]
		if ok {
			av := toStringArray(att)
			entries = append(entries, &ldap.EntryAttribute{Name: attrKey, Values: av})
		}
	}
	return entries
}

func buildDN(attrs map[string]interface{}, attr string, subDN string, baseDN string) string {
	return fmt.Sprintf("%s=%v,%s,%s", attr, attrs[attr], subDN, baseDN)
}

// Close - close
func (h *handler) Close(boundDn string, conn net.Conn) error {
	//stats_frontend.Add("closes", 1)
	return nil
}

func extractAttrs(f *ber.Packet, m map[string]string) error {
	switch ldap.FilterMap[f.Tag] {
	case "Equality Match":
		if len(f.Children) != 2 {
			return errors.New("Equality match must have only two children")
		}
		attribute := strings.ToLower(f.Children[0].Value.(string))
		value := f.Children[1].Value.(string)
		m[attribute] = strings.ToLower(value)
	case "And":
		for _, child := range f.Children {
			err := extractAttrs(child, m)
			if err != nil {
				return err
			}
		}
	case "Or":
		for _, child := range f.Children {
			err := extractAttrs(child, m)
			if err != nil {
				return err
			}
		}
	case "Not":
		if len(f.Children) != 1 {
			return errors.New("Not filter must have only one child")
		}
		err := extractAttrs(f.Children[0], m)
		if err != nil {
			return err
		}
	}
	return nil
}
