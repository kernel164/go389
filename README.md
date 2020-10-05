## go389
A Simple LDAP Server

### Features
- Pluggable backend
  - LDAP
  - YAML
- Pluggable Auth
  - Simple (hash based - sha256,md5,bcrypt)
  - PAM

### Build
```
task build
```

##### Supported Build Tags
- Backend
  - yaml
- Auth
  - pam
- App Config
  - yaml
  - ini

### Sample App Config
```yaml
server: ldap
backend: yaml

servers:
  ldap:
    Bind: "localhost:8389"

backends:
  yaml:
    Path: db.yml
    BaseDN: "dc=example,dc=com"

auths:
  pam:
    Service: go389
```

### Sample Backend YAML DB
```yaml

settings:
  sa:
    dn: cn=sa
    bindAttr: sa
  user:
    dn: cn=users
    objectclass: user
    bindAttr: mail
    searchAttr: mail
  group:
    dn: cn=groups
    objectclass: group
    bindAttr: cn
    searchAttr: name
    alias:
      cn: name

serviceAccounts:
  sa1:
    auths:
      - "bcrypt:$2a$10$dsfdsfdsf.dsfdsfdsfsfdf"

users:
  test@example.com:
    attrs:
      name: test
      groups:
        - admin
    alias:
      uid: name
    auths:
      - "sha256:asdfgewy45645645645"
```

### References
- https://github.com/nmcclain/ldap
- https://github.com/nmcclain/glauth
