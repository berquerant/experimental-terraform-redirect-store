package main

import (
	"errors"
	"experimental-terraform-redirect-store/api"
	"flag"
	"log/slog"
	"os"
)

func main() {
	var (
		addr = flag.String("addr", "127.0.0.1:8030", "")
		db   = flag.String("db", "api.db", "DB file")
	)
	flag.Parse()

	if err := touchFile(*db); err != nil {
		panic(err)
	}
	dbFile := api.NewDatabaseFile(*db)
	database := api.NewDatabaseImpl(dbFile)
	server := api.NewServerImpl(database)
	slog.Info("listen", slog.String("addr", *addr), slog.String("db", *db))
	panic(api.ListenAndServe(*addr, server, server))
}

func touchFile(name string) error {
	_, err := os.Stat(name)
	switch {
	case errors.Is(err, os.ErrNotExist):
		f, err := os.Create(name)
		if err != nil {
			return err
		}
		return f.Close()
	case err != nil:
		return err
	default:
		return nil
	}
}
