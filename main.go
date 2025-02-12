package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"ProjectAWSStore-ReadCustomer/config"
	"ProjectAWSStore-ReadCustomer/routes"

	"github.com/gorilla/handlers"
	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("🚀 Iniciando CustomerService en Golang...")

	// 📌 Cargar variables de entorno
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ Advertencia: No se pudo cargar el archivo .env, verificando variables de entorno...")
	}

	fmt.Println("🛠️ MONGO_URI desde Go:", os.Getenv("MONGO_URI"))
	fmt.Println("🛠️ MONGO_DB_NAME desde Go:", os.Getenv("MONGO_DB_NAME"))

	// ✅ Conectar a MongoDB antes de configurar las rutas
	db, err := config.ConnectDB()
	if err != nil {
		log.Fatal("❌ Error al conectar con MongoDB:", err)
	}
	fmt.Println("✅ Conexión exitosa a MongoDB")

	// ✅ Configurar rutas después de conectar a MongoDB
	router := routes.SetupRoutes(db)

	// ✅ Aplicar middleware de CORS al router
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"http://localhost:3000"}), // ✅ Permitir solicitudes desde el frontend
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"Content-Type"}),
	)(router)

	// 📌 Obtener el puerto desde `.env`
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082" // 🔥 Valor por defecto corregido a 8082
	}

	// ✅ Iniciar el servidor con CORS habilitado
	fmt.Println("✅ Servidor corriendo en el puerto", port)
	log.Fatal(http.ListenAndServe(":"+port, corsHandler))
}
