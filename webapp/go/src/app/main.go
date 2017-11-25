package main

import (
	"crypto/md5"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/go-redis/redis"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/jmoiron/sqlx"
)

var (
	db      *sqlx.DB
	rediCli *redis.Client
)

var localServers = []string{
	"192.168.12.1",
	"192.168.12.2",
	"192.168.12.3",
	"192.168.12.4",
}

func init() {
	debug.SetGCPercent(-1)
}

func initDB() {
	db_host := os.Getenv("ISU_DB_HOST")
	db_host = "127.0.0.1"
	db_port := os.Getenv("ISU_DB_PORT")
	if db_port == "" {
		db_port = "3306"
	}
	db_user := os.Getenv("ISU_DB_USER")
	if db_user == "" {
		db_user = "root"
	}
	db_password := os.Getenv("ISU_DB_PASSWORD")
	if db_password != "" {
		db_password = ":" + db_password
	}

	dsn := fmt.Sprintf("%s%s@tcp(%s:%s)/isudb?parseTime=true&loc=Local&charset=utf8mb4",
		db_user, db_password, db_host, db_port)

	log.Printf("Connecting to db: %q", dsn)
	db, _ = sqlx.Connect("mysql", dsn)
	for {
		err := db.Ping()
		if err == nil {
			break
		}
		log.Println(err)
		time.Sleep(time.Second * 3)
	}

	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(5 * time.Minute)
	log.Printf("Succeeded to connect db.")

	rediCli = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
}

func getInitializeHandler(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("broadcast") == "" {
		var wg sync.WaitGroup
		for _, s := range localServers {
			wg.Add(1)
			go func(s string) {
				defer wg.Done()

				_, err := http.Get("http://" + s + "/initialize?broadcast=1")
				if err != nil {
					panic(err)
				}
			}(s)
		}
		wg.Wait()
		w.WriteHeader(204)
	} else {
		db.MustExec("TRUNCATE TABLE adding")
		db.MustExec("TRUNCATE TABLE buying")
		db.MustExec("TRUNCATE TABLE room_time")
		runtime.GC()
		w.WriteHeader(204)
	}
}

var servers = [...]string{
	"app0121.isu7f.k0y.org",
	//	"app0122.isu7f.k0y.org",
	//	"app0123.isu7f.k0y.org",
	//	"app0124.isu7f.k0y.org",
}

func getRoomServer(room string) string {
	hashed := md5.Sum([]byte(room))
	var s []byte = hashed[:4]
	l := len(servers)
	idx := int(binary.BigEndian.Uint32(s)) % l
	return servers[idx]
}

func getRoomHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomName := vars["room_name"]
	path := "/ws/" + url.PathEscape(roomName)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Host string `json:"host"`
		Path string `json:"path"`
	}{
		Host: getRoomServer(roomName),
		Path: path,
	})
}

func wsGameHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	roomName := vars["room_name"]

	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)
	if _, ok := err.(websocket.HandshakeError); ok {
		log.Println("Failed to upgrade", err)
		return
	}
	go serveGameConn(ws, roomName)
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	initDB()

	r := mux.NewRouter()
	r.HandleFunc("/initialize", getInitializeHandler)
	r.HandleFunc("/room/", getRoomHandler)
	r.HandleFunc("/room/{room_name}", getRoomHandler)
	r.HandleFunc("/ws/", wsGameHandler)
	r.HandleFunc("/ws/{room_name}", wsGameHandler)
	r.HandleFunc("/debug/pprof", pprof.Index)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("../public/")))

	log.Fatal(http.ListenAndServe(":5000", handlers.LoggingHandler(os.Stderr, r)))
}
