# Logistics API

![Go Version](https://img.shields.io/badge/Go-1.21+-blue)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![API Status](https://img.shields.io/badge/API-operational-green)

API REST para gestión de órdenes de logística construida con **Go** y **arquitectura hexagonal**, diseñada para ser escalable, mantenible y lista para producción.

## 🚀 Quick Start

### 1. Instalar y ejecutar

```bash
git clone <repository-url>
cd logistics-api
cp .env.example .env  # Configurar variables
go mod download
make run              # http://localhost:8080
```

### 2. **Probar con Postman (RECOMENDADO)**

El método más eficiente para probar la API:

1. **Importar colección**: Abrir Postman → Import → `docs/logistics_api_postman_collection.json`
2. **Variables automáticas**: La colección ya incluye todas las variables necesarias
3. **Ejecutar en orden**:
   - `Health Checks` → Verificar que la API funciona
   - `Autenticación` → Registrar cliente y admin
   - `Órdenes - Cliente` → Crear y listar órdenes
   - `Órdenes - Administrador` → Gestionar estados
   - `Casos de Error` → Validar manejo de errores

**La colección incluye scripts automáticos que:**

- ✅ Guardan tokens JWT automáticamente
- ✅ Validan respuestas y tiempos
- ✅ Muestran logs útiles en la consola
- ✅ Prueban 25+ escenarios diferentes

## 📋 Características Principales

- **Arquitectura Hexagonal** (Ports & Adapters) para máxima testabilidad
- **Autenticación JWT** con roles diferenciados (Cliente/Administrador)
- **Validación de reglas de negocio** (pesos, coordenadas, transiciones de estado)
- **Rate Limiting** (10 req/sec, burst 20) para prevenir abuso
- **Logging estructurado** con Logrus para observabilidad
- **Health Checks** para monitoreo de infraestructura
- **Graceful Shutdown** para deployments sin downtime

## 🏛️ Reglas de Negocio

### Clasificación por Peso

- **S**: ≤ 5kg
- **M**: ≤ 15kg
- **L**: ≤ 25kg
- **>25kg**: Requiere convenio especial

### Estados de Órdenes

```
creado → recolectado → en_estacion → en_ruta → entregado
   ↓         ↓            ↓           ↓
cancelado  cancelado   cancelado   cancelado
```

### Control de Acceso

- **Clientes**: Crear órdenes + ver las propias
- **Administradores**: Acceso completo + cambiar estados

## 🔌 API Endpoints

| Método | Endpoint                    | Descripción                       | Auth |
| ------ | --------------------------- | --------------------------------- | ---- |
| `GET`  | `/health`                   | Health check completo             | No   |
| `POST` | `/api/v1/auth/register`     | Registro de usuario               | No   |
| `POST` | `/api/v1/auth/login`        | Inicio de sesión                  | No   |
| `POST` | `/api/v1/orders/`           | Crear orden                       | JWT  |
| `GET`  | `/api/v1/orders/`           | Listar órdenes (filtrado por rol) | JWT  |
| `PUT`  | `/api/v1/orders/:id/status` | Actualizar estado (solo admin)    | JWT  |

## 📝 Ejemplos Rápidos

### Registro y Login

```bash
# Registrar cliente
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email": "cliente@test.com", "password": "password123", "role": "client"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email": "cliente@test.com", "password": "password123"}'
```

### Crear Orden Simple

```bash
curl -X POST http://localhost:8080/api/v1/orders/ \
  -H "Authorization: Bearer <JWT_TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "origin_coordinates": {"latitude": 19.4326, "longitude": -99.1332},
    "destination_coordinates": {"latitude": 19.4969, "longitude": -99.1276},
    "origin_address": {
      "street": "Av. Reforma", "zipcode": "06600", "ext_num": "222",
      "city": "CDMX", "state": "CDMX", "country": "México"
    },
    "destination_address": {
      "street": "Av. Insurgentes", "zipcode": "07700", "ext_num": "1500",
      "city": "CDMX", "state": "CDMX", "country": "México"
    },
    "product_quantity": 2,
    "total_weight": 5.5
  }'
```

## 🏗️ Stack Tecnológico

- **Backend**: Go 1.21+ con Gin Framework
- **Base de Datos**: PostgreSQL (Neon)
- **Autenticación**: JWT con bcrypt
- **ORM**: GORM con Auto-migrations
- **Logging**: Logrus (JSON structured)
- **Validación**: go-playground/validator

## 🔧 Configuración

### Variables de Entorno Requeridas

```bash
DATABASE_URL=postgresql://user:pass@host:port/db?sslmode=require
JWT_SECRET=your-super-secret-jwt-key-32-chars-min
JWT_EXPIRY_HOUR=24
SERVER_HOST=localhost
SERVER_PORT=8080
LOG_LEVEL=info
LOG_FORMAT=json
```

### Comandos Disponibles

```bash
make run          # Ejecutar aplicación
make build        # Construir binario
make test         # Ejecutar tests
make clean        # Limpiar artifacts
make gen-jwt      # Generar JWT secret
make db-status    # Verificar BD
```

## 🏛️ Arquitectura Hexagonal

### Estructura del Proyecto

```
logistics-api/
├── cmd/server/                 # Entry point
├── internal/
│   ├── core/                   # Domain Layer
│   │   ├── domain/            # Entities & Value Objects
│   │   ├── ports/             # Interfaces
│   │   └── usecases/          # Application Logic
│   ├── adapters/              # Infrastructure Layer
│   │   ├── primary/           # HTTP, CLI
│   │   └── secondary/         # DB, External APIs
│   ├── config/                # Configuration
│   └── pkg/                   # Shared utilities
```

### ¿Por qué Arquitectura Hexagonal?

- **Testabilidad**: Fácil creación de mocks y tests
- **Flexibilidad**: Cambio de infraestructura sin afectar negocio
- **Mantenibilidad**: Separación clara de responsabilidades
- **Escalabilidad**: Preparado para evolución a microservicios

## 🔒 Seguridad Implementada

- ✅ **JWT Authentication** con tokens firmados
- ✅ **Password hashing** con bcrypt (cost 12)
- ✅ **Rate limiting** (10 req/sec, burst 20)
- ✅ **CORS** configurado apropiadamente
- ✅ **Input validation** en todos los endpoints
- ✅ **SQL injection prevention** via GORM

## 📊 Monitoreo

### Health Checks

- `/health` - Estado completo del sistema
- `/health/live` - Liveness probe (aplicación responde)
- `/health/ready` - Readiness probe (dependencias OK)

### Logging Estructurado

- Formato JSON para agregación
- Contexto de request (IP, User-Agent, latencia)
- Correlation IDs para trazabilidad

## 🧪 Testing

### Con Postman (Recomendado)

1. Importar `docs/logistics_api_postman_collection.json`
2. Ejecutar carpetas en orden secuencial
3. Verificar logs en consola de Postman

### Con Go Tests

```bash
go test ./...                    # Todos los tests
go test -cover ./...            # Con coverage
go test ./internal/core/...     # Solo domain tests
```

## 🚀 Deploy

### Variables Críticas

```bash
DATABASE_URL=postgresql://...      # PostgreSQL connection
JWT_SECRET=<32-char-minimum>       # JWT signing key
SERVER_PORT=8080                   # Server port
LOG_LEVEL=info                     # Logging level
```

### Opciones Recomendadas

- **Base de Datos**: Neon.tech (PostgreSQL serverless, gratuito)
- **Hosting**: Railway, Render, o similar
- **CI/CD**: GitHub Actions

## 🐛 Troubleshooting

**❌ Error de conexión a BD**

```bash
echo $DATABASE_URL  # Verificar variable
go run cmd/server/main.go  # Ver logs detallados
```

**❌ "Authentication required"**

```bash
# Verificar formato del header
Authorization: Bearer <token>
# Obtener nuevo token via /auth/login
```

**❌ "Rate limit exceeded"**

```bash
# Esperar 1 minuto o cambiar IP
# Límite: 10 req/sec con burst de 20
```

**❌ Coordenadas inválidas**

- Latitude: [-90, 90]
- Longitude: [-180, 180]

## 📚 Documentación Adicional

- 📁 [Colección Postman](./docs/logistics_api_postman_collection.json) - **25+ requests con scripts automáticos**
- 📄 [Casos de Uso Detallados](./docs/use_cases.md)
- 🏗️ [Diagrama de Arquitectura](./docs/architecture.md)
- 🤝 [Guía de Contribución](./docs/contributing.md)

## 🎯 Flujo de Testing Recomendado

### Método 1: Postman (Más Eficiente)

1. Importar colección
2. Ejecutar `Health Checks`
3. Ejecutar `Autenticación` (auto-guarda tokens)
4. Ejecutar `Órdenes - Cliente`
5. Ejecutar `Órdenes - Administrador`
6. Revisar `Casos de Error`

### Método 2: curl Manual

Ver ejemplos en las secciones anteriores o usar la colección de Postman como referencia.

---

**💡 Tip**: Para una experiencia óptima, usa la **colección de Postman** que incluye scripts automáticos, validaciones y casos de prueba completos.

**📧 Contacto**: [tu-email@example.com] | **🐙 GitHub**: [tu-usuario]
