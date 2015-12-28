// +build yaml

package yaml

import (
	auth "../../auth"
	cfg "../../cfg"
	log "../../log"
	model "../../model"
	"fmt"
	"github.com/nmcclain/ldap"
	"gopkg.in/yaml.v1"
	"io/ioutil"
	"net"
	"strings"
)

type YUser struct {
	Name        string
	UnixId      int      `yaml:"uid"`
	GroupId     int      `yaml:"gid"`
	OtherGroups []int    `yaml:"other_groups"`
	Auths       []string `yaml:"auths"`
	SshKeys     []string `yaml:"ssh_keys"`
}

type YGroup struct {
	Name   string
	UnixId int `yaml:"uid"`
}

type YamlLdapDB struct {
	Groups []YGroup
	Users  []YUser
}

type YamlBackendHandler struct {
	model.BackendHandler
	db  YamlLdapDB
	cfg YamlBackendSettings
}

type YamlBackendSettings struct {
	model.BaseBackend
	Path string
}

func NewYamlBackendHandler(name string) (model.BackendHandler, error) {
	settings := YamlBackendSettings{}
	cfg.GetBackendCfg(name, &settings)
	db := YamlLdapDB{}
	content, err := ioutil.ReadFile(settings.Path)
	if err == nil {
		if unmarshalerr := yaml.Unmarshal(content, &db); unmarshalerr != nil {
			return nil, unmarshalerr
		}
	} else {
		return nil, err
	}
	return YamlBackendHandler{db: db, cfg: settings}, nil
}

func (h YamlBackendHandler) Bind(bindDN, bindSimplePw string, conn net.Conn) (resultCode ldap.LDAPResultCode, err error) {
	bindDN = strings.ToLower(bindDN)
	log.Info("Bind Request", "BindDN", bindDN, "BaseDN", h.cfg.BaseDN, "Remote", conn.RemoteAddr().String())
	//stats_frontend.Add("bind_reqs", 1)

	// parse the bindDN
	if !strings.HasSuffix(bindDN, h.cfg.BaseDN) {
		log.Warn(fmt.Sprintf("Bind Error: BindDN %s not our BaseDN %s", bindDN, h.cfg.BaseDN))
		return ldap.LDAPResultInvalidCredentials, nil
	}
	parts := strings.Split(strings.TrimSuffix(bindDN, ","+h.cfg.BaseDN), ",")
	groupName := ""
	userName := ""
	if len(parts) == 1 {
		userName = strings.TrimPrefix(parts[0], "cn=")
	} else if len(parts) == 2 {
		userName = strings.TrimPrefix(parts[0], "cn=")
		groupName = strings.TrimPrefix(parts[1], "ou=")
	} else {
		log.Warn(fmt.Sprintf("Bind Error: BindDN %s should have only one or two parts (has %d)", bindDN, len(parts)))
		return ldap.LDAPResultInvalidCredentials, nil
	}
	// find the user
	user := YUser{}
	found := false
	for _, u := range h.db.Users {
		if u.Name == userName {
			found = true
			user = u
		}
	}
	if !found {
		log.Warn(fmt.Sprintf("Bind Error: User %s not found.", user))
		return ldap.LDAPResultInvalidCredentials, nil
	}
	// find the group
	group := YGroup{}
	found = false
	for _, g := range h.db.Groups {
		if g.Name == groupName {
			found = true
			group = g
		}
	}
	if !found {
		log.Warn(fmt.Sprintf("Bind Error: Group %s not found.", group))
		return ldap.LDAPResultInvalidCredentials, nil
	}
	// validate group membership
	if user.GroupId != group.UnixId {
		log.Warn(fmt.Sprintf("Bind Error: User %s primary group is not %s.", userName, groupName))
		return ldap.LDAPResultInvalidCredentials, nil
	}

	authx := false
	for _, authstr := range user.Auths {
		aparts := strings.Split(authstr, ":")
		backendPassword := ""
		if len(aparts) > 1 {
			backendPassword = aparts[1]
		}
		a, _ := auth.GetAuthHandler(aparts[0])
		if a.Auth(user.Name, backendPassword, bindSimplePw) == nil { // if auth successful
			authx = true
			break
		}
	}
	if !authx {
		log.Warn(fmt.Sprintf("Bind Error: invalid credentials as %s from %s", bindDN, conn.RemoteAddr().String()))
		return ldap.LDAPResultInvalidCredentials, nil
	}

	//stats_frontend.Add("bind_successes", 1)
	log.Info("Bind Success")
	return ldap.LDAPResultSuccess, nil
}

//
func (h YamlBackendHandler) Search(bindDN string, searchReq ldap.SearchRequest, conn net.Conn) (result ldap.ServerSearchResult, err error) {
	bindDN = strings.ToLower(bindDN)
	//baseDN := strings.ToLower("," + h.baseDN)
	searchBaseDN := strings.ToLower(searchReq.BaseDN)
	log.Info("Search Request", "BindDN", bindDN, "From", conn.RemoteAddr().String(), "Query", searchReq.Filter)
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
	entries := []*ldap.Entry{}
	filterEntity, err := ldap.GetFilterObjectClass(searchReq.Filter)
	if err != nil {
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, fmt.Errorf("Search Error: error parsing filter: %s", searchReq.Filter)
	}
	switch filterEntity {
	default:
		return ldap.ServerSearchResult{ResultCode: ldap.LDAPResultOperationsError}, fmt.Errorf("Search Error: unhandled filter type: %s [%s]", filterEntity, searchReq.Filter)
	case "posixgroup":
		for _, g := range h.db.Groups {
			attrs := []*ldap.EntryAttribute{}
			attrs = append(attrs, &ldap.EntryAttribute{"cn", []string{g.Name}})
			attrs = append(attrs, &ldap.EntryAttribute{"description", []string{fmt.Sprintf("%s via LDAP", g.Name)}})
			attrs = append(attrs, &ldap.EntryAttribute{"gidNumber", []string{fmt.Sprintf("%d", g.UnixId)}})
			attrs = append(attrs, &ldap.EntryAttribute{"objectClass", []string{"posixGroup"}})
			attrs = append(attrs, &ldap.EntryAttribute{"uniqueMember", h.getGroupMembers(g.UnixId)})
			attrs = append(attrs, &ldap.EntryAttribute{"memberUid", h.getGroupMemberIDs(g.UnixId)})
			dn := fmt.Sprintf("cn=%s,ou=groups,%s", g.Name, h.cfg.BaseDN)
			entries = append(entries, &ldap.Entry{dn, attrs})
		}
	case "posixaccount", "":
		for _, u := range h.db.Users {
			attrs := []*ldap.EntryAttribute{}
			attrs = append(attrs, &ldap.EntryAttribute{"cn", []string{u.Name}})
			attrs = append(attrs, &ldap.EntryAttribute{"uid", []string{u.Name}})
			attrs = append(attrs, &ldap.EntryAttribute{"ou", []string{h.getGroupName(u.GroupId)}})
			attrs = append(attrs, &ldap.EntryAttribute{"uidNumber", []string{fmt.Sprintf("%d", u.UnixId)}})
			attrs = append(attrs, &ldap.EntryAttribute{"accountStatus", []string{"active"}})
			attrs = append(attrs, &ldap.EntryAttribute{"objectClass", []string{"posixAccount"}})
			attrs = append(attrs, &ldap.EntryAttribute{"homeDirectory", []string{"/home/" + u.Name}})
			attrs = append(attrs, &ldap.EntryAttribute{"loginShell", []string{"/bin/bash"}})
			attrs = append(attrs, &ldap.EntryAttribute{"description", []string{fmt.Sprintf("%s via LDAP", u.Name)}})
			attrs = append(attrs, &ldap.EntryAttribute{"gecos", []string{fmt.Sprintf("%s via LDAP", u.Name)}})
			attrs = append(attrs, &ldap.EntryAttribute{"gidNumber", []string{fmt.Sprintf("%d", u.GroupId)}})
			attrs = append(attrs, &ldap.EntryAttribute{"memberOf", h.getGroupDNs(u.OtherGroups)})
			if len(u.SshKeys) > 0 {
				attrs = append(attrs, &ldap.EntryAttribute{"sshPublicKey", u.SshKeys})
			}
			dn := fmt.Sprintf("cn=%s,ou=%s,%s", u.Name, h.getGroupName(u.GroupId), h.cfg.BaseDN)
			entries = append(entries, &ldap.Entry{dn, attrs})
		}
	}
	//stats_frontend.Add("search_successes", 1)
	log.Info("Search OK")
	return ldap.ServerSearchResult{entries, []string{}, []ldap.Control{}, ldap.LDAPResultSuccess}, nil
}

//
func (h YamlBackendHandler) Close(boundDn string, conn net.Conn) error {
	//stats_frontend.Add("closes", 1)
	return nil
}

//
func (h YamlBackendHandler) getGroupMembers(gid int) []string {
	members := make(map[string]bool)
	for _, u := range h.db.Users {
		if u.GroupId == gid {
			dn := fmt.Sprintf("cn=%s,ou=%s,%s", u.Name, h.getGroupName(u.GroupId), h.cfg.BaseDN)
			members[dn] = true
		} else {
			for _, othergid := range u.OtherGroups {
				if othergid == gid {
					dn := fmt.Sprintf("cn=%s,ou=%s,%s", u.Name, h.getGroupName(u.GroupId), h.cfg.BaseDN)
					members[dn] = true
				}
			}
		}
	}
	m := []string{}
	for k, _ := range members {
		m = append(m, k)
	}
	return m
}

//
func (h YamlBackendHandler) getGroupMemberIDs(gid int) []string {
	members := make(map[string]bool)
	for _, u := range h.db.Users {
		if u.GroupId == gid {
			members[u.Name] = true
		} else {
			for _, othergid := range u.OtherGroups {
				if othergid == gid {
					members[u.Name] = true
				}
			}
		}
	}
	m := []string{}
	for k, _ := range members {
		m = append(m, k)
	}
	return m
}

//
func (h YamlBackendHandler) getGroupDNs(gids []int) []string {
	groups := make(map[string]bool)
	for _, gid := range gids {
		for _, g := range h.db.Groups {
			if g.UnixId == gid {
				dn := fmt.Sprintf("cn=%s,ou=groups,%s", g.Name, h.cfg.BaseDN)
				groups[dn] = true
			}
		}
	}
	g := []string{}
	for k, _ := range groups {
		g = append(g, k)
	}
	return g
}

//
func (h YamlBackendHandler) getGroupName(gid int) string {
	for _, g := range h.db.Groups {
		if g.UnixId == gid {
			return g.Name
		}
	}
	return ""
}
