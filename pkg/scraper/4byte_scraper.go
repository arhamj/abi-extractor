package scraper

import (
	"context"
	"database/sql"
	"errors"
	"github.com/arhamj/abi-extractor/pkg/external"
	"github.com/arhamj/abi-extractor/pkg/util"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
	"strings"
)

type MappingKind string

const (
	Function MappingKind = "function"
	Event    MappingKind = "event"
)

var (
	FourByteMigrations = `
CREATE TABLE IF NOT EXISTS sign_mapping_fourbyte
(
    id          INTEGER						PRIMARY KEY,
    kind        VARCHAR(16)              	NOT NULL,
    hex_sign    VARCHAR(64)              	NOT NULL,
    string_sign TEXT                     	NOT NULL,
    created_at  TIMESTAMP WITH TIME ZONE 	NOT NULL
);

CREATE TABLE IF NOT EXISTS sync_status_fourbyte
(
    kind				VARCHAR(16)	NOT NULL,
    last_synced_page 	INT			NOT NULL
);

CREATE UNIQUE INDEX IF NOT EXISTS sync_status_fourbyte__kind_index ON sync_status_fourbyte (kind);

CREATE UNIQUE INDEX IF NOT EXISTS sign_mapping_fourbyte__unique_index ON sign_mapping_fourbyte (kind, hex_sign, string_sign);

CREATE INDEX IF NOT EXISTS sign_mapping_fourbyte__kind_hex_sign_index ON sign_mapping_fourbyte (kind, hex_sign);
`
)

type FourByteScraper struct {
	logger  *zap.Logger
	db      *sql.DB
	gateway external.FourByteGateway

	ctx context.Context
}

func NewFourByteScraper(ctx context.Context, fourByteGateway external.FourByteGateway) (*FourByteScraper, error) {
	db, err := util.NewSQLiteDB("db/scraper.db", FourByteMigrations)
	if err != nil {
		return nil, err
	}
	s := FourByteScraper{
		logger:  zap.L().With(zap.String("loc", "FourByteScraper")),
		gateway: fourByteGateway,
		db:      db,
		ctx:     ctx,
	}
	return &s, nil
}

func (s *FourByteScraper) Start(kind MappingKind) error {
	for {
		select {
		case <-s.ctx.Done():
			s.logger.Info("Stopping sync!", zap.String("kind", string(kind)))
			s.Stop()
			return nil
		default:
			err := s.sync(kind)
			if err != nil && strings.Contains(err.Error(), "sync completed") {
				s.logger.Error("Sync completed and up to date!", zap.String("kind", string(kind)))
				return nil
			} else if err != nil {
				return err
			}
		}
	}
}

func (s *FourByteScraper) sync(kind MappingKind) error {
	lastPageSynced, err := s.fetchLastSyncedPage(kind)
	if err != nil {
		return err
	}
	pageToSync := lastPageSynced + 1
	var resp *external.FourBytesResp
	if kind == Function {
		resp, err = s.gateway.GetFunctionSignatures(pageToSync)
		if err != nil {
			return err
		}
	} else if kind == Event {
		resp, err = s.gateway.GetEventSignatures(pageToSync)
		if err != nil {
			return err
		} else if len(resp.Results) == 0 {
			return errors.New("sync completed and up to date")
		}
	}
	err = s.bulkInsertRecords(kind, resp)
	if err != nil {
		s.logger.Error("Error when bulk inserting Function records to SQLite", zap.Error(err))
		return err
	}
	err = s.updateLastSyncedPage(pageToSync, kind)
	if err != nil {
		s.logger.Error("Error when updating the last synced page for functions", zap.Error(err))
		return err
	}
	s.logger.Info("Sync info", zap.String("kind", string(kind)), zap.Int("count", len(resp.Results)))
	return nil
}

func (s *FourByteScraper) bulkInsertRecords(recordType MappingKind, resp *external.FourBytesResp) error {
	for _, sign := range resp.Results {
		stmt, err := s.db.Prepare("INSERT OR IGNORE INTO sign_mapping_fourbyte (kind,hex_sign,string_sign,created_at) VALUES (?,?,?,?)")
		if err != nil {
			return err
		}
		_, err = stmt.Exec(recordType, sign.HexSignature, sign.TextSignature, sign.CreatedAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *FourByteScraper) fetchLastSyncedPage(kind MappingKind) (int, error) {
	var lastPageSynced int
	row := s.db.QueryRow("SELECT last_synced_page FROM sync_status_fourbyte WHERE kind = ?", kind)
	if row.Err() != nil {
		return 0, row.Err()
	}
	err := row.Scan(&lastPageSynced)
	if err == sql.ErrNoRows {
		lastPageSynced = 0
	} else if err != nil && err != sql.ErrNoRows {
		return 0, err
	}
	return lastPageSynced, nil
}

func (s *FourByteScraper) updateLastSyncedPage(pageNo int, kind MappingKind) error {
	stmt, err := s.db.Prepare("UPDATE sync_status_fourbyte set last_synced_page = ? where kind = ?")
	if err != nil {
		return err
	}
	if pageNo == 1 {
		stmt, err = s.db.Prepare("INSERT INTO sync_status_fourbyte(last_synced_page, kind) VALUES (?,?)")
		if err != nil {
			return err
		}
	}
	_, err = stmt.Exec(pageNo, kind)
	if err != nil {
		return err
	}
	return nil
}

func (s *FourByteScraper) Stop() {
	s.db.Close()
}
