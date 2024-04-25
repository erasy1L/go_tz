package models

type CarRequest struct {
	RegNum string        `json:"regNum,omitempty" example:"AA1234AA"`
	Mark   string        `json:"mark,omitempty" example:"Toyota"`
	Model  string        `json:"model,omitempty" example:"Corolla"`
	Year   int           `json:"year,omitempty" example:"2010"`
	Owner  PersonRequest `json:"owner,omitempty"`
} // @name Car

type CarResponse struct {
	RegNum string         `json:"regNum,omitempty" example:"AA1234AA"`
	Mark   string         `json:"mark,omitempty" example:"Toyota"`
	Model  string         `json:"model,omitempty" example:"Corolla"`
	Year   int            `json:"year,omitempty" example:"2010"`
	Owner  PersonResponse `json:"owner,omitempty"`
	ID     string         `json:"id,omitempty"`
} // @name CarResponse

// CarFilter = "make" | "model" | "owner" | "reg_num"
type CarFilter string

var ID CarFilter = "id"
var RegNum CarFilter = "reg_number"
var Mark CarFilter = "mark"
var Model CarFilter = "model"
var Year CarFilter = "year"
var Owner CarFilter = "owner"

func (c *CarResponse) ValueToUpdate() map[CarFilter]interface{} {
	values := make(map[CarFilter]interface{})
	if c.RegNum != "" {
		values[RegNum] = c.RegNum
	}
	if c.Mark != "" {
		values[Mark] = c.Mark
	}
	if c.Model != "" {
		values[Model] = c.Model
	}
	if c.Year != 0 {
		values[Year] = c.Year
	}
	if c.Owner != (PersonResponse{}) {
		values[Owner] = c.Owner
	}
	return values
}
