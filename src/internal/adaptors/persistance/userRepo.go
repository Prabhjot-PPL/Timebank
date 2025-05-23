package persistance

import (
	"context"
	"fmt"
	"log"
	"timebank/src/internal/adaptors/ports"
	"timebank/src/internal/core/dto"
	"timebank/src/pkg"
)

type UserRepo struct {
	db *Database
}

func NewUserRepo(d *Database) ports.UserRepository {
	return &UserRepo{db: d}
}

// --------------------------------AUTH----------------------------------------

func (u *UserRepo) CreateUser(ctx context.Context, user dto.UserDetails) error {
	hpassword, e := pkg.HashPassword(user.Password)
	if e != nil {
		fmt.Println("Unable to hash password : ", e)
	}

	var userId int
	err := u.db.db.QueryRowContext(ctx, `
		INSERT INTO user_details (username, email, password)
		VALUES ($1, $2, $3)
		RETURNING id
	`, user.Username, user.Email, hpassword).Scan(&userId)

	if err != nil {
		log.Println("Error inserting user_details: ", err)
		return err
	}

	// Insert skills_offered
	for _, skill := range user.SkillsOffered {
		_, err := u.db.db.ExecContext(ctx, `
			INSERT INTO skills_offered (user_id, skill_offered) VALUES ($1, $2)
		`, userId, skill)

		if err != nil {
			log.Println("Error inserting skill into skills_offered: ", err)
			return err
		}
	}

	// Insert skills_needed
	for _, skill := range user.SkillsNeeded {
		_, err := u.db.db.ExecContext(ctx, `
			INSERT INTO skills_needed (user_id, skill) VALUES ($1, $2)
		`, userId, skill)

		if err != nil {
			log.Println("Error inserting skill into skills_needed: ", err)
			return err
		}
	}

	return nil
}

func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (dto.UserDetails, error) {
	var user dto.UserDetails

	query := `SELECT id, username, email, password FROM user_details WHERE email=$1`
	row := u.db.db.QueryRowContext(ctx, query, email)

	err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return user, err
	}

	return user, nil
}

// --------------------------------SESSION----------------------------------------

func (u *UserRepo) CreateSession(ctx context.Context, session dto.Session) error {
	_, err := u.db.db.ExecContext(ctx, `
		INSERT INTO sessions (
			session_date, helper, recipient, skill, hours, session_status, feedback
		)
		VALUES ($1, $2, $3, $4, $5, $6, NULL)
	`,
		session.SessionDate,
		session.Helper,
		session.Recipient,
		session.Skill,
		session.Hours,
		session.SessionStatus,
	)

	if err != nil {
		log.Println("Error inserting session: ", err)
		return err
	}

	return nil
}

func (u *UserRepo) CompleteSessionTx(ctx context.Context, sessionID int, feedback, status string) error {
	tx, err := u.db.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get helper, recipient, and hours from the session
	var helperID, recipientID int
	var hours float64
	err = tx.QueryRowContext(ctx, `
		SELECT helper, recipient, hours FROM sessions WHERE session_id = $1
	`, sessionID).Scan(&helperID, &recipientID, &hours)
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Update the session
	_, err = tx.ExecContext(ctx, `
		UPDATE sessions
		SET session_status = $1, feedback = $2
		WHERE session_id = $3
	`, status, feedback, sessionID)
	if err != nil {
		return fmt.Errorf("failed to update session: %w", err)
	}

	// Increment helper's balance_hours
	_, err = tx.ExecContext(ctx, `
		UPDATE user_details SET balance_hours = balance_hours + $1 WHERE id = $2
	`, hours, helperID)
	if err != nil {
		return fmt.Errorf("failed to update helper balance: %w", err)
	}

	// Decrement recipient's balance_hours
	_, err = tx.ExecContext(ctx, `
		UPDATE user_details SET balance_hours = balance_hours - $1 WHERE id = $2
	`, hours, recipientID)
	if err != nil {
		return fmt.Errorf("failed to update recipient balance: %w", err)
	}

	return tx.Commit()
}
