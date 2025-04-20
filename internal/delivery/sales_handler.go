package delivery

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "time"

    "github.com/dibsnvas/golang-2025/internal/models"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

type SalesHandler struct {
    DB *gorm.DB
}

func NewSalesHandler(db *gorm.DB) *SalesHandler {
    return &SalesHandler{DB: db}
}

type createSaleRequest struct {
    EmployeeID    uint    `json:"employee_id"`
    ShopID        uint    `json:"shop_id"`
    PaymentMethod string  `json:"payment_method"`
    Items         []struct {
        ItemID      uint    `json:"item_id"`
        Quantity    int     `json:"quantity"`
        PriceAtSale float64 `json:"price_at_sale"`
    } `json:"items"`
}

// CreateSale registers a new sales transaction
// @Summary Create a sales transaction
// @Description Register a new sales transaction for an employee
// @Tags Sales
// @Accept json
// @Produce json
// @Param createSaleRequest body createSaleRequest true "Sale data"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sales [post]
func (h *SalesHandler) CreateSale(c *gin.Context) {
    var req createSaleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    tx := models.SalesTransaction{
        EmployeeID:      req.EmployeeID,
        ShopID:          req.ShopID,
        TransactionTime: time.Now(),
        PaymentMethod:   req.PaymentMethod,
    }

    var total float64
    var saleItems []models.SaleItem
    for _, item := range req.Items {
        total += float64(item.Quantity) * item.PriceAtSale
        saleItems = append(saleItems, models.SaleItem{
            ItemID:      item.ItemID,
            Quantity:    item.Quantity,
            PriceAtSale: item.PriceAtSale,
        })
    }
    tx.TotalAmount = total
    tx.SaleItems = saleItems

    if err := h.DB.Create(&tx).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    for _, item := range saleItems {
        go func(it models.SaleItem) {
            payload := map[string]interface{}{
                "item_id":  it.ItemID,
                "quantity": it.Quantity,
            }

            body, err := json.Marshal(payload)
            if err != nil {
                fmt.Printf("Failed to marshal payload: %v\n", err)
                return
            }

            resp, err := http.Post(
                "http://catalog-service/inventory/deduct",
                "application/json",
                bytes.NewBuffer(body),
            )
            if err != nil {
                fmt.Printf("Failed to notify catalog service: %v\n", err)
                return
            }
            defer resp.Body.Close()
            fmt.Printf("Notified catalog service for item_id=%d, status=%s\n", it.ItemID, resp.Status)
        }(item)
    }

    c.JSON(http.StatusCreated, gin.H{"transaction_id": tx.ID})
}
// GetSalesByEmployeeAndDate returns sales for a specific employee on a specific date
// @Summary Get sales by employee and date
// @Description Get total sales count and amount by employee ID and date
// @Tags Sales
// @Accept json
// @Produce json
// @Param employee_id path int true "Employee ID"
// @Param date query string true "Date in YYYY-MM-DD format"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /sales/employee/{employee_id} [get]
func (h *SalesHandler) GetSalesByEmployeeAndDate(c *gin.Context) {
    employeeIDStr := c.Param("employee_id")
    employeeID, err := strconv.ParseUint(employeeIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
        return
    }

    dateStr := c.Query("date")
    if dateStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "date query param is required, e.g. ?date=2025-04-10"})
        return
    }

    parsedDate, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
        return
    }

    startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, parsedDate.Location())
    endOfDay := startOfDay.Add(24 * time.Hour)

    // Выбираем все транзакции, где employee_id=? и transaction_time в этот день
    var sales []models.SalesTransaction
    if err := h.DB.Where(
        "employee_id = ? AND transaction_time >= ? AND transaction_time < ?",
        employeeID, startOfDay, endOfDay,
    ).Find(&sales).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    countChecks := len(sales)
    var totalAmount float64
    for _, s := range sales {
        totalAmount += s.TotalAmount
    }

    c.JSON(http.StatusOK, gin.H{
        "employee_id":  employeeID,
        "date":         dateStr,
        "count_checks": countChecks,
        "total_amount": totalAmount,
    })
}
