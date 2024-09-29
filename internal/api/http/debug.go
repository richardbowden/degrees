package apihttp

import (
	"encoding/json"
	"net/http"
)

type DebugHandler struct{}

func (e *DebugHandler) Debug(w http.ResponseWriter, r *http.Request) {
	f := map[string]int32{}
	f["rando_string"] = 402

	// f["total_cons"] = e.pg.Stat().TotalConns()
	// f["idel_cons"] = e.pg.Stat().IdleConns()
	// f["max_idle_destroy_count"] = int32(e.pg.Stat().MaxIdleDestroyCount())
	//
	jd, err := json.Marshal(f)

	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jd)
}
