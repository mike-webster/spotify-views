package spotify

import (
	"context"
	"fmt"
	"net/http"
	"reflect"

	"github.com/mike-webster/spotify-views/data"
	"github.com/mike-webster/spotify-views/keys"
)

func GetDependencies(ctx context.Context) *Dependencies {
	ideps := keys.GetContextValue(ctx, keys.ContextDependencies)
	if ideps == nil {
		fmt.Println("missing deps")
		return nil
	}

	deps, ok := ideps.(*Dependencies)
	if !ok {
		fmt.Println("invalid deps", reflect.TypeOf(ideps))
		return nil
	}

	return deps
}

type Dependencies struct {
	Client HttpClient
	DB     data.DB
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
