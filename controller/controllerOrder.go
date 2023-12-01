package controller

import (
	"assignment-2/config"
	"assignment-2/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Create(c echo.Context) error {
	db := config.GetDB()

	var order models.Order
	if err := c.Bind(&order); err != nil {
		return err
	}

	err := db.Create(&order).Error
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"status": "error create data",
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"status": "success create data",
	})
}

func Update(c echo.Context) error {
	db := config.GetDB()

	updateOrder := models.UpdateOrder{}
	if err := c.Bind(&updateOrder); err != nil {
		return err
	}

	order := models.Order{}

	// find order
	err := db.Preload(clause.Associations).Find(&order, "id = ?", updateOrder.ID).Error
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success":  false,
			"messages": "Not found",
		})
	}

	items := []models.Item{}

	for _, v := range updateOrder.UpdateItem {
		// find item
		err := db.Preload(clause.Associations).Find(&items, "order_id = ? and item_code = ?", updateOrder.ID, v.ItemCode).Error
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success":  false,
				"messages": "Not found",
			})
		}

		var idItem uint
		for _, v := range items {
			idItem = v.ID
		}

		// set model Item
		item := getModelItem(idItem, v, updateOrder.ID)

		// update if exist create when not exist
		err = db.Save(&item).Error

		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]interface{}{
				"success":  false,
				"messages": "Update data failed",
			})
		}
	}

	// update order
	err = db.Debug().Model(&order).Where("id = ?", updateOrder.ID).Updates(models.Order{
		CustomerName: updateOrder.CustomerName,
	}).Error

	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success":  false,
			"messages": "Update data failed",
		})
	}

	return c.JSON(http.StatusNotFound, map[string]interface{}{
		"success":  true,
		"messages": "Update data success",
		"data":     updateOrder,
	})

}

func getModelItem(idItem uint, v models.UpdateItem, orderId uint) models.Item {
	var item models.Item
	if idItem == 0 {
		item = models.Item{
			ItemCode:    v.ItemCode,
			Description: v.Description,
			Quantity:    v.Quantity,
			OrderID:     int(orderId),
		}
		return item
	}
	item = models.Item{
		Model: gorm.Model{
			ID: idItem,
		},
		ItemCode:    v.ItemCode,
		Description: v.Description,
		Quantity:    v.Quantity,
		OrderID:     int(orderId),
	}
	return item
}

func Delete(c echo.Context) error {
	db := config.GetDB()

	orderId := c.Param("id")

	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return err
	}

	if err := tx.Where("id = ?", orderId).Delete(&models.Order{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("order_id = ?", orderId).Delete(&models.Item{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":  true,
		"messages": "Delete data success",
	})
}

func Get(c echo.Context) error {
	db := config.GetDB()

	var order models.Order
	orderId := c.Param("id")

	err := db.Preload(clause.Associations).Find(&order, "id = ?", orderId).Error

	if err != nil {
		return err
	}

	if order.ID == 0 {
		return c.JSON(http.StatusNotFound, map[string]interface{}{
			"success":  false,
			"messages": "Not found",
		})

	}

	return c.JSON(http.StatusOK, order)
}
