package tables

import (
	"github.com/Masterminds/squirrel"
	"github.com/mohammedaouamri5/JDM-back/db"
	"github.com/sirupsen/logrus"
)

type Packet struct {
	ID_Packet       int32
	Start           int32
	End             int32
	ID_Download     int32
	ID_Packet_State int32
}

type Packets []Packet
type Chunck []Packets

// TODO : creating the tree of packets

func (me Packets) NULL() Packet {
	return Packet{
		ID_Packet:       0,
		Start:           0,
		End:             0,
		ID_Download:     0,
		ID_Packet_State: 0,
	}
}
func (me Packet) IsNULL() bool {
	return me.ID_Packet == 0 && me.Start == 0 && me.End == 0 && me.ID_Download == 0 && me.ID_Packet_State == 0
}
func (me Packets) Select(limit int8, ID_Download int) (Packets, error) {
	// Ensure that the ID_State is retrieved correctly
	state := State{}.GET("dow")

	sql, args, err := squirrel.Select(
		"ID_Packet",
		"Start",
		"End",
		"ID_Packet_State",
		"ID_Download",
	).From("Packet").Where(squirrel.Eq{
		"ID_Packet_State": state.ID_State,
		"ID_Download":     ID_Download,
	}).Limit(uint64(limit)).ToSql()

	logrus.Info(sql, args)
	if err != nil {
		logrus.Error("Error building SQL query: ", err)
		return nil, err
	}

	result, err := db.DB().Query(sql, args...)
	if err != nil {
		logrus.Error("Error executing query: ", err)
		return nil, err
	}
	defer result.Close()

	var packets = make(Packets, 0)
	for result.Next() {
		var packet Packet
		if err := result.Scan(
			&packet.ID_Packet,
			&packet.Start,
			&packet.End,
			&packet.ID_Packet_State,
			&packet.ID_Download,
		); err != nil {
			logrus.Error("Error scanning result: ", err)
			return nil, err
		}
		packets = append(packets, packet)
	}

	if err := result.Err(); err != nil {
		logrus.Error("Error with result iteration: ", err)
		return nil, err
	}

	return packets, nil
}
