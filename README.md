# harborw


### Required environment variables
- LDAP_USERNAME
- LDAP_PASSWORD
- HARBOR_BASEURL
- PORTAINER_BASEURL

```bash
DEBUG=1 LDAP_USERNAME=username LDAP_PASSWORD=password HARBOR_BASEURL=http://localhost:3000 PORTAINER_BASEURL=http://localhost:3000 go run ./...
```
