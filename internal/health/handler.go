package health

import (
	"github.com/nayefradwi/nayef_go_common/httputil"
	"net/http"
)

func HealthCheck(w http.ResponseWriter, r *http.Request) {
	jw := httputil.NewJsonResponseWriter(w)
	jw.WriteSuccessMessage("Healthy", nil)
}
