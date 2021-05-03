package service

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	_ "errors"
	"fmt"
	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
	"log"
	"strconv"
	"strings"
	"unicode/utf8"
)
func Decrypt(x string, key []byte) string {
	fmt.Print("Start of decryption without Vector IV ####### decryption with key : ")
	fmt.Println(key)
	fmt.Print("String to decrypt : ")
	fmt.Println(string(x))
	ct ,_:=base64.StdEncoding.DecodeString(x)
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(ct))
	mode.CryptBlocks(pt, ct)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Unpad(pt) // unpad plaintext after decryption
	if err != nil {
		panic(err.Error())
	}
	fmt.Print(" Result : ")
	fmt.Println(string(pt))
	fmt.Println("End of decryption without Vector IV ###########")
	return string(pt)
}

func Encrypt(pt, key []byte) string {
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Pad(pt) // padd last block of plaintext if block size less than block cipher size
	if err != nil {
		panic(err.Error())
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	encryptedString := base64.StdEncoding.EncodeToString(ct)
	fmt.Print("Start of encryption without Vector IV |||||| Encryption with key : ")
	fmt.Println(key)
	fmt.Print("String to encrypt : ")
	fmt.Println(string(pt))
	fmt.Print(" Result : ")
	fmt.Println(encryptedString)
	fmt.Println("End of encryption without Vector IV |||||||||")
	return encryptedString
}


func AESEncryptWithIV( iv[]byte,src,key string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println("key error1", err)
	}
	if src == "" {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCEncrypter(block, iv)
	content := []byte(src)
	content = PKCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)
	encryptedString := base64.StdEncoding.EncodeToString(crypted)
	fmt.Print("Start of encryption -------- Encryption with key : ")
	fmt.Print(key)
	fmt.Print(" and vector IV : ")
	fmt.Println(iv)
	fmt.Print("String to encrypt : ")
	fmt.Print(src)
	fmt.Print(" Result : ")
	fmt.Println(encryptedString)
	fmt.Println("End of encryption -------------")
	return encryptedString
}

func DecryptWithIV(crypt , key string,iv []byte) (string) {
	enc,_ := base64.StdEncoding.DecodeString(crypt)
	btenc := []byte(enc)
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println("key error1", err)
	}
	if len(btenc) == 0 {
		fmt.Println("plain content empty")
	}
	ecb := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(btenc))
	ecb.CryptBlocks(decrypted, btenc)

	bytetrim := PKCS5Trimming(decrypted)
	if(bytetrim==nil){
		return ""
	}
	decryptedstring := base64.StdEncoding.EncodeToString(bytetrim)
	data, _ := base64.StdEncoding.DecodeString(decryptedstring)

	fmt.Print("Start of Decryption ++++++++ Decryption with key : ")
	fmt.Print(key)
	fmt.Print(" and vector IV : ")
	fmt.Println(iv)
	fmt.Print("String to Decrypt : ")
	fmt.Print(crypt)
	fmt.Print(" Result : ")
	fmt.Println(string(data))
	fmt.Println("End of Decryption +++++++++++++")

	return string(data)
}
func DecodeHex(input []byte) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(input)))
	_, err := hex.Decode(dst, input)
	if err != nil {
		return nil, err
	}
	return dst, nil
}
func Base64Encode(input []byte) ([]byte) {
	eb := make([]byte, base64.StdEncoding.EncodedLen(len(input)))
	base64.StdEncoding.Encode(eb, input)

	return eb
}
func GetByteArray(in string)([]byte){
	fmt.Printf("String that will be converted to byte array Vector IV : %+q\n",in)
	x, err := parseUsers([]byte(in));
	if(err!=nil){
		panic(err)
	}
	var s []byte;
	fmt.Printf("Vector IV Type Byte[] as byte value : %+q\n", x)

	s = []byte{byte(x[0]),byte(x[1]),byte(x[2]),byte(x[3]),byte(x[4]),byte(x[5]),byte(x[6]),byte(x[7]),byte(x[8]),byte(x[9]),byte(x[10]),byte(x[11]),byte(x[12]),byte(x[13]),byte(x[14]),byte(x[15])};
	myString := string(s[:]);
	fmt.Print("Vector IV Type Byte[] as string : ")
	fmt.Println(myString);
	fmt.Print("Vector IV Type Byte[] as array int : ")
	fmt.Println(s)
	return s;
}
func parseUsers(jsonBuffer []byte) ([]int8, error) {

	// We create an empty array
	users := []int8{}

	// Unmarshal the json into it. this will use the struct tag
	err := json.Unmarshal(jsonBuffer, &users)
	if err != nil {
		return nil, err
	}

	// the array is now filled with users
	return users, nil

}

func GenerateSecureIVVector()[]byte{
	var x string
	x ="["
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		log.Fatal(err)
	}
	fmt.Print("Generated Random Byte [] for Vector IV : ")
	fmt.Println(iv)
	fmt.Print("After conversion Vector IV Result : ")
	for i := 0; i < len(iv); i++ {
		s ,_ := strconv.Atoi(fmt.Sprintf("%d",iv[i]))
		if(s>127){
			s = s-128
		}
		x = x + fmt.Sprintf("%d",s)
		if(i!=len(iv)-1){
			x = x + " ,"
		}
	}
	x = x+"]"
	fmt.Println(x)
	return GetByteArray(x)
}
func GetVectorIVAsString(iv []byte) string{
	var x string
	x = "["
	for i := 0; i < len(iv); i++ {
		s ,_ := strconv.Atoi(fmt.Sprintf("%d",iv[i]))
		if(s>127){
			s = s-128
		}
		x = x + fmt.Sprintf("%d",s)
		if(i!=len(iv)-1){
			x = x + " ,"
		}
	}
	x = x+"]"
	return x
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func PKCS5Trimming(encrypt []byte) []byte{
	if(len(encrypt)-1<0){
		return nil
	}
	padding := encrypt[len(encrypt)-1]
	if(padding<0){
		return nil
	}
	if(len(encrypt)-int(padding)<0){
		return nil
	}
	return encrypt[:len(encrypt)-int(padding)]
}
func DecrypRespBodyAPI(js string)string{
	var x = "{"
	var bt []byte
	c := make(map[string]json.RawMessage)
	e := json.Unmarshal([]byte(js), &c)
	if e != nil {
		panic(e)
	}
	//k := make([]string, len(c))
	kv := make([]string,len(c))
	i := 0
	for s, v := range c {
		if(s=="val"){
			bt = GetByteArray(Decrypt(trimFirstRune(strings.TrimSuffix( string(v),"\"")),[]byte(secretkey)))
		}
	}
	i = 0
	var strval ,strkey string
	for s, v := range c {
		if(s!="val") {
				fmt.Println("asdasd", v)
				if (v != nil) {
					strval = "\"" + DecryptWithIV(trimFirstRune(strings.TrimSuffix(string(v), "\"")), secretkey, bt) + "\""
					if (strval == "") {
						strval = "null"
					}
				} else {
					strval = ""
				}
				if (s != "val") {
					strkey = DecryptWithIV(s, secretkey, bt)
				} else {
					strkey = "val"
				}
				x = x + "\"" + strkey + "\":" + strval + ""

				kv[i] = string(v)

			i++
			if (i < len(c)-1) {
				x = x + ","
			}
		}
	}
	x = x+"}"
	fmt.Println(x)

	//fmt.Printf("%#v\n", kv)
	return x
}
func trimFirstRune(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}