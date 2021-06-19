package main

import (
	"crypto/tls"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	resty "gopkg.in/resty.v1"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	//session       *sessions.Session
	templateCache map[string]*template.Template
	fileServer    http.Handler
	webSvcHost    string
}

var (
	app *application
)

func loadEnv(envFilename string) (*string, tls.Certificate) {
	if err := godotenv.Load(envFilename); err != nil {
		app.errorLog.Fatal("ERROR loading .env file")
	}

	conn := os.Getenv("CONN_HOST") + ":" + os.Getenv("CONN_PORT")
	addr := &conn

	// Initialize Certificates
	cert, err := tls.LoadX509KeyPair(os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"))
	if err != nil {
		app.errorLog.Fatalf("ERROR server certificate: %s", err)
	}

	// Initialize a new template cache...
	tCache, err := newTemplateCache(os.Getenv("TEMPLATE"))
	if err != nil {
		log.Fatalf("ERROR loading templates: %s", err)
	}

	app.templateCache = tCache
	app.fileServer = http.FileServer(http.Dir(os.Getenv("STATIC")))
	app.webSvcHost = os.Getenv("WEB_SERVICE_HOST")

	//session := sessions.New([]byte(os.Getenv("SECRET")))
	//session.Lifetime = 12 * time.Hour
	//session.Secure = true

	return addr, cert
}

func main() {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	app = &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	addr, cert := loadEnv("webclient.env")

	tlsConfig := &tls.Config{
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true, // Disable security check (https)
	}
	resty.SetTLSClientConfig(tlsConfig)

	// Initialize a new http.Server struct.
	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		TLSNextProto: map[string]func(s *http.Server, c *tls.Conn, h http.Handler){},
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("INFO https server starting @ %s", srv.Addr)
	if err := srv.ListenAndServeTLS("", ""); err != nil {
		log.Fatal("error starting https server :: ", err)
		return
	}
}
