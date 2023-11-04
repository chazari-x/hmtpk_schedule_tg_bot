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
	getActiveChats     = `SELECT count(group_name) FROM schedule_bot WHERE last_activity >= $1`
	insertChat         = `INSERT INTO schedule_bot (id, last_activity) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET last_activity = $2`
	createTable        = `CREATE TABLE IF NOT EXISTS schedule_bot (
							id 				VARCHAR PRIMARY KEY NOT NULL, 
							group_name		VARCHAR 			NOT NULL DEFAULT 0,
							last_activity 	DATE 				NOT NULL)`
)

type Storage struct {
	DB  *sql.DB
	ctx context.Context
	dsn string
}

func NewStorage(dbCfg *config.DataBase, ctx context.Context) (*Storage, *sql.DB, error) {
	ctxNew, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbCfg.Host, dbCfg.Port, dbCfg.User, dbCfg.Pass, dbCfg.Name)
	db, err := connect(dsn, ctxNew)
	if err != nil {
		return nil, nil, err
	}

	return &Storage{DB: db, ctx: ctx, dsn: dsn}, db, nil
}

func connect(dsn string, ctx context.Context) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, err
	}

	if _, err = db.ExecContext(ctx, createTable); err != nil {
		return nil, err
	}

	return db, nil
}

func (s *Storage) Ping() (*sql.DB, error) {
	ctx, cancel := context.WithTimeout(s.ctx, time.Second*2)
	defer cancel()

	if err := s.DB.PingContext(ctx); err != nil {
		db, err := connect(s.dsn, ctx)
		return db, err
	}

	return nil, nil
}

func (s *Storage) InsertChat(chatID int) error {
	if s.DB == nil {
		return errors.New("не установленно подключение к базе данных")
	}

	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	if _, err := s.DB.ExecContext(ctx, insertChat, strconv.Itoa(chatID), time.Now()); err != nil {
		return err
	}

	return nil
}

func (s *Storage) UpdateLastActivity(chatID int) error {
	if s.DB == nil {
		return errors.New("не установленно подключение к базе данных")
	}

	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	if _, err := s.DB.ExecContext(ctx, updateLastActivity, strconv.Itoa(chatID), time.Now()); err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetActiveChats() (int, int, error) {
	if s.DB == nil {
		return 0, 0, errors.New("не установленно подключение к базе данных")
	}

	ctx, cancel := context.WithTimeout(s.ctx, time.Second*2)
	defer cancel()

	var monthNum string
	if err := s.DB.QueryRowContext(ctx, getActiveChats, time.Now().AddDate(0, -1, 0)).Scan(&monthNum); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, errors.New("storage is nil")
		}

		return 0, 0, err
	}

	var dayNum string
	if err := s.DB.QueryRowContext(ctx, getActiveChats, time.Now()).Scan(&dayNum); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, errors.New("storage is nil")
		}

		return 0, 0, err
	}

	day, err := strconv.Atoi(dayNum)
	if err != nil {
		return 0, 0, err
	}

	month, err := strconv.Atoi(monthNum)
	if err != nil {
		return 0, 0, err
	}

	return day, month, nil
}

func (s *Storage) SelectGroupID(chatID int) (string, error) {
	if s.DB == nil {
		return "", errors.New("не установленно подключение к базе данных")
	}

	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	var dbItem string
	if err := s.DB.QueryRowContext(ctx, selectGroupID, strconv.Itoa(chatID)).Scan(&dbItem); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", errors.New("storage is nil")
		}

		return "", err
	}

	return dbItem, nil
}

func (s *Storage) ChangeGroupID(chatID int, group string) error {
	if s.DB == nil {
		return errors.New("не установленно подключение к базе данных")
	}

	ctx, cancel := context.WithTimeout(s.ctx, time.Second)
	defer cancel()

	if _, err := s.DB.ExecContext(ctx, changeGroupID, strconv.Itoa(chatID), group); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errors.New("storage is nil")
		}

		return err
	}

	return nil
}
