package tables

import (
	"github.com/Masterminds/squirrel"
	// "github.com/fatih/structs"
	"github.com/mohammedaouamri5/JDM-back/db"
	"github.com/sahilm/fuzzy"
	"github.com/sirupsen/logrus"
)

type State struct {
	ID_State int8
	Name     string
}

type States []State

func (e States) String(i int) string {
	return e[i].Name
}

func (e States) Len() int {
	return len(e)
}

var states = make(States, 0)

func (State) Pull() error {

	sql, args, err := squirrel.Select("*").From("STATE").ToSql()

	if err != nil {
		logrus.Fatal(err.Error())
		return err
	}

	rows, err := db.DB().Query(sql, args...)
	defer rows.Close()

	for rows.Next() {
		var row State
		if err := rows.Scan(
			&row.ID_State,
			&row.Name,
		); err != nil {
			logrus.Fatal(err.Error())
			return err
		}
		states = append(states, row)
	}

	logrus.Info(states)
	return nil
}

func (State) GET(__str string) State {
	matches := fuzzy.FindFrom(__str, states)
	return states[matches[0].Index]
}
