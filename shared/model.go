package shared

type User struct {
    UserID           string `json:"user_id"`
    Name             string `json:"name"`
    Email            string `json:"email"`
    MonthlyIncome    int    `json:"monthly_income"`
    CreditScore      int    `json:"credit_score"`
    EmploymentStatus string `json:"employment_status"`
    Age              int    `json:"age"`
}


