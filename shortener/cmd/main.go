package main

import (
	"context"
	"github.com/basterrus/Go_users_catalog_app/shortener/internal/api/chiRouter"
	"github.com/basterrus/Go_users_catalog_app/shortener/internal/api/handler"
	"github.com/basterrus/Go_users_catalog_app/shortener/internal/api/server"
	"github.com/basterrus/Go_users_catalog_app/shortener/internal/app/redirectBL"
	"github.com/basterrus/Go_users_catalog_app/shortener/internal/app/repository/followingBL"
	"github.com/basterrus/Go_users_catalog_app/shortener/internal/app/repository/shortenerBL"
	"github.com/basterrus/Go_users_catalog_app/shortener/internal/app/starter"
	"github.com/basterrus/Go_users_catalog_app/shortener/internal/db/inmemoryDB"
	"github.com/basterrus/Go_users_catalog_app/shortener/internal/db/sql/postgresDB"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	if tz := os.Getenv("TZ"); tz != "" {
		var err error
		time.Local, err = time.LoadLocation(tz)
		if err != nil {
			log.Printf("error loading location '%s': %v\n", tz, err)
		}
	}

	// output current time zone
	tnow := time.Now()
	tz, _ := tnow.Zone()
	log.Printf("Local time zone %s. Service started at %s", tz,
		tnow.Format("2006-01-02T15:04:05.000 MST"))

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)

	srvPort := os.Getenv("SRV_PORT")
	if srvPort == "" {
		log.Fatal("unknown SRV_PORT = ", srvPort)
	}

	var shortdb shortenerBL.ShortenerStore
	var followdb followingBL.FollowingStore
	store := os.Getenv("SHORTENER_STORE")

	switch store {
	case "mem":
		shortdb = inmemoryDB.NewShortenerMapDB()
		followdb = inmemoryDB.NewFollowingMapDB()
	case "pg":
		pgDB := postgresDB.NewPostgresDB()
		defer pgDB.Close()
		shortdb = pgDB
		followdb = pgDB
	default:
		log.Fatal("unknown SHORTENER_STORE = ", store)
	}

	shortBL := shortenerBL.NewShotenerBL(shortdb)
	followBL := followingBL.NewFollowingBL(followdb)
	redirBL := redirectBL.NewRedirect(shortBL, followBL)

	app := starter.NewApp(redirBL)
	handlers := handler.NewHandlers(redirBL)
	chi := chiRouter.NewChiRouter(handlers)
	srv := server.NewServer(":"+srvPort, chi)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go app.Serve(ctx, wg, srv)

	<-ctx.Done()
	cancel()
	wg.Wait()
}
