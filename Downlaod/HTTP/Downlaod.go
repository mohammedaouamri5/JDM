package http

import (
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/mohammedaouamri5/JDM-back/db"
	"github.com/mohammedaouamri5/JDM-back/tables"
	"github.com/sirupsen/logrus"
	"io"
	netHttp "net/http"
	"os"
	"sync"
)

func Downlaod(p_downlaod tables.Downlaod) {
	_, err := tables.GetTheHead(p_downlaod.Remote)
	if errors.Is(err, tables.ECantGetHeader) {
		downlaod_without_range(p_downlaod)
	}
	downlaod_with_range(p_downlaod)
}

var done_id int8

func stilldownloading(p_downlaod tables.Downlaod) bool {
	var state int
	sql, args, _ := squirrel.Select("ID_Download_State").From("Download").Where(squirrel.Eq{"ID_Download": p_downlaod.IdDownlaod}).ToSql()
	db.DB().QueryRow(sql, args...).Scan(&state)
	down_id := int(tables.State{}.GET("dow").ID_State)
	if state !=	 down_id{
	logrus.Info("PAUSE")
	println()
	println()
	println()
	println()
	println()
	println()
	}
	return state == down_id
}
func downlaod_with_range(p_downlaod tables.Downlaod) error {
	NbChunks, ChunkSize := 2, 3
	logrus.WithFields(logrus.Fields{
		"curent id":     p_downlaod.IdDownlaodStatus,
		"download id":   (tables.State{}.GET("dow").ID_State),
		"download name": (tables.State{}.GET("dow").Name),
	}).Info()
	if (stilldownloading(p_downlaod) ) {
		for chunks, err := getBunchOfChunks(p_downlaod.IdDownlaod, NbChunks, ChunkSize); (chunks != nil) && (stilldownloading(p_downlaod)); chunks, err = getBunchOfChunks(p_downlaod.IdDownlaod, NbChunks, ChunkSize) {

			logrus.Info("len = ", len(chunks))
			logrus.Infof("chunk = %++v", (chunks))
			logrus.Infof("chunk = %++v", (chunks))
			done_id = tables.State{}.GET("done").ID_State

			if err != nil {
				logrus.Error(err.Error())
				return err
			}
			var wg sync.WaitGroup
			wg.Add(NbChunks)
			for _, packets := range chunks {
				go download_packets(packets, p_downlaod, &wg)
			}
			wg.Wait()
			logrus.Info("Done waiting")
		}
	}
	return nil
}

func download_packets(p_packets tables.Packets, p_download tables.Downlaod, worker *sync.WaitGroup) {
	defer worker.Done()
	logrus.Info("worder")
	for _, packet := range p_packets {
		if packet.IsNULL() {
			continue
		}
		err := download_packet(packet, p_download)
		if err != nil {
			logrus.Error(err.Error())
		}

		sql, args, err := squirrel.Update("Packet").Set("ID_Packet_State", done_id).Where(squirrel.Eq{
			"ID_Packet": packet.ID_Packet,
		}).ToSql()

		_, err = db.DB().Exec(sql, args...)
		if err != nil {
			logrus.Error(err.Error())
		}
	}
}

func download_packet(p_packet tables.Packet, p_download tables.Downlaod) error {
	// Open the output file in write-only mode
	out, err := os.OpenFile(p_download.WorkingFilePath, os.O_WRONLY, 0777)
	if err != nil {
		logrus.Error("Failed to open output file: ", err)
		return err
	}
	defer out.Close()

	// Create an HTTP request to download the byte range
	client := &netHttp.Client{}
	req, err := netHttp.NewRequest("GET", p_download.Remote, nil)
	if err != nil {
		logrus.Error("Failed to create request: ", err)
		return err
	}

	// Set the Range header for partial download
	rangeHeader := fmt.Sprintf("bytes=%d-%d", p_packet.Start, p_packet.End)
	req.Header.Set("Range", rangeHeader)
	logrus.Info("Requesting byte range: ", rangeHeader)

	// Make the HTTP request
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Failed to make HTTP request: ", err)
		return err
	}
	defer resp.Body.Close()

	// Ensure that the server supports range requests
	if resp.StatusCode != netHttp.StatusPartialContent {
		logrus.Error("Server does not support range requests or invalid range")
		return err
	}

	// Move the file pointer to the appropriate position
	if _, err := out.Seek(int64(p_packet.Start), 0); err != nil {
		logrus.Error("Failed to seek in the output file: ", err)
		return err
	}

	// Copy the downloaded bytes to the file
	written, err := io.Copy(out, resp.Body)
	if err != nil {
		logrus.Error("Failed to write data to file: ", err)
		return err
	}

	logrus.Infof("Downloaded %d bytes for packet (%d) [%d-%d] ", written, p_packet.ID_Packet, p_packet.Start, p_packet.End)
	return nil
}

func getBunchOfChunks(ID_Downlaod int, pNbChunks int, pChunkSize int) (tables.Chunck, error) {
	if pNbChunks <= 0 || pChunkSize <= 0 {
		return nil, fmt.Errorf("both pNbChunks and pChunkSize must be positive integers")
	}

	result := make(tables.Chunck, pNbChunks)
	for i := 0; i < pNbChunks; i++ {
		result[i] = make(tables.Packets, pChunkSize)
	}

	limit := pNbChunks * pChunkSize

	__packets, err := tables.Packets{}.Select(int8(limit), ID_Downlaod)
	if err != nil {
		logrus.Fatal(err.Error())
		return nil, err
	}
	_len := len(__packets)

	if _len == 0 {
		return nil, nil
	}
	for i := 0; i < _len; i++ {
		result[int8(i/pChunkSize)][i%pChunkSize] = __packets[i]
	}
	for i := _len; i < limit; i++ {
		result[int8(i/pChunkSize)][i%pChunkSize] = tables.Packets{}.NULL()
	}
	logrus.WithFields(logrus.Fields{
		"result":    result,
		"__packets": __packets,
		"_len":      _len,
		"limit":     limit,
	}).Info()
	return result, nil
}
func downlaod_without_range(p_downlaod tables.Downlaod) error {
	// Create the file
	out, err := os.OpenFile(p_downlaod.WorkingFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := netHttp.Get(p_downlaod.Remote)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != netHttp.StatusOK {
		err := (fmt.Errorf("bad status: %s", resp.Status))
		return err
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}
