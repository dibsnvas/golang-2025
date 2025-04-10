# Sales & Operations Service

A Go-based microservice handling sales transactions, employee attendance (clock-in/clock-out), and salary payments. This project was designed under the assumption that each domain (Sales, Inventory, etc.) is managed by separate microservices.

## Features

1. **Sales**
   - **POST** `/sales`  
     Registers a new sales transaction (including sale items).
     - Stores records in `sales_transactions` and `sale_items`.
     - Asynchronously notifies the external Catalog/Inventory service to deduct stock.
   - **GET** `/sales/employee/:employee_id?date=YYYY-MM-DD`  
     Retrieves how many transactions (checks) and the total sold amount for a given employee on a specific date.
   
2. **Employee Attendance**
   - **POST** `/attendance/clock-in`  
     Marks the time when an employee starts work (clock-in).
   - **POST** `/attendance/clock-out`  
     Marks the time when an employee ends work (clock-out).
   
3. **Salary**
   - **POST** `/salary/pay`  
     Records a salary payment to an employee.
   - **GET** `/salary/:id`  
     Retrieves details of a specific salary payment by ID.

## Entities & Database Structure

- **`sales_transactions`**  
  - Columns: `id`, `employee_id`, `shop_id`, `transaction_time`, `total_amount`, `payment_method`  
  - Represents the "header" of a sale.

- **`sale_items`**  
  - Columns: `id`, `transaction_id`, `item_id`, `quantity`, `price_at_sale`  
  - Stores each sold item in a single transaction.

- **`employee_attendance`**  
  - Columns: `id`, `employee_id`, `clock_in`, `clock_out`  
  - Tracks the working hours for each employee.

- **`salary_payments`**  
  - Columns: `id`, `employee_id`, `pay_period_start`, `pay_period_end`, `amount`, `paid_at`  
  - Records salary payments to employees.

## Installation & Setup

1. **Clone the repository**:
   `
   git clone https://github.com/<username>/sales-operations-service.git
   cd sales-operations-service`
2. **Install Go dependencies**:
  `
go mod tidy
Configure PostgreSQL`

3. **Ensure you have a PostgreSQL instance running (locally or via Docker)**.

`Create a database named sales_ops or update the DSN in cmd/main.go if needed.`

4. **Run the service**:
`
go run ./cmd/main.go
The service will listen on http://localhost:8080.`
