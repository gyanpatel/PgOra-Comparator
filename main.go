package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	_ "github.com/godror/godror"
	_ "github.com/lib/pq"
)

//global variables hold db connections and application secrets
var dbConPG *sql.DB
var dbConOra *sql.DB
var secretDetails SecretDetails
var tableQueryList []TableQueryList

// Initialize db connections
func init() {

	secretDetApp, err := readK8LocalSecret("appconfig")
	secretDetails = secretDetApp
	if err != nil {
		log.Fatal("FATAL : UNABLE TO RETRIEVE THE K8LocalSecret appconfig SECRET, APPLICATION STARTUP FAILURE ", err)
	}
	secretDetOra, err := readK8LocalSecret("orasecret")

	if err != nil {
		log.Fatal("FATAL : UNABLE TO RETRIEVE THE K8LocalSecret SECRET, APPLICATION STARTUP FAILURE ", err)
	}
	oraConnectString := secretDetOra.User + "/" + secretDetOra.Password + "@" + secretDetOra.DbSslMode + "://" + secretDetOra.Host + ":" + secretDetOra.Port + "/" + secretDetOra.Dbname
	dbConSetupOra, err := sql.Open("godror", oraConnectString)
	if err != nil {
		log.Fatal("FATAL : UNABLE TO connect to Oracle  brdb,  FAILURE ", err)
	}

	dbConOra = dbConSetupOra
	dbConOra.Ping()
	err = dbConOra.Ping()
	if err != nil {
		log.Fatal("FATAL : UNABLE TO ping to Oracle BRDB db,  FAILURE ", err)
	}
	dbConOra.SetMaxIdleConns(5)
	secretDetPG, err := readK8LocalSecret("pgsecret")
	if err != nil {
		log.Fatal("FATAL : UNABLE TO RETRIEVE THE K8LocalSecret SECRET, APPLICATION STARTUP FAILURE ", err)
	}
	bdbConnectString := "dbname=" + secretDetPG.Dbname + " port=" + secretDetPG.Port + " user=" + secretDetPG.User + " password=" + secretDetPG.Password + " host=" + secretDetPG.Host + " sslmode=" + secretDetPG.DbSslMode
	dbConSetupPG, err := sql.Open("postgres", bdbConnectString)
	if err != nil {
		log.Fatal("FATAL : UNABLE TO connect to PG BRDB db,  FAILURE ", err)
	}
	dbConPG = dbConSetupPG
	err = dbConPG.Ping()
	if err != nil {
		log.Fatal("FATAL : UNABLE TO ping to PG BRDB db,  FAILURE ", err)
	}
	dbConPG.SetMaxIdleConns(5)

	err = getTableQuery()
	if err != nil {
		log.Fatal("FATAL : calling getTableQuery() UNABLE TO GET query list,  FAILURE ", err)
	}
	/*
		queryFile, err := ioutil.ReadFile("ComparisonQuery.json")
		if err != nil {
			log.Println("ERROR : UNABLE TO READ ComparisonQuery.json => ", err)
		}
		err = json.Unmarshal(queryFile, &tableQueryList)
		if err != nil {
			log.Println("ERROR : UNABLE TO PARSE THE ComparisonQuery.json ", err)
		} */
}
func readK8LocalSecret(envVarName string) (SecretDetails, error) {
	var secretDet SecretDetails
	secretFilePath := os.Getenv(envVarName)
	secretFilePath = filepath.Clean(secretFilePath)
	if len(secretFilePath) == 0 {
		log.Fatal("FATAL : UNABLE TO LOCATE THE K8LocalSecret SECRET FILE CHECK ENV VARIABLE " + envVarName + " FOR FILE PATH, APPLICATION STARTUP FAILURE ")
	}
	secretsFile, err := ioutil.ReadFile(filepath.Clean(secretFilePath))
	if err != nil {
		log.Fatal("FATAL : UNABLE TO READ THE rK8LocalSecret SECRET FILE => "+secretFilePath+" , APPLICATION STARTUP FAILURE ", err)
	}
	err = json.Unmarshal(secretsFile, &secretDet)
	if err != nil {
		log.Fatal("FATAL : UNABLE TO PARSE THE K8LocalSecret SECRET, FILE => "+secretFilePath+", APPLICATION STARTUP FAILURE ", err)
	}
	return secretDet, err
}

func main() {
	log.Println("Info :", "Starting Ora PG Comparator...")
	log.Println("Current Query list ", tableQueryList)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("templates/assets"))))

	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/", handleLogin)
	http.HandleFunc("/home", handleHome)
	http.HandleFunc("/hist", handleHist)
	http.HandleFunc("/comp", handleCompare)
	http.HandleFunc("/compall", handleCompare)
	http.HandleFunc("/loginval", handleLoginVal)
	http.HandleFunc("/logout", handleLogout)
	http.HandleFunc("/orapgcompawsssologin", handleAwsSsoLogin)

	log.Fatal(http.ListenAndServe(":"+secretDetails.WebAppPort, nil))
}
