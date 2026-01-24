package jobs

import (
	"time"

	"github.com/SaltaGet/ecommerce-fiber-ms/internal/cache"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/dependencies"

	"github.com/rs/zerolog/log"
)

// Job de actualización
func ReloadTenants(tenantsStore *cache.TenantStore, deps *dependencies.ContainerGrpc) {
    ticker := time.NewTicker(30 * time.Minute)
    
    // Función local para recargar
    reload := func() {
        log.Info().Msg("🔄 Recargando caché de tenants desde gRPC...")
        // Llamas a tu servicio gRPC existente a través de deps
        tenants, err := deps.Services.TenantService.TenantList()
        if err != nil {
            log.Error().Err(err).Msg("Fallo al recargar tenants")
            return
        }
        
        tenantsStore.Update(tenants)
        log.Info().Int("count", len(tenants)).Msg("Caché de tenants actualizada")
    }

    // Carga inicial
    reload()

    // Bucle infinito para el ticker
    for range ticker.C {
        reload()
    }
}