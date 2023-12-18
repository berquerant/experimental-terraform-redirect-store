package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
)

type Record struct {
	Name string `json:"name"`
	To   string `json:"to"`
}

var (
	ErrRecordNotFound = errors.New("RecordNotFound")

	_ Database = NewDatabaseImpl(nil)
)

// NOTE: make a simple implementation for verification purposes

type Database interface {
	Scan(ctx context.Context) ([]*Record, error)
	Get(ctx context.Context, name string) (*Record, error)
	Put(ctx context.Context, record *Record) error
	Delete(ctx context.Context, name string) error
}

func NewDatabaseImpl(dbFile DatabaseFile) *DatabaseImpl {
	return &DatabaseImpl{
		dbFile: dbFile,
		mux:    sync.RWMutex{},
	}
}

type DatabaseImpl struct {
	dbFile DatabaseFile
	mux    sync.RWMutex
}

func (db *DatabaseImpl) Scan(ctx context.Context) ([]*Record, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	records, err := db.dbFile.Read()
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, ErrRecordNotFound
	}
	return records, nil
}

func (db *DatabaseImpl) Get(ctx context.Context, name string) (*Record, error) {
	db.mux.RLock()
	defer db.mux.RUnlock()

	records, err := db.dbFile.Read()
	if err != nil {
		return nil, err
	}
	for _, r := range records {
		if name == r.Name {
			return r, nil
		}
	}

	return nil, fmt.Errorf("%w, %s", ErrRecordNotFound, name)
}

func (db *DatabaseImpl) Put(ctx context.Context, record *Record) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	records, err := db.dbFile.Read()
	if err != nil {
		return err
	}

	var found *Record
	for _, r := range records {
		if record.Name == r.Name {
			found = r
			break
		}
	}

	if found != nil {
		found.To = record.To
	} else {
		records = append(records, record)
	}
	return db.dbFile.Write(records)
}

func (db *DatabaseImpl) Delete(ctx context.Context, name string) error {
	db.mux.Lock()
	defer db.mux.Unlock()

	records, err := db.dbFile.Read()
	if err != nil {
		return err
	}

	var (
		rs    []*Record
		found bool
	)
	for _, r := range records {
		if r.Name == name {
			found = true
			continue
		}
		rs = append(rs, r)
	}

	if !found {
		return fmt.Errorf("%w, %s", ErrRecordNotFound, name)
	}
	return db.dbFile.Write(rs)
}

var (
	ErrConnectDatabase = errors.New("ConnectDatabase")
	ErrReadDatabase    = errors.New("ReadDatabase")
	ErrWriteDatabase   = errors.New("WriteDatabase")
)

type DatabaseFile interface {
	Write(records []*Record) error
	Read() ([]*Record, error)
}

type databaseFile struct {
	filename string
	mux      sync.RWMutex
}

func NewDatabaseFile(filename string) DatabaseFile {
	return &databaseFile{
		filename: filename,
		mux:      sync.RWMutex{},
	}
}

func (f *databaseFile) Write(records []*Record) error {
	f.mux.Lock()
	defer f.mux.Unlock()

	b, err := json.Marshal(records)
	if err != nil {
		return fmt.Errorf("%w, marshal", ErrWriteDatabase)
	}
	if err := os.WriteFile(f.filename, b, 0666); err != nil {
		return fmt.Errorf("%w, write", ErrWriteDatabase)
	}
	return nil
}

func (f *databaseFile) Read() ([]*Record, error) {
	f.mux.RLock()
	defer f.mux.RUnlock()

	b, err := os.ReadFile(f.filename)
	if err != nil {
		return nil, fmt.Errorf("%w, read", ErrReadDatabase)
	}

	if len(b) == 0 {
		return nil, nil
	}

	var records []*Record
	if err := json.Unmarshal(b, &records); err != nil {
		return nil, fmt.Errorf("%w, unmarshal", ErrReadDatabase)
	}
	return records, nil
}
