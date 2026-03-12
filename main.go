package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)




// puerto con el número de carnet
const puerto = ":24750"
const archivoJSON = "./data/canciones.json"




// Estructura de una cancion
type Cancion struct {
	ID       int    `json:"id"`
	Nombre   string `json:"cancion"`
	Artista  string `json:"artista"`
	Album    string `json:"album"`
	Genero   string `json:"genero"`
	Duracion string `json:"duracion"`
	Vistas   int    `json:"vistas"`
}

// respuesta de error
type RespError struct {
	Status  int    `json:"status"`
	Error   string `json:"error"`
	Detalle string `json:"detalle"`
}

var canciones []Cancion

func main() {
	cargarDatos()

	log.Println("Servidor corriendo en", puerto)
	log.Fatal(http.ListenAndServe(puerto, nil))
}




//DATOS  
func cargarDatos() {
	archivo, err := os.ReadFile(archivoJSON)
	if err != nil {
		log.Fatal("Error leyendo archivo:", err)
	}

	err = json.Unmarshal(archivo, &canciones)
	if err != nil {
		log.Fatal("Error parseando JSON:", err)
	}
}


func guardarDatos() {
	datos, err := json.MarshalIndent(canciones, "", "  ")
	if err != nil {
		log.Println("Error al convertir a JSON:", err)
		return
	}

	err = os.WriteFile(archivoJSON, datos, 0644)
	if err != nil {
		log.Println("Error escribiendo archivo:", err)
	}
}

func siguienteID() int {
	max := 0
	for _, c := range canciones {
		if c.ID > max {
			max = c.ID
		}
	}
	return max + 1
}





//RESPUESTAS JSON  

func respJSON(w http.ResponseWriter, status int, datos interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(datos)
}

func respError(w http.ResponseWriter, status int, error string, detalle string) {
	resp := RespError{
		Status:  status,
		Error:   error,
		Detalle: detalle,
	}
	respJSON(w, status, resp)
}




// VALIDACION  

func validarCancion(c Cancion) (bool, string) {
	if c.Nombre == "" {
		return false, "El campo 'cancion' es requerido"
	}
	if c.Artista == "" {
		return false, "El campo 'artista' es requerido"
	}
	if c.Genero == "" {
		return false, "El campo 'genero' es requerido"
	}
	if c.Duracion == "" {
		return false, "El campo 'duracion' es requerido"
	}
	if c.Vistas < 0 {
		return false, "El campo 'vistas' no puede ser negativo"
	}
	return true, ""
}
