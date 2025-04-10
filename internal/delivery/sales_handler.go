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

// ----- createSaleRequest: входная структура для POST /sales -----
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

// CreateSale обрабатывает POST /sales
func (h *SalesHandler) CreateSale(c *gin.Context) {
    var req createSaleRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Формируем шапку транзакции
    tx := models.SalesTransaction{
        EmployeeID:      req.EmployeeID,
        ShopID:          req.ShopID,
        TransactionTime: time.Now(),
        PaymentMethod:   req.PaymentMethod,
    }

    // Считаем итоговую сумму
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

    // Сохраняем в базе
    if err := h.DB.Create(&tx).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // Уведомляем сервис "Каталог" об уменьшении склада (асинхронно)
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

    // Возвращаем ID транзакции
    c.JSON(http.StatusCreated, gin.H{"transaction_id": tx.ID})
}

// ----- GetSalesByEmployeeAndDate: GET /sales/employee/:employee_id?date=YYYY-MM-DD -----
// Показывает, сколько продаж сделал сотрудник за конкретную дату.
func (h *SalesHandler) GetSalesByEmployeeAndDate(c *gin.Context) {
    // Параметр пути: /sales/employee/:employee_id
    employeeIDStr := c.Param("employee_id")
    employeeID, err := strconv.ParseUint(employeeIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid employee_id"})
        return
    }

    // Query-параметр: ?date=2025-04-10
    dateStr := c.Query("date")
    if dateStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "date query param is required, e.g. ?date=2025-04-10"})
        return
    }

    // Парсим дату
    parsedDate, err := time.Parse("2006-01-02", dateStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format, use YYYY-MM-DD"})
        return
    }

    // Определяем границы дня: [startOfDay, endOfDay)
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

    // Считаем кол-во чеков и сумму total_amount
    countChecks := len(sales)
    var totalAmount float64
    for _, s := range sales {
        totalAmount += s.TotalAmount
    }

    // Возвращаем краткий отчёт
    c.JSON(http.StatusOK, gin.H{
        "employee_id":  employeeID,
        "date":         dateStr,
        "count_checks": countChecks,
        "total_amount": totalAmount,
    })
}
