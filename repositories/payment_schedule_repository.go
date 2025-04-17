package repositories

import (
    "database/sql"
    "time"

    "bank-api/models"
)

type PaymentScheduleRepository interface {
    Create(ps *models.PaymentSchedule) error
    GetByID(id int) (*models.PaymentSchedule, error)
    // Получить все просроченные неоплаченные платежи
    GetOverdueUnpaid(cutoff time.Time) ([]*models.PaymentSchedule, error)
    // Новый метод: получить график по кредиту
    GetByCreditID(creditID int) ([]*models.PaymentSchedule, error)
    Update(ps *models.PaymentSchedule) error
}

type paymentScheduleRepository struct {
    db *sql.DB
}

func NewPaymentScheduleRepository(db *sql.DB) PaymentScheduleRepository {
    return &paymentScheduleRepository{db: db}
}

func (r *paymentScheduleRepository) Create(ps *models.PaymentSchedule) error {
    _, err := r.db.Exec(
        `INSERT INTO payment_schedules (credit_id, due_date, amount, is_paid, created_at)
         VALUES ($1, $2, $3, $4, NOW())`,
        ps.CreditID, ps.DueDate, ps.Amount, ps.IsPaid,
    )
    return err
}

func (r *paymentScheduleRepository) GetByID(id int) (*models.PaymentSchedule, error) {
    row := r.db.QueryRow(
        `SELECT id, credit_id, due_date, amount, is_paid, created_at
         FROM payment_schedules WHERE id = $1`, id,
    )
    ps := &models.PaymentSchedule{}
    if err := row.Scan(
        &ps.ID, &ps.CreditID, &ps.DueDate, &ps.Amount, &ps.IsPaid, &ps.CreatedAt,
    ); err != nil {
        return nil, err
    }
    return ps, nil
}

func (r *paymentScheduleRepository) GetOverdueUnpaid(cutoff time.Time) ([]*models.PaymentSchedule, error) {
    rows, err := r.db.Query(
        `SELECT id, credit_id, due_date, amount, is_paid, created_at
         FROM payment_schedules WHERE due_date < $1 AND is_paid = false`,
        cutoff,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []*models.PaymentSchedule
    for rows.Next() {
        ps := &models.PaymentSchedule{}
        if err := rows.Scan(&ps.ID, &ps.CreditID, &ps.DueDate, &ps.Amount, &ps.IsPaid, &ps.CreatedAt); err != nil {
            return nil, err
        }
        list = append(list, ps)
    }
    return list, nil
}

func (r *paymentScheduleRepository) GetByCreditID(creditID int) ([]*models.PaymentSchedule, error) {
    rows, err := r.db.Query(
        `SELECT id, credit_id, due_date, amount, is_paid, created_at
         FROM payment_schedules WHERE credit_id = $1 ORDER BY due_date`,
        creditID,
    )
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var list []*models.PaymentSchedule
    for rows.Next() {
        ps := &models.PaymentSchedule{}
        if err := rows.Scan(&ps.ID, &ps.CreditID, &ps.DueDate, &ps.Amount, &ps.IsPaid, &ps.CreatedAt); err != nil {
            return nil, err
        }
        list = append(list, ps)
    }
    return list, nil
}

func (r *paymentScheduleRepository) Update(ps *models.PaymentSchedule) error {
    _, err := r.db.Exec(
        `UPDATE payment_schedules SET is_paid=$1 WHERE id=$2`,
        ps.IsPaid, ps.ID,
    )
    return err
}