package ali

import (
	"context"

	"github.com/hashicorp/go-cleanhttp"
	"github.com/hashicorp/vault/logical"
	"github.com/hashicorp/vault/logical/framework"
	"net/http"
)

func Factory(ctx context.Context, conf *logical.BackendConfig) (logical.Backend, error) {
	client := cleanhttp.DefaultClient()
	client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	b := newBackend(client)
	if err := b.Setup(ctx, conf); err != nil {
		return nil, err
	}
	return b, nil
}

// newBackend exists for testability. It allows us to inject a fake client.
func newBackend(client *http.Client) *backend {
	b := &backend{
		getCallerIdentityClient: client,
		roleMgr:                 NewRoleManager(),
	}
	b.Backend = &framework.Backend{
		AuthRenew: b.pathLoginRenew,
		Help:      backendHelp,
		PathsSpecial: &logical.Paths{
			Unauthenticated: []string{
				"login",
			},
		},
		Paths: []*framework.Path{
			pathLogin(b),
			pathListRole(b),
			pathListRoles(b),
			pathRole(b),
		},
		BackendType: logical.TypeCredential,
	}
	return b
}

type backend struct {
	*framework.Backend

	getCallerIdentityClient *http.Client
	roleMgr                 *RoleManager
}

const backendHelp = `
That Alibaba RAM auth method allows entities to authenticate based on their
identity and pre-configured roles.
`
