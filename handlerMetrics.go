package main

import (
	"fmt"
	"net/http"
)

func handlerMetrics(w http.ResponseWriter, r *http.Request, cfg *apiConfig) {
	displayString := `<html>
<body>
<h1>Welcome, Chirpy Admin</h1>
<p>Chirpy has been visited %d times!</p>
</body>
</html>`

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(displayString, cfg.fileserverHits.Load())))
}
