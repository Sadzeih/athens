package download

import (
	"net/http"

	"github.com/gomods/athens/pkg/download/mode"
	"github.com/gomods/athens/pkg/errors"
	"github.com/gomods/athens/pkg/log"
)

// PathVersionModule URL.
const PathVersionModule = "/{module:.+}/@v/{version}.mod"

// ModuleHandler implements GET baseURL/module/@v/version.mod
func ModuleHandler(dp Protocol, lggr log.Entry, df *mode.DownloadFile) http.Handler {
	const op errors.Op = "download.VersionModuleHandler"
	f := func(w http.ResponseWriter, r *http.Request) {
		mod, ver, err := getModuleParams(r, op)
		if err != nil {
			err = errors.E(op, errors.M(mod), errors.V(ver), err)
			lggr.SystemErr(err)
			w.WriteHeader(errors.Kind(err))
			return
		}
		modBts, err := dp.GoMod(r.Context(), mod, ver)
		if err != nil {
			lggr.SystemErr(err)
			if errors.Kind(err) == errors.KindRedirect {
				http.Redirect(w, r, getRedirectURL(df.URL(mod), r.URL.Path), errors.KindRedirect)
				return
			}
			w.WriteHeader(errors.Kind(err))
			return
		}

		w.Write(modBts)
	}
	return http.HandlerFunc(f)
}
