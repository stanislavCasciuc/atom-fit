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

func (resp *Responser) BadRequestError(w http.ResponseWriter, r *http.Request, err error) {
	resp.logger.Warnw("bad request error", "method", r.Method, "path", r.URL.Path, "error", err)

	WriteJSONError(w, http.StatusBadRequest, err.Error())
}

func (resp *Responser) NotFoundErorr(w http.ResponseWriter, r *http.Request, err error) {
	resp.logger.Warnw("not found error", "method", r.Method, "path", r.URL.Path, "error", err)

	WriteJSONError(w, http.StatusNotFound, err.Error())
}

func (resp *Responser) ConflictError(w http.ResponseWriter, r *http.Request, err error) {
	resp.logger.Warnw("conflict error", "method", r.Method, "path", r.URL.Path, "error", err)

	WriteJSONError(w, http.StatusConflict, err.Error())
}

func (resp *Responser) UnauthorizedError(w http.ResponseWriter, r *http.Request, err error) {
	resp.logger.Warnw("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err)

	WriteJSONError(w, http.StatusUnauthorized, err.Error())
}
