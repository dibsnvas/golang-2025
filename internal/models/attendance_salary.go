package models

import "time"

type EmployeeAttendance struct {
    ID         uint       `gorm:"primaryKey;column:id"`
    EmployeeID uint       `gorm:"column:employee_id"`
    ClockIn    time.Time  `gorm:"column:clock_in"`
    ClockOut   *time.Time `gorm:"column:clock_out"`
    CreatedAt  time.Time  `gorm:"column:created_at"`
    UpdatedAt  time.Time  `gorm:"column:updated_at"`
}
