package database

import (
	"github.com/imdario/mergo"
	"github.com/jinzhu/gorm"
	"github.com/pmdcosta/treasure-coin"
	log "github.com/sirupsen/logrus"
)

// UserService represents a service for managing user persistence.
type UserService struct {
	client *Client
}

// Create creates a new user record in the database.
func (s *UserService) Add(user coin.User) (coin.User, error) {
	u := NewUser(user)
	if dbc := s.client.db.Create(&u); dbc.Error != nil {
		s.client.logger.WithFields(log.Fields{"err": dbc.Error}).Info(ErrCreate)
		return user, ErrCreate
	}
	return u.toDomain(), nil
}

// Find retrieves a user from the database by the ID.
func (s *UserService) Find(id uint) (coin.User, error) {
	u := User{}
	if dbc := s.client.db.First(&u, id); dbc.Error != nil {
		s.client.logger.WithFields(log.Fields{"err": dbc.Error, "id": id}).Info(ErrRetrieving)
		return coin.User{}, ErrRetrieving
	}
	return u.toDomain(), nil
}

// FindByEmail retrieves a user from the database by the email address.
func (s *UserService) FindByEmail(email string) (coin.User, error) {
	u := User{}
	if dbc := s.client.db.Where(&User{Email: email}).First(&u); dbc.Error != nil {
		s.client.logger.WithFields(log.Fields{"err": dbc.Error, "email": email}).Info(ErrRetrieving)
		return coin.User{}, ErrRetrieving
	}
	return u.toDomain(), nil
}

// FindByUsername retrieves a user from the database by the username.
func (s *UserService) FindByUsername(username string) (coin.User, error) {
	u := User{}
	if dbc := s.client.db.Where(&User{Username: username}).First(&u); dbc.Error != nil {
		s.client.logger.WithFields(log.Fields{"err": dbc.Error, "username": username}).Info(ErrRetrieving)
		return coin.User{}, ErrRetrieving
	}
	return u.toDomain(), nil
}

// Update updated a user record in the database.
func (s *UserService) Update(user coin.User) (coin.User, error) {
	// begin a transaction
	tx := s.client.db.Begin()
	if tx.Error != nil {
		s.client.logger.WithFields(log.Fields{"err": tx.Error}).Info("failed to start transaction")
		return user, ErrUpdate
	}

	// retrieve the user from the database.
	u := User{}
	if dbc := tx.First(&u, user.ID); dbc.Error != nil {
		s.client.logger.WithFields(log.Fields{"err": dbc.Error, "id": user.ID}).Info("failed to get the current user state from the database")
		tx.Rollback()
		return user, ErrRetrieving
	}

	// update the user structure.
	if err := u.fromDomain(user); err != nil {
		s.client.logger.WithFields(log.Fields{"err": err, "id": user.ID}).Info("failed to merge structures")
		tx.Rollback()
		return user, ErrUpdate
	}

	// update the record.
	if dbc := tx.Save(&u); dbc.Error != nil {
		s.client.logger.WithFields(log.Fields{"err": dbc.Error, "id": user.ID}).Info(ErrUpdate)
		tx.Rollback()
		return user, ErrUpdate
	}

	// commit transaction.
	if dbc := tx.Commit(); dbc.Error != nil {
		s.client.logger.WithFields(log.Fields{"err": dbc.Error, "id": user.ID}).Info("failed to commit the transaction to the database")
		tx.Rollback()
		return user, ErrUpdate
	}

	return u.toDomain(), nil
}

// Delete removes a user from the database.
func (s *UserService) Delete(user coin.User) error {
	// retrieve the user from the database.
	u := User{}
	if dbc := s.client.db.First(&u, user.ID); dbc.Error != nil {
		s.client.logger.WithFields(log.Fields{"err": dbc.Error, "id": user.ID}).Info("failed to get the current user state from the database")
		return ErrRetrieving
	}

	// delete user.
	if dbc := s.client.db.Delete(u); dbc.Error != nil {
		s.client.logger.WithFields(log.Fields{"err": dbc.Error}).Info(ErrDelete)
		return ErrDelete
	}
	return nil
}

// PurgeAll deletes all the users from the database.
func (s *UserService) PurgeAll() {
	s.client.db.DropTableIfExists(&User{})
}

// User represents the user structure to be persisted.
type User struct {
	Email    string `gorm:"type:varchar(100);unique_index"`
	Username string `gorm:"type:varchar(100);unique_index"`
	Password string
	gorm.Model
}

// NewUser creates a new database user from a domain user.
func NewUser(user coin.User) User {
	u := User{
		Email:    user.Email,
		Username: user.Username,
		Password: user.Password,
	}
	return u
}

// toDomain builds a domain user from the database representation.
func (u *User) toDomain() coin.User {
	return coin.User{
		ID:       u.ID,
		Email:    u.Email,
		Username: u.Username,
		Password: u.Password,
	}
}

// fromDomain updates the user structure from a domain user data.
func (u *User) fromDomain(user coin.User) error {
	if err := mergo.Merge(u, NewUser(user), mergo.WithOverride); err != nil {
		return err
	}
	return nil
}
