package main

type DateTime struct {
	time [15]rune
}

type Email struct {
	Email string
}

type FrequencySettings struct {
	Settings [20]rune
}

type Account struct {
	ID    int8
	Email Email
}

type Illness struct {
	ID    int8
	Title [100]rune
}

type Procedure struct {
	ID    int8
	Title [100]rune
}

type Drug struct {
	ID    int8
	Title [100]rune
}

type Scheme struct {
	ID                int8
	Drug              Drug
	Procedure         Procedure
	FrequencySettings FrequencySettings
}

type TreatmentScheme struct {
	ID      int8
	Illness Illness
	Schemes []Scheme
}

type Treatment struct {
	ID              int8
	Account         Account
	TreatmentScheme TreatmentScheme
	BegunAt         DateTime
	EndedAt         DateTime
}

type Schedule struct {
	ID        int8
	Treatment Treatment
	Procedure Procedure
	Drug      Drug
	PlannedAt DateTime
}

type Notification struct {
	ID         int8
	Schedule   Schedule
	NotifiedAt DateTime
}

type MedicationSchedule struct {
	ID                 int8
	Schedule           Schedule
	TakingMedicationAt DateTime
}
