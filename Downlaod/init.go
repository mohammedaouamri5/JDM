package downlaod

import (
	"errors"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/fatih/structs"
	http "github.com/mohammedaouamri5/JDM-back/Downlaod/HTTP"
	"github.com/mohammedaouamri5/JDM-back/db"
	"github.com/mohammedaouamri5/JDM-back/tables"
	"github.com/sirupsen/logrus"
)

var downloding = []int{}

func isStarted(p_id int) bool {
	for _, val := range downloding {
		if val == p_id {
			return true
		}
	}
	return false
}

func DownloadLiciner() error {

	for ; true; time.Sleep(time.Second / 10) {
		sql, args, err := squirrel.
			Select(
				"ID_Download",
				"ID_Download_Type",
				"ID_Download_State",
				"ID_File_Type",
				"Working_file_path",
				"Output_file_path",
				"Remote",
			).
			From("Download").
			Where(squirrel.Eq{"ID_Download_State": tables.State{}.GET("down").ID_State}). // Ensure the WHERE clause is valid

			ToSql()
		if err != nil {
			logrus.Error(err.Error())
			return err
		}

		result, err := db.DB().Query(sql, args...)
		if err != nil {
			logrus.Error(err.Error())
			return err
		}
		defer result.Close() // Always close the result set

		// Check if there are any rows

		for result.Next() {
			var row tables.Downlaod
			if err = result.Scan(
				&row.IdDownlaod,
				&row.IdDownlaodType,
				&row.IdDownlaodStatus,
				&row.IdFileType,
				&row.WorkingFilePath,
				&row.OutputFilePath,
				&row.Remote,
			); err != nil {
				logrus.Error(err.Error())
				return err
			}
			if !isStarted(row.IdDownlaod) {
				logrus.WithFields(structs.Map(row)).Info()
				go http.Downlaod(row)
				downloding = append(downloding, row.IdDownlaod)
			}
		}

		// Check if there is any error in the iteration
		if err = result.Err(); err != nil {
			logrus.Error(err.Error())
			return err
		}
	}
	return nil
}

func popByValue(slice []int, value int) ([]int, int, error) {
    // Iterate over the slice to find the index of the value
    index := -1
    for i, v := range slice {
        if v == value {
            index = i
            break
        }
    }
    
    // Check if the value was found
    if index == -1 {
        return slice, 0, errors.New("value not found in slice") // Handle value not found
    }
    
    // Remove the value from the slice
    slice = append(slice[:index], slice[index+1:]...)
    
    // Return the modified slice and the removed value
    return slice, value, nil
}
func Pause(id int) {
	if isStarted(id) {
		sql, args, err := squirrel.
			Update("Download").Set("ID_Download_State" , tables.State{}.GET("paus").ID_State).Where(squirrel.Eq{"ID_Download" : id }).ToSql()
		_ , err  = db.DB().Exec(sql , args...)
		if err != nil {
			logrus.Fatal(err.Error())
		}
		downloding , _ , _ = popByValue(downloding,id)
	} else {

		sql, args, err := squirrel.
			Update("Download").Set("ID_Download_State" , tables.State{}.GET("dow").ID_State).Where(squirrel.Eq{"ID_Download" : id }).ToSql()
		_ , err  = db.DB().Exec(sql , args...)
		if err != nil {
			logrus.Fatal(err.Error())
		}
	}
}
