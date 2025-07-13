package routers

import (
	"personal-diary/controllers"
	"personal-diary/middleware"

	"github.com/gorilla/mux"
)

func DiaryRouters(routers *mux.Router) {
	dairyRouter := routers.PathPrefix("/diary").Subrouter()
	dairyRouter.Use(middleware.JwtVerify)

	dairyRouter.HandleFunc("", controllers.CreateDiary).Methods("POST")
	dairyRouter.HandleFunc("", controllers.GetAllDiaries).Methods("GET")
	dairyRouter.HandleFunc("/{id}", controllers.UpdateDiary).Methods("PUT")
	dairyRouter.HandleFunc("/{id}", controllers.DeleteDiary).Methods("DELETE")
	dairyRouter.HandleFunc("/{id}", controllers.RefineTextHandler).Methods("POST")
	

}
