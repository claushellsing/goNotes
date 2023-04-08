package main

import (
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

var defaultSubjectID = uint(1) //first subject

var rootCmd = &cobra.Command{
	Use: "note",
}

func initDB() (*gorm.DB, error) {
	var count int

	db, err := gorm.Open(sqlite.Open("notes.db"), &gorm.Config{})
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

	return db, err
}

func main() {

	db, _ := initDB()

	cmdAddNote := &cobra.Command{
		Use:   "add [note]",
		Short: "Add a new note",
		Run: func(cmd *cobra.Command, args []string) {
			noteText := args[0]
			note := Note{
				Text:      noteText,
				SubjectID: defaultSubjectID,
			}

			db.Create(&note)
		},
	}

	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdAddNote)
	rootCmd.Execute()
}
