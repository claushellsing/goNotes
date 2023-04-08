package main

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
)

type Subject struct {
	ID   uint `gorm:"primarykey"`
	Name string
}

type Note struct {
	ID        uint `gorm:"primarykey"`
	Text      string
	SubjectID uint
	Subject   Subject
}

func main() {
	var count int
	defaultSubjectID := uint(1) //first subject

	if len(os.Args) < 2 {
		fmt.Println("Need a message")
		return
	}

	noteText := os.Args[1]

	db, _ := gorm.Open(sqlite.Open("notes.db"), &gorm.Config{})
	db.
		Raw("SELECT count(*) FROM sqlite_master where type='table'").
		Scan(&count)

	if count < 2 {
		db.AutoMigrate(&Subject{}, &Note{})
		db.Create(&Subject{
			ID:   defaultSubjectID,
			Name: "default",
		})
	}

	note := Note{
		Text:      noteText,
		SubjectID: defaultSubjectID,
	}

	db.Create(&note)
}
