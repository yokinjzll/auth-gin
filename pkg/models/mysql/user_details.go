package mysql

import (
	"database/sql"
	"errors"
	"gin-chat/pkg/models"
)

type UserDetailsModel struct {
	DB *sql.DB
}

func (udm *UserDetailsModel) Insert(ud models.UserDetails) (int64, error) {

	stmp := `INSERT INTO user_details (user_id, first_name, last_name, dob, created_at)
	VALUES(?, ?, ?, ?, UTC_TIMESTAMP())`

	result, err := udm.DB.Exec(stmp, &ud.UserID, &ud.First_name, &ud.Last_name, &ud.Dob)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (udm *UserDetailsModel) Get(id string) (*models.UserDetails, error) {
	stmt := `SELECT user_id, first_name, last_name, dob, created_at FROM user_details WHERE user_id = ?`

	row := udm.DB.QueryRow(stmt, id)
	u := &models.UserDetails{}
	err := row.Scan(&u.UserID, &u.First_name, &u.Last_name, &u.Dob, &u.Created_at)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}
	return u, nil
}

func (udm *UserDetailsModel) SetByID(id string, ud models.UserDetails) bool {
	req := `UPDATE user_details SET first_name = ?, last_name = ?, dob = ? WHERE user_id = ?`

	_, err := udm.DB.Exec(req, &ud.First_name, &ud.Last_name, &ud.Dob, &ud.UserID)
	if err != nil {
		return false
	}

	return true
}
