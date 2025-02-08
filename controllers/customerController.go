package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"ProjectAWSStore-ReadCustomer/config"
	"ProjectAWSStore-ReadCustomer/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var customerCollection *mongo.Collection

// 📌 Función para establecer la colección
func SetCustomerCollection(db *mongo.Database) {
	customerCollection = db.Collection("customers")
}

// 📌 Obtener todos los clientes
func GetAllCustomers(w http.ResponseWriter, r *http.Request) {
	if customerCollection == nil {
		http.Error(w, "❌ Error: la base de datos no está inicializada", http.StatusInternalServerError)
		return
	}

	var customers []models.Customer

	cursor, err := customerCollection.Find(context.TODO(), bson.M{})
	if err != nil {
		http.Error(w, "❌ Error obteniendo clientes", http.StatusInternalServerError)
		return
	}
	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var customer models.Customer
		if err := cursor.Decode(&customer); err != nil {
			http.Error(w, "❌ Error decodificando cliente", http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(customers)
}
func SyncCreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer models.Customer
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		http.Error(w, "❌ Entrada inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("📌 Recibida solicitud de sincronización:", customer.Email)

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	// ✅ Verificar si el cliente ya existe en ReadCustomer
	var existingCustomer models.Customer
	err = customerCollection.FindOne(context.TODO(), bson.M{"email": customer.Email}).Decode(&existingCustomer)
	if err == nil {
		fmt.Println("⚠️ Cliente ya existe en ReadCustomer:", customer.Email)
		w.WriteHeader(http.StatusOK)
		return
	}

	// ✅ Insertar nuevo cliente
	_, err = customerCollection.InsertOne(context.TODO(), customer)
	if err != nil {
		http.Error(w, "❌ Error al sincronizar cliente", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente sincronizado correctamente en ReadCustomer:", customer.Email)
	w.WriteHeader(http.StatusCreated)
}

// 📌 **Sincronizar actualización de clientes desde `UpdateCustomer`**
func SyncUpdateCustomer(w http.ResponseWriter, r *http.Request) {
	var updatedCustomer models.Customer
	err := json.NewDecoder(r.Body).Decode(&updatedCustomer)
	if err != nil {
		http.Error(w, "❌ Entrada inválida", http.StatusBadRequest)
		return
	}

	fmt.Println("📌 Recibida solicitud de sincronización para:", updatedCustomer.Email)

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "Database not initialized", http.StatusInternalServerError)
		return
	}

	// ✅ Actualizar cliente en MongoDB
	_, err = customerCollection.UpdateOne(
		context.TODO(),
		bson.M{"email": updatedCustomer.Email},
		bson.M{"$set": updatedCustomer},
	)
	if err != nil {
		http.Error(w, "❌ Error al sincronizar actualización", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente sincronizado correctamente en ReadCustomer/CreateCustomer:", updatedCustomer.Email)
	w.WriteHeader(http.StatusOK)
}
