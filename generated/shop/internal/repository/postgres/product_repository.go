package postgres

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/domain"
	"github.com/makoto-developer/go_microservice_example/generated/shop/internal/repository"
)

type productRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) repository.ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(ctx context.Context, product *domain.Product) error {
	query := `
		INSERT INTO products (id, shop_id, name, description, price, category, stock_quantity,
		                     weight, size, jan_code, published, deleted, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
	`
	_, err := r.db.ExecContext(ctx, query,
		product.ID, product.ShopID, product.Name, product.Description, product.Price,
		product.Category, product.StockQuantity, product.Weight, product.Size, product.JANCode,
		product.Published, product.Deleted, product.CreatedAt, product.UpdatedAt,
	)
	return err
}

func (r *productRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, shop_id, name, description, price, category, stock_quantity,
		       weight, size, jan_code, published, deleted, created_at, updated_at
		FROM products WHERE id = $1
	`
	product := &domain.Product{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&product.ID, &product.ShopID, &product.Name, &product.Description, &product.Price,
		&product.Category, &product.StockQuantity, &product.Weight, &product.Size, &product.JANCode,
		&product.Published, &product.Deleted, &product.CreatedAt, &product.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, domain.ErrProductNotFound
	}
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (r *productRepository) Update(ctx context.Context, product *domain.Product) error {
	query := `
		UPDATE products SET name = $2, description = $3, price = $4, category = $5, stock_quantity = $6,
		                   weight = $7, size = $8, jan_code = $9, updated_at = $10
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query,
		product.ID, product.Name, product.Description, product.Price, product.Category,
		product.StockQuantity, product.Weight, product.Size, product.JANCode, product.UpdatedAt,
	)
	return err
}

func (r *productRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE products SET deleted = true, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *productRepository) UpdatePublished(ctx context.Context, productID uuid.UUID, published bool) error {
	query := `UPDATE products SET published = $2, updated_at = NOW() WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, productID, published)
	return err
}

func (r *productRepository) List(ctx context.Context, shopID uuid.UUID, includeDeleted bool) ([]*domain.Product, error) {
	query := `
		SELECT id, shop_id, name, description, price, category, stock_quantity,
		       weight, size, jan_code, published, deleted, created_at, updated_at
		FROM products WHERE shop_id = $1
	`
	if !includeDeleted {
		query += ` AND deleted = false`
	}
	query += ` ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, shopID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*domain.Product
	for rows.Next() {
		p := &domain.Product{}
		if err := rows.Scan(
			&p.ID, &p.ShopID, &p.Name, &p.Description, &p.Price, &p.Category, &p.StockQuantity,
			&p.Weight, &p.Size, &p.JANCode, &p.Published, &p.Deleted, &p.CreatedAt, &p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, rows.Err()
}

func (r *productRepository) AddImage(ctx context.Context, image *domain.ProductImage) error {
	query := `
		INSERT INTO product_images (id, product_id, url, display_order, thumbnail_200_url,
		                           thumbnail_400_url, thumbnail_800_url, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(ctx, query,
		image.ID, image.ProductID, image.URL, image.DisplayOrder,
		image.Thumbnail200URL, image.Thumbnail400URL, image.Thumbnail800URL, image.CreatedAt,
	)
	return err
}

func (r *productRepository) GetImages(ctx context.Context, productID uuid.UUID) ([]*domain.ProductImage, error) {
	query := `
		SELECT id, product_id, url, display_order, thumbnail_200_url, thumbnail_400_url,
		       thumbnail_800_url, created_at
		FROM product_images WHERE product_id = $1 ORDER BY display_order ASC
	`
	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []*domain.ProductImage
	for rows.Next() {
		img := &domain.ProductImage{}
		if err := rows.Scan(
			&img.ID, &img.ProductID, &img.URL, &img.DisplayOrder,
			&img.Thumbnail200URL, &img.Thumbnail400URL, &img.Thumbnail800URL, &img.CreatedAt,
		); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, rows.Err()
}

func (r *productRepository) DeleteImage(ctx context.Context, imageID uuid.UUID) error {
	query := `DELETE FROM product_images WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, imageID)
	return err
}

func (r *productRepository) CountImages(ctx context.Context, productID uuid.UUID) (int, error) {
	query := `SELECT COUNT(*) FROM product_images WHERE product_id = $1`
	var count int
	err := r.db.QueryRowContext(ctx, query, productID).Scan(&count)
	return count, err
}

func (r *productRepository) AddTag(ctx context.Context, tag *domain.ProductTag) error {
	query := `INSERT INTO product_tags (id, product_id, tag_name, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(ctx, query, tag.ID, tag.ProductID, tag.TagName, tag.CreatedAt)
	return err
}

func (r *productRepository) GetTags(ctx context.Context, productID uuid.UUID) ([]*domain.ProductTag, error) {
	query := `SELECT id, product_id, tag_name, created_at FROM product_tags WHERE product_id = $1`
	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []*domain.ProductTag
	for rows.Next() {
		tag := &domain.ProductTag{}
		if err := rows.Scan(&tag.ID, &tag.ProductID, &tag.TagName, &tag.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func (r *productRepository) DeleteTags(ctx context.Context, productID uuid.UUID) error {
	query := `DELETE FROM product_tags WHERE product_id = $1`
	_, err := r.db.ExecContext(ctx, query, productID)
	return err
}

func (r *productRepository) CreateVariation(ctx context.Context, variation *domain.ProductVariation) error {
	query := `
		INSERT INTO product_variations (id, product_id, sku, attribute_name, attribute_value,
		                               price, stock_quantity, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(ctx, query,
		variation.ID, variation.ProductID, variation.SKU, variation.AttributeName,
		variation.AttributeValue, variation.Price, variation.StockQuantity,
		variation.CreatedAt, variation.UpdatedAt,
	)
	return err
}

func (r *productRepository) GetVariations(ctx context.Context, productID uuid.UUID) ([]*domain.ProductVariation, error) {
	query := `
		SELECT id, product_id, sku, attribute_name, attribute_value, price, stock_quantity,
		       created_at, updated_at
		FROM product_variations WHERE product_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var variations []*domain.ProductVariation
	for rows.Next() {
		v := &domain.ProductVariation{}
		if err := rows.Scan(
			&v.ID, &v.ProductID, &v.SKU, &v.AttributeName, &v.AttributeValue,
			&v.Price, &v.StockQuantity, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, err
		}
		variations = append(variations, v)
	}
	return variations, rows.Err()
}

func (r *productRepository) UpdateVariation(ctx context.Context, variation *domain.ProductVariation) error {
	query := `
		UPDATE product_variations SET price = $2, stock_quantity = $3, updated_at = $4
		WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, variation.ID, variation.Price, variation.StockQuantity, variation.UpdatedAt)
	return err
}

func (r *productRepository) DeleteVariations(ctx context.Context, productID uuid.UUID) error {
	query := `DELETE FROM product_variations WHERE product_id = $1`
	_, err := r.db.ExecContext(ctx, query, productID)
	return err
}
