package main

import (
	"TA/service"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)
var (
	//db,_ = service.MySQL()
	ctx          context.Context
	cookie  *http.Cookie;
	itemlistcookie string;
	itemlisturl string;
	cookie2  *http.Cookie;
	itemlistcookie2 string;
	itemlisturl2 string;
)

type VARequest struct {
	Amount string
	Description string
	Name string
	VA string
}
type OrderStatus struct {
	Id string `json:"id"`
	Status string `json:"status"`

}

type M map[string]interface{}

var cookieName = "CookieData"
var cookieName2 = "CookieData2"

func main() {
	fmt.Println("")

	fmt.Printf("asdsad %f","1.1E7")
	handleRequests()
	//bt := service.GetByteArray("[102 ,3 ,121 ,64 ,61 ,125 ,29 ,53 ,80 ,55 ,41 ,1 ,34 ,86 ,9 ,7]")

	//rs:= service.AESEncryptWithIV( bt,"87654321","1suio2930ekdmcxk");
	//_ = service.DecryptWithIV(rs, "1suio2930ekdmcxk", bt)
	//service.GenerateSecureIVVector()
	//service.DecrypRespBodyAPI(service.GetVA("100013"))

	//var x = service.GeneratReqBodyForCreateVA("10000000","Beli hp cash " ,"Belanjaonline.com - hp samsung ","100012")
	//service.DecrypRespBodyAPI(service.CreateVA(x))

	//var x ="{\"val\":\"5XaO8oQTK5AihKV65RMMetYRjISrscjjwC9h2cPgAJ/xDr9pZcBS5EMvWbMyf9H1w8Pw0BhAXVYyU3T9SUZUkA==\",\"b+KvZKwF1TlIMq4r8WEuNw==\":\"fvB8Ez8BF8RnJnVw19ka1vR4LSLQH2whpVkuKIUcZG4=\"}"
	//service.DecrypRespBodyAPI(x)
}
func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/delete", ActionDelete)
	myRouter.HandleFunc("/deleteredirect", ActionDelete2)
	myRouter.HandleFunc("/deletecookieorder",deletecookieorder)
	myRouter.HandleFunc("/getitem", getListItemForPage)
	myRouter.HandleFunc("/getitemdetail", getItemDetail)
	myRouter.HandleFunc("/additemcookie",addItemOnCookie)
	myRouter.HandleFunc("/purchase",purchase)
	myRouter.HandleFunc("/ord",tezzz)
	myRouter.HandleFunc("/order",order)
	myRouter.HandleFunc("/callapi/createva", callAPICreateVA)
	myRouter.HandleFunc("/callapi/getva/{va}", callAPIGetVA)
	myRouter.HandleFunc("/va/notif", notificationPaymentVA)
	myRouter.HandleFunc("/insert/item", insertItem)
	myRouter.HandleFunc("/insert/billing", insertBilling)
	myRouter.HandleFunc("/css/{file}", getCss)
	myRouter.HandleFunc("/js/{file}", getJs)
	myRouter.HandleFunc("/img/{file}",getImage)
	ctx = context.Background()
	log.Println("starting...")
	log.Fatal(http.ListenAndServe(":10000", myRouter))
}
func getImage(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["file"]
	//fmt.Println(key)
	fmt.Fprintf(w,"%v",ReadFile(key))
}
func getJs(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["file"]
	//fmt.Println(key)
	fmt.Fprintf(w,"%v",ReadFile(key))
}
func getCss(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "text/css")
	vars := mux.Vars(r)
	key := vars["file"]
	//fmt.Println(key)
	fmt.Fprintf(w,"%v",ReadFile(key))
}
func getItemDetail(w http.ResponseWriter,r *http.Request){
	checkCookie(w,r)
	key := r.URL.Query()["id"]
	keyint,_ := strconv.Atoi(key[0])
	db ,_ := service.MySQL()
	var v1,v2,v3,v4,v5 string
	item := service.GetItemById(keyint,db)
	v1 = "Nama Barang : "+item.Name
	t := strings.Split(strconv.Itoa(item.Price),"")
	stramount := ""
	i:=0
	for i = 0; i < len(t) ; i++{
		if((len(t) - i - 1) % 3 == 0){
			stramount += t[i] + ",";
		}else{
			stramount += t[i]
		}
	}

	v2 = fmt.Sprintf("Harga Barang RP.%v %v %v",stramount , " stok : " ,item.Stock)
	v3 = "Nama Toko : " + item.NamaToko + " </br>Alamat " + item.AddressSuplier +" - "+ item.City
	v4 = "Deskripsi : " + item.Description
	v5 = item.Img
	var strprint string
	strprint=fmt.Sprintf(ReadFile("itemdetail.html"),v1,v2,v3,v4,key[0],key[0],v5)
	fmt.Sprintln(strprint)
	fmt.Fprintf(w, fmt.Sprintf(ReadFile("base.html"),itemlistcookie2,itemlistcookie,itemlisturl,strprint))


}
func getListItemForPage(w http.ResponseWriter, r *http.Request){
	checkCookie(w,r)
	var url string;
	url = "getitem?"
	name:= r.URL.Query()["name"]
	page:=r.URL.Query()["page"]
	total:=r.URL.Query()["total"]
	category:=r.URL.Query()["category"]
	order:=r.URL.Query()["order"]
	var urlselect = "getitem?page=1&"
	//reqBody, _ := ioutil.ReadAll(r.Body)
	var req service.PageFilterItem
	if(len(name)>0) {
		req.Name = name[0]
		url = url + strings.ReplaceAll("name="+req.Name+"&"," ","+")
		urlselect =urlselect + strings.ReplaceAll("name="+req.Name+"&"," ","+")
	}
	if(len(page)>0){
		req.Page,_ = strconv.Atoi(page[0])
		url = url + "page="+page[0]+"&"
	}else{
		url = url + "page=1&"
		req.Page = 1
	}
	if(len(total)>0){
		req.Total,_ = strconv.Atoi(total[0])
		url = url + "total="+total[0]+"&"
	}else{
		req.Total = 20
	}
	if(len(category)>0) {
		req.Category = category[0]
		url = url + "category="+category[0]+"&"
	}
	if(len(order)>0) {
		req.Order = order[0]
		url = url + strings.ReplaceAll("order="+order[0]+"&"," ","+")
	}else{
		req.Order = "id desc";
	}
	db,_:=service.MySQL()
	var str string
	str = "<tr>"
	count :=1
	backgroundcolour := ""
	for i,v:=range service.GetItemList(req,db){
		if(i%2==0){
			backgroundcolour = "style = \"background-color: #D1F2EB\""
		}else{
			backgroundcolour = "style = \"background-color: #D7DBDD;\""
		}
		if(count%5==0){
			str = str +"</tr><tr>"
		}
		str = str +"<td "+backgroundcolour+"><a href=\"../getitemdetail?id="+fmt.Sprintf("%v\"",v.ID)+">"+v.Name+"</br>"+v.NamaToko+"</br>"+v.City+"</br>"+formatCurrency(v.Price)+"</a></br>" +
			"<a href=\"../getitemdetail?id="+fmt.Sprintf("%v",v.ID)+"\"><img src=\""+v.Img+"\" alt=\""+v.Name+"\" style=\"width: 100%;height: auto;\" width=\"125\" height=\"150\"></a></td>"
		count++
	}
	str = str+"</tr>"
	db,_=service.MySQL()
	totaldata := service.CountItemByFilter(req,db)

	var recordperpage string
	recordperpage = ""
	if(req.Total==1){
		recordperpage = recordperpage + "<option selected=\"selected\" value=\"1\">1</option>"
	}else {
		recordperpage = recordperpage + "<option value=\"1\">1</option>"
	}
	if(req.Total==4){
		recordperpage = recordperpage + "<option selected=\"selected\" value=\"4\">4</option>"
	}else {
		recordperpage = recordperpage + "<option value=\"4\">4</option>"
	}
	if(req.Total==8){
		recordperpage = recordperpage + "<option selected=\"selected\" value=\"8\">8</option>"
	}else {
		recordperpage = recordperpage + "<option value=\"8\">8</option>"
	}
	if(req.Total==20){
		recordperpage = recordperpage + "<option selected=\"selected\" value=\"20\">20</option>"
	}else {
		recordperpage = recordperpage + "<option value=\"20\">20</option>"
	}
	if(req.Total==40){
		recordperpage = recordperpage + "<option selected=\"selected\" value=\"40\">40</option>"
	}else {
		recordperpage = recordperpage + "<option value=\"40\">40</option>"
	}

	var orderbyid string
	orderbyid = ""
	if(req.Order=="id desc"){
		orderbyid = orderbyid + "<option selected=\"selected\" value=\"id desc\">terbaru</option>"
	}else {
		orderbyid = orderbyid + "<option value=\"id desc\">terbaru</option>"
	}
	if(req.Order=="id asc"){
		orderbyid = orderbyid + "<option selected=\"selected\" value=\"id asc\">terlama</option>"
	}else {
		orderbyid = orderbyid + "<option value=\"id asc\">terlama</option>"
	}
	if(req.Order=="price desc"){
		orderbyid = orderbyid + "<option selected=\"selected\" value=\"price desc\">harga tertinggi</option>"
	}else {
		orderbyid = orderbyid + "<option value=\"price desc\">harga tertinggi</option>"
	}
	if(req.Order=="price asc"){
		orderbyid = orderbyid + "<option selected=\"selected\" value=\"price asc\">harga terendah</option>"
	}else {
		orderbyid = orderbyid + "<option value=\"price asc\">harga terendah</option>"
	}
	if(req.Order=="name asc"){
		orderbyid = orderbyid + "<option selected=\"selected\" value=\"name asc\">nama barang</option>"
	}else {
		orderbyid = orderbyid + "<option value=\"name asc\">nama barang</option>"
	}
	if(req.Order=="city asc"){
		orderbyid = orderbyid + "<option selected=\"selected\" value=\"city asc\">kota</option>"
	}else {
		orderbyid = orderbyid + "<option value=\"name asc\">city barang</option>"
	}

	fmt.Fprintf(w, fmt.Sprintf(ReadFile("base.html"),itemlistcookie2,itemlistcookie,itemlisturl,fmt.Sprintf(ReadFile("listitem.html"),recordperpage,orderbyid,str)))
	totalint64:= int64(req.Total)
	maxpage := (totaldata/totalint64)
	if(totaldata%totalint64>0){
		maxpage++
	}
	i := int64(1)
	if(maxpage>1) {
		fmt.Fprintf(w,"<div  class=\"pagination\"> ")
		fmt.Fprintf(w,"<a href="+strings.ReplaceAll(url, "page="+fmt.Sprintf("%d",req.Page), "page=1")+">&laquo;</a>")
		for i = 1; i <= maxpage; i++ {
			if (i == int64(req.Page)) {
				fmt.Fprintf(w, "<a class=\"active\" href="+strings.ReplaceAll(url, "page="+fmt.Sprintf("%d", req.Page), "page="+fmt.Sprintf("%d", i))+">"+fmt.Sprintf("%d", i)+"</a>")
			} else {
				fmt.Fprintf(w, "<a href="+strings.ReplaceAll(url, "page="+fmt.Sprintf("%d", req.Page), "page="+fmt.Sprintf("%d", i))+">"+fmt.Sprintf("%d", i)+"</a>")
			}
		}
		fmt.Fprintf(w, "<a href="+strings.ReplaceAll(url, "page="+fmt.Sprintf("%d", req.Page), "page="+fmt.Sprintf("%d", maxpage))+">&raquo;</a>")
		fmt.Fprintf(w, "</div>")
	}else{
		fmt.Fprintf(w,"<div  class=\"pagination\"> ")
		fmt.Fprintf(w,"<a class=\"disabled\" href="+strings.ReplaceAll(url, "page="+fmt.Sprintf("%d",req.Page), "page=1")+">&laquo;</a>")
		fmt.Fprintf(w, "<a class=\"disabled\" href="+strings.ReplaceAll(url, "page="+fmt.Sprintf("%d", req.Page), "page="+fmt.Sprintf("%d", maxpage))+">&raquo;</a>")
		fmt.Fprintf(w, "</div>")
	}


	fmt.Fprintf(w,"<script> function myFunction(x) { window.location = \"../"+urlselect+"total=\"+x; }" +
		"function myFunction2(x) { window.location = \"../"+urlselect+"total="+fmt.Sprintf("%v",req.Total)+"&order=\"+x; }</script>\n")
}

func callAPIGetVA(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	key := vars["va"]
	fmt.Fprintf(w, "%+v", (service.DecrypRespBodyAPI(service.GetVA(key))))
}
func callAPICreateVA(w http.ResponseWriter, r *http.Request){
	reqBody, _ := ioutil.ReadAll(r.Body)
	var req VARequest
	json.Unmarshal([]byte(reqBody), &req)
	var x = service.GeneratReqBodyForCreateVA(req.Amount,req.Description ,req.Name,req.VA)
	fmt.Fprintf(w, "%+v",(service.DecrypRespBodyAPI(service.CreateVA(x))))
}
func insertItem(w http.ResponseWriter, r *http.Request){
	reqBody, _ := ioutil.ReadAll(r.Body)
	var req service.Item
	json.Unmarshal([]byte(reqBody), &req)
	if(ctx==nil){
		fmt.Println("ctx nya null")
	}
	service.InsertItem(ctx,req)
}
func insertBilling(w http.ResponseWriter, r *http.Request){
	reqBody, _ := ioutil.ReadAll(r.Body)
	var req service.Billing
	json.Unmarshal([]byte(reqBody), &req)
	if(ctx==nil){
		fmt.Println("ctx nya null")
	}
	fmt.Fprint(w,service.InsertBilling(ctx,req))
}
func insertBillingwithouthttp(req service.Billing)string{
	if(ctx==nil){
		fmt.Println("ctx nya null")
	}
	return service.InsertBilling(ctx,req)
}
func notificationPaymentVA(w http.ResponseWriter, r *http.Request){
	var req service.VaResp
	reqBody, _ := ioutil.ReadAll(r.Body)
	fmt.Println(string(reqBody))
	fmt.Println("after decrypt ")
	x := service.DecrypRespBodyAPI(string(reqBody))
	json.Unmarshal([]byte(x), &req)
	service.UpdatePayment(ctx,req.VaNo,req.Amount)
}
func homePage(w http.ResponseWriter, r *http.Request){
	checkCookie(w,r)
	strbaseanditemcookie := fmt.Sprintf(ReadFile("base.html"),itemlistcookie2,itemlistcookie,itemlisturl,ReadFile("homepage.html"))
	fmt.Fprintf(w, fmt.Sprintf(strbaseanditemcookie))
	fmt.Println("Endpoint Hit: homePage")
}
func ActionDelete(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{}
	c.Name = cookieName
	c.Expires = time.Unix(0, 0)
	c.MaxAge = -1
	http.SetCookie(w, c)

	//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
func ActionDelete2(w http.ResponseWriter, r *http.Request) {
	c := &http.Cookie{}
	c.Name = cookieName
	c.Expires = time.Unix(0, 0)
	c.MaxAge = -1
	http.SetCookie(w, c)

	//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
func ReadFile(filename string) (string){
	b, err := ioutil.ReadFile("ui/"+filename) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	str := string(b) // convert content to a 'string'
	return str;
}
func formatCurrency(price int)string{
t := strings.Split(strconv.Itoa(price),"")
stramount := ""
i:=0
for i = 0; i < len(t) ; i++{
if((len(t) - i - 1) % 3 == 0){
stramount += t[i] + ",";
}else{
stramount += t[i]
}
}
return stramount;
}

func checkCookie(w http.ResponseWriter,r *http.Request){
	var erz error;
	var erz2 error;

	cookie ,erz=r.Cookie("CookieData")
	cookie2,erz2 = r.Cookie(cookieName2)

	if erz != nil {
		fmt.Printf("Cant find cookie :/\r\n")
		cookieName := "CookieData"
		c := &http.Cookie{}
		if storedCookie, _ := r.Cookie(cookieName); storedCookie != nil {
			c = storedCookie
		}
		if c.Value == "" {
			c = &http.Cookie{}
			c.Name = cookieName
			//c.Expires = time.Now().Add(5 * time.Minute)
			http.SetCookie(w, c)
		}
		cookie = c;
	}
	var list []service.Item
	json.Unmarshal([]byte(strings.ReplaceAll(cookie.Value,"@#","\"")), &list)
	itemlistcookie = ""
	itemlisturl=""
	stritem := "<li class=\"clearfix\"><img width=\"70\" height=\"70\" src=\"%v\" alt=\"%v\" /><span class=\"item-name\">%v</span><span class=\"item-price\"></br>Rp.%v</span></li>"
	for _,v := range(list){
		itemlistcookie = itemlistcookie + fmt.Sprintf(stritem,v.Img,v.Name,v.Name,formatCurrency(v.Price))
		itemlisturl = itemlisturl+fmt.Sprintf("%v",v.ID)+","
	}
	itemlisturl = strings.TrimSuffix(itemlisturl,",")


	if erz2 != nil {
		fmt.Printf("Cant find cookie :/\r\n")
		cookieName2 := "CookieData2"
		c := &http.Cookie{}
		if storedCookie2, _ := r.Cookie(cookieName2); storedCookie2 != nil {
			c = storedCookie2
		}
		if c.Value == "" {
			c = &http.Cookie{}
			c.Name = cookieName2
			http.SetCookie(w, c)
		}
		cookie2 = c;
	}
	var list2 []OrderStatus
	json.Unmarshal([]byte(strings.ReplaceAll(cookie2.Value,"@#","\"")), &list2)
	itemlistcookie2 = ""
	itemlisturl2=""
	stritem2 := "<li <span class=\"item-name\"><a href=\"../order?id=%v&c=1\" >Nomer transaksi %v <a href=\"../deletecookieorder?id=%v\">\n<img border=\"0\" alt=\"W3Schools\" src=\"../img/closebutton.jpg\" width=\"15\" height=\"15\"></a></span><span class=\"item-price\"></br>%v</span></li>"
	statusorder:=""
	var ccc []service.Item
	for indexlist2,v := range(list2){
			db, _ := service.MySQL()
			ccc = service.GetItemByOrderId(db, v.Id)
		if(ccc[0].Status=="U"){
			statusorder = "Belum Dibayar"
		}else{
			statusorder = "Lunas"
		}
		itemlistcookie2 = itemlistcookie2 + fmt.Sprintf(stritem2,v.Id,v.Id,indexlist2,statusorder)
		itemlisturl2 = itemlisturl2+fmt.Sprintf("%v",v.Id)+","
	}
	fmt.Printf("%v\n",itemlistcookie2)
	itemlisturl2 = strings.TrimSuffix(itemlisturl2,",")
}
func addItemOnCookie(w http.ResponseWriter, r *http.Request){
	checkCookie(w,r);
	key := r.URL.Query()["iddetail"]
	db,_:= service.MySQL()
	strkey,_:= strconv.Atoi(fmt.Sprintf("%v",key[0]))
	item:=service.GetItemById(strkey,db)
	var list []service.Item
	json.Unmarshal([]byte(strings.ReplaceAll(cookie.Value,"@#","\"")), &list)
	res:=append(list, item)
	pagesJson, _ := json.Marshal(res)
	cookie.Value=strings.ReplaceAll(string(pagesJson),"\"","@#");
	//fmt.Printf("%v",strings.ReplaceAll(cookie.Value,"@#","\""))
	cookie.Name= "CookieData"
	cookie.Expires = time.Now().Add(60 * time.Minute)
	http.SetCookie(w,cookie)
	fmt.Fprintf(w,"<html><head></head><body>Redirecting..<script>  history.go(-2);</script></body></html")
}
func addItemOnCookie2(w http.ResponseWriter, r *http.Request,id string){
	checkCookie(w,r);
	var list []OrderStatus
	json.Unmarshal([]byte(strings.ReplaceAll(cookie2.Value,"@#","\"")), &list)
	var item OrderStatus
	item.Id = id
	res:=append(list, item)
	pagesJson, _ := json.Marshal(res)
	cookie2.Value=strings.ReplaceAll(string(pagesJson),"\"","@#");
	//fmt.Printf("%v",strings.ReplaceAll(cookie.Value,"@#","\""))
	cookie2.Name= "CookieData2"

	http.SetCookie(w,cookie2)
}
func purchase(w http.ResponseWriter, r *http.Request){

	key := r.URL.Query()["idpurchase"]
	checkCookie(w,r)
	//fmt.Fprintf(w, fmt.Sprintf(ReadFile("base.html"),fmt.Sprintf(ReadFile("purchase.html"),fmt.Sprintf("%v",key))))
	strbaseanditemcookie := fmt.Sprintf(ReadFile("base.html"),itemlistcookie2,itemlistcookie,itemlisturl,fmt.Sprintf(ReadFile("purchase.html"),fmt.Sprintf("%v",key)))
	fmt.Fprintf(w, fmt.Sprintf(strbaseanditemcookie))
}
func tezzz(w http.ResponseWriter,r *http.Request){
	//reqBody, _ := ioutil.ReadAll(r.Body)
	ActionDelete(w,r)
	address:=r.PostFormValue("lat") + "\n Detail Alamat : " + r.PostFormValue("address");

	var req service.Billing;
	var detail []service.DetailBilling;
	var listiddetail []int
	json.Unmarshal([]byte(r.PostFormValue("item")), &listiddetail)
	db,_:=service.MySQL()
	strdetail := ""
	TotalAmount := 0
	str:=""
	str2:=""
	background:=""
	for i,v := range (listiddetail){
		if(i%2==0){
			background ="style=\"background-color:#D7DBDD;\""
		}else {
			background ="style=\"background-color:#D1F2EB;\""
		}
		var c service.DetailBilling;
		item:= service.GetItemById(v,db)
		TotalAmount = TotalAmount+item.Price
		str2 ="<td "+background+"><a href=\"../getitemdetail?id="+fmt.Sprintf("%v\"",item.ID)+">"+item.Name+"</br>"+item.NamaToko+"</br>"+item.City+"</br>"+formatCurrency(item.Price)+"</a></td>"
		str ="<td style=\"vertical-align:top\">" +
			"<a href=\"../getitemdetail?id="+fmt.Sprintf("%v",item.ID)+"\"><img src=\""+item.Img+"\" alt=\""+item.Name+"\" style=\"width: 100%;height: auto;\" width=\"125\" height=\"150\"></a></td>"
		strdetail = strdetail + "<tr>"+str2+str+"</tr>"
		c.IDItem=v
		detail = append(detail,c)
	}
	req.ExpiredDate = 3;
	req.PaymentTools = r.PostFormValue("payment");
	req.StatusPayment = "U"
	req.Address =address
	req.PhoneNumber = r.PostFormValue("nohp");
	req.Email = r.PostFormValue("email");
	req.Name = r.PostFormValue("name");
	req.DetailBilling = detail;
	x:=insertBillingwithouthttp(req)
	idorder:= strings.ReplaceAll(x,"100010","")
	//addItemOnCookie2(w,r,strings.ReplaceAll(x,"100010",""))

	q := r.URL.Query()
	q.Add("id",x)
	x = "  : " + x
	strdetail = "  : " +strdetail
	//checkCookie(w,r)
	page := fmt.Sprintf(ReadFile("resultpurchase.html"),"  : Digital Banking TA",x,"  : 3 Hari",strdetail,TotalAmount)
	base:=ReadFile("base.html")
	baseandpage:=fmt.Sprintf(base,itemlistcookie2,itemlistcookie,itemlisturl,page)
	addItemOnCookie2(w,r,idorder)
	fmt.Fprintf(w,baseandpage)
}
func order(w http.ResponseWriter,r *http.Request){
	checkCookie(w,r)
	key := r.URL.Query()["id"]
	id := key[0]
	db,_ := service.MySQL()
	itemlist := service.GetItemByOrderId(db,id)
	strdetail := ""
	TotalAmount := 0
	str:=""
	str2:=""
	background:=""
	for i,item:= range (itemlist){
		if(i%2==0){
			background ="style=\"background-color:#D7DBDD;\""
		}else {
			background ="style=\"background-color:#D1F2EB;\""
		}
		TotalAmount = TotalAmount+item.Price
		str2 ="<td "+background+"><a href=\"../getitemdetail?id="+fmt.Sprintf("%v\"",item.ID)+">"+item.Name+"</br>"+item.NamaToko+"</br>"+item.City+"</br>"+formatCurrency(item.Price)+"</a></td>"
		str ="<td style=\"vertical-align:top\">" +
			"<a href=\"../getitemdetail?id="+fmt.Sprintf("%v",item.ID)+"\"><img src=\""+item.Img+"\" alt=\""+item.Name+"\" style=\"width: 100%;height: auto;\" width=\"125\" height=\"150\"></a></td>"
		strdetail = strdetail + "<tr>"+str2+str+"</tr>"
	}
	var statustrx string
	if(itemlist[0].Status=="U"){
		statustrx = " Belum dibayar"
	}else{
		statustrx = " Lunas"
	}
	x:= ": 100010"+id+statustrx ;
	strdetail = "  : " +strdetail
	page := fmt.Sprintf(ReadFile("resultpurchase.html"),"  : Digital Banking TA",x,": "+itemlist[0].ExpiredDate.Format(time.RFC1123),strdetail,TotalAmount)
	base:=ReadFile("base.html")
	baseandpage:=fmt.Sprintf(base,itemlistcookie2,itemlistcookie,itemlisturl,page)
	fmt.Fprintf(w,baseandpage)
}
func deletecookieorder(w http.ResponseWriter,r *http.Request)  {
	key := r.URL.Query()["id"]
	checkCookie(w,r)
	var list2 []OrderStatus
	json.Unmarshal([]byte(strings.ReplaceAll(cookie2.Value,"@#","\"")), &list2)
	intkey ,_:= strconv.Atoi(key[0])
	final:=RemoveIndexOrder(list2,intkey)
	pagesJson, _ := json.Marshal(final)
	cookie2.Value=strings.ReplaceAll(string(pagesJson),"\"","@#");
	http.SetCookie(w,cookie2)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
func RemoveIndexOrder(s []OrderStatus, index int) []OrderStatus {
	return append(s[:index], s[index+1:]...)
}