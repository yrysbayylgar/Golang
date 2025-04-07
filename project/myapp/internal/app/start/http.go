package start

import (
	"fmt"
	"net/http"
)

func StartHTTPServer(port string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	fmt.Println("HTTP сервер запущен на порту", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}
