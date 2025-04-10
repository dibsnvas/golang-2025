package delivery

import (
    "net/http"
    "strconv"
    "time"

    "github.com/dibsnvas/golang-2025/internal/models"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type SalaryHandler struct {
    DB *gorm.DB
}

func NewSalaryHandler(db *gorm.DB) *SalaryHandler {
    return &SalaryHandler{DB: db}
}

// PaySalaryRequest — структура для JSON-запроса при POST /salary/pay
type PaySalaryRequest struct {
    EmployeeID     uint    `json:"employee_id"`
    PayPeriodStart string  `json:"pay_period_start"` // строка, чтобы потом распарсить "YYYY-MM-DD"
    PayPeriodEnd   string  `json:"pay_period_end"`
    Amount         float64 `json:"amount"`
    PaidAt         string  `json:"paid_at"` // можно не указывать и взять time.Now()
}

// PaySalary — POST /salary/pay
func (h *SalaryHandler) PaySalary(c *gin.Context) {
    var req PaySalaryRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    start, err := time.Parse("2006-01-02", req.PayPeriodStart)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pay_period_start"})
        return
    }

    end, err := time.Parse("2006-01-02", req.PayPeriodEnd)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid pay_period_end"})
        return
    }

    paidAt := time.Now()
    if req.PaidAt != "" {
        paidAt, err = time.Parse("2006-01-02", req.PaidAt)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "invalid paid_at"})
            return
        }
    }

    salary := models.SalaryPayment{
        EmployeeID:     req.EmployeeID,
        PayPeriodStart: start,
        PayPeriodEnd:   end,
        Amount:         req.Amount,
        PaidAt:         paidAt,
    }

    if err := h.DB.Create(&salary).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"salary_id": salary.ID})
}

// GetSalaryByID — GET /salary/:id
func (h *SalaryHandler) GetSalaryByID(c *gin.Context) {
    idStr := c.Param("id")
    id, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
        return
    }

    var salary models.SalaryPayment
    if err := h.DB.First(&salary, id).Error; err != nil {
        if err == gorm.ErrRecordNotFound {
            c.JSON(http.StatusNotFound, gin.H{"error": "salary not found"})
        } else {
            c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        }
        return
    }

    c.JSON(http.StatusOK, salary)
}
