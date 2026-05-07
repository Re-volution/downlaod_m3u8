package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

func (dsr *dirstr) getUrlTs2() int {
	f, _ := os.OpenFile("2", os.O_RDONLY, 0666)
	b := bufio.NewReader(f)
	var i = 0
	//var fileData []byte
	var i2 = 0

	for {
		i++
		line, _, err := b.ReadLine()
		if err == io.EOF {
			break
		}

		if string(line)[:4] == "http" {
			i2++
			url := string(line)
			//url = url[0:4] + url[5:]
			//fmt.Println(url)
			//r, e := http.Get(url)
			//if e != nil {
			//	time.Sleep(1e9)
			//	r, e = http.Get(url)
			//	if e != nil {
			//		fmt.Println(e)
			//		i2--
			//		continue
			//	}
			//}
			//ds, _ := io.ReadAll(r.Body)
			//r.Body.Close()
			ds := dsr.httppoxy(url)
			f2, e := os.Create(dsr.dir + strconv.Itoa(i2) + ".ts")
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
	return i2
}

var client = &http.Client{Transport: &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}}

func gethttps(addr string) ([]byte, error) {
	resp, err := client.Get(addr)
	if err != nil {
		time.Sleep(1e9)
		resp, err = client.Get(addr)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		fmt.Println("error:", err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	return body, err
}

func (dsr *dirstr) getUrlTs3() int {
	f, _ := os.OpenFile(dsr.onceName+"1.txt", os.O_RDONLY, 0666)
	b := bufio.NewReader(f)
	var i = 0
	//var fileData []byte
	var i2 = 0
	//var addr = strings.Split(dsr.wangzhhi, "index.m3u8")[0]
	for {
		i++
		line, _, err := b.ReadLine()
		if err == io.EOF {
			break
		}
		if (len(strings.Split(string(line), ".ts")) > 1 || len(strings.Split(string(line), "https")) > 1) && len(strings.Split(string(line), "adjump")) == 1 {
			i2++
			url := string(line)
			//url = url[0:4] + url[5:]
			fmt.Println(url)
			ds, e := gethttps(url)
			if e != nil {
				fmt.Println(e)
				i2--
				continue
			}
			f2, e := os.Create(dsr.dir + strconv.Itoa(i2) + ".ts")
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

	return i2
}

// 加密全连接
func (dsr *dirstr) getUrlTs4() (int, []byte, []byte) {
	f, _ := os.OpenFile(dsr.onceName+"1.txt", os.O_RDONLY, 0666)
	b := bufio.NewReader(f)
	var i = 0
	//var fileData []byte
	var i2 = 0
	var key = []byte{}
	var vi = ""
	//var addr = strings.Split(dsr.wangzhhi, "index.m3u8")[0]
	for {
		i++
		line, _, err := b.ReadLine()
		if err == io.EOF {
			break
		}
		if i == 3 {

			data := strings.Split(string(line), ",")

			addr := strings.Split(data[1], "\"")
			bd, _ := http.Get(addr[1])
			key, err = io.ReadAll(bd.Body)
			if err != nil {
				fmt.Println(err)
				return 0, nil, nil
			}
			bd.Body.Close()
			data2 := strings.Split(data[2], "=")
			vi = data2[1]
		}
		if i < 5 {
			continue
		}

		if (len(strings.Split(string(line), ".ts")) > 1 || len(strings.Split(string(line), "https")) > 1) && len(strings.Split(string(line), "adjump")) == 1 {
			i2++
			url := string(line)
			//url = url[0:4] + url[5:]
			fmt.Println(url)
			ds, e := gethttps(url)
			if e != nil {
				fmt.Println(e)
				i2--
				continue
			}
			f2, e := os.Create(dsr.dir + strconv.Itoa(i2) + ".ts")
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
				return 0, key, nil
			}
			iv = append(iv, byte(nudd))
		}
	}

	return i2, key, iv
}

func (dsr *dirstr) getUrlTs() (int, string, []byte) {
	f, _ := os.OpenFile(dsr.onceName+"1.txt", os.O_RDONLY, 0666)
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
		if string(line)[:4] == "http" {
			i2++
			url := string(line)
			//url = url[0:4] + url[5:]
			fmt.Println(url)
			r, e := http.Get(url)
			if e != nil {
				time.Sleep(1e9)
				r, e = http.Get(url)
				if e != nil {
					fmt.Println(e)
					i2--
					continue
				}
			}
			ds, _ := io.ReadAll(r.Body)
			r.Body.Close()
			f2, e := os.Create(dsr.dir + strconv.Itoa(i2) + ".ts")
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

var formatStr = "2006_01_02_15_04_05"
var movieDir = "movie"
var proxyUrl = "http://127.0.0.1:1080"

// https://hls.vdtuzv.com/videos3/ad866431958933e7a668624a999859d8/ad866431958933e7a668624a999859d8.m3u8?auth_key=1711024782-65fc2a8e746f5-0-1a0ae9ab99cb9633dbddb69b16fff052&v=2
func (dsr *dirstr) httppoxy(wz string) []byte {

	/*
		1. 代理请求
		2. 跳过https不安全验证
		3. 自定义请求头 User-Agent
	*/
	// webUrl := "http://ip.gs/"
	// proxyUrl := "http://171.215.227.125:9000"

	request, _ := http.NewRequest("GET", wz, nil)
	request.Header.Set("Connection", "keep-alive")
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/56.0.2924.87 Safari/537.36")
	request.AddCookie(&http.Cookie{
		Name:       "FC2_GDPR",
		Value:      "ture",
		Path:       "/",
		Domain:     "liaoningmovie.net",
		Expires:    time.Now().AddDate(1, 0, 0),
		RawExpires: "",
		MaxAge:     0,
		Secure:     false,
		HttpOnly:   false,
		SameSite:   0,
		Raw:        "",
		Unparsed:   nil,
	})

	proxy, _ := url.Parse(proxyUrl)
	tr := &http.Transport{
		Proxy:           http.ProxyURL(proxy),
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 5, //超时时间
	}

	resp, err := client.Do(request)
	if err != nil {
		fmt.Println("出错了", err)
		return nil
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
	return body
}

type dirstr struct {
	dir      string
	wangzhhi string
	onceName string
	addrs    []byte
	data     []byte
}

func initdsr(wz string) *dirstr {
	var now = time.Now()
	var t = new(dirstr)
	t.wangzhhi = wz
	t.onceName = now.Format(formatStr)
	t.dir = "./" + t.onceName + "/"
	os.Mkdir(t.dir, os.ModeDir)
	return t
}

func run2() {
	var dsr = initdsr("wz")
	dsr.getdataByAddrs2()
}

func run(wz string) {
	if len(os.Args) == 1 {
		fmt.Println("plase input m3u8 addr")
		fmt.Println("请输入 m3u8的地址作为参数运行")
		time.Sleep(10e9)
		os.Exit(1)
	}
	var dsr = initdsr(wz)
	var addrs []byte
	if len(os.Args) == 2 {
		if wz == "file" {
			f, e := os.Open("./file.txt")
			if e != nil {
				fmt.Println("open filetxt fial.", e)
				return
			}
			data, _ := io.ReadAll(f)
			for _, adrz := range strings.Split(string(data), "\r\n") {
				run(adrz)
			}
			return
		} else {
			var e error
			addrdb, e := http.Get(dsr.wangzhhi)
			if e != nil {
				fmt.Println(e)
				return
			}
			addrs, e = io.ReadAll(addrdb.Body)
			if e != nil {
				fmt.Println(e)
				return
			}
		}

	} else if len(os.Args) == 3 {
		addrs = dsr.httppoxy(dsr.wangzhhi)
	}
	if addrs == nil {
		fmt.Println("未读取到数据")
		return
	}
	dsr.getdataByAddrs(addrs)
}

func (dsr *dirstr) getdataByAddrs2() {

	num := dsr.getUrlTs2()
	if num == 0 {
		fmt.Println("num 0")
		return
	}

	var alldata []byte
	for i := 1; i <= num; i++ {
		f, _ := os.Open(dsr.dir + strconv.Itoa(i) + ".ts")
		wd, e := io.ReadAll(f)
		if e != nil {
			fmt.Println(e)
			return
		}
		fmt.Println(len(wd))
		//
		plaintext := wd
		//plaintext := CBCDecrypt(wd, key, iv)
		//	fmt.Println(plaintext)
		f.Close()

		alldata = append(alldata, []byte(plaintext)...)
		//writeFile(dir+strconv.Itoa(i)+"_un.ts", []byte(plaintext))
	}
	writeFile(dsr.dir+"all.ts", alldata)
	dsr.ffe()
	dsr.clear()
}

func (dsr *dirstr) getdataByAddrs(addrs []byte) {
	writeFile(dsr.onceName+"1.txt", addrs)
	dsr.data = addrs
	//不加密全连接
	//num := dsr.getUrlTs3()
	//if num == 0 {
	//	fmt.Println("num 0")
	//	return
	//}

	//加密全连接
	num, key, iv := dsr.getUrlTs4()
	if num == 0 {
		fmt.Println("num 0")
		return
	}
	//bd, _ := http.Get(adr)
	//key, e := io.ReadAll(bd.Body)
	//if e != nil {
	//	fmt.Println(e)
	//	return
	//}
	//bd.Body.Close()

	var alldata []byte
	for i := 1; i <= num; i++ {
		f, _ := os.Open(dsr.dir + strconv.Itoa(i) + ".ts")
		wd, e := io.ReadAll(f)
		if e != nil {
			fmt.Println(e)
			return
		}
		fmt.Println(len(wd))
		//
		//plaintext := wd
		plaintext := CBCDecrypt(wd, []byte(key), iv)
		//	fmt.Println(plaintext)
		f.Close()

		alldata = append(alldata, []byte(plaintext)...)
		//writeFile(dir+strconv.Itoa(i)+"_un.ts", []byte(plaintext))
	}
	writeFile(dsr.dir+"all.ts", alldata)
	dsr.ffe()
	dsr.clear()
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

func (dsr *dirstr) ffe() {
	now := time.Now().Format(formatStr)
	os.Mkdir("../"+movieDir, os.ModeDir)
	arg := []string{"-i", dsr.dir + "all.ts", "-c", "copy", "-map", "0:v", "-map", "0:a", "../" + movieDir + "/" + now + ".mp4"}
	cmd := exec.Command(ffeName, arg...)
	_, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}

func (dsr *dirstr) clear() {
	os.RemoveAll(dsr.dir)
	os.Remove(dsr.onceName + "1.txt")
}

func main() {
	run(os.Args[1])
	//	run2()
}
