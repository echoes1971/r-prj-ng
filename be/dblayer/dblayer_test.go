package dblayer

import (
	"log"
	"testing"
)

func TestCreateTableStrings(t *testing.T) {
	factory := NewDBEFactory(true)
	user := NewDBUser()
	factory.Register(user)
	group := NewDBGroup()
	factory.Register(group)
	userGroup := NewUserGroup()
	factory.Register(userGroup)

	// Loop through all registered entities and print their CREATE TABLE strings
	for _, className := range factory.GetAllClassNames() {
		entity := factory.GetInstanceByClassName(className)
		createTableSQL := entity.GetCreateTableSQL("rprj")
		log.Printf("CREATE TABLE SQL for %s:\n%s\n", className, createTableSQL)
	}
}
func TestEnsureDBSchema(t *testing.T) {

	InitDBLayer("mysql", "root:mysecret@tcp(localhost:3306)/rproject", "rprj")

	// Call EnsureDBSchema to test table creation logic
	EnsureDBSchema()
}
