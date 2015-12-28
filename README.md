## go389
A Simple LDAP Server

### Features
- Pluggable backend
  - LDAP
  - YAML
- Pluggable Auth
  - Simple (hash based - currently sha256)
  - PAM

### Build
```
go get
go build -i -tags "yaml pam"
```

##### Supported Build Tags
- Backend
  - yaml
- Auth
  - pam
- App Config
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
groups:
  - name: sadmin
    uid:  9001
  - name: admin
    uid:  9002

users:
  - name: suser1
    uid: 9001
    gid: 9001
    auths:
      - "sha256:06004fd4f328fab028833d5156da66649be95afd61d41b690de33c1e3e3941a6" # suser1
      - pam
  - name: user1
    uid: 9002
    gid: 9002
    auths:
      - "sha256:0a041b9462caa4a31bac3567e0b6e6fd9100787db2ab433d96f6d178cabfce90" # user1
```

### References
- https://github.com/nmcclain/ldap
- https://github.com/nmcclain/glauth
