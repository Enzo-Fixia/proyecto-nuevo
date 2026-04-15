# Agent Handoff — Backend Go + MySQL (FIXIA)

> **Propósito de este documento**
> Dar a otro agente de IA (Claude Code, Cursor, Copilot Workspace, etc.) todo el contexto necesario para **entender el estado actual** del proyecto `golang-api`, **levantarlo en una máquina nueva**, y **extenderlo** sin romper las convenciones ya establecidas.
>
> El proyecto **ya está construido y funcional** — este NO es un prompt de construcción desde cero. Si necesitás el prompt original de construcción, mirá el historial de git.

---

## 1. Estado actual del proyecto

- **Directorio raíz del código:** `golang-api/`
- **Stack:** Go 1.22+, Gin, GORM, MySQL 8, Flyway (migraciones), JWT, bcrypt, godotenv, testify
- **Arquitectura:** capas `model → repository → service → handler → routes` por módulo
- **Módulos implementados:** `auth`, `user`, `product`, `order` (con transacción GORM para crear órdenes)
- **Tests:** unit tests en `internal/user/service_test.go` e `internal/product/service_test.go`
- **API disponible en:** `http://localhost:8080`
- **Collection Postman:** `golang-api/postman_collection.json` (importable directo)

### Decisiones de diseño clave (NO reinventar)

1. **MySQL corre LOCAL en el host, NO en Docker.** El `docker-compose.yml` solo levanta la app Go y Flyway (bajo profile `tools`). La app se conecta a la instancia MySQL del sistema operativo del desarrollador vía `host.docker.internal` (desde contenedor) o `localhost` (cuando se corre `go run` local).
   - **Razón:** el dev usa su MySQL ya instalado; persistencia fuera de contenedores; más rápido iterar.
2. **El esquema es 100 % responsabilidad de Flyway.** `database/db.go` NO llama `AutoMigrate`. Las tablas existen porque corrieron las migraciones `migrations/V*.sql`, no porque GORM las creó.
   - **Razón:** schema versionado, auditable, reversible por migración inversa. Dos devs pueden crear migraciones en paralelo sin conflictos (timestamp único).
3. **Convención Flyway: timestamp.** Archivos con formato `V{YYYYMMDDHHmmss}__descripcion.sql`. Header SQL estándar (ver sección 5).
4. **Módulo Go:** `github.com/fixia/golang-api` (NO cambiar sin actualizar todos los imports).
5. **Config:** `godotenv` (no Viper). El `.env` se carga una vez en `config/config.go`.

---

## 2. Estructura del proyecto (actual, real)

```
golang-api/
├── cmd/main.go                          # Entry point: carga config, conecta DB, registra rutas, arranca Gin
├── config/config.go                     # Lee .env con godotenv
├── database/db.go                       # Conexión GORM a MySQL (SIN AutoMigrate)
├── internal/
│   ├── auth/
│   │   ├── jwt.go                       # Generación/validación de JWT (HS256)
│   │   └── middleware.go                # AuthRequired + RequireRole("admin")
│   ├── user/
│   │   ├── model.go                     # User, RegisterRequest, LoginRequest, UpdateUserRequest, LoginResponse
│   │   ├── repository.go                # Interfaz + impl GORM
│   │   ├── service.go                   # Lógica: register, login, CRUD
│   │   ├── handler.go                   # Handlers Gin
│   │   ├── routes.go                    # Registro de rutas /auth/* y /users/*
│   │   └── service_test.go              # Unit tests con testify
│   ├── product/
│   │   ├── model.go                     # Product, CreateProductRequest, UpdateProductRequest
│   │   ├── repository.go, service.go, handler.go, routes.go
│   │   └── service_test.go
│   └── order/
│       ├── model.go                     # Order, OrderItem, CreateOrderRequest, OrderItemRequest, UpdateStatusRequest
│       └── repository.go, service.go, handler.go, routes.go
├── middleware/logger.go                 # Logger HTTP simple
├── utils/
│   ├── response.go                      # APIResponse + SuccessResponse/ErrorResponse helpers
│   └── errors.go                        # ErrDuplicateEmail, ErrInvalidCredentials, ErrUnauthorized, ErrInsufficientStock
├── migrations/                          # Flyway SQL versionado
│   ├── V20260415100000__create_users_table.sql
│   ├── V20260415100100__create_products_table.sql
│   ├── V20260415100200__create_orders_table.sql
│   └── V20260415100300__create_order_items_table.sql
├── .env                                 # NO commitear (está en .gitignore)
├── .env.example                         # Plantilla de variables
├── Dockerfile                           # Multi-stage: golang:1.22-alpine → alpine:latest
├── docker-compose.yml                   # Solo servicios `app` y `flyway` (profile tools)
├── postman_collection.json              # Collection con todos los endpoints + auto-save de JWT
├── README.md
├── go.mod                               # module github.com/fixia/golang-api
└── go.sum
```

---

## 3. Levantar el proyecto desde cero (en máquina nueva)

### 3.1 Requisitos

| Herramienta | Uso | ¿Obligatorio? |
|---|---|---|
| Docker Desktop | Correr la app y Flyway | Sí |
| MySQL 8 local | Base de datos | Sí |
| MySQL Workbench (o equivalente) | Crear la BD inicial y administrar | Recomendado |
| Go 1.22+ | Solo si se quiere correr `go run` o tests locales | Opcional |

> En Windows, Docker Desktop debe estar corriendo (ícono verde) antes de cualquier `docker-compose`. Si aparece `failed to connect to the docker API at npipe:////./pipe/...`, es que el daemon no arrancó.

### 3.2 Pasos en orden

1. **Clonar el repo** y posicionarse en `golang-api/`.
2. **Crear la BD en MySQL local** (una sola vez):
   ```sql
   CREATE DATABASE golang_api CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
   ```
   La forma más simple es MySQL Workbench → conectar a la instancia local → nueva query → ejecutar.
3. **Configurar `.env`**:
   ```bash
   cp .env.example .env
   ```
   Editar `DB_PASSWORD` con la password real del root MySQL local. El resto de defaults suele andar:
   ```env
   APP_PORT=8080
   APP_ENV=development
   DB_HOST=localhost
   DB_PORT=3306
   DB_USER=root
   DB_PASSWORD=<tu_password_real>
   DB_NAME=golang_api
   JWT_SECRET=supersecreto_cambialo_en_produccion
   JWT_EXPIRATION_HOURS=24
   ```
   > Cuando la app corre dentro del contenedor Docker, el compose override (ver `docker-compose.yml`) fuerza `DB_HOST=host.docker.internal` automáticamente. No hay que editar `.env` para eso.
4. **Correr las migraciones**:
   ```bash
   docker-compose run --rm flyway info       # ver estado
   docker-compose run --rm flyway migrate    # aplicar pendientes
   ```
5. **Levantar la app**:
   ```bash
   docker-compose up --build -d app
   docker-compose logs -f app                # ver logs
   ```
   La API queda en `http://localhost:8080`. El health check responde en `GET /health`.
6. **(Opcional) Promover un usuario a admin** para poder probar endpoints de admin:
   ```sql
   UPDATE users SET role = 'admin' WHERE email = 'tu@email.com';
   ```
   Después **volver a loguearse** para que el JWT incluya el rol nuevo.

### 3.3 Levantar con `go run` (sin Docker para la app)

Útil para iteración rápida en Go. Requiere Go 1.22+ instalado.

```bash
cd golang-api
go mod tidy
go run ./cmd/main.go
```

En este modo, `DB_HOST=localhost` del `.env` funciona tal cual.

### 3.4 Correr tests

```bash
go test ./...
```

---

## 4. Endpoints y bodies (referencia canónica)

Todas las rutas bajo `/api/v1`, excepto `/health`. Auth por header `Authorization: Bearer <jwt>`.

### Auth

| Método | Ruta | Auth | Body |
|---|---|---|---|
| POST | `/auth/register` | ❌ | `{ first_name, last_name, email, password }` |
| POST | `/auth/login` | ❌ | `{ email, password }` → `{ token, user }` |
| GET | `/auth/me` | ✅ | — |

Reglas de validación (ver `internal/user/model.go`):
- `first_name`, `last_name`: required, min 2, max 50
- `email`: required, email válido
- `password` (register): required, min 8

### Users (admin)

| Método | Ruta | Body |
|---|---|---|
| GET | `/users` | — |
| GET | `/users/:id` | — |
| PUT | `/users/:id` | `{ first_name?, last_name?, role?, is_active? }` (role ∈ user\|admin) |
| DELETE | `/users/:id` | — (soft delete vía GORM `DeletedAt`) |

### Products

| Método | Ruta | Auth | Body |
|---|---|---|---|
| GET | `/products` | ❌ | — |
| GET | `/products/:id` | ❌ | — |
| POST | `/products` | ✅ admin | `{ name, description?, price, stock?, category?, image_url? }` |
| PUT | `/products/:id` | ✅ admin | cualquier subset de los campos anteriores |
| DELETE | `/products/:id` | ✅ admin | — |

### Orders

| Método | Ruta | Auth | Body |
|---|---|---|---|
| POST | `/orders` | ✅ | `{ items: [{ product_id, quantity }, ...] }` |
| GET | `/orders` | ✅ | — (solo las del usuario autenticado) |
| GET | `/orders/:id` | ✅ | — (el handler valida que sea propia o admin) |
| PUT | `/orders/:id/status` | ✅ admin | `{ status }` (∈ pending\|confirmed\|shipped\|delivered\|cancelled) |

`POST /orders` corre dentro de una **transacción GORM**: valida stock, descuenta inventario, graba `Order` + `OrderItem`s y calcula `total`. Si cualquier paso falla → rollback.

### Formato de respuesta estándar

```json
{ "success": true, "message": "...", "data": { ... } }
{ "success": false, "error": "..." }
```

Helpers en `utils/response.go` (`SuccessResponse`, `ErrorResponse`).

### Collection Postman

`postman_collection.json` en la raíz del proyecto:
- Import directo en Postman.
- Variables: `baseUrl`, `token`, `userId`, `productId`, `orderId`.
- El request `Auth > Login` tiene un test script que guarda el JWT automáticamente en `{{token}}`.
- Auth a nivel colección (Bearer `{{token}}`); las requests públicas lo sobrescriben con `noauth`.

---

## 5. Cómo agregar una migración nueva (Flyway)

**Nunca editar un archivo de migración ya aplicado.** Siempre crear uno nuevo con timestamp posterior.

### 5.1 Generar timestamp

- Git Bash / Linux / Mac: `date +V%Y%m%d%H%M%S` → `V20260415143027`
- PowerShell: `Get-Date -Format "\VyyyyMMddHHmmss"`

### 5.2 Crear archivo

`migrations/V20260415143027__add_phone_to_users.sql`:

```sql
-- =============================================================================
-- Migration: add phone column to users
-- Version:   V20260415143027
-- Date:      2026-04-15 14:30:27
-- Author:    tu@email.com
-- =============================================================================

ALTER TABLE users ADD COLUMN phone VARCHAR(30);
```

### 5.3 Aplicar

```bash
docker-compose run --rm flyway migrate
```

### 5.4 Convenciones Flyway

- Separador doble underscore: `V{timestamp}__{descripcion}.sql`
- Timestamp único (por eso incluye hora/min/seg)
- Sin rollback automático en Community — cambios destructivos requieren su migración inversa explícita
- Si editaste sin querer un archivo ya aplicado, `docker-compose run --rm flyway repair`

---

## 6. Cómo agregar un módulo nuevo (p.ej. `category`)

Seguir el patrón de `product/`. Pasos:

1. **Crear migración Flyway** para la tabla nueva (ver sección 5).
2. **Crear directorio `internal/category/`** con los 5 archivos:
   - `model.go` — struct GORM + Request DTOs con tags `binding:`
   - `repository.go` — `type Repository interface { ... }` + `type repository struct { db *gorm.DB }` + `func NewRepository(db *gorm.DB) Repository`
   - `service.go` — `type Service interface { ... }` + impl que recibe `Repository`
   - `handler.go` — recibe `Service`, usa `utils.SuccessResponse/ErrorResponse`, valida con `c.ShouldBindJSON`
   - `routes.go` — `func RegisterRoutes(r *gin.RouterGroup, h *Handler)` registra bajo `/api/v1/categories` con middlewares apropiados (`auth.AuthRequired()`, `auth.RequireRole("admin")`)
3. **Cablear en `cmd/main.go`:** instanciar repo → service → handler y llamar `category.RegisterRoutes(...)`.
4. **Agregar tests** en `service_test.go` siguiendo el estilo de `user` y `product`.
5. **Actualizar la collection Postman** (`postman_collection.json`) con los nuevos endpoints.
6. **Actualizar el README** si el endpoint es parte del contrato público.

---

## 7. Variables de entorno

| Variable | Default | Notas |
|---|---|---|
| `APP_PORT` | `8080` | Puerto donde escucha Gin |
| `APP_ENV` | `development` | En `production` baja el log level de GORM a `Error` |
| `DB_HOST` | `localhost` | Se sobrescribe a `host.docker.internal` en el compose del servicio `app` |
| `DB_PORT` | `3306` | |
| `DB_USER` | `root` | |
| `DB_PASSWORD` | — | Password del MySQL local del host |
| `DB_NAME` | `golang_api` | Debe existir en MySQL antes de correr Flyway |
| `JWT_SECRET` | — | Cambiar en producción; firma HS256 |
| `JWT_EXPIRATION_HOURS` | `24` | |

---

## 8. Troubleshooting

| Síntoma | Causa probable | Solución |
|---|---|---|
| `failed to connect to the docker API at npipe:...` | Docker Desktop no está corriendo | Abrir Docker Desktop y esperar ícono verde |
| Flyway: `Access denied for user 'root'@...` | Password del `.env` no coincide con el MySQL local | Editar `DB_PASSWORD` en `.env` y reintentar |
| Flyway: `Unknown database 'golang_api'` | No se creó la BD | `CREATE DATABASE golang_api;` en Workbench |
| App: `Error 1146: Table 'golang_api.users' doesn't exist` | No se corrieron las migraciones | `docker-compose run --rm flyway migrate` |
| `POST /auth/login` retorna 401 siempre | User no existe o password mal hasheada | Registrar primero con `/auth/register` |
| Endpoints de admin retornan 403 | JWT tiene `role=user` | `UPDATE users SET role='admin' WHERE ...` + volver a loguearse |
| `host.docker.internal` no resuelve | Linux sin `host-gateway` configurado | El compose ya incluye `extra_hosts: host.docker.internal:host-gateway` — asegurar Docker ≥ 20.10 |
| Warning `the attribute 'version' is obsolete` | Compose v2+ ignora `version:` | Cosmético, se puede borrar la línea del `docker-compose.yml` |

---

## 9. Qué cambió respecto del prompt original de construcción

Si encontrás referencias viejas en otros docs, sabé que:

- ❌ **Ya no se usa `GORM AutoMigrate`.** Todo el esquema vive en `migrations/*.sql` y se aplica con Flyway.
- ❌ **Ya no hay servicio `mysql` en `docker-compose.yml`.** MySQL es local del host.
- ❌ **Ya no se menciona Viper.** El único loader de env es `godotenv`.
- ✅ Se agregó el servicio `flyway` bajo profile `tools` en el compose.
- ✅ Se agregó `postman_collection.json`.
- ✅ El módulo Go se llama `github.com/fixia/golang-api`.

---

## 10. Checklist para un agente que toma el proyecto

Antes de tocar código, confirmar:

- [ ] Leí `README.md` y este archivo completos
- [ ] Entiendo que MySQL corre en el host, no en Docker
- [ ] Entiendo que las tablas se crean con Flyway, NUNCA con AutoMigrate
- [ ] Sé cómo generar un timestamp Flyway y crear una migración nueva
- [ ] Vi la collection `postman_collection.json` para entender los contratos
- [ ] Si voy a agregar un módulo, seguiré el patrón de `product/` (5 archivos, interfaces, capas)
- [ ] Si toco el esquema, lo hago en una migración nueva — nunca editando una ya aplicada
- [ ] Si cambio un endpoint, actualizo la collection Postman y el README

---

*Documento de handoff para FIXIA S.A.S. — Proyecto `golang-api` alineado con curso de Golang en Udemy. Última actualización: 2026-04-15.*
