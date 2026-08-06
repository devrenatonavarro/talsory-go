# Matrix QR Factorization & Rotation API

API REST construida en **Go 1.22+** utilizando **Fiber v2** para el cálculo de descomposición QR ($A = Q \cdot R$) de matrices rectangulares, rotación de matrices en 90 grados, e integración HTTP resiliente con un servicio externo en Node.js.

---

## 🚀 Características Principales

1. **Factorización QR mediante Gram-Schmidt Modificado (MGS)**:
   - Descompone una matriz $A_{m \times n}$ ($m \ge n$) en una matriz ortogonal $Q_{m \times n}$ ($Q^T Q = I_n$) y una matriz triangular superior $R_{n \times n}$.
   - Algoritmo numéricamente estable implementado sin librerías externas de álgebra lineal.
2. **Rotación de Matriz (90°)**:
   - Rota matrices rectangulares en dirección `clockwise` (sentido horario) y `counterclockwise` (sentido antihorario).
3. **Cliente HTTP Resiliente hacia Node.js**:
   - Envía las matrices resultantes $Q$ y $R$ al endpoint `POST {NODE_API_URL}/api/v1/stats`.
   - Timeout configurable vía `context.WithTimeout`.
   - En caso de falla o inalcanzabilidad de la API de Node.js, responde con **502 Bad Gateway** y un mensaje estructurado claro sin detener la ejecución del servidor Go.
4. **Validación Rigurosa de Entradas**:
   - Detecta matrices vacías, filas de longitud inconsistente (matrices irregulares) y valores no numéricos (`NaN` o `Infinito`), devolviendo respuestas `400 Bad Request`.
5. **Middlewares Incorporados**:
   - Logging estructurado con `X-Request-ID`, método HTTP, ruta, código de respuesta y tiempo de ejecución.
   - Habilitación de CORS (Cross-Origin Resource Sharing).
   - Autenticación opcional mediante **JWT** en `/api/v1/*` (activable vía variable de entorno `AUTH_ENABLED=true`).

---

## 🧮 Algoritmo QR Elegido y Razón

Se eligió el algoritmo de **Gram-Schmidt Modificado (MGS)** (*Modified Gram-Schmidt*).

### ¿Por qué MGS?
En la versión clásica de Gram-Schmidt, el error de redondeo de punto flotante acumula rápidamente una pérdida de ortogonalidad entre los vectores calculados. **Modified Gram-Schmidt** corrige esto orthogonalizando los vectores restantes **inmediatamente** contra el nuevo vector ortonormal $q_k$ recién calculado en cada paso de la iteración.

### Complejidad Computacional
- **Tiempo**: $\mathcal{O}(m \cdot n^2)$ operaciones flotantes para una matriz de dimensiones $m \times n$.
- **Espacio**: $\mathcal{O}(m \cdot n + n^2)$ memoria requerida para almacenar las matrices de salida $Q$ ($m \times n$) y $R$ ($n \times n$).

---

## 📁 Estructura del Proyecto

```text
.
├── cmd/
│   └── server/
│       └── main.go           # Punto de entrada de la aplicación y cableado de rutas
├── internal/
│   ├── client/
│   │   └── node_client.go    # Cliente HTTP hacia Node API con timeout y manejo 502
│   ├── config/
│   │   └── config.go         # Carga de variables de entorno (.env via godotenv)
│   ├── handlers/
│   │   ├── matrix_handler.go # Controladores de los endpoints REST (Fiber handlers)
│   │   └── matrix_handler_test.go
│   ├── middleware/
│   │   ├── auth.go           # Middleware JWT opcional
│   │   └── logger.go         # Middleware de logging y Request ID
│   ├── models/
│   │   └── matrix.go         # Estructuras DTO de entrada, salida y errores
│   └── services/
│       ├── matrix.go         # Validación y rotación de matrices (90°)
│       ├── matrix_test.go
│       ├── qr.go             # Algoritmo de Factorización QR (GoDoc + MGS)
│       └── qr_test.go
├── .env.example              # Plantilla de variables de entorno
├── Dockerfile                # Build multi-stage (golang:1.22-alpine -> alpine:3.19)
├── go.mod
├── go.sum
└── README.md
```

---

## ⚙️ Variables de Entorno

Puedes configurar la aplicación creando un archivo `.env` basado en `.env.example`:

| Variable | Descripción | Valor por Defecto |
| :--- | :--- | :--- |
| `PORT` | Puerto HTTP del servidor Fiber | `3000` |
| `NODE_API_URL` | URL base del microservicio en Node.js | `http://localhost:4000` |
| `HTTP_TIMEOUT_MS` | Timeout de peticiones hacia Node.js (en ms) | `5000` |
| `AUTH_ENABLED` | Habilita/deshabilita middleware JWT (`true`/`false`) | `false` |
| `JWT_SECRET` | Clave secreta para firma y verificación de JWT | `supersecretkey` |

---

## 🛠️ Cómo Ejecutar Localmente

### Requisitos Previos
- Go 1.22 o superior.

### Pasos

1. Clonar o descargar el repositorio.
2. Instalar dependencias:
   ```bash
   go mod download
   ```
3. Ejecutar las pruebas unitarias y de integración:
   ```bash
   go test -v -cover ./...
   ```
4. Iniciar el servidor:
   ```bash
   go run cmd/server/main.go
   ```

El servidor estará corriendo en `http://localhost:3000`.

---

## 🐳 Cómo Ejecutar con Docker

### 1. Construir la Imagen Docker (Multi-stage)

```bash
docker build -t matrix-qr-service .
```

### 2. Ejecutar el Contenedor

```bash
docker run -d \
  -p 3000:3000 \
  -e PORT=3000 \
  -e NODE_API_URL="http://host.docker.internal:4000" \
  --name matrix-service \
  matrix-qr-service
```

---

## 📡 Ejemplos de cURL para cada Endpoint

### 1. Health Check

```bash
curl -X GET http://localhost:3000/health
```

**Respuesta Exitosa (200 OK):**
```json
{
  "status": "ok"
}
```

---

### 2. Factorización QR

```bash
curl -X POST http://localhost:3000/api/v1/matrix/qr \
  -H "Content-Type: application/json" \
  -d '{
    "matrix": [
      [1, 2],
      [3, 4],
      [5, 6]
    ]
  }'
```

**Respuesta Exitosa (200 OK):**
```json
{
  "q": [
    [0.1690308509457033, 0.8970852271427503],
    [0.50709255283711, 0.2760262237362308],
    [0.8451542547285165, -0.3450327796702888]
  ],
  "r": [
    [5.916079783099616, 7.437357443329971],
    [0, 0.8280786712086925]
  ],
  "stats": {
    "norm_q": 1.0,
    "status": "processed"
  }
}
```

**Respuesta de Error en Validación (400 Bad Request):**
```bash
curl -X POST http://localhost:3000/api/v1/matrix/qr \
  -H "Content-Type: application/json" \
  -d '{
    "matrix": [
      [1, 2],
      [3]
    ]
  }'
```
```json
{
  "error": "Matriz de entrada no válida",
  "details": "todas las filas de la matriz deben tener la misma longitud: fila 1 tiene longitud 1, se esperaba 2"
}
```

**Respuesta si Node API está inalcanzable (502 Bad Gateway):**
```json
{
  "error": "Error de comunicación con el servicio externo Node.js (502 Bad Gateway)",
  "details": "no se pudo conectar con Node API (http://localhost:4000): dial tcp 127.0.0.1:4000: connect: connection refused"
}
```

---

### 3. Rotación de Matriz (90°)

#### En sentido horario (Clockwise - Default)

```bash
curl -X POST http://localhost:3000/api/v1/matrix/rotate \
  -H "Content-Type: application/json" \
  -d '{
    "matrix": [
      [1, 2],
      [3, 4],
      [5, 6]
    ],
    "direction": "clockwise"
  }'
```

**Respuesta (200 OK):**
```json
{
  "rotated": [
    [5, 3, 1],
    [6, 4, 2]
  ]
}
```

#### En sentido antihorario (Counterclockwise)

```bash
curl -X POST http://localhost:3000/api/v1/matrix/rotate \
  -H "Content-Type: application/json" \
  -d '{
    "matrix": [
      [1, 2],
      [3, 4],
      [5, 6]
    ],
    "direction": "counterclockwise"
  }'
```

**Respuesta (200 OK):**
```json
{
  "rotated": [
    [2, 4, 6],
    [1, 3, 5]
  ]
}
```

---

### 4. Con Autenticación JWT (si `AUTH_ENABLED=true`)

Para realizar solicitudes cuando la autenticación está activa, incluya el encabezado `Authorization`:

```bash
curl -X POST http://localhost:3000/api/v1/matrix/rotate \
  -H "Authorization: Bearer <TU_JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "matrix": [[1, 2], [3, 4]]
  }'
```
