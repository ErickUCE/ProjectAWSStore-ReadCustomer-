package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Variable global para la conexión a la base de datos
var DB *mongo.Database

// ConnectDB establece la conexión con MongoDB
func ConnectDB() (*mongo.Database, error) {
	fmt.Println("📌 Ejecutando ConnectDB...")

	// 📌 Intentar cargar .env primero
	err := godotenv.Load()
	if err != nil {
		fmt.Println("⚠️ No se pudo cargar el archivo .env, verificando variables de entorno...")
	}

	// 📌 Obtener MONGO_URI desde las variables de entorno
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("❌ Error: No se encontró MONGO_URI en las variables de entorno")
	}

	// 📌 Configurar opciones de cliente
	clientOptions := options.Client().ApplyURI(mongoURI)

	// 📌 Crear un cliente de MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal("❌ Error al crear el cliente de MongoDB:", err)
		return nil, err
	}

	// 📌 Verificar la conexión con MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("❌ Error: No se pudo conectar a MongoDB:", err)
		return nil, err
	}

	// 📌 Obtener el nombre de la base de datos desde las variables de entorno
	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "ReadCustomerDB" // Nombre por defecto
	}

	fmt.Println("🛢️ Base de datos seleccionada:", dbName)

	// 📌 Guardar la base de datos en la variable global
	DB = client.Database(dbName)
	return DB, nil
}

// GetDB devuelve la instancia de la base de datos ya conectada
func GetDB() *mongo.Database {
	if DB == nil {
		log.Fatal("❌ Error: La base de datos no está inicializada. Asegúrate de llamar a ConnectDB() primero.")
	}
	return DB
}
