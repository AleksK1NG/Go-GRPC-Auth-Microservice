package repository

const (
	createUserQuery = `INSERT INTO users (first_name, last_name, email, password, role, avatar) 
		VALUES ($1, $2, $3, $4, $5 ,$6) 
		RETURNING user_id, first_name, last_name, email, password, avatar, created_at, updated_at`

	findByEmailQuery = `SELECT user_id, email, first_name, last_name, role, avatar, password, created_at, updated_at FROM users WHERE email = $1`

	findByIDQuery = `SELECT user_id, email, first_name, last_name, role, avatar, created_at, updated_at FROM users WHERE user_id = $1`
)
