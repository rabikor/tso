package server

import (
	"github.com/labstack/echo/v4"
	"treatment-scheme-organizer/database"
)

type createDrugRequest struct {
	Drug struct {
		Title string `json:"title" validate:"required"`
	} `json:"drug" validate:"required"`
}

func (r *createDrugRequest) bind(c echo.Context, d *database.Drug) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	d.Title = r.Drug.Title

	return nil
}

func (r *createDrugRequest) populate(d database.Drug) {
	r.Drug.Title = d.Title
}

type createIllnessRequest struct {
	Illness struct {
		Title string `json:"title" validate:"required"`
	} `json:"illness" validate:"required"`
}

func (r *createIllnessRequest) bind(c echo.Context, i *database.Illness) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	i.Title = r.Illness.Title

	return nil
}

func (r *createIllnessRequest) populate(i database.Illness) {
	r.Illness.Title = i.Title
}

type createProcedureRequest struct {
	Procedure struct {
		Title string `json:"title" validate:"required"`
	} `json:"procedure" validate:"required"`
}

func (r *createProcedureRequest) bind(c echo.Context, p *database.Procedure) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	p.Title = r.Procedure.Title

	return nil
}

func (r *createProcedureRequest) populate(p database.Procedure) {
	r.Procedure.Title = p.Title
}

type schemeDayData struct {
	DrugID      uint `json:"drugId" validate:"required"`
	ProcedureID uint `json:"procedureId" validate:"required"`
	Order       uint `json:"order" validate:"required"`
	Times       uint `json:"times" validate:"required"`
	Frequency   uint `json:"frequency" validate:"required"`
}

type createSchemeDayRequest struct {
	SchemeDay schemeDayData `json:"schemeDay" validate:"required"`
}

func (r *createSchemeDayRequest) bind(c echo.Context, sd *database.SchemeDay) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	sd.DrugID = r.SchemeDay.DrugID
	sd.ProcedureID = r.SchemeDay.ProcedureID
	sd.Order = r.SchemeDay.Order
	sd.Times = r.SchemeDay.Times
	sd.Frequency = r.SchemeDay.Frequency

	return nil
}

func (r *createSchemeDayRequest) populate(sd database.SchemeDay) {
	r.SchemeDay.DrugID = sd.Drug.ID
	r.SchemeDay.ProcedureID = sd.Procedure.ID
	r.SchemeDay.Order = sd.Order
	r.SchemeDay.Times = sd.Times
	r.SchemeDay.Frequency = sd.Frequency
}

type createSchemeRequest struct {
	Scheme struct {
		IllnessID uint `json:"illness" validate:"required"`
		Length    uint `json:"length" validate:"required"`
	} `json:"scheme" validate:"required"`
	Days []schemeDayData `json:"days" validate:"dive,required"`
}

func (r *createSchemeRequest) bind(c echo.Context, s *database.Scheme) error {
	if err := c.Bind(r); err != nil {
		return err
	}

	if err := c.Validate(r); err != nil {
		return err
	}

	s.IllnessID = r.Scheme.IllnessID
	s.Length = r.Scheme.Length

	for _, sd := range r.Days {
		s.Days = append(
			s.Days,
			database.SchemeDay{
				Scheme:      *s,
				SchemeID:    s.ID,
				ProcedureID: sd.ProcedureID,
				DrugID:      sd.DrugID,
				Order:       sd.Order,
				Times:       sd.Times,
				Frequency:   sd.Frequency,
			},
		)
	}

	return nil
}

func (r *createSchemeRequest) populate(s database.Scheme) {
	r.Scheme.IllnessID = s.Illness.ID
	r.Scheme.Length = s.Length
}
