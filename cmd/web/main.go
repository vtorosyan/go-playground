package main

import (
	"database/sql"
	"flag"
	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"
	"vtorosyan.learning/internal/models"

	_ "github.com/go-sql-driver/mysql"
)

type application struct {
	logger         *slog.Logger
	snippets       *models.SnippetModel
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
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

	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	templateCache, err := newTemplateCache()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	formDecoder := form.NewDecoder()
	snippets := models.SnippetModel{DB: db}
	app := &application{
		logger:         logger,
		snippets:       &snippets,
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
	}

	// Note to self: even if you don't use the address of the struct,
	//Go will automatically get the address when calling ListenAndServe, as the method is a pointer receiver
	server := &http.Server{
		Addr:     *addr,
		Handler:  app.routes(),
		ErrorLog: slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("Starting the server.", "address", server.Addr)
	err = server.ListenAndServe()

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
