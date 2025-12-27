package dependencies

import (
	"github.com/DanielChachagua/ecommerce-noagestion-protos/pb"
	"github.com/SaltaGet/ecommerce-fiber-ms/cmd/server/controllers"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/repositories"
	"github.com/SaltaGet/ecommerce-fiber-ms/internal/services"
	"google.golang.org/grpc"
)

type ContainerGrpc struct {
	Controllers struct {
		TenantController *controllers.TenantController
		ProductController *controllers.ProductController
		CategoryController *controllers.CategoryController
	}
	Services struct {
		TenantService *services.TenantService
		ProductService *services.ProductService
		CategoryService *services.CategoryService
	}
	Repositories struct {
		TenantClient *repositories.TenantRepository
		ProductClient *repositories.ProductRepository
		CategoryClient *repositories.CategoryRepository
	}
}

func NewContainerGrpc(conn *grpc.ClientConn) *ContainerGrpc {
	c := &ContainerGrpc{}
	
	// Repositorios
	repoTenant := pb.NewTenantServiceClient(conn)
	repoProduct := pb.NewProductServiceClient(conn)
	repoCategory := pb.NewCategoryServiceClient(conn)

	c.Repositories.TenantClient = &repositories.TenantRepository{
		Client: repoTenant,
	}
	c.Repositories.ProductClient = &repositories.ProductRepository{
		Client: repoProduct,
	}
	c.Repositories.CategoryClient = &repositories.CategoryRepository{
		Client: repoCategory,
	}

	c.Services.TenantService = &services.TenantService{
		Repo: c.Repositories.TenantClient,
	}
	c.Services.ProductService = &services.ProductService{
		Repo: c.Repositories.ProductClient,
	}
	c.Services.CategoryService = &services.CategoryService{
		Repo: c.Repositories.CategoryClient,
	}

	c.Controllers.TenantController = &controllers.TenantController{
		TenantService: c.Services.TenantService,
	}
	c.Controllers.ProductController = &controllers.ProductController{
		ProductService: c.Services.ProductService,
	}
	c.Controllers.CategoryController = &controllers.CategoryController{
		CategoryService: c.Services.CategoryService,
	}

	return c
}