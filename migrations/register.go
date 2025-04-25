package migrations

import (
	"os"
	"path/filepath"
)

// IMPORTANT INFO: when you wanna add a new table or sql file, you should add it to migration list for up and down
func init() {
	migrations := []Migration{
		{
			Version: "000001",
			Up:      readSQLFile("000001_create_enums.sql"),
			Down: `
                DROP TYPE IF EXISTS user_role;
                DROP TYPE IF EXISTS user_status;
            `,
		},
		{
			Version: "000002",
			Up:      readSQLFile("000002_create_users.sql"),
			Down: `
                DROP TRIGGER IF EXISTS update_users_updated_at ON users;
                DROP FUNCTION IF EXISTS update_updated_at_column();
                DROP TABLE IF EXISTS users CASCADE;
            `,
		},
		{
			Version: "000003",
			Up:      readSQLFile("000003_create_auth_tables.sql"),
			Down: `
                DROP TABLE IF EXISTS sessions CASCADE;
                DROP TABLE IF EXISTS token_blacklists CASCADE;
                DROP TABLE IF EXISTS tokens CASCADE;
            `,
		},
	}

	Migrations = append(Migrations, migrations...)
}

func readSQLFile(filename string) string {
	content, err := os.ReadFile(filepath.Join("migrations", "sql", filename))
	if err != nil {
		panic(err)
	}
	return string(content)
}
