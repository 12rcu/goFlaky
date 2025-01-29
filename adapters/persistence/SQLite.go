package persistence

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"goFlaky/core/framework"
	"strings"
)

type SQLite struct {
	db *sql.DB
}

func CreateSQLiteConnection() (SQLite, error) {
	const file string = "default.db"
	db, err := sql.Open("sqlite3", file)
	return SQLite{db}, err
}

// CreateNewRun Create a new run in the db and return the run id
func (s SQLite) CreateNewRun(projects []string) (int, error) {
	query := `INSERT INTO runs (start_time, projects) VALUES (CURRENT_TIMESTAMP, ?)`
	res, err := s.db.Exec(query, strings.Join(projects, ","))
	if err != nil {
		return 0, err
	}
	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return 0, err
	}
	return int(id), nil
}

func (s SQLite) CreateTestResult(runId int, runType string, project string, result framework.TestResult, testOrder string) error {
	t := "PRE_RUN"
	if testOrder == "" {
		t = "OD_RUN"
	}
	query := `INSERT INTO test_results(run_id, run_type, project, test_suite, test_id, result, test_order) VALUES (?,?,?,?,?,?,?)`
	_, err := s.db.Exec(query, runId, t, project, result.TestSuite, result.TestName, result.TestOutcome, testOrder)
	if err != nil {
		return err
	}
	return nil
}

func (s SQLite) CreateClassification(runId int, project string, testSuite string, testName string, classification string) error {
	query := `INSERT INTO test_classifications(run_id, project, test_suite, test_id, classification) VALUES (?,?,?,?,?)`
	_, err := s.db.Exec(query, runId, project, testSuite, testName, classification)
	if err != nil {
		return err
	}
	return nil
}

func (s SQLite) GetProjectTestResults(runId int, project string) ([]framework.TestResult, error) {
	query := `SELECT test_suite, test_id, result FROM test_results WHERE project=? AND run_id=?`
	rows, err := s.db.Query(query, project, runId)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	var results []framework.TestResult
	for rows.Next() {
		result := framework.TestResult{}
		err := rows.Scan(&result.TestSuite, &result.TestName, &result.TestOutcome)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}
	return results, nil
}
