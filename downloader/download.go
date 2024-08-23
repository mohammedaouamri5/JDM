package downloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

	"github.com/google/uuid"
	"github.com/mohammedaouamri5/JDM-back/utile"
	log "github.com/sirupsen/logrus"

	"io"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	// "sync"
	"time"
)

const DEFAULT_DOWNLOAD_PATH = "/home/mohammedaouamri/Download/"

type Packet struct {
	Start int  `json:"start"`
	End   int  `json:"end"`
	Done  bool `json:"done"`
}

type FILE struct {
	Id      uuid.UUID
	Output  string `json:"output" binding:"required"`
	Url     string `json:"url" binding:"required"`
	Packets []Packet
}

func Unfiniched(me FILE) string {
	return "./data/" + me.Id.String() + ".unfiniched"
}

func Cfgjson(me FILE) string {
	return "./data/" + me.Id.String() + ".cfg.json"
}

// TODO: change the formate of the MetaData

func (me *FILE) ReadFromMetaData(path string) error {

	file, err := os.Open(path)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	// Read the file's content

	if err != nil {
		return err
	}

	// Unmarshal the JSON data into the Person struct
	err = json.Unmarshal(content, &me)
	if err != nil {
		return err
	}

	return nil
}

// Method to marshal the FILE struct to JSON and write it to a file
func (me *FILE) writeToAMetaData(path string) error {

	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Errorf("Error opening file: %v", err)
		return err
	}
	defer file.Close()

	content, err := json.MarshalIndent(*me, "", "  ")
	if err != nil {
		log.Errorln("\n\t", "Error marshalling JSON:", err.Error())
		return err
	}
	_, err = file.Write(content)
	if err != nil {
		log.Errorln("\n\t", "Error writing to file:", err.Error())
		return err
	}
	log.Infoln("\n\twrite into ", me.Output)
	return nil
}

func (me *FILE) mkeConfig(numThreads int) error {
	log.Infoln("\n\tmkeConfig")
	cfg, err := os.Create(Cfgjson(*me))

	if err != nil {
		log.Errorln("\n\terror creating output file: ", err.Error())
		return err
	}

	Unf, err := os.Create(Unfiniched(*me))

	if err != nil {
		log.Errorln("\n\terror creating output file: ", err.Error())
		return err

	}
	if err := Unf.Close(); err != nil {
		log.Errorln(err.Error())
		return err
	}

	resp, err := http.Head(me.Url)
	if err != nil {
		log.Errorln("\n\terror getting head request: ", err.Error())
		return err

	}
	defer resp.Body.Close()

	if resp.Header.Get("Accept-Ranges") != "bytes" {
		log.Errorln("\n\tserver does not support range requests")
		return errors.New("server does not support range requests")

	}

	contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	if err != nil {
		log.Errorln("\n\tfailed to get content length: ", err.Error())
		return err
	}

	// 	segmentSize :=    int(math.Ceil(float64(contentLength) / float64(numThreads)))

	segmentSize := 1_000_000
	segmentNb_ := int(math.Ceil(math.Ceil(float64(contentLength) / float64(segmentSize))))

	for i := 0; i < segmentNb_; i++ {
		start := segmentSize * i
		end := start + segmentSize - 1
		if i == segmentNb_-1 {
			end = contentLength - 1
		}
		me.Packets = append(me.Packets, Packet{Start: start, End: end, Done: false})
	}

	log.Infoln("The Packes Are Created", len(me.Packets))

	if err := me.writeToAMetaData(Cfgjson(*me)); err != nil {
		log.Errorln("\n\t", err.Error())
	}
	if err := cfg.Close(); err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}

// Constructor initializes the FILE struct
func (me *FILE) Constructor(url_p string, name_p string, path_p *string) error {
	if !strings.HasPrefix(url_p, "http://") && !strings.HasPrefix(url_p, "https://") {
		return errors.New("invalid URL")
	}
	me.Url = url_p
	if path_p == nil {
		me.Output = utile.InfoS.PATH + name_p
	} else {
		me.Output = (*path_p) + "/" + name_p
	}
	// Ensure the directory exists
	dir := strings.TrimSuffix(me.Output, "/"+name_p)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}
	me.Id = uuid.New()
	return nil
}

func (me *FILE) finishIt() error {

	os.Rename(Unfiniched(*me), me.Output)
	// os.Remove(Cfgjson(*me))
	return nil
}

// func (me *FILE) downloadRange(out *os.File, start, end int, wg *sync.WaitGroup, retries int) {
func (me *FILE) downloadRange(start, end int) error {

	log.Infoln(fmt.Sprintf("start download range %d -> %d", start, end))
	out, err := os.OpenFile(Unfiniched(*me), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer out.Close()
	retries := 5
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	for attempts := 0; attempts <= retries; attempts++ {
		req, err := http.NewRequest("GET", me.Url, nil)
		if err != nil {
			log.Println("Error creating request:", err)
			return errors.New(fmt.Sprint("Error creating request:", err))
		}
		rangeHeader := fmt.Sprintf("bytes=%d-%d", start, end)
		req.Header.Set("Range", rangeHeader)
		resp, err := client.Do(req)
		if err != nil {
			log.Warn(fmt.Sprintf("Error executing request (attempt %d/%d): %v\n", attempts+1, retries, err))
			time.Sleep(time.Second * time.Duration(attempts+1+utile.RandomIntByRange(-5, 5)))
			continue
		}

		defer resp.Body.Close()
		if resp.StatusCode != http.StatusPartialContent {
			log.Printf("Unexpected status code: %d\n", resp.StatusCode)
			return errors.New(fmt.Sprintf("Unexpected status code: %d", resp.StatusCode))
		}
		if _, err := out.Seek(int64(start), io.SeekStart); err != nil {
			log.Println("Error seeking in FILE:", err)
			return errors.New(fmt.Sprint("Error seeking in FILE:", err))
		}

		max_i := 5
		done := false
		for i := 0; i < max_i; i++ {
			if _, err := io.Copy(out, resp.Body); err != nil {
				attempts := i + 1
				waiting_time := (time.Second * time.Duration(attempts+utile.RandomIntByRange(-max_i, max_i)))
				log.Infoln("waiting", waiting_time)
				time.Sleep(waiting_time)
				log.Warnln(fmt.Sprintf("Error copying data (attempt %d/%d): %v\n", attempts, max_i, err))
			} else {
				i = max_i + 2
				done = true
			}
		}
		if done {
			log.Infoln(fmt.Sprintf("Downloaded bytes %d-%d\n", start, end))
		} else {
			log.Trace(fmt.Sprintf("can't download bytes %d-%d\n", start, end))
		}
		return nil
	}
	errmsg := fmt.Sprintf("Failed to download bytes %d-%d after %d attempts\n", start, end, retries+1)
	log.Errorln(err)
	return errors.New(errmsg)
}

func (me *FILE) downloadchunk(chunk []int, wg *sync.WaitGroup) error {
	defer wg.Done()
	log.Infoln(fmt.Sprintf("Downloading chunk of %d lenght ", len(chunk)))
	for _, index := range chunk {
		packet := me.Packets[index]
		if !packet.Done {
			if err := me.downloadRange(packet.Start, packet.End); err != nil {
				log.Errorln(err.Error())
				return err
			} else {
				me.Packets[index].Done = true
				go me.writeToAMetaData(Cfgjson(*me))
			}
		} else {
			log.Warnln("The packe is samhow downloaded")
		}
	}
	return nil
}
func (me *FILE) download(numThreads int) error {

	if err := me.writeToAMetaData(Cfgjson(*me)); err != nil {
		log.Errorln("\n\t", err.Error())
		return err
	}

	ReminderPacket := me.getReminderPacket()

	if len(ReminderPacket) == 0 {
		me.finishIt()
		return nil
	}

	chunks, err := utile.SplitSlice(ReminderPacket, numThreads)

	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	// the actual downlowd
	var wg sync.WaitGroup
	count := 0
	for _, chunk := range chunks {
		if len(chunk) > 0 {
			count++
		}
	}
	wg.Add(utile.Min(numThreads, count ))

	for _, chunk := range chunks {
		go me.downloadchunk(chunk, &wg)
	}
	wg.Wait()
	log.Infoln(14)
	me.writeToAMetaData(Cfgjson(*me))

	return me.finishIt()
}

func (me *FILE) MkeConfig(numThreads int) error {
	if IsExist, err := utile.PathIsExist(Cfgjson(*me)); err != nil {
		log.Errorln("\n\t", err.Error())
		return err
	} else if !IsExist {
		if err := me.mkeConfig(numThreads); err != nil {
			log.Errorln("\n\t", err.Error())
			return err
		}
	}
	me.writeToAMetaData(Cfgjson(*me))
	return nil
}

func (me *FILE) Download(numThreads int) error {

	me.MkeConfig(numThreads)

	if err := me.download(numThreads); err != nil {

		log.Errorln("\n\t", err.Error())

		return err
	}
	return nil
} 
func (me *FILE) getReminderPacket() []int {
	result := make([]int, 0)
	for index, packet := range me.Packets {
		if !packet.Done {
			result = append(result, index)
		}
	}
	return result
}
