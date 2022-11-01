package starter

import (
	"context"
	"github.com/basterrus/Go_users_catalog_app/shortener/internal/app/redirectBL"
	"sync"
)

type App struct {
	redirectBL *redirectBL.Redirect
}

func NewApp(redirect *redirectBL.Redirect) *App {
	app := &App{
		redirectBL: redirect,
	}
	return app
}

type APIServer interface {
	Start(redirectBL *redirectBL.Redirect)
	Stop()
}

func (a *App) Serve(ctx context.Context, wg *sync.WaitGroup, hs APIServer) {
	defer wg.Done()
	hs.Start(a.redirectBL)
	<-ctx.Done()
	hs.Stop()
}
