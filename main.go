package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	http.HandleFunc("/api/canciones/", handlerPorID)
	http.HandleFunc("/api/canciones", handlerCanciones)

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





// HANDLER PRINCIPAL  

func handlerCanciones(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		handleGet(w, r)
	case http.MethodPost:
		handlePost(w, r)
	default:
		respError(w, 405, "Metodo no permitido", "Usa GET o POST en /api/canciones")
	}
}




// HANDLER POR ID

func handlerPorID(w http.ResponseWriter, r *http.Request) {
	partes := strings.Split(r.URL.Path, "/")

	if len(partes) < 4 || partes[3] == "" {
		respError(w, 400, "ID faltante", "Incluye un ID en la URL: /api/canciones/1")
		return
	}

	id, err := strconv.Atoi(partes[3])
	if err != nil {
		respError(w, 400, "ID invalido", "El ID debe ser un numero entero")
		return
	}

	switch r.Method {
	case http.MethodGet:
		handleGetPorID(w, id)
	case http.MethodPut:
		handlePut(w, r, id)
	case http.MethodPatch:
		handlePatch(w, r, id)
	default:
		respError(w, 405, "Metodo no permitido", "Usa GET, PUT o PATCH en /api/canciones/{id}")
	}
}


//   GET  

func handleGet(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	// busqueda por query param id
	idParam := query.Get("id")
	if idParam != "" {
		id, err := strconv.Atoi(idParam)
		if err != nil {
			respError(w, 400, "ID invalido", "El query param 'id' debe ser un numero")
			return
		}
		for _, c := range canciones {
			if c.ID == id {
				respJSON(w, 200, c)
				return
			}
		}
		respError(w, 404, "No encontrado", "No existe una cancion con id "+idParam)
		return
	}

	respJSON(w, 200, canciones)
}

func handleGetPorID(w http.ResponseWriter, id int) {
	for _, c := range canciones {
		if c.ID == id {
			respJSON(w, 200, c)
			return
		}
	}
	respError(w, 404, "No encontrado", "No existe una cancion con id "+strconv.Itoa(id))
}


//   POST  

func handlePost(w http.ResponseWriter, r *http.Request) {
	var nueva Cancion

	err := json.NewDecoder(r.Body).Decode(&nueva)
	if err != nil {
		respError(w, 400, "JSON invalido", "El body no es JSON valido")
		return
	}

	if ok, msg := validarCancion(nueva); !ok {
		respError(w, 400, "Validacion fallida", msg)
		return
	}

	nueva.ID = siguienteID()
	canciones = append(canciones, nueva)
	guardarDatos()

	respJSON(w, 201, nueva)
}


//   PUT  

func handlePut(w http.ResponseWriter, r *http.Request, id int) {
	var actualizada Cancion

	err := json.NewDecoder(r.Body).Decode(&actualizada)
	if err != nil {
		respError(w, 400, "JSON invalido", "El body no es JSON valido")
		return
	}

	if ok, msg := validarCancion(actualizada); !ok {
		respError(w, 400, "Validacion fallida", msg)
		return
	}

	for i, c := range canciones {
		if c.ID == id {
			actualizada.ID = id
			canciones[i] = actualizada
			guardarDatos()
			respJSON(w, 200, actualizada)
			return
		}
	}

	respError(w, 404, "No encontrado", "No existe una cancion con id "+strconv.Itoa(id))
}


//   PATCH  

func handlePatch(w http.ResponseWriter, r *http.Request, id int) {
	var campos map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&campos)
	if err != nil {
		respError(w, 400, "JSON invalido", "El body no es JSON valido")
		return
	}

	for i, c := range canciones {
		if c.ID == id {
			if val, ok := campos["cancion"]; ok {
				canciones[i].Nombre = val.(string)
			}
			if val, ok := campos["artista"]; ok {
				canciones[i].Artista = val.(string)
			}
			if val, ok := campos["album"]; ok {
				canciones[i].Album = val.(string)
			}
			if val, ok := campos["genero"]; ok {
				canciones[i].Genero = val.(string)
			}
			if val, ok := campos["duracion"]; ok {
				canciones[i].Duracion = val.(string)
			}
			if val, ok := campos["vistas"]; ok {
				canciones[i].Vistas = int(val.(float64))
			}

			guardarDatos()
			respJSON(w, 200, canciones[i])
			return
		}
	}

	respError(w, 404, "No encontrado", "No existe una cancion con id "+strconv.Itoa(id))
}
