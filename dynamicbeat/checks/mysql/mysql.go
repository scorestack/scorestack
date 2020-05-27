package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"strconv"

	// MySQL driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/s-newman/scorestack/dynamicbeat/checks/schema"
)

// The Definition configures the behavior of the MySQL check
// it implements the "check" interface
type Definition struct {
	Config       schema.CheckConfig // generic metadata about the check
	Host         string             `optiontype:"required"`                      // IP of Hostname for the MySQL server
	Username     string             `optiontype:"required"`                      // Username for the database
	Password     string             `optiontype:"required"`                      // Password for the user
	Database     string             `optiontype:"required"`                      // Name of the database to access
	Table        string             `optiontype:"required"`                      // Name of the table to access
	Column       string             `optiontype:"required"`                      // Name of the column to access
	MatchContent string             `optiontype:"optional"`                      // Whether to perform a regex content match on the results of the query
	ContentRegex string             `optiontype:"optional" optiondefault:".*"`   // Regex to match on
	Port         string             `optiontype:"optional" optiondefault:"3306"` // Port for the server
}

// Run a single instance of the check
func (d *Definition) Run(ctx context.Context) schema.CheckResult {
	// Initialize empty result
	result := schema.CheckResult{}

	// Create DB handle
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", d.Username, d.Password, d.Host, d.Port, d.Database))
	if err != nil {
		result.Message = fmt.Sprintf("Creating database handle failed : %s", err)
		return result
	}
	defer db.Close()

	// Set connection parameters
	db.SetMaxIdleConns(-1)
	db.SetMaxOpenConns(1)

	// Check db connection
	err = db.PingContext(ctx)
	if err != nil {
		result.Message = fmt.Sprintf("Failed to ping database : %s", err)
	}

	// Query the DB
	// TODO: This is SQL injectable. Figure out Paramerterized queries
	rows, err := db.QueryContext(ctx, fmt.Sprintf("SELECT %s FROM %s;", d.Column, d.Table))
	if err != nil {
		result.Message = fmt.Sprintf("Could not query database : %s", err)
		return result
	}
	defer rows.Close()

	// Store the value from the column
	var val string

	// Perform regex matching, if necessary
	if matchContent, _ := strconv.ParseBool(d.MatchContent); matchContent {
		// Compile the regex
		regex, err := regexp.Compile(d.ContentRegex)
		if err != nil {
			result.Message = fmt.Sprintf("Error compiling regex string %s : %s", d.ContentRegex, err)
			return result
		}

		// Check the rows
		for rows.Next() {
			// Grab a value
			err := rows.Scan(&val)
			if err != nil {
				result.Message = fmt.Sprintf("Could not scan row values : %s", err)
				return result
			}
			// Check value with regex
			if regex.MatchString(val) {
				// If we reach here the check passes
				result.Passed = true
				return result
			}

		}
	}

	// Check for error in the rows
	if rows.Err() != nil {
		result.Message = fmt.Sprintf("Something happened to the rows : %s", err)
		return result
	}

	// Check fails if we reach here
	return result
}

// GetConfig returns the current CheckConfig struct this check has been
// configured with
func (d *Definition) GetConfig() schema.CheckConfig {
	return d.Config
}

// SetConfig reconfigures this check with a new CheckConfig struct
func (d *Definition) SetConfig(config schema.CheckConfig) {
	d.Config = config
}
