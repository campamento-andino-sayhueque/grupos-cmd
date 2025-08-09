# Grupos Command Service

API de comandos para el servicio de grupos, construida con Go, Gin, Google Cloud Firestore y Google Cloud Pub/Sub.

## 📋 Requisitos Previos

### Software Necesario

- **Go 1.24.3+**: [Descargar Go](https://golang.org/dl/)
- **Docker y Docker Compose**: [Descargar Docker](https://www.docker.com/get-started)
- **Google Cloud SDK** (para despliegue): [Instalar gcloud CLI](https://cloud.google.com/sdk/docs/install)
- **Git**: Para clonar el repositorio

### Dependencias del Proyecto

Las dependencias principales incluyen:
- `gin-gonic/gin`: Framework web
- `cloud.google.com/go/firestore`: Cliente de Firestore
- `cloud.google.com/go/pubsub`: Cliente de Pub/Sub
- `ThreeDotsLabs/watermill`: Biblioteca de mensajería
- `swaggo/gin-swagger`: Documentación API con Swagger

## 🚀 Instalación y Configuración

### 1. Clonar el Repositorio

```bash
git clone https://github.com/campamento-andino-sayhueque/grupos-cmd.git
cd grupos-cmd
```

### 2. Instalar Dependencias

```bash
go mod download
```

### 3. Configuración de Variables de Entorno

Crea un archivo `.env` en el directorio raíz con las siguientes variables:

```env
# Para desarrollo local con emuladores
GOOGLE_CLOUD_PROJECT_ID=test-project
FIRESTORE_EMULATOR_HOST=localhost:8081
PUB_SUB_EMULATOR_HOST=localhost:8085
PUB_SUB_TOPIC_ID=grupos
FIRESTORE_COLLECTION=eventos
PORT=8080

# Para producción (Google Cloud Run)
# GOOGLE_CLOUD_PROJECT_ID=tu-proyecto-gcp
# PUB_SUB_TOPIC_ID=grupos
# FIRESTORE_COLLECTION=eventos
# PORT=8080
```

## 🔧 Desarrollo Local

### Opción 1: Usando Docker Compose (Recomendado)

Este método levanta automáticamente los emuladores de Firestore y Pub/Sub:

```bash
# Levantar todos los servicios (app + emuladores)
docker-compose up

# Levantar en background
docker-compose up -d

# Ver logs
docker-compose logs -f

# Detener servicios
docker-compose down
```

### Opción 2: Desarrollo Manual

#### 1. Levantar Emuladores de Google Cloud

**Emulador de Firestore:**
```bash
gcloud beta emulators firestore start --host-port=localhost:8081
```

**Emulador de Pub/Sub:**
```bash
gcloud beta emulators pubsub start --project=test-project --host-port=localhost:8085
```

#### 2. Ejecutar la Aplicación

```bash
# Compilar y ejecutar
go run cmd/server/main.go

# O compilar primero y luego ejecutar
go build -o grupos-cmd cmd/server/main.go
./grupos-cmd
```

### 🌐 Acceso a la Aplicación

- **API**: http://localhost:8080
- **Documentación Swagger**: http://localhost:8080/swagger/index.html
- **Emulador Firestore**: http://localhost:8081
- **Emulador Pub/Sub**: http://localhost:8085
- **Firebase Emulator Suite UI**: http://localhost:4000
- **Pub/Sub UI Dedicada**: http://localhost:4200

## 🔍 Herramientas de Monitoreo

### Firebase Emulator Suite UI (Puerto 4000)
Interfaz oficial de Google que incluye:
- Vista general de todos los emuladores
- Navegación entre Firestore y Pub/Sub
- Datos básicos de tópicos y suscripciones

### Pub/Sub UI Dedicada (Puerto 4200)
Interfaz especializada para Pub/Sub que permite:
- ✅ Crear y gestionar tópicos
- ✅ Crear y gestionar suscripciones
- ✅ Enviar mensajes de prueba
- ✅ Visualizar mensajes en las colas
- ✅ Monitorear flujo de eventos en tiempo real
- ✅ Navegar historial de mensajes

Para usar la UI de Pub/Sub:
1. Accede a http://localhost:4200
2. Agrega el proyecto `test-project`
3. Explora tópicos, suscripciones y mensajes

## 🧪 Ejecución de Tests

### Tests Unitarios

```bash
# Ejecutar todos los tests
go test ./...

# Tests con verbose output
go test -v ./...

# Tests con coverage
go test -cover ./...

# Coverage detallado
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Tests de Integración

Los tests de integración requieren que los emuladores estén ejecutándose:

```bash
# Opción 1: Con Docker Compose
docker-compose up -d firestore pubsub

# Ejecutar tests de integración
go test -v ./cmd/server/

# Opción 2: Emuladores manuales
gcloud beta emulators firestore start --host-port=localhost:8081 &
gcloud beta emulators pubsub start --project=test-project --host-port=localhost:8085 &

# Configurar variables de entorno para tests
export FIRESTORE_EMULATOR_HOST=localhost:8081
export PUBSUB_EMULATOR_HOST=localhost:8085
export GOOGLE_CLOUD_PROJECT_ID=test-project

# Ejecutar tests
go test -v ./cmd/server/
```

### Tests Específicos

```bash
# Test específico por nombre
go test -v -run TestCreateGrupo_Integration ./cmd/server/

# Tests de un paquete específico
go test -v ./internal/handler/
```

## 🏗️ Compilación

### Compilación Local

```bash
# Compilación básica
go build -o grupos-cmd cmd/server/main.go

# Compilación optimizada para producción
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o grupos-cmd cmd/server/main.go
```

### Compilación con Docker

```bash
# Construir imagen Docker
docker build -t grupos-cmd .

# Ejecutar contenedor
docker run -p 8080:8080 grupos-cmd
```

## ☁️ Despliegue en Google Cloud Run

### Requisitos Previos para Despliegue

1. **Cuenta de Google Cloud** con facturación habilitada
2. **Proyecto de GCP** creado
3. **APIs habilitadas**:
   - Cloud Run API
   - Cloud Build API
   - Firestore API
   - Pub/Sub API

### Configuración Inicial

```bash
# Autenticarse con Google Cloud
gcloud auth login

# Configurar proyecto
gcloud config set project TU-PROYECTO-ID

# Habilitar APIs necesarias
gcloud services enable run.googleapis.com
gcloud services enable cloudbuild.googleapis.com
gcloud services enable firestore.googleapis.com
gcloud services enable pubsub.googleapis.com
```

### Crear Recursos en GCP

```bash
# Crear tópico de Pub/Sub
gcloud pubsub topics create grupos

# Inicializar Firestore (si no está creado)
gcloud firestore databases create --region=us-central1
```

### Despliegue

#### Opción 1: Usando Cloud Build

```bash
# Desplegar directamente desde el código fuente
gcloud run deploy grupos-cmd \
  --source . \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars="GOOGLE_CLOUD_PROJECT_ID=TU-PROYECTO-ID,PUB_SUB_TOPIC_ID=grupos,FIRESTORE_COLLECTION=eventos"
```

#### Opción 2: Usando Docker

```bash
# Construir y subir imagen a Container Registry
gcloud builds submit --tag gcr.io/TU-PROYECTO-ID/grupos-cmd

# Desplegar desde la imagen
gcloud run deploy grupos-cmd \
  --image gcr.io/TU-PROYECTO-ID/grupos-cmd \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars="GOOGLE_CLOUD_PROJECT_ID=TU-PROYECTO-ID,PUB_SUB_TOPIC_ID=grupos,FIRESTORE_COLLECTION=eventos"
```

#### Opción 3: Usando Artifact Registry (Recomendado)

```bash
# Crear repositorio en Artifact Registry
gcloud artifacts repositories create grupos-cmd-repo \
  --repository-format=docker \
  --location=us-central1

# Configurar Docker para usar Artifact Registry
gcloud auth configure-docker us-central1-docker.pkg.dev

# Construir y subir imagen
docker build -t us-central1-docker.pkg.dev/TU-PROYECTO-ID/grupos-cmd-repo/grupos-cmd .
docker push us-central1-docker.pkg.dev/TU-PROYECTO-ID/grupos-cmd-repo/grupos-cmd

# Desplegar
gcloud run deploy grupos-cmd \
  --image us-central1-docker.pkg.dev/TU-PROYECTO-ID/grupos-cmd-repo/grupos-cmd \
  --platform managed \
  --region us-central1 \
  --allow-unauthenticated \
  --set-env-vars="GOOGLE_CLOUD_PROJECT_ID=TU-PROYECTO-ID,PUB_SUB_TOPIC_ID=grupos,FIRESTORE_COLLECTION=eventos"
```

### Configuración de Variables de Entorno en Cloud Run

```bash
# Actualizar variables de entorno
gcloud run services update grupos-cmd \
  --region us-central1 \
  --set-env-vars="GOOGLE_CLOUD_PROJECT_ID=TU-PROYECTO-ID,PUB_SUB_TOPIC_ID=grupos,FIRESTORE_COLLECTION=eventos,PORT=8080"
```

### Configurar Permisos IAM

```bash
# Obtener la cuenta de servicio de Cloud Run
SERVICE_ACCOUNT=$(gcloud run services describe grupos-cmd --region=us-central1 --format="value(spec.template.spec.serviceAccountName)")

# Dar permisos de Firestore
gcloud projects add-iam-policy-binding TU-PROYECTO-ID \
  --member="serviceAccount:${SERVICE_ACCOUNT}" \
  --role="roles/datastore.user"

# Dar permisos de Pub/Sub
gcloud projects add-iam-policy-binding TU-PROYECTO-ID \
  --member="serviceAccount:${SERVICE_ACCOUNT}" \
  --role="roles/pubsub.publisher"
```

## 📊 Monitoreo y Logs

### Ver Logs de Cloud Run

```bash
# Logs en tiempo real
gcloud run services logs tail grupos-cmd --region=us-central1

# Logs recientes
gcloud run services logs read grupos-cmd --region=us-central1 --limit=50
```

### Métricas en Google Cloud Console

- Ve a [Cloud Run Console](https://console.cloud.google.com/run)
- Selecciona tu servicio `grupos-cmd`
- Revisa métricas de CPU, memoria, latencia y errores

## 🛠️ Comandos Útiles

```bash
# Generar documentación Swagger
swag init -g cmd/server/main.go

# Limpiar módulos Go
go mod tidy

# Verificar vulnerabilidades
go list -json -m all | nancy sleuth

# Formatear código
go fmt ./...

# Linting
golangci-lint run

# Ejecutar benchmarks
go test -bench=. ./...
```

## 📁 Estructura del Proyecto

```
grupos-cmd/
├── cmd/
│   └── server/          # Punto de entrada de la aplicación
├── internal/
│   ├── config/          # Configuración de la aplicación
│   ├── domain/          # Entidades de dominio
│   ├── event/           # Publisher de eventos
│   ├── handler/         # Handlers HTTP
│   └── repository/      # Repositorio de datos
├── docs/                # Documentación Swagger generada
├── docker-compose.yml   # Configuración para desarrollo local
├── Dockerfile          # Imagen Docker para producción
└── README.md           # Este archivo
```

## 🤝 Contribución

1. Fork el proyecto
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Commit tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 📄 Licencia

Este proyecto está bajo la licencia MIT. Ver el archivo `LICENSE` para más detalles.