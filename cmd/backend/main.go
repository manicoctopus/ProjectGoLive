package main

import (
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"ProjectGoLive/pkg/models"
	"ProjectGoLive/pkg/models/mysql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type application struct {
	errorLog *log.Logger
	infoLog  *log.Logger
	users    interface {
		Create(*models.User) (int, error)
		Update(*models.User) error
		Delete(uint32) error
		Retrieve(uint32) (*models.User, error)
		RetrieveAll() ([]*models.User, error)
		AuthenticateUser(string, string) (int, error)
	}
	pdtsvcs interface {
		Create(*models.Pdtsvc) (int, error)
		Update(*models.Pdtsvc) error
		Delete(uint32) error
		Retrieve(uint32) (*models.Pdtsvc, error)
		RetrieveAll() ([]*models.Pdtsvc, error)
		RetrieveAllByID(uint32) ([]*models.Pdtsvc, error)
	}
	listings interface {
		Create(*models.Listing) (int, error)
		Update(*models.Listing) error
		Delete(uint32) error
		Retrieve(uint32) (*models.Listing, error)
		RetrieveAll() ([]*models.Listing, error)
	}
	reviews interface {
		Create(*models.Review) (int, error)
		Update(*models.Review) error
		Delete(uint32) error
		Retrieve(uint32) (*models.Review, error)
		RetrieveAll() ([]*models.Review, error)
		RetrieveAllByID(uint32) ([]*models.Review, error)
	}
	categories interface {
		Create(*models.Category) (int, error)
		Update(*models.Category) error
		Delete(uint32) error
		Retrieve(uint32) (*models.Category, error)
		RetrieveAll() ([]*models.Category, error)
	}
}

var (
	app *application
)

func loadEnv(envFilename string) (*string, *string, *string, *string, *string, tls.Certificate, []byte) {
	if err := godotenv.Load(envFilename); err != nil {
		app.errorLog.Fatal("ERROR loading .env file")
	}

	conn := os.Getenv("CONN_HOST")
	addr := &conn

	portString := os.Getenv("CONN_PORT")
	port := &portString

	sslString := os.Getenv("SSL")
	ssl := &sslString

	sslPortString := os.Getenv("SSL_CONN_PORT")
	sslPort := &sslPortString

	dsnString := os.Getenv("DSN")
	dsn := &dsnString

	// Initialize Certificates
	cert, err := tls.LoadX509KeyPair(os.Getenv("CERT_FILE"), os.Getenv("KEY_FILE"))
	if err != nil {
		app.errorLog.Fatalf("ERROR server certificate: %s", err)
	}

	caCert, err := ioutil.ReadFile(os.Getenv("CACERT_FILE"))
	if err != nil {
		app.errorLog.Fatalf("ERROR ca certificate: %s", err)
	}

	return addr, port, ssl, sslPort, dsn, cert, caCert
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)

	app = &application{
		errorLog: errorLog,
		infoLog:  infoLog,
	}

	addr, port, ssl, sslPort, dsn, cert, caCert := loadEnv("apiserver.env")

	db, err := openDB(*dsn)
	if err != nil {
		app.errorLog.Fatal(err)
	}
	defer db.Close()

	app.users = &mysql.UserModel{DB: db}
	app.pdtsvcs = &mysql.PdtsvcModel{DB: db}
	app.listings = &mysql.ListingModel{DB: db}
	app.reviews = &mysql.ReviewModel{DB: db}
	app.categories = &mysql.CategoryModel{DB: db}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	// Initialize a tls.Config struct to hold the non-default TLS settings we want
	// the server to use.
	tlsConfig := &tls.Config{
		ClientCAs:    caCertPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		Certificates: []tls.Certificate{cert},
	}
	tlsConfig.BuildNameToCertificate()

	// Initialize a new http.Server struct.
	srv := &http.Server{
		Addr:         *addr + ":" + *sslPort,
		ErrorLog:     app.errorLog,
		Handler:      app.routes(),
		TLSConfig:    tlsConfig,
		TLSNextProto: map[string]func(s *http.Server, c *tls.Conn, h http.Handler){},
		// Add Idle, Read and Write timeouts to the server.
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if *ssl == "true" {
		app.infoLog.Printf("INFO starting https api server @ %s\n", srv.Addr)
		if err := srv.ListenAndServeTLS("", ""); err != nil {
			app.errorLog.Fatal("ERROR starting https api server :: ", err)
			return
		}

		httpSrv := &http.Server{
			Addr: *addr + ":" + *port,
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				host, _, _ := net.SplitHostPort(r.Host)
				u := r.URL
				u.Host = net.JoinHostPort(host, *sslPort)
				u.Scheme = "https"
				log.Println(u.String())
				http.Redirect(w, r, u.String(), http.StatusMovedPermanently)
			}),
			IdleTimeout:  time.Minute,
			ReadTimeout:  5 * time.Second,
			WriteTimeout: 10 * time.Second,
		}
		// Start the HTTP server and redirect all incoming connections to HTTPS
		go httpSrv.ListenAndServe()

	} else if *ssl == "false" {
		srv.Addr = *addr + ":" + *port
		app.infoLog.Printf("INFO starting http api server @ %s\n", srv.Addr)
		if err := srv.ListenAndServe(); err != nil {
			app.errorLog.Fatal("ERROR starting http api server :: ", err)
			return
		}
	}
}
