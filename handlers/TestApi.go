package handlers

import (
	"net/http"
)

func TestApi(w http.ResponseWriter, r *http.Request) {
	//params := mux.Vars(r)
	//
	//podName := params["podName"]
	//
	//// Get business layer
	//business, err := getBusiness(r)
	//if err != nil {
	//	RespondWithError(w, http.StatusInternalServerError, "Workloads initialization error: "+err.Error())
	//	return
	//}
	//
	//version := business.Istiod.GetIstiodVersion(podName)

	RespondWithJSON(w, http.StatusOK, "test")
}
