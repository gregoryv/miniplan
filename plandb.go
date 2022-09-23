package miniplan

import (
	"database/sql"
	"fmt"
	"strings"
)

func NewPlanDB(filename string) (*PlanDB, error) {
	db, err := sql.Open("sqlite3", filename)
	if _, err = db.Exec(changesTbl.CREATE); err != nil {
		if !strings.Contains(err.Error(), "already exists") {
			return nil, err
		}
	}

	mdb := &PlanDB{DB: db}
	stmt, err := db.Prepare(changesTbl.INSERT)
	mdb.InsertChange = stmt

	stmt, err = db.Prepare(changesTbl.DELETE)
	mdb.DeleteChange = stmt

	stmt, err = db.Prepare(changesTbl.UPDATE)
	mdb.UpdateChange = stmt
	return mdb, err
}

type PlanDB struct {
	*sql.DB

	InsertChange *sql.Stmt
	DeleteChange *sql.Stmt
	UpdateChange *sql.Stmt
}

func (me *PlanDB) insert(v interface{}) (err error) {
	switch v := v.(type) {
	case *Change:
		_, err = me.InsertChange.Exec(v.UUID, v.Title, v.Description)
	default:
		err = fmt.Errorf("insert %T", v)
	}
	return

}

var changesTbl = struct {
	CREATE, INSERT, DELETE, UPDATE string
}{
	CREATE: `CREATE TABLE changes (
        uuid VARCHAR(36) NULL,
        title VARCHAR(64) NULL,
        description VARCHAR(2048) NULL
    )`,
	INSERT: "INSERT INTO changes(uuid, title, description) values(?,?,?)",
	DELETE: "DELETE FROM changes WHERE uuid LIKE ?",
	UPDATE: "UPDATE changes SET title = ?, description = ? WHERE uuid LIKE ?",
}
