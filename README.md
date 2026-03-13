# API de Canciones

API REST de canciones hecha en Go usando unicamente la libreria estandar. Los datos se guardan en un archivo JSON.

---

## Estructura del proyecto

```
.
├── main.go
├── data/
│   └── canciones.json
├── Dockerfile
└── README.md
```

---

## Como ejecutar

```bash
go run main.go
```

El servidor corre en el puerto `24750`.

### Con Docker

```bash
docker build -t api-canciones .
docker run -p 24750:24750 api-canciones
```

---

## Endpoints

### GET /api/canciones

Devuelve todas las canciones.

```bash
curl http://localhost:24750/api/canciones
```

### GET /api/canciones?id=1

Busca una cancion por ID usando query parameter.

```bash
curl http://localhost:24750/api/canciones?id=1
```

### GET /api/canciones con filtros combinados

Se pueden combinar los filtros `artista`, `genero`, `album` y `cancion`.

```bash
curl "http://localhost:24750/api/canciones?genero=Rock"
curl "http://localhost:24750/api/canciones?artista=d4vd&genero=Alternative"
curl "http://localhost:24750/api/canciones?cancion=fire"
```

### GET /api/canciones/{id}

Busca una cancion por ID usando path parameter.

```bash
curl http://localhost:24750/api/canciones/3
```

### POST /api/canciones

Crea una nueva cancion. Campos requeridos: `cancion`, `artista`, `genero`, `duracion`.

```bash
curl -X POST http://localhost:24750/api/canciones \
  -H "Content-Type: application/json" \
  -d '{"cancion":"Blinding Lights","artista":"The Weeknd","album":"After Hours","genero":"Pop","duracion":"3:20","vistas":4200000000}'
```

Respuesta (201 Created):
```json
{
  "id": 17,
  "cancion": "Blinding Lights",
  "artista": "The Weeknd",
  "album": "After Hours",
  "genero": "Pop",
  "duracion": "3:20",
  "vistas": 4200000000
}
```

### PUT /api/canciones/{id}

Reemplaza completamente una cancion existente.

```bash
curl -X PUT http://localhost:24750/api/canciones/1 \
  -H "Content-Type: application/json" \
  -d '{"cancion":"Circles (Remix)","artista":"Will Paquin","album":"Circles","genero":"Indie Pop","duracion":"3:00","vistas":60000000}'
```

### PATCH /api/canciones/{id}

Actualiza parcialmente una cancion. Solo se envian los campos que se quieren cambiar.

```bash
curl -X PATCH http://localhost:24750/api/canciones/1 \
  -H "Content-Type: application/json" \
  -d '{"vistas":70000000}'
```

### DELETE /api/canciones/{id}

Elimina una cancion por ID.

```bash
curl -X DELETE http://localhost:24750/api/canciones/1
```

---

## Casos de error

### Cancion no encontrada (404)

```bash
curl http://localhost:24750/api/canciones/999
```

```json
{
  "status": 404,
  "error": "No encontrado",
  "detalle": "No existe una cancion con id 999"
}
```

### JSON invalido en POST (400)

```bash
curl -X POST http://localhost:24750/api/canciones \
  -H "Content-Type: application/json" \
  -d 'esto no es json'
```

```json
{
  "status": 400,
  "error": "JSON invalido",
  "detalle": "El body no es JSON valido"
}
```

### Validacion fallida (400)

```bash
curl -X POST http://localhost:24750/api/canciones \
  -H "Content-Type: application/json" \
  -d '{"cancion":"Test"}'
```

```json
{
  "status": 400,
  "error": "Validacion fallida",
  "detalle": "El campo 'artista' es requerido"
}
```

---

## Campos de cada cancion

| Campo    | Tipo   | Requerido | Descripcion                        |
|----------|--------|-----------|------------------------------------|
| id       | int    | auto      | Se genera automaticamente          |
| cancion  | string | si        | Nombre de la cancion               |
| artista  | string | si        | Artista o banda                    |
| album    | string | no        | Album al que pertenece             |
| genero   | string | si        | Genero musical                     |
| duracion | string | si        | Duracion en formato M:SS           |
| vistas   | int    | no        | Numero de reproducciones           |

---

## Persistencia

Todos los cambios (POST, PUT, PATCH, DELETE) se guardan automaticamente en `data/canciones.json`.
