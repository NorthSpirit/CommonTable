package hub

import (
	"net/http"
)

// ServeAPK serves the waiter app APK with an explicit, correct
// Content-Type — without this, some browsers misidentify the file as a
// plain ZIP archive (which, structurally, an APK actually is) and offer
// to extract it instead of installing it.
func ServeAPK(apkPath string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.android.package-archive")
		w.Header().Set("Content-Disposition", `attachment; filename="waiter_app.apk"`)
		http.ServeFile(w, r, apkPath)
	}
}
