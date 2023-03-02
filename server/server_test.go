package server

import (
	"os"
	"testing"

	"treatment-scheme-organizer/database"
)

var (
	d *database.DB
)

func TestMain(t *testing.M) {
	setup()
	tearDown()
	code := t.Run()
	os.Exit(code)
}

func setup() {
	d = database.TestDB()
	loadFixtures()
}

func tearDown() {
	database.TruncateTables(d.DB)
}

var SchemeExample *database.Scheme

func loadFixtures() {
	// drugs
	_ = d.Drugs.Add(&database.Drug{Title: "Drug 1"})
	_ = d.Drugs.Add(&database.Drug{Title: "Drug 2"})

	// illnesses
	_ = d.Illnesses.Add(&database.Illness{Title: "Illness 1"})
	_ = d.Illnesses.Add(&database.Illness{Title: "Illness 2"})
	_ = d.Illnesses.Add(&database.Illness{Title: "Illness 3"})

	// procedures
	_ = d.Procedures.Add(&database.Procedure{Title: "Procedure 1"})
	_ = d.Procedures.Add(&database.Procedure{Title: "Procedure 2"})

	// schemes
	SchemeExample = &database.Scheme{IllnessID: 1, Length: 3}
	_ = d.Schemes.Add(SchemeExample)
	_ = d.Schemes.Add(&database.Scheme{IllnessID: 1, Length: 2})

	// scheme days
	_ = d.SchemeDays.Add(&database.SchemeDay{SchemeID: SchemeExample.ID, DrugID: 1, ProcedureID: 1, Order: 1, Times: 1, Frequency: 0})
	_ = d.SchemeDays.Add(&database.SchemeDay{SchemeID: SchemeExample.ID, DrugID: 1, ProcedureID: 1, Order: 2, Times: 2, Frequency: 4})
}
