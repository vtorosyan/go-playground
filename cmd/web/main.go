package main

import (
	"database/sql"
	"flag"
	"github.com/go-playground/form/v4"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"vtorosyan.learning/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger        *slog.Logger
	snippets      *models.SnippetModel
	templateCache map[string]*template.Template
	formDecoder   *form.Decoder
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP port that the server needs to run")
	dsn := flag.String("dsn", "user:password@/snippetbox?parseTime=true", "Database connection string")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelDebug,
		AddSource: true,
	}))

	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer func() {
		err = db.Close()
		if err != nil {
			logger.Error(err.Error())
		}
	}()

	templateCache, err := newTemplateCache()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()
	snippets := models.SnippetModel{DB: db}
	app := &application{
		logger:        logger,
		snippets:      &snippets,
		templateCache: templateCache,
		formDecoder:   formDecoder,
	}

	logger.Info("Starting the server.", "address", *addr)
	err = http.ListenAndServe(*addr, app.routes())
	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}
