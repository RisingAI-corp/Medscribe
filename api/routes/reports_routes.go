package routes

import (
	"Medscribe/api/handlers/reportsHandler"

	"github.com/go-chi/chi/v5"
)

func ReportRoutes(handler reportsHandler.ReportsHandler) *chi.Mux {
	r := chi.NewRouter()

	r.Post("/generate", handler.GenerateReport)

	r.Patch("/regenerate", handler.RegenerateReport)

	r.Patch("/changeName", handler.ChangeReportName)

	r.Patch("/updateContentSection", handler.UpdateContentSection)

	r.Patch("/learn-style", handler.LearnStyle)

	r.Delete("/delete", handler.DeleteReport)

	r.Post("/get", handler.GetReport)

	r.Post("/getTranscript", handler.GetTranscript)

	return r
}
