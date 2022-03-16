package api

import (
	"github.com/brunomdev/digital-account/app/api/handlers"
	"github.com/brunomdev/digital-account/app/api/routes"
)

func (s *Server) router() {
	routes.DocRoutes(s.httpServer)
	routes.AccountRoutes(s.httpServer, handlers.NewAccountHandler(s.service.Account))
	routes.TransactionRoutes(s.httpServer, handlers.NewTransactionHandler(s.service.Transaction))
}
