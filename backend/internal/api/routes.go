package api

import (
	"database/sql"
	"net/http"
)

func NewRouter(db *sql.DB) http.Handler {
    mux := http.NewServeMux()

    // Define routes
    mux.HandleFunc("/expenses", testHandler)
    mux.HandleFunc("/incomes", testHandler)
    mux.HandleFunc("/categories", testHandler)

    return mux
}

// Test handler for now
func testHandler(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Hello, World!"))
}
