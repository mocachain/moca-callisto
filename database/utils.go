package database

import (
	"fmt"
	"strings"
)

// InsertEnableModules allows to save enabled module into the database
func (db *DB) InsertEnableModules(modules []string) error {
	if len(modules) == 0 {
		return nil
	}

	// Remove existing modules
	stmt := "DELETE FROM modules WHERE TRUE"
	_, err := db.SQL.Exec(stmt)
	if err != nil {
		return fmt.Errorf("error while deleting modules: %s", err)
	}

	if len(modules) == 0 {
		return nil
	}

	var values []interface{}
	stmt = `INSERT INTO modules (module_name) VALUES`
	for key, value := range modules {
		stmt += fmt.Sprintf("($%d),", key+1)
		values = append(values, value)
	}
	stmt = stmt[:len(stmt)-1] // remove tailing ","
	stmt += " ON CONFLICT DO NOTHING"
	_, err = db.SQL.Exec(stmt, values...)
	if err != nil {
		return fmt.Errorf("error while storing modules: %s", err)
	}

	return nil
}

func ExtractEvent(event string) string {
	lastDotIndex := strings.LastIndex(event, ".")
	if lastDotIndex == -1 {
		return ""
	}
	afterDot := event[lastDotIndex+1:]
	prefix := "Event"
	if strings.HasPrefix(afterDot, prefix) {
		return afterDot[len(prefix):]
	}
	return afterDot
}
