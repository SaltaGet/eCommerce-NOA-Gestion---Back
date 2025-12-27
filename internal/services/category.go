package services

import "github.com/SaltaGet/ecommerce-fiber-ms/internal/schemas"

func (s *CategoryService) CategoryGetAll(tenantID string) ([]schemas.Category, error) {
  pbCategories, err := s.Repo.CategoryGetAll(tenantID)
  if err != nil {
    return nil, err
  }

  categories := make([]schemas.Category, 0, len(pbCategories))

  for i := range pbCategories {
    pbCategory := pbCategories[i] 

    category := schemas.Category{
      ID:   pbCategory.Id,
      Name: pbCategory.Name,
    }
    categories = append(categories, category)
  }

  return categories, nil
}