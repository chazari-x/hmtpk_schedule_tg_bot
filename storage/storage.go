package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/chazari-x/hmtpk_schedule/config"
	"strconv"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/net/context"
)

const (
	selectGroupID      = `SELECT group_name FROM schedule_bot WHERE id = $1`
	changeGroupID      = `UPDATE schedule_bot SET group_name = $2 WHERE id = $1`
	updateLastActivity = `UPDATE schedule_bot SET last_activity = $2 WHERE id = $1`
	getActiveChats     = `SELECT count(group_name) FROM schedule_bot WHERE last_activity > $1`
	insertChat         = `INSERT INTO schedule_bot (id, last_activity) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET last_activity = $2`
	createTable        = `CREATE TABLE IF NOT EXISTS schedule_bot (
							id 				VARCHAR PRIMARY KEY NOT NULL, 
							group_name		VARCHAR 			NOT NULL DEFAULT 0,
							last_activity 	DATE 				NOT NULL)`
)

type Storage struct {
	db  *sql.DB
	ctx context.Context
}

func NewStorage(dbCfg *config.DataBase, ctx context.Context) (*Storage, *sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Pass, dbCfg.Name)
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, nil, err
	}

	ctxNew, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	if err = db.PingContext(ctxNew); err != nil {
		return nil, nil, err
	}

	if _, err = db.ExecContext(ctxNew, createTable); err != nil {
		return nil, nil, err
	}

	return &Storage{db: db, ctx: ctx}, db, nil
}

func (s *Storage) InsertChat(chatID int) error {
	if s.db == nil {
		return errors.New("не установленно подключение к базе данных")
	}

	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, insertChat, strconv.Itoa(chatID), time.Now()); err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateLastActivity(chatID int) error {
	if s.db == nil {
		return errors.New("не установленно подключение к базе данных")
	}

	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, updateLastActivity, strconv.Itoa(chatID), time.Now()); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetActiveChats() (int, error) {
	if s.db == nil {
		return 0, errors.New("не установленно подключение к базе данных")
	}

	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	var chatsNum string
	if err := s.db.QueryRowContext(ctx, getActiveChats, time.Duration(time.Now().Month())).Scan(&chatsNum); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("storage is nil")
		}

		return 0, err
	}

	return strconv.Atoi(chatsNum)
}

func (s *Storage) SelectGroupID(chatID int) (string, error) {
	if s.db == nil {
		return "", errors.New("не установленно подключение к базе данных")
	}

	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	var dbItem string
	if err := s.db.QueryRowContext(ctx, selectGroupID, strconv.Itoa(chatID)).Scan(&dbItem); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("storage is nil")
		}

		return "", err
	}

	return dbItem, nil
}

func (s *Storage) ChangeGroupID(chatID int, group string) error {
	if s.db == nil {
		return errors.New("не установленно подключение к базе данных")
	}

	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	if _, err := s.db.ExecContext(ctx, changeGroupID, strconv.Itoa(chatID), group); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("storage is nil")
		}

		return err
	}

	return nil
}
