package models

import "time"

type SalaryPayment struct {
    ID             uint      `gorm:"primaryKey;column:id"`
    EmployeeID     uint      `gorm:"column:employee_id"`
    PayPeriodStart time.Time `gorm:"column:pay_period_start"`
    PayPeriodEnd   time.Time `gorm:"column:pay_period_end"`
    Amount         float64   `gorm:"column:amount"`
    PaidAt         time.Time `gorm:"column:paid_at"`
}

func (SalaryPayment) TableName() string {
    return "salary_payments"
}
