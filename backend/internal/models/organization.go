package models

type Organization struct {
	ID            string
	Name          string
	Address       *string
	City          *string
	State         *string
	ZipCode       *string
	VATNumber     *string
	ContactEmail  *string
	ContactPhone  *string
	PecEmail      *string
	SdiCode       *string
	PersonnelInfo *string
}
