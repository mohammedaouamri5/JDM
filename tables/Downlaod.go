package tables

import (
	"context"
	"errors"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/mattn/go-sqlite3"
	"github.com/mohammedaouamri5/JDM-back/db"
	. "github.com/mohammedaouamri5/JDM-back/db"
	"github.com/sirupsen/logrus"
)

var (
	ECantGetHeader    = errors.New("Cant Get The Header")
	ECantCreatePacket = errors.New("Cant Get The Header")
)

type Downlaod struct {
	IdDownlaod       int    `json:"id-downlaod"`
	IdDownlaodType   int    `json:"id-downlaod-type"`
	IdDownlaodStatus int    `json:"id-downlaod-status"`
	IdFileType       int    `json:"id-file-type"`
	WorkingFilePath  string `json:"working-file-type"`
	OutputFilePath   string `json:"output-file-path"`
	Remote           string `json:"remote"`
}

// FIXME : The New function just handele the cas where the pretocol is ```http``` and can't determend the type Yet
func (me *Downlaod) New(p_remote string, p_working_dir *string, p_output_dir *string) Downlaod {
	me.IdDownlaodStatus = 1
	me.Remote = p_remote

	if p_output_dir == nil {
		p_output_dir = &Settings().Output_dir
	}
	if p_working_dir == nil {
		p_working_dir = &Settings().Working_dir
	}

	header, err := getTheHead(p_remote)
	if errors.Is(err, ECantGetHeader) {
		logrus.Error(err.Error())
		// TODO : Downlaod directory
	}

	// FIXME: get the actule type
	var FileType, DownlaodType string
	me.IdFileType, FileType = 1, "undefined"
	me.IdDownlaodType, DownlaodType = 1, "http"
	_, _ = FileType, DownlaodType

	if err != nil {
		logrus.Error(err.Error())
	}

	name, err := getTheName(p_remote, header)
	logrus.Info("the name is : ", name)

	me.WorkingFilePath = (*p_working_dir) + "/" + FileType + "/" + name + ".work"
	me.OutputFilePath = (*p_output_dir) + "/" + FileType + "/" + name

	/*
		CREATE TABLE IF NOT EXISTS Download (
		ID_Download INTEGER PRIMARY KEY AUTOINCREMENT,
		ID_Download_Type INTEGER,
		ID_Download_State INTEGER,
		ID_File_Type INTEGER,
		Working_file_path TEXT,
		Output_file_path TEXT,
		Remote TEXT,
		FOREIGN KEY (ID_Download_Type) REFERENCES DownloadType(ID_Download_Type),
		FOREIGN KEY (ID_Download_State) REFERENCES State(ID_State),
		FOREIGN KEY (ID_File_Type) REFERENCES FileType(ID_File_Type)
		);
	*/
	sql, args, err := sq.
		Insert("Download").
		Columns(
			"ID_Download_Type",
			"ID_Download_State",
			"ID_File_Type",
			"Working_file_path",
			"Output_file_path",
			"Remote").
		Values(
			me.IdDownlaodType,
			me.IdDownlaodStatus,
			me.IdFileType,
			me.WorkingFilePath,
			me.OutputFilePath,
			me.Remote).
		Suffix("RETURNING ID_Download").
		ToSql()

	err = DB().QueryRow(sql, args...).Scan(&me.IdDownlaod)

	if err != nil {
		logrus.Fatalf("Failed to insert and return ID: %v", err)
	} else {
		logrus.Info("The Id : ", me.IdDownlaod)
	}

	DB().Exec(sql, args...)
	return *me
}

func (me *Downlaod) Init() error {
	if err := _MkDirForFile(me.WorkingFilePath); err != nil {
		logrus.Fatal(err.Error())
		return err
	}

	if err := _MkDirForFile(me.OutputFilePath); err != nil {
		logrus.Fatal(err.Error())
		return err
	}

	if _, err := os.Create(me.WorkingFilePath); err != nil {
		logrus.Error(err.Error())
		return err
	}

	sql, args, err := sq.
		Update("Download").
		Set("ID_Download_State",
			sq.Select("ID_State").
				From("State").
				Where("Name LIKE ?", "%downl__d%"),
		).
		Where(sq.Eq{"ID_Download": me.IdDownlaod}).
		ToSql()

	if err != nil {
		logrus.Fatal(err.Error())
		return err
	}

	_, err = DB().Exec(sql, args...)

	if err != nil {
		logrus.Fatal(err.Error())
		return err
	}

	me.creat_packages_for_http_protochoe()
	return nil

}

func _MkDirForFile(p_file_path string) error {
	p_file_pathSplited := strings.Split(p_file_path, "/")
	WorkingiDirPath := p_file_pathSplited[:len(p_file_pathSplited)-1]

	return os.MkdirAll(strings.Join(WorkingiDirPath, "/"), 0777)
}

var GetTheHead = getTheHead

func (me *Downlaod) creat_packages_for_http_protochoe() error {
	header, err := GetTheHead(me.Remote)

	if err != nil {
		logrus.Error(err.Error())
		return err
	}

	ranges, size := header.Get("Accept-Ranges"), header.Get("Content-Length")
	logrus.WithFields(logrus.Fields{"range": ranges, "size": size}).Info()

	if ranges == "" || size == "" {
		logrus.Error(ECantCreatePacket.Error())
		return ECantCreatePacket
	}
	logrus.Infof("%++v", Settings())
	size_int, err := strconv.Atoi(size)

	queryBuilder := sq.Insert("Packet").
		Columns("Start", "End", "ID_Packet_State", "ID_Download")
	/*
		CREATE TABLE IF NOT EXISTS Packet (
			ID_Packet INTEGER PRIMARY KEY AUTOINCREMENT,
			Start INTEGER,
			End INTEGER,
			ID_Packet_State INTEGER,
			FOREIGN KEY (ID_Packet_State) REFERENCES State(ID_State)
		);
	*/

	downloading_id := State{}.GET("dow").ID_State
	for i := 0; i < size_int; i += Settings().PacketSize {
		queryBuilder = queryBuilder.Values(i, min(i-1+Settings().PacketSize, size_int), downloading_id, me.IdDownlaod)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return err
	}

	_, err = db.DB().Exec(query, args...)

	if err != nil {
		logrus.Info(err.Error())
		return err
	}

	return err
}

func getTheHead(p_remote string) (http.Header, error) {
	// See If The Head Is Supported
	client := &http.Client{}

	// Create a new HEAD request
	req, err := http.NewRequest("HEAD", p_remote, nil)
	if err != nil {
		logrus.Println("Error creating request:", err)
		return nil, err
	}

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		logrus.Error("Error executing request:", err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 400 {
		return resp.Header, nil
	}

	logrus.Info("Test The time out :")
	// Set up a context with a timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // Ensure the cancel function is called to free resources

	// Create a new GET request with the context
	req, err = http.NewRequestWithContext(ctx, "GET", p_remote, nil)
	if err != nil {
		logrus.Error("Error creating request:", err)
		return nil, err
	}

	// Add the Range header
	req.Header.Add("Range", "bytes=0-0")

	// Create an HTTP client
	client = &http.Client{}

	// Execute the request
	resp, err = client.Do(req)
	if err != nil {
		// If the context was canceled (due to timeout), handle it here
		if errors.Is(err, context.DeadlineExceeded) {
			logrus.Error("Request timed out")
			return nil, errors.Join(err, ECantGetHeader)
		}
		logrus.Error("Error executing request:", errors.Join(err, ECantGetHeader))
		return nil, err
	}

	// Ensure the response body is closed
	defer resp.Body.Close()

	return resp.Header, nil
}

func getTheName(p_remote string, p_header http.Header) (string, error) {
	// see if the name in on the headerHeader
	if p_header != nil {
		if ContentDisposition := p_header.Get("Content-Disposition"); ContentDisposition != "" {
			ContentDispositionSplites := strings.Split(ContentDisposition, "=")
			result := ContentDispositionSplites[len(ContentDispositionSplites)-1]
			return result, nil
		}
	}
	// if not just get the last part of the URL
	remote_splited := strings.Split(p_remote, "/")
	result := remote_splited[len(remote_splited)-1]
	return result, nil
}
