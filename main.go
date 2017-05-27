package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"time"
	"log"
	"encoding/json"
	"strconv"
)


type Note struct {
	Title string `json:"title"`
	Description string `json:"description"`
	CreatedAt time.Time        `json:"created_at"`
}

var noteStore = make(map[string]Note)
var id int

//GETNoteHandler - GET . api/notes
func GETNoteHandler(w http.ResponseWriter, r *http.Request){
	var notes []Note
	for _, v:= range noteStore {
		notes = append(notes,v)
	}

	//Hacer la cabecera http
	w.Header().Set("Content-Type","application/json")
	j,err:=json.Marshal(notes)
	if err != nil {
		panic(err)
	}
	//Codigo 200
	w.WriteHeader(http.StatusOK)
	//escribimos el body en json
	w.Write(j)

}

func POSTNoteHandler(w http.ResponseWriter, r *http.Request){
	var note Note
	//Decodificando el json que llega a la estructura
	err:=json.NewDecoder(r.Body).Decode(&note)
	if err != nil {
		panic(err)
	}
	//Asignamos la fecha al objeto
	note.CreatedAt = time.Now()
	//Incrementamos id
	id++
	//Entero a string para la llave del noteStore
	k:= strconv.Itoa(id)
	//Asignamos el objeto con el id
	noteStore[k]=note

	//Regresamos el mismo objeto para confirmar
	//Hacer la cabecera http
	w.Header().Set("Content-Type","application/json")
	j,err:=json.Marshal(note)
	if err != nil {
		panic(err)
	}
	//Codigo 201 objeto Creado
	w.WriteHeader(http.StatusCreated)
	//escribimos el body en json
	w.Write(j)


}

//Actualizar la nota
func PutNoteHandler(w http.ResponseWriter, r *http.Request){
	//Extraer variables del request
	vars:=mux.Vars(r)
	//valor del Id de la nota
	k:= vars["id"]
	var noteUpdate Note
	err:= json.NewDecoder(r.Body).Decode(&noteUpdate)
	if err != nil {
		panic(err)
	}

	//Checar si el id existe en nuestro Map
	//Trata de extraer el dato y guarda un booleano en ok y lo checa
	if note, ok := noteStore[k]; ok {
		//Asignarle la fecha
		noteUpdate.CreatedAt = note.CreatedAt
		//      map, id
		delete(noteStore,k)
		noteStore[k] =  noteUpdate
	} else{
		//Si no existe la nota
		log.Printf("No encontramos el id %s ", k)
	}

	//codigo 204
	w.WriteHeader(http.StatusNoContent)
}

// DELETE
func DeleteNoteHandler(w http.ResponseWriter, r *http.Request) {
	//Extraer variables del request
	vars:=mux.Vars(r)
	//valor del Id de la nota
	k:= vars["id"]
	//Nos llega un id y se borra
	if _, ok := noteStore[k]; ok {
		//Se borra la nota
		delete(noteStore,k)

	} else{
		//Si no existe la nota
		log.Printf("No encontramos el id %s ", k)
	}

	//codigo 204
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter().StrictSlash(false)
	r.HandleFunc("/api/notes", GETNoteHandler).Methods("GET")
	r.HandleFunc("/api/notes", POSTNoteHandler).Methods("POST")
	r.HandleFunc("/api/notes/{id}", PutNoteHandler).Methods("PUT")
	r.HandleFunc("/api/notes/{id}", DeleteNoteHandler).Methods("DELETE")

	server:=&http.Server{
		Addr:":8080",
		Handler: r,
		ReadTimeout: 10 * time.Second,
		WriteTimeout: 10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Println("Listening...")
	server.ListenAndServe()

}
