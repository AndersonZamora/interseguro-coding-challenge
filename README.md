# Reto técnico — Interseguro

Dos APIs comunicadas por HTTP:

- **`go-api`** (Go + Fiber): recibe una matriz, calcula su **factorización QR**
  (Householder), y reenvía Q y R a `node-api`.
- **`node-api`** (Node.js + Express + TypeScript): recibe Q y R, calcula
  estadísticas (máximo, mínimo, promedio, suma total, si alguna matriz es
  diagonal) y las devuelve.
- **`frontend`** (React + Vite + TypeScript + Tailwind): consume ambas APIs y
  muestra el resultado.

Las decisiones ante ambigüedades del enunciado (QR vs. rotación, algoritmo de
factorización, diseño de JWT, topología de despliegue) están documentadas en
la sección [Decisiones de arquitectura](#decisiones-de-arquitectura).

## Arquitectura

```
Usuario ──▶ frontend (React/Vite)
              │  1. POST /auth/token         ──▶ node-api
              │  2. POST /api/v1/qr (Bearer) ──▶ go-api
              │                                    │
              │                                    │ JWT de servicio
              │                                    ▼
              │                              POST /api/v1/stats ──▶ node-api
              │                                    │
              │◀────────── { q, r, stats } ────────┘
```

## Estructura del repo

```
go-api/       API en Go (Fiber) — factorización QR
node-api/     API en Node.js (Express + TS) — estadísticas + emisión de JWT
frontend/     SPA en React + Vite + TS + Tailwind
docker-compose.yml   Levanta las 3 apps localmente
render.yaml           Blueprint de despliegue en Render
```

## Requisitos

- Docker + Docker Compose (forma recomendada de correr todo junto).
- Para desarrollo sin Docker: Go 1.22+, Node.js 20+, npm.

## Desarrollo local con Docker Compose

```bash
cp .env.example .env   # editar JWT_SECRET / DEMO_CLIENT_SECRET si se desea
docker compose up --build
```

- Frontend: http://localhost:8081
- API Go: http://localhost:8080
- API Node: http://localhost:3000

### Probar el flujo con curl

```bash
# 1. Obtener un token
TOKEN=$(curl -s -X POST http://localhost:3000/auth/token \
  -H "Content-Type: application/json" \
  -d '{"client_id":"interseguro-demo","client_secret":"<DEMO_CLIENT_SECRET del .env>"}' \
  | node -pe "JSON.parse(require('fs').readFileSync(0,'utf8')).access_token")

# 2. Calcular QR + estadísticas
curl -s -X POST http://localhost:8080/api/v1/qr \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"matrix": [[12,-51,4],[6,167,-68],[-4,24,-41]]}'
```

## Desarrollo de cada servicio por separado

### go-api

Go no es un requisito estricto si se usa Docker; para correrlo sin instalar Go
localmente:

```bash
docker run --rm -v "$(pwd)/go-api:/app" -w /app golang:1.22 go test ./... -v
```

Con Go instalado localmente:

```bash
cd go-api
go run ./cmd/api      # requiere JWT_SECRET y NODE_API_URL en el entorno
go test ./...
```

### node-api

```bash
cd node-api
npm install
npm run dev     # requiere .env con JWT_SECRET y DEMO_CLIENT_SECRET
npm test
npm run build && npm start   # build de producción
```

### frontend

```bash
cd frontend
npm install
cp .env.example .env.local
npm run dev
npm run build   # type-check + build de producción a dist/
```

## Variables de entorno

| Variable | Servicio | Descripción |
|---|---|---|
| `JWT_SECRET` | go-api, node-api | Secreto HMAC compartido para firmar/verificar JWT. **Debe ser idéntico en ambos.** |
| `NODE_API_URL` | go-api | URL base de node-api (ej. `http://node-api:3000` en compose). |
| `ALLOWED_ORIGIN` | go-api, node-api | Origen permitido por CORS (URL del frontend). |
| `DEMO_CLIENT_ID` / `DEMO_CLIENT_SECRET` | node-api | Credenciales de demo para `POST /auth/token`. |
| `PORT` | go-api, node-api | Puerto HTTP (default 8080 / 3000). |
| `VITE_GO_API_URL` / `VITE_NODE_API_URL` | frontend (build-time) | URLs públicas de las APIs, inyectadas al build de Vite. |

## Contrato de las APIs

### `POST /auth/token` (node-api, sin auth)

```json
// Request
{ "client_id": "interseguro-demo", "client_secret": "..." }
// Response 200
{ "access_token": "ey...", "expires_in": 3600 }
```

### `POST /api/v1/qr` (go-api, requiere `Authorization: Bearer <token>`)

```json
// Request
{ "matrix": [[12,-51,4],[6,167,-68],[-4,24,-41]] }
// Response 200
{
  "q": [[...]], "r": [[...]],
  "stats": { "max": 70, "min": -175, "average": -5.22, "sum": -93.92, "anyDiagonal": false }
}
```

Errores: `400` (matriz vacía/irregular/`m<n`/rango deficiente/JSON inválido),
`401` (token faltante/inválido), `502` (node-api no disponible).

### `POST /api/v1/stats` (node-api, requiere `Authorization: Bearer <token>`)

```json
// Request
{ "q": [[...]], "r": [[...]] }
// Response 200
{ "max": 70, "min": -175, "average": -5.22, "sum": -93.92, "anyDiagonal": false }
```

## Tests

- `go-api`: `go test ./...` (o vía Docker, ver arriba) — cubre la
  factorización QR (reconstrucción `A≈QR`, ortogonalidad de Q, casos de
  error), JWT (emisión/expiración/secreto inválido) y el handler HTTP con un
  `StatsFetcher` mockeado.
- `node-api`: `npm test` (Jest + Supertest) — cubre el servicio de
  estadísticas (casos borde: 1x1, negativos, diagonal real, rectangular) y
  las rutas HTTP (401 sin token, 200 caso feliz, 400 payload inválido).

## Despliegue en Render

El repo incluye [`render.yaml`](render.yaml) (Blueprint) con 3 servicios:
`interseguro-go-api`, `interseguro-node-api` (ambos Docker) e
`interseguro-frontend` (Static Site, build de Vite).

Como las URLs públicas de cada servicio solo se conocen después de crearlo,
algunas variables quedan marcadas para completar a mano. Checklist:

1. En Render, "New" → "Blueprint" → conectar este repo. Render crea los 3
   servicios a partir de `render.yaml`.
2. Esperar el primer deploy de `interseguro-node-api` y copiar:
   - Su `JWT_SECRET` generado (Dashboard → Environment) → pegarlo también en
     `interseguro-go-api` (deben ser idénticos).
   - Su URL pública (ej. `https://interseguro-node-api.onrender.com`) →
     pegarla en `NODE_API_URL` de `interseguro-go-api` y en
     `VITE_NODE_API_URL` de `interseguro-frontend`.
3. Esperar el deploy de `interseguro-go-api` y copiar su URL pública →
   `VITE_GO_API_URL` de `interseguro-frontend`.
4. Esperar el deploy de `interseguro-frontend` y copiar su URL pública →
   `ALLOWED_ORIGIN` en `interseguro-go-api` y `interseguro-node-api`.
5. Redeploy manual de `interseguro-go-api` y `interseguro-node-api` para que
   tomen el `ALLOWED_ORIGIN` actualizado.

Alternativa sin Blueprint: crear los 3 servicios a mano en el dashboard de
Render (2 "Web Service" con runtime Docker apuntando a
`go-api/Dockerfile` / `node-api/Dockerfile`, 1 "Static Site" con build
command `cd frontend && npm ci && npm run build` y publish directory
`frontend/dist`), configurando las mismas variables de entorno de la tabla
de arriba.

## Decisiones de arquitectura

Este apartado sustenta las decisiones tomadas ante ambigüedades del
enunciado, tal como el propio reto pide justificar en la entrevista.

### 1. Factorización QR vs. "rotación de la matriz"

El enunciado es inconsistente entre dos secciones: la descripción general de
arquitectura pide que la API de Go "realice la rotación de la matriz",
mientras que la sección de "Funcionalidad requerida" (más específica) pide
que "reciba una matriz rectangular y devuelva la factorización QR de dicha
matriz".

**Decisión: se implementó la factorización QR.** La sección de
funcionalidad requerida define un contrato de entrada/salida explícito y es
la fuente más específica y operativa; la descripción general es un resumen
narrativo, más propenso a imprecisión de redacción. Ante una contradicción
entre un resumen general y una especificación detallada, se prioriza la
especificación detallada.

### 2. Householder vs. Gram-Schmidt para la factorización QR

Se evaluaron dos algoritmos clásicos:

| Criterio | Householder | Gram-Schmidt modificado (MGS) |
|---|---|---|
| Estabilidad numérica | Alta (estándar en LAPACK/numpy) | Media (mejor que GS clásico, pero inferior a Householder) |
| Complejidad de implementación | ~120 líneas con casos borde | Similar |
| Manejo de rango deficiente | Se detecta por norma casi nula del vector reflejado | Se detecta por norma casi nula al normalizar |

**Decisión: Householder reflections.** Mismo costo de implementación que
MGS, pero con mejor estabilidad numérica — es la elección estándar en
implementaciones de referencia (LAPACK, NumPy `linalg.qr`).

**Caso `m < n` (menos filas que columnas):** se rechaza con un error 400
(`ErrTooFewRows`) en vez de transponer la matriz silenciosamente. Transponer
cambiaría la semántica de lo que el usuario envió; el enunciado no pide ese
comportamiento, así que un input inválido debe fallar explícitamente, no
transformarse en otra cosa.

### 3. Diseño de autenticación JWT

- Un único emisor de tokens: `POST /auth/token` en la API de Node, con
  credenciales demo fijas por variable de entorno (`DEMO_CLIENT_ID` /
  `DEMO_CLIENT_SECRET`), ya que el enunciado no define un modelo de usuarios
  (el requisito literal es solo "aplicar un nivel de seguridad utilizando
  JWT para proteger las consultas a las APIs").
- La API de Go **no** llama a ese endpoint para su comunicación interna con
  Node: firma su propio JWT de servicio (`IssueServiceToken`) con el mismo
  `JWT_SECRET` compartido. Evita una llamada HTTP adicional y un punto de
  fallo extra en cada request de usuario.
- Secreto simétrico HS256 vía variable de entorno — no se justifica un par
  de claves RS256 para este alcance (no hay múltiples emisores ni necesidad
  de verificación pública del token).
- Expiración de 1 hora para tokens de usuario/demo, 5 minutos para tokens de
  servicio (de vida más corta porque se emiten y consumen en el mismo
  request).
- Mensaje de error 401 uniforme (`"missing or invalid authorization token"`)
  en ambas APIs para que el frontend maneje un solo caso.
- El token del frontend se cachea en `localStorage` (con verificación de
  expiración al cargar) para sobrevivir a un refresh de página — mejora de
  UX, no un requisito del reto.

### 4. Topología de despliegue en Render

Se desplegaron 3 servicios independientes en vez de empaquetar el frontend
dentro del contenedor de Node:

- `interseguro-go-api` y `interseguro-node-api`: servicios web Docker.
- `interseguro-frontend`: **Static Site** nativo de Render (build de Vite,
  sin Docker).

Empaquetar el frontend como archivos estáticos servidos por Express habría
evitado un tercer dominio, pero mezcla la responsabilidad de servir una SPA
con la de una API de negocio, y los Static Sites de Render son gratuitos y
más simples de configurar (build command + publish dir, sin Dockerfile). El
costo es que ambas APIs necesitan CORS habilitado hacia el origen del
frontend — un middleware que de todas formas era necesario.

### 5. Sin capa de repositorio ni gestor de estado global en el frontend

Ambas APIs son completamente stateless (no hay persistencia), por lo que no
se agregó una capa de repositorio/DAO — habría sido una abstracción sin
propósito. De la misma forma, el frontend usa `useState` local en `App.tsx`
en vez de Redux/Zustand/Context: hay dos piezas de estado (`token`,
`result`) en una sola pantalla, lo cual no justifica un gestor de estado
global.

### 6. Tipado y stack de cada servicio

Decisión explícita: tipar todo el stack. Go es tipado por naturaleza; la API
de Node se escribió en TypeScript (compilada a `dist/` para producción) y el
frontend usa el template Vite `react-ts`. Los tipos de dominio (`Matrix`,
`StatsResponse`, `QrResponse`, `TokenResponse`) se duplican entre
`node-api/src/types` y `frontend/src/types` sin un paquete compartido — una
decisión consciente para no introducir un monorepo con gestión de paquetes
(ej. workspaces) solo para 4 interfaces, dado el alcance del challenge.
