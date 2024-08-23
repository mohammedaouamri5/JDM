package downloader_test

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	// "os/exec"
	T "testing"

	"github.com/mohammedaouamri5/JDM-back/downloader"
	. "github.com/stretchr/testify/assert"
)

/*
"00df751d287d01821e0123c10056d020" : "archlinux-2010.05-netinstall-i686.iso.torrent"
"a08ef236184006c9c735436883fa6b66" : "archlinux-2010.05-netinstall-dual.iso.torrent"
"32e126ab42422faba7e4811be85e5992" : "archlinux-2010.05-netinstall-x86_64.iso.torrent"
"5253edaf4afe0472f1a96c2985814cc5" : "archlinux-2010.05-core-x86_64.iso.torrent"
"47164714e4b19da795eb16270a1fd354" : "archlinux-2010.05-core-dual.iso.torrent"
"2754567a712fc0e4f9c71dbc31279821" : "archlinux-2010.05-core-i686.iso.torrent"
"577ca8026a9997f7cecb430276d78793" : "archlinux-2010.05-netinstall-x86_64.iso"
"327f3865c7d1d785e62ff7f566a4fe7b" : "archlinux-2010.05-netinstall-i686.iso"
"5db5fd11713635cff208b11a498c59ef" : "archlinux-2010.05-core-i686.iso"
"9e9057702af5826a3b924233bf44fe66" : "archlinux-2010.05-core-x86_64.iso"
"b272d4fdfff702b84c5a2065a9b0737b" : "archlinux-2010.05-netinstall-dual.iso"
"6b0fec50e4895eaecd58a6157d1b949a" : "archlinux-2010.05-core-dual.iso"
*/


func md5Checksum(filePath string) (string, error) {
    // Open the file
    file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    // Create a new MD5 hash object
    hash := md5.New()

    // Copy the file content to the hash object
    if _, err := io.Copy(hash, file); err != nil {
        return "", err
    }

    // Get the MD5 checksum as a byte slice
    checksum := hash.Sum(nil)

    // Encode the checksum as a hexadecimal string
    checksumString := hex.EncodeToString(checksum)

    return checksumString, nil
}

 


func TestDonwload(t *T.T)  {
	fileChecksumMap := map[string]string{
        "archlinux-2010.05-netinstall-i686.iso.torrent" : "00df751d287d01821e0123c10056d020" , 
        "archlinux-2010.05-netinstall-dual.iso.torrent" : "a08ef236184006c9c735436883fa6b66" , 
        "archlinux-2010.05-netinstall-x86_64.iso.torrent" : "32e126ab42422faba7e4811be85e5992" , 
        "archlinux-2010.05-core-x86_64.iso.torrent" : "5253edaf4afe0472f1a96c2985814cc5" , 
        "archlinux-2010.05-core-dual.iso.torrent" : "47164714e4b19da795eb16270a1fd354" , 
        "archlinux-2010.05-core-i686.iso.torrent" : "2754567a712fc0e4f9c71dbc31279821" , 
        "archlinux-2010.05-netinstall-x86_64.iso" : "577ca8026a9997f7cecb430276d78793" , 
        "archlinux-2010.05-netinstall-i686.iso" : "327f3865c7d1d785e62ff7f566a4fe7b" , 
        "archlinux-2010.05-core-i686.iso" : "5db5fd11713635cff208b11a498c59ef" , 
        "archlinux-2010.05-core-x86_64.iso" : "9e9057702af5826a3b924233bf44fe66" , 
        "archlinux-2010.05-netinstall-dual.iso" : "b272d4fdfff702b84c5a2065a9b0737b" , 
        "archlinux-2010.05-core-dual.iso" : "6b0fec50e4895eaecd58a6157d1b949a" , 
    }

    ResultMap := make(map[string]string)
    path := "./tmp" 
    for name, hash := range fileChecksumMap {
        file := downloader.FILE{}
        file.Constructor(
            "https://archive.archlinux.org/iso/2010.05/"+name, 
            name,
            &path,
        )
        assert := New(t)
        assert.NoError(file.Download(10))
        checksum, err := md5Checksum(file.Output)
        assert.NoError(err)
        ResultMap[hash] = checksum
    }
    
    for True , Pred := range ResultMap {
        assert := New(t) 
        assert.Equal(Pred, ResultMap[True])
    }

}

