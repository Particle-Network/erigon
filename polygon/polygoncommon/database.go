// Copyright 2024 The Erigon Authors
// This file is part of Erigon.
//
// Erigon is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// Erigon is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with Erigon. If not, see <http://www.gnu.org/licenses/>.

package polygoncommon

import (
	"context"
	"path/filepath"
	"sync"

	"github.com/c2h5oh/datasize"

	"github.com/ledgerwatch/erigon-lib/log/v3"

	"github.com/ledgerwatch/erigon-lib/kv"
	"github.com/ledgerwatch/erigon-lib/kv/mdbx"
)

type Database struct {
	db kv.RwDB

	dataDir  string
	openOnce sync.Once

	logger log.Logger
}

func NewDatabase(
	dataDir string,
	logger log.Logger,
) *Database {
	return &Database{dataDir: dataDir, logger: logger}
}

func (db *Database) open(ctx context.Context, label kv.Label, tableCfg kv.TableCfg) error {
	dbPath := filepath.Join(db.dataDir, label.String())
	db.logger.Info("Opening Database", "label", label.String(), "path", dbPath)

	var err error
	db.db, err = mdbx.NewMDBX(db.logger).
		Label(label).
		Path(dbPath).
		WithTableCfg(func(_ kv.TableCfg) kv.TableCfg { return tableCfg }).
		MapSize(16 * datasize.GB).
		GrowthStep(16 * datasize.MB).
		Open(ctx)
	return err
}

func (db *Database) OpenOnce(ctx context.Context, label kv.Label, tableCfg kv.TableCfg) error {
	var err error
	db.openOnce.Do(func() {
		err = db.open(ctx, label, tableCfg)
	})
	return err
}

func (db *Database) Close() {
	if db.db != nil {
		db.db.Close()
	}
}

func (db *Database) BeginRo(ctx context.Context) (kv.Tx, error) {
	return db.db.BeginRo(ctx)
}

func (db *Database) BeginRw(ctx context.Context) (kv.RwTx, error) {
	return db.db.BeginRw(ctx)
}
