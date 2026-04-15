# Golang API — E-commerce Backend

API REST en **Go + MySQL** para gestión de usuarios, productos y órdenes, con autenticación JWT, arquitectura en capas limpias y dockerización completa.

> Proyecto alineado con el curso de Go en Udemy (FIXIA S.A.S.).

---

## 🛠️ Stack

- **Go 1.22+**
- **Gin** — framework HTTP
- **GORM** — ORM
- **MySQL 8.0**
- **JWT v5** + **bcrypt**
- **godotenv** — variables de entorno
- **Docker** + **docker-compose**
- **testify** — testing

---

## 📋 Requisitos previos

- Go 1.22 o superior
- Docker y Docker Compose (recomendado)
- MySQL 8.0 (si se corre localmente sin Docker)

---

## ⚙️ Configuración

### 1. Crear la base de datos en tu MySQL local

Este proyecto **usa tu instancia MySQL local** (no levanta MySQL en Docker). Creá la BD vacía una sola vez:

```sql
CREATE DATABASE golang_api;
```

Las tablas se crean mediante **migraciones Flyway** (ver sección siguiente).

### 2. Configurar variables de entorno

```bash
cp .env.example .env
```

Editar `.env` y poner **tu password de MySQL local**:

```env
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=tu_password_mysql_local
DB_NAME=golang_api
```

---

## 🗃️ Migraciones con Flyway

El esquema de la BD está versionado en `migrations/` con la convención Flyway **timestamp** (`V{YYYYMMDDHHmmss}__descripcion.sql`):

```
migrations/
├── V20260415100000__create_users_table.sql
├── V20260415100100__create_products_table.sql
├── V20260415100200__create_orders_table.sql
└── V20260415100300__create_order_items_table.sql
```

Cada archivo incluye un header con fecha y autor:

```sql
-- =============================================================================
-- Migration: create users table
-- Version:   V20260415100000
-- Date:      2026-04-15 10:00:00
-- Author:    autor@ejemplo.com
-- =============================================================================

CREATE TABLE users (...);
```

**Ventajas del formato timestamp:**
- La fecha del cambio queda en el nombre del archivo (visible en `git log`, `ls`, `flyway info`).
- Ordenamiento cronológico natural (Flyway ordena alfabéticamente → coincide con el tiempo).
- Dos devs pueden crear migraciones en paralelo sin conflictos de merge (cada timestamp es único).

Flyway además mantiene la tabla `flyway_schema_history` en la BD con el tracking de qué versiones se aplicaron, cuándo y con qué checksum.

### Correr migraciones (vía Docker, recomendado)

No necesitás instalar Flyway localmente — el compose lo provee como servicio bajo demanda:

```bash
# Aplicar migraciones pendientes
docker-compose run --rm flyway migrate

# Ver estado actual (qué se aplicó, qué falta)
docker-compose run --rm flyway info

# Validar que archivos SQL coincidan con los checksums en BD
docker-compose run --rm flyway validate

# Reparar history (útil si editaste un archivo ya aplicado)
docker-compose run --rm flyway repair
```

### Crear una nueva migración

#### 1. Generar el timestamp actual

**Git Bash / WSL / Linux / Mac:**
```bash
date +V%Y%m%d%H%M%S
# salida: V20260415143027
```

**PowerShell (Windows):**
```powershell
Get-Date -Format "\VyyyyMMddHHmmss"
# salida: V20260415143027
```

**CMD (Windows):**
```cmd
powershell -Command "Get-Date -Format VyyyyMMddHHmmss"
```

#### 2. Crear el archivo con ese prefijo

```
migrations/V20260415143027__add_phone_to_users.sql
```

#### 3. Contenido con header estándar

```sql
-- =============================================================================
-- Migration: add phone column to users
-- Version:   V20260415143027
-- Date:      2026-04-15 14:30:27
-- Author:    tu@email.com
-- =============================================================================

ALTER TABLE users ADD COLUMN phone VARCHAR(30);
```

#### 4. Aplicar

```bash
docker-compose run --rm flyway migrate
```

### Reglas importantes de Flyway

- **Nunca edites un archivo ya aplicado**. Creá uno nuevo con timestamp posterior.
- El timestamp debe ser **único** — por eso incluye hora/minuto/segundo.
- El separador es **doble underscore**: `V{timestamp}__{descripcion}.sql`.
- Flyway Community no tiene rollback automático — los cambios destructivos requieren su propia migración inversa.

### Alternativa: Flyway local

Si preferís no usar Docker para Flyway, instalalo nativo:

```bash
# Windows (scoop)
scoop install flyway

# Mac
brew install flyway
```

Y correlo apuntando a la BD local:

```bash
flyway -url=jdbc:mysql://localhost:3306/golang_api \
       -user=root -password=tu_password \
       -locations=filesystem:./migrations \
       migrate
```

---

## 🚀 Ejecución de la app

**⚠️ Antes de arrancar la app por primera vez, correr las migraciones** (`flyway migrate`). La app ya no crea tablas automáticamente.

### Local con `go run` (recomendado para desarrollo)

```bash
go mod tidy
go run ./cmd/main.go
```

### Con Docker (la app en contenedor)

```bash
docker-compose up --build app
```

API disponible en `http://localhost:8080`.

### Build del binario

```bash
go build -o server ./cmd/main.go
./server
```

---

## 🧪 Tests

```bash
go test ./...
```

---

## 📁 Estructura del proyecto

```
golang-api/
├── cmd/main.go                 # Entry point
├── config/config.go            # Carga de variables de entorno
├── database/db.go              # Conexión MySQL + AutoMigrate
├── internal/
│   ├── auth/                   # JWT + middleware de auth/roles
│   ├── user/                   # Módulo de usuarios (CRUD + register/login)
│   ├── product/                # Módulo de productos (CRUD)
│   └── order/                  # Módulo de órdenes (con transacciones)
├── middleware/logger.go        # Logger HTTP
├── utils/                      # Respuestas y errores estandarizados
├── .env.example
├── Dockerfile
├── docker-compose.yml
└── go.mod
```

Cada módulo en `internal/` sigue el patrón: `model → repository → service → handler → routes`.

---

## 🌐 Endpoints

Todas las rutas están bajo `/api/v1`.

### Auth

| Método | Ruta | Auth |
|--------|------|------|
| POST | `/auth/register` | ❌ |
| POST | `/auth/login` | ❌ |
| GET | `/auth/me` | ✅ |

**POST `/auth/register`**
```json
{
  "first_name": "Ada",
  "last_name": "Lovelace",
  "email": "ada@example.com",
  "password": "supersecret"
}
```

**POST `/auth/login`** → retorna `{ "token": "<jwt>", "user": {...} }`.

Incluir el token en requests protegidos:
```
Authorization: Bearer <jwt>
```

### Users (admin)

| Método | Ruta | Descripción |
|--------|------|-------------|
| GET | `/users` | Listar usuarios |
| GET | `/users/:id` | Obtener por ID |
| PUT | `/users/:id` | Actualizar |
| DELETE | `/users/:id` | Soft-delete |

### Products

| Método | Ruta | Auth |
|--------|------|------|
| GET | `/products` | ❌ |
| GET | `/products/:id` | ❌ |
| POST | `/products` | ✅ admin |
| PUT | `/products/:id` | ✅ admin |
| DELETE | `/products/:id` | ✅ admin |

**POST `/products`**
```json
{
  "name": "Teclado mecánico",
  "description": "Switches rojos",
  "price": 120.50,
  "stock": 30,
  "category": "periféricos"
}
```

### Orders

| Método | Ruta | Auth |
|--------|------|------|
| POST | `/orders` | ✅ |
| GET | `/orders` | ✅ (propias) |
| GET | `/orders/:id` | ✅ (propia o admin) |
| PUT | `/orders/:id/status` | ✅ admin |

**POST `/orders`**
```json
{
  "items": [
    { "product_id": 1, "quantity": 2 },
    { "product_id": 3, "quantity": 1 }
  ]
}
```

La creación se hace en una **transacción GORM**: valida stock, descuenta inventario, calcula total y graba orden + items atómicamente.

---

## 🔐 Formato de respuesta estandarizado

```json
{
  "success": true,
  "message": "operation description",
  "data": { ... },
  "error": ""
}
```

En caso de error:
```json
{ "success": false, "error": "descripción del error" }
```

---

## 🧪 Crear un admin para testear

Tras registrar un usuario, promoverlo manualmente en MySQL:

```sql
UPDATE users SET role = 'admin' WHERE email = 'tu@email.com';
```

---

## 📄 Licencia

Proyecto educativo — FIXIA S.A.S.
