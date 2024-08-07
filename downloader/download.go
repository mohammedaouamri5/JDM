package downloader

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"sync"

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
	Url     string `json:"url"`
	Output  string `json:"output"`
	Packets []Packet
}

func (me *FILE) readFromJson(path string) error {

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
	log.Infoln(fmt.Sprintf("\n\t%+v", me))

	return nil
}

// Method to marshal the FILE struct to JSON and write it to a file
func (me *FILE) writeToAJson(path string) error {

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

// Constructor initializes the FILE struct
func (me *FILE) Constructor(url_p string, name_p string, path_p *string) error {
	if !strings.HasPrefix(url_p, "http://") && !strings.HasPrefix(url_p, "https://") {
		return errors.New("invalid URL")
	}

	me.Url = url_p

	if path_p == nil {
		me.Output = DEFAULT_DOWNLOAD_PATH + name_p
	} else {
		me.Output = (*path_p) + "/" + name_p
	}

	// Ensure the directory exists
	dir := strings.TrimSuffix(me.Output, "/"+name_p)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("error creating directory: %v", err)
	}

	return nil
}

func (me *FILE) finishIt() error {
	os.Rename(utile.Unfiniched(me.Output), me.Output)
	os.Remove(utile.Cfgjson(me.Output))
	return nil
}

func (me *FILE) mkeConfig(numThreads int) error {
	log.Infoln("\n\tmkeConfig")
	cfg, err := os.Create(utile.Cfgjson(me.Output))

	if err != nil {
		log.Errorln("\n\terror creating output file: ", err.Error())
		return err
	}

	Unf, err := os.Create(utile.Unfiniched(me.Output))

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

	segmentSize := int(math.Ceil(float64(contentLength) / float64(numThreads)))

	for i := 0; i < numThreads; i++ {

		start := segmentSize * i
		end := start + segmentSize - 1
		if i == numThreads-1 {
			end = contentLength - 1
		}
		me.Packets = append(me.Packets, Packet{Start: start, End: end, Done: false})

	}

	if err := me.writeToAJson(utile.Cfgjson(me.Output)); err != nil {
		log.Errorln("\n\t", err.Error())
	}
	if err := cfg.Close(); err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}

// func (me *FILE) downloadRange(out *os.File, start, end int, wg *sync.WaitGroup, retries int) {
func (me *FILE) downloadRange(start, end int) error {

	log.Infoln(fmt.Sprintf("start download range %d -> %d", start, end))
	out, err := os.OpenFile(utile.Unfiniched(me.Output), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)

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
			return errors.New(fmt.Sprint("Unexpected status code: %d\n", resp.StatusCode))
		}

		if _, err := out.Seek(int64(start), io.SeekStart); err != nil {
			log.Println("Error seeking in FILE:", err)
			return errors.New(fmt.Sprint("Error seeking in FILE:", err))
		}

		for i := 0; i < 5; i++ {
			if _, err := io.Copy(out, resp.Body); err != nil {
				attempts := i + 1
				waiting_time := attempts + utile.RandomIntByRange(-attempts/2, attempts/2)
				log.Infoln("waiting", waiting_time, "\n")
				time.Sleep(time.Second * time.Duration(attempts+utile.RandomIntByRange(-5, 5)))
				log.Errorln(fmt.Sprintf("Error copying data (attempt %d/%d): %v\n", attempts, 5, err))
			} else {
				i = 10
			}
		}
		

		log.Infoln(fmt.Sprintf("Downloaded bytes %d-%d\n", start, end))

		return nil
	}
	errmsg := fmt.Sprintf("Failed to download bytes %d-%d after %d attempts\n", start, end, retries+1)
	log.Errorln(err)
	return errors.New(errmsg)
}

func (me *FILE) downloadchunk(chunk []int, wg *sync.WaitGroup) error {
	defer wg.Done()
	for _, index := range chunk {
		packet := me.Packets[index]
		log.Infoln(fmt.Sprintf("The index %d  of %v ", index, chunk))
		if !packet.Done {
			if err := me.downloadRange(packet.Start, packet.End); err != nil {
				log.Errorln(err.Error())
				return err
			} else {
				me.Packets[index].Done = true
			}
		} else {
			log.Warnln("The packe is samhow downloaded")
		}
	}
	// me.writeToAJson(utile.Cfgjson(me.Output))
	return nil
}
func (me *FILE) download(numThreads int) error {

	if err := me.readFromJson(utile.Cfgjson(me.Output)); err != nil {
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
	wg.Add(numThreads)

	for _, chunk := range chunks {
		go me.downloadchunk(chunk, &wg)

	}
	wg.Wait()
	me.writeToAJson(utile.Cfgjson(me.Output))

	return me.finishIt()

	// for index, packet := range me.Packets {
	// 	if !me.Packets[index].Done {
	// 		if err := me.downloadRange(
	// 			unf,
	// 			packet.Start,
	// 			packet.End,
	// 		); err != nil {
	// 			log.Errorln("\n\t", err.Error())
	// 			return err
	// 		} else {
	// 			me.Packets[index].Done = true
	// 			if err := me.writeToAJson(utile.Cfgjson(me.Output)); err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// }

	return nil
}
func (me *FILE) Download(numThreads int) error {

	if IsExist, err := utile.PathIsExist(utile.Cfgjson(me.Output)); err != nil {
		log.Errorln("\n\t", err.Error())
		return err
	} else if !IsExist {
		if err := me.mkeConfig(numThreads); err != nil {
			log.Errorln("\n\t", err.Error())
			return err
		}
	}

	me.readFromJson(utile.Cfgjson(me.Output))

	if err := me.download(numThreads); err != nil {
		log.Errorln("\n\t", err.Error())
		return err
	}
	return nil

	// if _, err := os.Stat(me.Output); err == nil {
	//     fmt.Println("File exists")
	//     return nil
	// } else if os.IsNotExist(err) {
	//     if  _, err := os.Stat(utile.Unfiniched(me.Output)); err == nil {
	//         fmt.Println("File exists")
	//     } else if os.IsNotExist(err) {
	//         fmt.Println("File does not exist")
	//     } else {
	//         fmt.Println("Error checking file:", err)
	//     }
	//     me.mkeConfig(numThreads)
	//     } else {
	//     fmt.Println("Error checking file:", err)
	// }
	// out, err := os.Create(me.Output)
	// if err != nil {
	// 	return fmt.Errorf("error creating output file: %v", err)
	// }
	// defer out.Close()
	// resp, err := http.Head(me.Url)
	// if err != nil {
	// 	return fmt.Errorf("error getting head request: %v", err)
	// }
	// defer resp.Body.Close()
	// if resp.Header.Get("Accept-Ranges") != "bytes" {
	// 	return  errors.New("server does not support range requests")
	// }
	// contentLength, err := strconv.Atoi(resp.Header.Get("Content-Length"))
	// if err != nil {
	// 	return fmt.Errorf("failed to get content length: %v", err)
	// }
	// segmentSize := int(math.Ceil(float64(contentLength) / float64(numThreads)))
	// var wg sync.WaitGroup
	// wg.Add(numThreads)
	// var packet  []Packet
	// for i := 0; i < numThreads; i++ {
	// 	start := segmentSize * i
	// 	end := start + segmentSize - 1
	// 	if i == numThreads-1 {
	// 		end = contentLength - 1
	// 	}
	//     packet = append(packet, Packet{Start: start , End: end , Done : false})
	// }
	// // go me.downloadRange(out, start, end, &wg, 3)
	// wg.Wait()
	// fmt.Println("All workers done")

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
