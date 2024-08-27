package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/abc401/digital-logic-simulator/api/helpers"
	"github.com/abc401/digital-logic-simulator/api/types"
	"github.com/abc401/digital-logic-simulator/models"
	"github.com/go-sql-driver/mysql"
	gormDb "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const Dsn = "root:1234@tcp(localhost:3306)/dls"

var cfg = mysql.Config{
	User:      "root",
	Passwd:    "1234",
	Net:       "tcp",
	Addr:      "127.0.0.1:3306",
	DBName:    "dls",
	ParseTime: true,
}

var con *gorm.DB = nil

func GetDBCon() *sql.DB {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		panic(err.Error())
	}
	return db
}

func GetGormDBCon() *gorm.DB {
	var err error
	if con == nil {
		con, err = gorm.Open(gormDb.Open(cfg.FormatDSN()), &gorm.Config{})
		if err != nil {
			panic(err.Error())
		}

	}
	return con
}

func AutoMigrate() {
	db := GetGormDBCon()
	err := db.AutoMigrate(&models.Article{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not auto migrate articles: %s\n", err.Error())
		panic("")
	}
	err = db.AutoMigrate(&models.MCQ{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not auto migrate mcqs: %s\n", err.Error())
		panic("")
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not auto migrate users: %s\n", err.Error())
		panic("")
	}
	err = db.AutoMigrate(&models.MarshaledProject{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not auto migrate projects: %s\n", err.Error())
		panic("")
	}
	fmt.Println("Successfully auto migrated")

}

func ConfigNextAndPrevious(tutorials []*models.Article) {
	for i := 0; i < len(tutorials); i++ {
		if i == 0 {
			tutorials[i].Previous = sql.NullString{}
		} else {
			tutorials[i].Previous = sql.NullString{
				String: tutorials[i-1].LinkTitle,
				Valid:  true,
			}
		}

		if i == len(tutorials)-1 {
			tutorials[i].Next = sql.NullString{}
		} else {
			tutorials[i].Next = sql.NullString{
				String: tutorials[i+1].LinkTitle,
				Valid:  true,
			}

		}
	}

}

func GetUser(con *gorm.DB, email string) *models.User {
	var dest []models.User
	con.Model(&models.User{}).Where("email = ?", email).First(&dest)
	if len(dest) > 0 {
		return &dest[0]
	}
	return nil
}

func GetUserByID(con *gorm.DB, uid uint) *models.User {
	var dest []models.User
	con.Model(&models.User{}).Where("id = ?", uid).First(&dest)
	if len(dest) > 0 {
		return &dest[0]
	}
	return nil
}

func AddTutorials(tutorials []*models.Article) {
	fmt.Println("[Info] Adding tutorials")
	db := GetGormDBCon()
	db.Exec("delete from mcqs")
	db.Exec("delete from articles")

	ConfigNextAndPrevious(tutorials)
	for i := 0; i < len(tutorials); i++ {
		AddTutorial(tutorials[i])
	}
}

func AddTutorial(tutorial *models.Article) {
	if tutorial.LinkTitle == "" {
		panic("Title of tutorial is not defined")
	}
	if tutorial.Content == "" {
		panic("Content of tutorial is not defined")
	}

	fmt.Println("[Info] Creating record in db: ", tutorial.LinkTitle)
	db := GetGormDBCon()
	db.Create(&tutorial)

	content, err := os.ReadFile(fmt.Sprintf("quiz/%s.txt", tutorial.LinkTitle))
	if err != nil {
		content = []byte("")
		fmt.Println("[Error] ", err.Error())
	}

	fmt.Println("Printing Mcqs lines: ", tutorial.LinkTitle)
	lines := strings.Split(string(content), "\n")

	lineIdx := SkipEmptyLines(lines, 0)

	mcqs := []*models.MCQ{}
	for lineIdx < len(lines) {

		mcq := models.MCQ{}

		line := strings.Trim(lines[lineIdx], "\r \t")
		mcq.Statement = line

		lineIdx++
		lineIdx = SkipEmptyLines(lines, lineIdx)

		line = strings.Trim(lines[lineIdx], "\r \t")
		mcq.Option1 = line

		lineIdx++
		lineIdx = SkipEmptyLines(lines, lineIdx)

		line = strings.Trim(lines[lineIdx], "\r \t")
		mcq.Option2 = line

		lineIdx++
		lineIdx = SkipEmptyLines(lines, lineIdx)

		line = strings.Trim(lines[lineIdx], "\r \t")
		mcq.Option3 = line

		lineIdx++
		lineIdx = SkipEmptyLines(lines, lineIdx)

		line = strings.Trim(lines[lineIdx], "\r \t")
		mcq.Option4 = line

		lineIdx++
		lineIdx = SkipEmptyLines(lines, lineIdx)

		line = strings.Trim(lines[lineIdx], "\r \t")
		num, err := strconv.ParseInt(line, 10, 64)
		if err != nil {
			fmt.Printf("[Info] tutorial: %s\n", tutorial.LinkTitle)
			fmt.Printf("[Info] line number: %d\n", lineIdx)
			fmt.Printf("MCQ: %s\n", helpers.SPrettyPrint(mcq))
			panic("")
		}
		mcq.CorrectOption = uint(num)

		lineIdx++
		lineIdx = SkipEmptyLines(lines, lineIdx)

		mcq.ArticleID = tutorial.ID
		mcqs = append(mcqs, &mcq)
	}

	// fmt.Printf("MCQ: %s\n", helpers.SPrettyPrint(mcqs))
	db.Create(mcqs)
	fmt.Printf("MCQ: %s\n", helpers.SPrettyPrint(mcqs))

}

func SkipEmptyLines(lines []string, startLineIdx int) int {
	for startLineIdx < len(lines) {
		line := lines[startLineIdx]
		trimmed := strings.Trim(line, "\r \t")

		if trimmed != "" {
			break
		}
		startLineIdx++
	}
	return startLineIdx
}

func GetMarshaledProject(con *gorm.DB, uid uint, projectID uint) (*models.MarshaledProject, error) {
	var queryResult = []models.MarshaledProject{}
	con.Where("user_id = ? and id = ?", uid, projectID).First(&queryResult)
	if len(queryResult) == 0 {
		return nil, errors.New("Project does not exist")
	}
	return &queryResult[0], nil
}

func GetUnmarshaledProject(con *gorm.DB, uid uint, projectID uint) (*types.ProjectType, error) {
	var marshaledProject, err = GetMarshaledProject(con, uid, projectID)
	if err != nil {
		return nil, err
	}
	return marshaledProject.Unmarshal(), nil
}

type ProjectMetaData struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt`
}

func GetAllProjectsMetaData(uid uint) []ProjectMetaData {
	var con = GetGormDBCon()
	var res []ProjectMetaData
	con.Model(&models.MarshaledProject{}).Select([]string{"id", "name", "created_at"}).Where("user_id = ?", uid).Order("created_at DESC").Find(&res)
	return res
}
