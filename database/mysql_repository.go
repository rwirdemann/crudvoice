package database

type MySQLRepository struct {
}

func NewMySQLRepository() *MySQLRepository {
	r := MySQLRepository{}
	return &r
}
