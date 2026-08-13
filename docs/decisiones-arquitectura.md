# Decisiones de arquitectura

Este documento sustenta las decisiones tomadas ante ambigüedades del enunciado
(`docs/reto-tecnico.pdf`), tal como el propio reto pide justificar en la
entrevista.

## 1. Factorización QR vs. "rotación de la matriz"

El PDF es inconsistente entre dos secciones:

- Página 2 ("Arquitectura de la solución"): *"API en Go: recibirá la matriz
  original como entrada, realizará la rotación de la matriz..."*
- Página 3 ("Funcionalidad requerida", más específica): *"Una API en Go que
  reciba... una matriz rectangular y devuelva la factorización QR de dicha
  matriz."*

**Decisión: se implementó la factorización QR.** La sección "Funcionalidad
requerida" define un contrato de entrada/salida explícito y es la fuente más
específica y operativa; la página 2 es un resumen narrativo de arquitectura,
más propenso a imprecisión de redacción. Ante una contradicción entre un
resumen general y una especificación detallada, se prioriza la especificación
detallada.

## 2. Householder vs. Gram-Schmidt para la factorización QR

Se evaluaron dos algoritmos clásicos:

| Criterio | Householder | Gram-Schmidt modificado (MGS) |
|---|---|---|
| Estabilidad numérica | Alta (estándar en LAPACK/numpy) | Media (mejor que GS clásico, pero inferior a Householder) |
| Complejidad de implementación | ~120 líneas con casos borde | Similar |
| Manejo de rango deficiente | Se detecta por norma casi nula del vector reflejado | Se detecta por norma casi nula al normalizar |

**Decisión: Householder reflections.** Mismo costo de implementación que MGS,
pero con mejor estabilidad numérica — es la elección estándar en
implementaciones de referencia (LAPACK, NumPy `linalg.qr`).

**Caso `m < n` (menos filas que columnas):** se rechaza con un error 400
(`ErrTooFewRows`) en vez de transponer la matriz silenciosamente. Transponer
cambiaría la semántica de lo que el usuario envió; el enunciado no pide ese
comportamiento, así que un input inválido debe fallar explícitamente, no
transformarse en otra cosa.

## 3. Diseño de autenticación JWT

- Un único emisor de tokens: `POST /auth/token` en la API de Node, con
  credenciales demo fijas por variable de entorno (`DEMO_CLIENT_ID` /
  `DEMO_CLIENT_SECRET`), ya que el enunciado no define un modelo de usuarios.
- La API de Go **no** llama a ese endpoint para su comunicación interna con
  Node: firma su propio JWT de servicio (`IssueServiceToken`) con el mismo
  `JWT_SECRET` compartido. Evita una llamada HTTP adicional y un punto de
  fallo extra en cada request de usuario.
- Secreto simétrico HS256 vía variable de entorno — no se justifica un par de
  claves RS256 para este alcance (no hay múltiples emisores ni necesidad de
  verificación pública del token).
- Expiración de 1 hora para tokens de usuario/demo, 5 minutos para tokens de
  servicio (de vida más corta porque se emiten y consumen en el mismo
  request).
- Mensaje de error 401 uniforme (`"missing or invalid authorization token"`)
  en ambas APIs para que el frontend maneje un solo caso.

## 4. Topología de despliegue en Render

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

## 5. Sin capa de repositorio ni gestor de estado global en el frontend

Ambas APIs son completamente stateless (no hay persistencia), por lo que no
se agregó una capa de repositorio/DAO — habría sido una abstracción sin
propósito. De la misma forma, el frontend usa `useState` local en `App.tsx`
en vez de Redux/Zustand/Context: hay dos piezas de estado (`token`,
`result`) en una sola pantalla, lo cual no justifica un gestor de estado
global.

## 6. Tipado y stack de cada servicio

Decisión explícita del candidato: tipar todo el stack. Go es tipado por
naturaleza; la API de Node se escribió en TypeScript (compilada a `dist/`
para producción) y el frontend usa el template Vite `react-ts`. Los tipos de
dominio (`Matrix`, `StatsResponse`, `QrResponse`, `TokenResponse`) se
duplican entre `node-api/src/types` y `frontend/src/types` sin un paquete
compartido — una decisión consciente para no introducir un monorepo con
gestión de paquetes (ej. workspaces) solo para 4 interfaces, dado el alcance
del challenge.
