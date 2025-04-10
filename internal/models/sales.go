package models

import (
    "time"
)

type SalesTransaction struct {
    ID              uint           `gorm:"primaryKey;column:id"`
    EmployeeID      uint           `gorm:"column:employee_id"`
    ShopID          uint           `gorm:"column:shop_id"`
    TransactionTime time.Time      `gorm:"column:transaction_time"`
    TotalAmount     float64        `gorm:"column:total_amount"`
    PaymentMethod   string         `gorm:"column:payment_method"`
    CreatedAt       time.Time      `gorm:"column:created_at"`
    UpdatedAt       time.Time      `gorm:"column:updated_at"`

    SaleItems []SaleItem `gorm:"foreignKey:TransactionID"`
}

type SaleItem struct {
    ID            uint    `gorm:"primaryKey;column:id"`
    TransactionID uint    `gorm:"column:transaction_id"`
    ItemID        uint    `gorm:"column:item_id"`
    Quantity      int     `gorm:"column:quantity"`
    PriceAtSale   float64 `gorm:"column:price_at_sale"`
    CreatedAt     time.Time
    UpdatedAt     time.Time
}
