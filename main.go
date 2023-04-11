package main

import (
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"strconv"
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
		Use:   "add [note] [subject]",
		Short: "Add a new note",
		Run: func(cmd *cobra.Command, args []string) {
			noteText := args[0]
			var note = Note{}

			if len(args) > 1 {
				var subjectId uint
				subjectName := args[1]
				subject := Subject{}

				db.Find(&subject, "name = ?", subjectName)

				if subject.ID == 0 { //No subject found
					newSubject := Subject{
						Name: subjectName,
					}
					db.Create(&newSubject)

					subjectId = newSubject.ID
				} else {
					subjectId = subject.ID
				}

				note = Note{
					Text:      noteText,
					SubjectID: subjectId,
				}
			} else {
				note = Note{
					Text:      noteText,
					SubjectID: defaultSubjectID,
				}
			}

			db.Create(&note)
		},
	}

	cmdListNotes := &cobra.Command{
		Use:   "list [subject]",
		Short: "List all notes",
		Run: func(cmd *cobra.Command, args []string) {
			var notes []Note
			var dtTable pterm.TableData

			db := db.Joins("Subject")
			if len(args) > 0 {
				subjectText := args[0]
				db.Where("Subject.Name = ?", subjectText).Find(&notes)
			}

			db.Find(&notes)

			dtTable = append(dtTable, []string{"ID", "Note", "Subject"})
			for _, note := range notes {
				dtTable = append(dtTable, []string{
					strconv.Itoa(int(note.ID)),
					note.Text,
					note.Subject.Name,
				})
			}

			pterm.DefaultTable.WithHasHeader().WithBoxed().WithData(dtTable).Render()
		},
	}

	rootCmd := &cobra.Command{Use: "app"}
	rootCmd.AddCommand(cmdAddNote)
	rootCmd.AddCommand(cmdListNotes)
	rootCmd.Execute()
}
