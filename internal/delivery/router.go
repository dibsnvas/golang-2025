package delivery

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
    _ "github.com/dibsnvas/golang-2025/docs" 
)

func SetupRouter(db *gorm.DB) *gin.Engine {
    r := gin.Default()

    salesHandler := NewSalesHandler(db)
    attendanceHandler := NewAttendanceHandler(db)
    salaryHandler := NewSalaryHandler(db)

    r.POST("/sales", salesHandler.CreateSale)

	r.POST("/salary/pay", salaryHandler.PaySalary)
	r.GET("/salary/:id", salaryHandler.GetSalaryByID)

    r.POST("/attendance/clock-in", attendanceHandler.ClockIn)
    r.POST("/attendance/clock-out", attendanceHandler.ClockOut)

	r.GET("/sales/employee/:employee_id", salesHandler.GetSalesByEmployeeAndDate)
    
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    return r
}
