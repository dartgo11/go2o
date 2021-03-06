/**
 * Copyright 2014 @ S1N1 Team.
 * name :
 * author : jarryliu
 * date : 2013-12-23 23:15
 * description :
 * history :
 */

package dps

import (
	"errors"
	"github.com/atnet/gof/web/ui/tree"
	"go2o/src/core/domain/interface/sale"
	"go2o/src/core/dto"
	"strconv"
)

type saleService struct {
	_rep sale.ISaleRep
}

func NewSaleService(r sale.ISaleRep) *saleService {
	return &saleService{
		_rep: r,
	}
}

func (this *saleService) GetValueGoods(partnerId, goodsId int) *sale.ValueGoods {
	sl := this._rep.GetSale(partnerId)
	pro := sl.GetGoods(goodsId)
	v := pro.GetValue()
	return &v
}

func (this *saleService) SaveGoods(partnerId int, v *sale.ValueGoods) (int, error) {
	sl := this._rep.GetSale(partnerId)
	var pro sale.IGoods
	if v.Id > 0 {
		pro = sl.GetGoods(v.Id)
		if pro == nil {
			return 0, errors.New("产品不存在")
		}
		if err := pro.SetValue(v); err != nil {
			return 0, err
		}
	} else {
		pro = sl.CreateGoods(v)
	}
	return pro.Save()
}

func (this *saleService) GetOnShelvesGoodsByCategoryId(partnerId, cid, num int) []*dto.ListGoods {
	var goods = this._rep.GetOnShelvesGoodsByCategoryId(partnerId, cid, num)
	var listGoods []*dto.ListGoods = make([]*dto.ListGoods, len(goods))
	for i, v := range goods {
		listGoods[i] = &dto.ListGoods{
			Id:         v.Id,
			Name:       v.Name,
			SmallTitle: v.SmallTitle,
			Image:      v.Image,
			Price:      v.Price,
			SalePrice:  v.SalePrice,
		}
	}
	return listGoods
}

func (this *saleService) DeleteGoods(partnerId, goodsId int) error {
	sl := this._rep.GetSale(partnerId)
	return sl.DeleteGoods(goodsId)
}

func (this *saleService) GetCategory(partnerId, id int) *sale.ValueCategory {
	sl := this._rep.GetSale(partnerId)
	c := sl.GetCategory(id)
	if c != nil {
		cv := c.GetValue()
		return &cv
	}
	return nil
}

func (this *saleService) DeleteCategory(partnerId, id int) error {
	sl := this._rep.GetSale(partnerId)
	return sl.DeleteCategory(id)
}

func (this *saleService) SaveCategory(partnerId int, v *sale.ValueCategory) (int, error) {
	sl := this._rep.GetSale(partnerId)
	var ca sale.ICategory
	if v.Id > 0 {
		ca = sl.GetCategory(v.Id)
		if err := ca.SetValue(v); err != nil {
			return 0, err
		}
	} else {
		ca = sl.CreateCategory(v)
	}

	return ca.Save()
}

func (this *saleService) GetCategoryTreeNode(partnerId int) *tree.TreeNode {
	sl := this._rep.GetSale(partnerId)
	cats := sl.GetCategories()
	rootNode := &tree.TreeNode{
		Text:   "根节点",
		Value:  "",
		Url:    "",
		Icon:   "",
		Open:   true,
		Childs: nil}
	this.iterCategoryTree(rootNode, 0, cats)
	return rootNode
}

func (this *saleService) iterCategoryTree(node *tree.TreeNode, parentId int, categories []sale.ICategory) {
	node.Childs = []*tree.TreeNode{}
	for _, v := range categories {
		cate := v.GetValue()
		if cate.ParentId == parentId {
			cNode := &tree.TreeNode{
				Text:   cate.Name,
				Value:  strconv.Itoa(cate.Id),
				Url:    "",
				Icon:   "",
				Open:   true,
				Childs: nil}
			node.Childs = append(node.Childs, cNode)
			this.iterCategoryTree(cNode, cate.Id, categories)
		}
	}
}

func (this *saleService) GetCategories(partnerId int) []*sale.ValueCategory {
	sl := this._rep.GetSale(partnerId)
	cats := sl.GetCategories()
	var list []*sale.ValueCategory = make([]*sale.ValueCategory, len(cats))
	for i, v := range cats {
		vv := v.GetValue()
		list[i] = &vv
	}
	return list
}
