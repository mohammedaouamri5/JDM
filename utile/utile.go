package utile

import (
	fmt "fmt"
	log "github.com/sirupsen/logrus"
	os "os"
)

func Unfiniched(me string) string {
	return me + ".unfiniched"
}

func Cfgjson(me string) string {
	return me + ".cfg.json"
}

func PathIsExist(p_path string) (bool, error) {

	log.Infoln(fmt.Sprintf("\n\tTesting If the File %s exist ", p_path))
	_, err := os.Stat(p_path)
	if err == nil {
		log.Infoln("\n\tFile does not exist")
		return true, nil

	} else if os.IsNotExist(err) {
		log.Infoln("\n\tFile does not exist")
		return false, nil
	}
 
	return false, err
}
