package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func getUrlTs() (int, string, []byte) {
	f, _ := os.OpenFile("1.txt", os.O_RDONLY, 0666)
	b := bufio.NewReader(f)
	var i = 0
	//var fileData []byte
	var i2 = 0
	var addr = ""
	var vi = ""
	for {
		i++
		line, _, err := b.ReadLine()
		if err == io.EOF {
			break
		}
		if i == 3 {

			data := strings.Split(string(line), ",")

			data1 := strings.Split(data[1], "\"")
			addr = data1[1]
			data2 := strings.Split(data[2], "=")
			vi = data2[1]
		}
		if i >= 7 && i%2 == 1 {
			i2++
			url := string(line)
			url = url[0:4] + url[5:]
			fmt.Println(url)
			r, e := http.Get(url)
			if e != nil {
				r, e = http.Get(url)
				if e != nil {
					fmt.Println(e)
					i2--
					continue
				}
			}
			ds, _ := io.ReadAll(r.Body)
			r.Body.Close()
			f2, e := os.Create(dir + strconv.Itoa(i2) + ".ts")
			if e != nil {
				fmt.Println(e)
				i2--
				continue
			}
			f2.Write(ds)
			f2.Close()
			//fileData = append(fileData, ds...)
		}
	}

	f.Close()
	iv := []byte{}
	for k, _ := range vi {
		if k > 1 && k%2 == 1 {
			nudd, e := strconv.ParseInt(vi[k-1:k+1], 16, 16)
			if e != nil {
				fmt.Println(e)
				return 0, addr, nil
			}
			iv = append(iv, byte(nudd))
		}
	}

	return i2, addr, iv
}

var dir = "./"
var wangzhhi = ""

func run() {
	if len(os.Args) == 1 {
		fmt.Println("plase input m3u8 addr")
		fmt.Println("请输入 m3u8的地址作为参数运行")
		return
	}
	wangzhhi = os.Args[1]

	now := time.Now()
	dir = "./" + now.Format("2006_01_02_15_04_05") + "/"
	os.Mkdir(dir, os.ModeDir)

	addrdb, _ := http.Get(wangzhhi)
	addrs, e := io.ReadAll(addrdb.Body)
	if e != nil {
		fmt.Println(e)
		return
	}
	writeFile("1.txt", addrs)

	num, adr, iv := getUrlTs()
	if num == 0 {
		fmt.Println("num 0")
		return
	}
	bd, _ := http.Get(adr)
	key, e := io.ReadAll(bd.Body)
	if e != nil {
		fmt.Println(e)
		return
	}
	bd.Body.Close()

	var alldata []byte
	for i := 1; i <= num; i++ {
		f, _ := os.Open(dir + strconv.Itoa(i) + ".ts")
		wd, e := io.ReadAll(f)
		if e != nil {
			fmt.Println(e)
			return
		}
		fmt.Println(len(wd))
		//

		plaintext := CBCDecrypt(wd, key, iv)
		//	fmt.Println(plaintext)
		f.Close()

		alldata = append(alldata, []byte(plaintext)...)
		//writeFile(dir+strconv.Itoa(i)+"_un.ts", []byte(plaintext))
	}
	writeFile(dir+"all.ts", alldata)
}

func writeFile(name string, data []byte) {
	f, _ := os.Create(name)
	f.Write(data)
	f.Close()
}

// CBCDecrypt AES-CBC 解密
func CBCDecrypt(ciphercode []byte, key []byte, iv []byte) string {
	defer func() {
		if err := recover(); err != nil {
			fmt.Println("cbc decrypt err:", err)
		}
	}()

	block, err := aes.NewCipher(key)
	if err != nil {
		fmt.Println("cbc NewCipher err:", err)
		return ""
	}

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphercode, ciphercode)

	plaintext := string(ciphercode) // ↓ 减去 padding
	return plaintext[:len(plaintext)-int(plaintext[len(plaintext)-1])]
}

func ffe() {
	arg := []string{"-i", dir + "all.ts", "-c", "copy", "-map", "0:v", "-map", "0:a", "output.mp4"}
	cmd := exec.Command(ffeName, arg...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func clear() {
	os.RemoveAll(dir)
	os.Remove("1.txt")
}

func main() {
	run()
	ffe()
	clear()
}
