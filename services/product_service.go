package services

import (
	"mall-seckill/datamodels"
	"mall-seckill/repositories"
)

type IProductService interface {
	GetProductById(int64) (*datamodels.Product, error)
	GetAllProduct() ([]*datamodels.Product, error)
	DeleteProductByID(int64) bool
	InsertProduct(product *datamodels.Product) (int64, error)
	UpdateProduct(product *datamodels.Product) error
	SubNumberOne(productID int64) error
}

type ProductService struct {
	productRepository repositories.IProduct
}

func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{repository}
}

func (p *ProductService) GetProductById(productId int64) (*datamodels.Product, error) {
	return p.productRepository.SelectByKey(productId)
}

func (p *ProductService) GetAllProduct() ([]*datamodels.Product, error) {
	return p.productRepository.SelectAll()
}

func (p *ProductService) DeleteProductByID(productId int64) bool {
	return p.productRepository.Delete(productId)
}

func (p *ProductService) InsertProduct(product *datamodels.Product) (int64, error) {
	return p.productRepository.Insert(product)
}

func (p *ProductService) UpdateProduct(product *datamodels.Product) error {
	return p.productRepository.Update(product)
}

func (p *ProductService) SubNumberOne(productId int64) error {
	return p.productRepository.SubProductNum(productId)
}
