package repositories

import (
    "database/sql"
    "fmt"

    "bank-api/models"
)

type CreditRepository interface {
    Create(c *models.Credit) error
    GetByID(id int) (*models.Credit, error)
    // Новый метод: получить все кредиты пользователя
    GetByUserID(userID int) ([]*models.Credit, error)
    // Обновить сумму кредита после штрафа
    UpdateAmount(creditID int, newAmount float64) error
}

type creditRepository struct {
    db *sql.DB
}

func NewCreditRepository(db *sql.DB) CreditRepository {
    return &creditRepository{db: db}
}

func (r *creditRepository) Create(c *models.Credit) error {
    _, err := r.db.Exec(
        `INSERT INTO credits (user_id, account_id, amount, interest_rate, created_at)
         VALUES ($1, $2, $3, $4, NOW())`,
        c.UserID, c.AccountID, c.Amount, c.InterestRate,
    )
    return err
}

func (r *creditRepository) GetByID(id int) (*models.Credit, error) {
    row := r.db.QueryRow(
        `SELECT id, user_id, account_id, amount, interest_rate, created_at
         FROM credits WHERE id = $1`, id,
    )
    cr := &models.Credit{}
    if err := row.Scan(&cr.ID, &cr.UserID, &cr.AccountID, &cr.Amount, &cr.InterestRate, &cr.CreatedAt); err != nil {
        return nil, err
    }
    return cr, nil
}

func (r *creditRepository) GetByUserID(userID int) ([]*models.Credit, error) {
    rows, err := r.db.Query(
        `SELECT id, user_id, account_id, amount, interest_rate, created_at
         FROM credits WHERE user_id = $1`, userID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []*models.Credit
    for rows.Next() {
        cr := &models.Credit{}
        if err := rows.Scan(&cr.ID, &cr.UserID, &cr.AccountID, &cr.Amount, &cr.InterestRate, &cr.CreatedAt); err != nil {
            return nil, err
        }
        list = append(list, cr)
    }
    return list, nil
}

func (r *creditRepository) UpdateAmount(creditID int, newAmount float64) error {
    res, err := r.db.Exec(
        `UPDATE credits SET amount = $1 WHERE id = $2`,
        newAmount, creditID,
    )
    if err != nil {
        return err
    }
    if cnt, _ := res.RowsAffected(); cnt == 0 {
        return fmt.Errorf("credit %d not found", creditID)
    }
    return nil
}