package dblayer

import "testing"

func TestMain(m *testing.M) {
	InitDBLayer("mysql", "root:mysecret@tcp(localhost:3306)/rproject", "rprj")

	// Esegui i test
	m.Run()

	// Teardown: chiudi la connessione
	CloseDBConnection()
}

/* ***** Helper functions for tests ***** */

func hardDeleteForTests(repo *DBRepository, object DBObjectInterface) error {
	deletedObject, err := repo.Delete(object)
	if err != nil {
		return err
	}
	// Second time to force the hard delete
	deletedObject, err = repo.Delete(deletedObject)
	if err != nil {
		return err
	}
	return nil
}
