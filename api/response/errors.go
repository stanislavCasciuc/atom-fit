package response

import (
	"net/http"
)

func (resp *Responser) InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	resp.logger.Errorw(
		"internal server error",
		"method",
		r.Method,
		"path",
		r.URL.Path,
		"error",
		err,
	)

	WriteJSONError(w, http.StatusInternalServerError, "internal error")
}

func (resp *Responser) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	resp.logger.Warnw("bad request error", "method", r.Method, "path", r.URL.Path, "error", err)

	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (resp *Responser) notFoundErorr(w http.ResponseWriter, r *http.Request, err error) {
	resp.logger.Warnw("not found error", "method", r.Method, "path", r.URL.Path, "error", err)

	WriteJSONError(w, http.StatusNotFound, err.Error())
}

func (resp *Responser) conflictError(w http.ResponseWriter, r *http.Request, err error) {
	resp.logger.Warnw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err)

	WriteJSONError(w, http.StatusConflict, err.Error())
}
