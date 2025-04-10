package delivery

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "gorm.io/gorm"

    "github.com/dibsnvas/golang-2025/internal/models"
)

type AttendanceHandler struct {
    DB *gorm.DB
}

func NewAttendanceHandler(db *gorm.DB) *AttendanceHandler {
    return &AttendanceHandler{db}
}

type clockInRequest struct {
    EmployeeID uint `json:"employee_id"`
}

func (h *AttendanceHandler) ClockIn(c *gin.Context) {
    var req clockInRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    record := models.EmployeeAttendance{
        EmployeeID: req.EmployeeID,
        ClockIn:    time.Now(),
    }

    if err := h.DB.Create(&record).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"attendance_id": record.ID})
}

type clockOutRequest struct {
    EmployeeID uint `json:"employee_id"`
}

func (h *AttendanceHandler) ClockOut(c *gin.Context) {
    var req clockOutRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var record models.EmployeeAttendance
    if err := h.DB.Where("employee_id = ? AND clock_out IS NULL", req.EmployeeID).
        Order("clock_in desc").
        First(&record).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "No active clock-in found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    now := time.Now()
    record.ClockOut = &now

    if err := h.DB.Save(&record).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "attendance_id": record.ID,
        "clock_in":      record.ClockIn,
        "clock_out":     record.ClockOut,
    })
}
