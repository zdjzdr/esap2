package models

type Models interface {
	Insert(int) error
	Update() error
	Delete() error
}
