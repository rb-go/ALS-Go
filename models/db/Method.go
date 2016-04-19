package db

type Method struct {
	Id           int     `gorm:"primary_key"`
	Name         string  `sql:"size:255;not null;unique"` // Default size for string is 255, you could reset it with this tag
	//Users        []User  `gorm:"many2many:user2method;"` // Many-To-Many relationship, 'user_languages' is join table
}

func (c Method) TableName() string {
	return "methods"
}