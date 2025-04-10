package repository

import (
    "log"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"

    "github.com/dibsnvas/golang-2025/internal/models"
)

func NewDB(dsn string) (*gorm.DB, error) {
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    err = db.AutoMigrate(
        &models.SalesTransaction{},
        &models.SaleItem{},
        &models.EmployeeAttendance{},
        &models.SalaryPayment{},
    )
    if err != nil {
        return nil, err
    }

    log.Println("Database migrated successfully!")
    return db, nil
}
