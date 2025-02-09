package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"ProjectAWSStore-ReadCustomer/config"
	"ProjectAWSStore-ReadCustomer/models"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var customerCollection *mongo.Collection

// 📌 Función para establecer la colección
func SetCustomerCollection(db *mongo.Database) {
	customerCollection = db.Collection("customers")
}

// 📌 Obtener todos los clientes
func GetAllCustomers(w http.ResponseWriter, r *http.Request) {
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

// 📌 **Sincronizar actualización de clientes desde `UpdateCustomer`**
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

	// ✅ Validar si `ID` está vacío
	if updatedCustomer.ID == primitive.NilObjectID {
		fmt.Println("⚠️ Error: `ID` vacío en la sincronización de actualización.")
		http.Error(w, "⚠️ Error: `ID` vacío en la sincronización", http.StatusBadRequest)
		return
	}

	// 📌 Crear el filtro para buscar por `_id` (NO ES NECESARIO CONVERTIR)
	filter := bson.M{"_id": updatedCustomer.ID}
	update := bson.M{"$set": updatedCustomer}

	// 📌 Intentar actualizar el cliente en la base de datos
	result, err := customerCollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		fmt.Println("❌ Error al actualizar cliente en MongoDB:", err)
		http.Error(w, "❌ Error al sincronizar actualización", http.StatusInternalServerError)
		return
	}

	// 📌 Verificar si realmente se encontró y actualizó el cliente
	if result.MatchedCount == 0 {
		fmt.Println("⚠️ Cliente no encontrado en la base de datos durante sincronización.")
		http.Error(w, "⚠️ Cliente no encontrado en la base de datos durante sincronización.", http.StatusNotFound)
		return
	}

	// ✅ Cliente actualizado correctamente
	fmt.Println("✅ Cliente sincronizado correctamente en ReadCustomer:", updatedCustomer.Email)
	w.WriteHeader(http.StatusOK)
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

func SyncDeleteCustomer(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	// 🔄 Intentar convertir el ID de string a ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println("❌ Error: ID inválido en la sincronización de eliminación:", id)
		http.Error(w, "❌ ID inválido en la sincronización", http.StatusBadRequest)
		return
	}

	customerCollection := config.GetDB().Collection("customers")
	if customerCollection == nil {
		http.Error(w, "❌ Database not initialized", http.StatusInternalServerError)
		return
	}

	// 🔍 Verificar si el cliente existe antes de eliminarlo
	var existingCustomer models.Customer
	err = customerCollection.FindOne(context.TODO(), bson.M{"_id": objID}).Decode(&existingCustomer)
	if err != nil {
		fmt.Println("⚠️ Cliente no encontrado en la base de datos durante sincronización:", id)
		http.Error(w, "⚠️ Cliente no encontrado en la base de datos durante sincronización", http.StatusNotFound)
		return
	}

	// 🗑️ Eliminar el cliente
	_, err = customerCollection.DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		fmt.Println("❌ Error al eliminar cliente en sincronización:", err)
		http.Error(w, "❌ Error al eliminar cliente en sincronización", http.StatusInternalServerError)
		return
	}

	fmt.Println("✅ Cliente eliminado en sincronización:", existingCustomer.Email)
	w.WriteHeader(http.StatusOK)
}
