package handlers

import (
	"net/http"

	"github.com/AsCat/acorn/log"
)

func NamespaceList(w http.ResponseWriter, r *http.Request) {
	business, err := getBusiness(r)

	if err != nil {
		log.Error(err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	namespaces, err := business.Namespace.GetNamespaces()
	if err != nil {
		log.Error(err)
		RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	RespondWithJSON(w, http.StatusOK, namespaces)
}
