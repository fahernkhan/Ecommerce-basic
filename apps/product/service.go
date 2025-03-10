package product

import (
	"Ecommerce-basic/infra/response"
	"Ecommerce-basic/internal/log"
	"context"
	"time"
)

type Repository interface {
	CreateProduct(ctx context.Context, model Product) (err error)
	GetAllProductsWithPaginationCursor(ctx context.Context, model ProductPagination) (products []Product, err error)
	GetProductBySKU(ctx context.Context, sku string) (product Product, err error)
	GetProductByID(ctx context.Context, id int) (product Product, err error)                                                                            // Method baru
	UpdateProduct(ctx context.Context, model Product) (err error)                                                                                       // Method baru
	SoftDeleteProduct(ctx context.Context, id int) (err error)                                                                                          // Method baru
	SearchProducts(ctx context.Context, keyword string, pagination ProductPagination) (products []Product, err error)                                   // Method baru
	FilterProducts(ctx context.Context, minPrice, maxPrice int, minStock, maxStock int16, pagination ProductPagination) (products []Product, err error) // Method baru
	GetProductByName(ctx context.Context, name string) (product Product, err error)
}

type service struct {
	repo Repository
}

func newService(repo Repository) service {
	return service{
		repo: repo,
	}
}

func (s service) CreateProduct(ctx context.Context, req CreateProductRequestPayload) (err error) {
	productEntity := NewProductFromCreateProductRequest(req)

	if err = productEntity.Validate(); err != nil {
		log.Log.Errorf(ctx, "[CreateProduct, Validate] with error detail %v", err.Error())
		return
	}

	// Validasi keunikan nama produk
	existingProduct, err := s.repo.GetProductByName(ctx, productEntity.Name)
	if err != nil && err != response.ErrNotFound {
		return
	}
	if existingProduct.Id != 0 {
		return response.ErrProductAlreadyExists
	}

	if err = s.repo.CreateProduct(ctx, productEntity); err != nil {
		return
	}

	return
}

//func (s service) CreateProduct(ctx context.Context, req CreateProductRequestPayload) (err error) {
//	productEntity := NewProductFromCreateProductRequest(req)
//
//	if err = productEntity.Validate(); err != nil {
//		log.Log.Errorf(ctx, "[CreateProduct, Validate] with error detail %v", err.Error())
//		return
//	}
//
//	if err = s.repo.CreateProduct(ctx, productEntity); err != nil {
//		return
//	}
//
//	return
//}

func (s service) ListProducts(ctx context.Context, req ListProductRequestPayload) (products []Product, err error) {
	pagination := NewProductPaginationFromListProductRequest(req)
	log.Log.Infof(ctx, "Fetching products with pagination: %+v", pagination)

	products, err = s.repo.GetAllProductsWithPaginationCursor(ctx, pagination)
	if err != nil {
		log.Log.Errorf(ctx, "Failed to fetch products: %v", err)
		if err == response.ErrNotFound {
			return []Product{}, nil
		}
		return
	}

	log.Log.Infof(ctx, "Fetched %d products", len(products))
	return
}

func (s service) ProductDetail(ctx context.Context, sku string) (model Product, err error) {
	model, err = s.repo.GetProductBySKU(ctx, sku)
	if err != nil {
		if err == response.ErrNotFound {
			return Product{}, response.ErrNotFound
		}
		return
	}
	return
}

func (s service) UpdateProduct(ctx context.Context, id int, req UpdateProductRequestPayload) (err error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return
	}

	product.Name = req.Name
	product.Stock = req.Stock
	product.Price = req.Price
	product.UpdatedAt = time.Now()

	if err = product.Validate(); err != nil {
		return
	}

	// Validasi keunikan di service
	existingProduct, err := s.repo.GetProductByName(ctx, product.Name)
	if err != nil {
		if err == response.ErrNotFound {
			//product not found, so it is valid to update the current product.
		} else {
			return err //return other errors
		}
	} else {
		if existingProduct.Id != 0 && existingProduct.Id != product.Id {
			return response.ErrProductAlreadyExists
		}
	}

	return s.repo.UpdateProduct(ctx, product)
}

//func (s service) UpdateProduct(ctx context.Context, id int, req UpdateProductRequestPayload) (err error) {
//	product, err := s.repo.GetProductByID(ctx, id)
//	if err != nil {
//		return
//	}
//
//	product.Name = req.Name
//	product.Stock = req.Stock
//	product.Price = req.Price
//	product.UpdatedAt = time.Now()
//
//	if err = product.Validate(); err != nil {
//		return
//	}
//
//	if err = product.ValidateUnique(ctx, s.repo); err != nil {
//		return
//	}
//
//	return s.repo.UpdateProduct(ctx, product)
//}

func (s service) DeleteProduct(ctx context.Context, id int) (err error) {
	return s.repo.SoftDeleteProduct(ctx, id)
}

func (s service) SearchProducts(ctx context.Context, keyword string, pagination ProductPagination) (products []Product, err error) {
	return s.repo.SearchProducts(ctx, keyword, pagination)
}

func (s service) FilterProducts(ctx context.Context, minPrice, maxPrice int, minStock, maxStock int16, pagination ProductPagination) (products []Product, err error) {
	return s.repo.FilterProducts(ctx, minPrice, maxPrice, minStock, maxStock, pagination)
}
