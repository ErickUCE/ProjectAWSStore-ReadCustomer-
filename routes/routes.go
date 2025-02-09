package routes

import (
	"ProjectAWSStore-ReadCustomer/controllers"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// SetupRoutes configura las rutas de la API
func SetupRoutes(db *mongo.Database) *mux.Router {
	router := mux.NewRouter()

	// ✅ Inicializar la colección en el controlador
	controllers.SetCustomerCollection(db)

	// ✅ Configurar controladores
	router.HandleFunc("/customers", controllers.GetAllCustomers).Methods("GET")
	router.HandleFunc("/sync-customer", controllers.SyncCreateCustomer).Methods("POST")
	// Endpoint de sincronización desde `CreateCustomer`
	router.HandleFunc("/sync-create", controllers.SyncCreateCustomer).Methods("POST")
	router.HandleFunc("/sync-update", controllers.SyncUpdateCustomer).Methods("POST")
	//router.HandleFunc("/sync-delete", controllers.SyncDeleteCustomer).Methods("POST")

	return router
}
