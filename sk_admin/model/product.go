package model

import (
	"SecKill/sk_admin/config"
	"log"
)

type Product struct {
	ProductId   int    `json:"product_id"`   //商品Id
	ProductName string `json:"product_name"` //商品名称
	Total       int    `json:"total"`        //商品数量
	Status      int    `json:"status"`       //商品状态
}

type ProductModel struct {
}

func NewProductModel() *ProductModel {
	return &ProductModel{}
}

func (p *ProductModel) getTableName() string {
	return "product"
}

func (p *ProductModel) GetProductList() ([]map[string]interface{}, error) {
	conn := config.SecAdminConfCtx.DbConf.DbConn.GetInstance()
	list, err := conn.Table(p.getTableName()).Get()
	if err != nil {
		log.Printf("Error : %v", err)
		return nil, err
	}
	return list, nil
}

func (p *ProductModel) CreateProduct(product *Product) error {
	conn := config.SecAdminConfCtx.DbConf.DbConn.GetInstance()
	_, err := conn.Table(p.getTableName()).Data(map[string]interface{}{
		"product_name": product.ProductName,
		"total":        product.Total,
		"status":       product.Status,
	}).Insert()
	if err != nil {
		log.Printf("Error : %v", err)
		return err
	}
	return nil
}
