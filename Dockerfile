# ==========================
# STAGE 1: COMPILACIÓN
# ==========================
FROM golang:1.24-alpine AS builder

# build-base incluye gcc y librerías de C necesarias para compilar WebP
RUN apk add --no-cache build-base

WORKDIR /app

# Cacheamos módulos
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# OPTIMIZACIÓN CLAVE: 
# 1. Usamos --mount para que las librerías de C ya compiladas se queden en caché.
# 2. Eliminamos UPX (que es lo que te quitaba 3 minutos).
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    CGO_ENABLED=1 go build -ldflags="-s -w -extldflags '-static'" -trimpath -o server ./cmd/server/main.go

# ==========================
# STAGE 2: IMAGEN FINAL
# ==========================
# Usamos Debian Distroless porque es compatible con binarios que usaron CGO
FROM gcr.io/distroless/static-debian12

WORKDIR /app

# El usuario no-root por defecto en distroless es 'nonroot'
COPY --from=builder /app/server .

# Crear carpeta media para las imágenes (Distroless es muy limitado, 
# si necesitas crear carpetas dinámicamente asegúrate que el binario tenga permisos)
# Nota: Como tu código usa os.MkdirAll, funcionará si el volumen está montado.

EXPOSE 3030
CMD ["./server"]

# # ==========================
# # STAGE 1: COMPILACIÓN
# # ==========================
# FROM golang:1.24-alpine AS builder

# # Necesitamos build-base para compilar los componentes de C de WebP
# RUN apk add --no-cache build-base upx

# WORKDIR /app

# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .

# # CAMBIO CLAVE: CGO_ENABLED=1 para que WebP pueda compilar
# RUN CGO_ENABLED=1 go build -ldflags="-s -w -extldflags '-static'" -trimpath -o server ./cmd/server/main.go

# RUN upx --best server || echo "Skipping UPX"

# # ==========================
# # STAGE 2: IMAGEN FINAL
# # ==========================
# FROM alpine:3.21

# RUN apk add --no-cache ca-certificates tzdata

# RUN addgroup -S appgroup && adduser -S appuser -G appgroup \
#     && mkdir -p /app/backups && chown -R appuser:appgroup /app

# WORKDIR /app

# # El binario ya es estático gracias a -extldflags '-static'
# COPY --from=builder --chown=appuser:appgroup /app/server .

# USER appuser
# EXPOSE 3030

# CMD ["./server"]