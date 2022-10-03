package main

import (
	"database/sql"
	sq "github.com/Masterminds/squirrel"
	"github.com/arhamj/abi-extractor/pkg/external"
	"github.com/arhamj/abi-extractor/pkg/util"
	_ "github.com/mattn/go-sqlite3"
)

type MappingKind string

const (
	Function MappingKind = "function"
	Event    MappingKind = "event"
)

var (
	migrations = `
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
	db      *sql.DB
	gateway external.FourByteGateway
}

func NewFourByteScraper(fourByteGateway external.FourByteGateway) (*FourByteScraper, error) {
	db, err := util.NewSQLiteDB("db/scraper.db", migrations)
	if err != nil {
		return nil, err
	}
	s := FourByteScraper{
		gateway: fourByteGateway,
		db:      db,
	}
	return &s, nil
}

func (s *FourByteScraper) Start() error {
	var lastPageSynced int
	query := sq.
		Select("last_synced_page").
		From("sync_status_fourbyte").
		Where("kind = ?", Function).RunWith(s.db)
	err := query.QueryRow().Scan(&lastPageSynced)
	if err == sql.ErrNoRows {
		lastPageSynced = 1
	} else if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	resp, err := s.gateway.GetFunctionSignatures(lastPageSynced)
	if err != nil {
		return err
	}
	for _, sign := range resp.Results {
		nativeInsert := sq.Expr("INSERT OR IGNORE INTO sign_mapping_fourbyte (kind,hex_sign,string_sign,created_at) VALUES (?,?,?,?)", Function, sign.HexSignature, sign.TextSignature, sign.CreatedAt)
		_, err = sq.ExecWith(s.db, nativeInsert)
		if err != nil {
			return err
		}
	}
	return nil
	// fetch the last synced page number
	// start event sync (go)
	// start function sync (go)
}

func main() {
	util.SetupDevLogger()
	scraper, err := NewFourByteScraper(external.NewFourByteGateway())
	if err != nil {
		panic(err)
	}
	scraper.Start()
}

func (s *FourByteScraper) Stop() {

}
