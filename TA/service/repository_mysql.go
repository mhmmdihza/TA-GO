package service

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	username string = "TA1711520104"//"root"
	password string = "vam"//""
	database string = "tcp(127.0.0.1:3306)/ecomm_apps?parseTime=true"
)
type (
	JsonBirthDate time.Time
	// Mahasiswa
	Billing struct {
		ID        int       `json:"id"`
		CreatedDate time.Time `json:"created_date"`
		ExpiredDate int `json:"expired_date"`
		PaymentDate time.Time `json:"payment_date"`
		Amount int `json:"amount"`
		StatusPayment        string       `json:"status_payment"`
		PaymentTools        string       `json:"payment_tools"`
		Address        string       `json:"address"`
		PhoneNumber        string       `json:"phone_number"`
		Email        string       `json:"email"`
		Name        string       `json:"name"`
		DetailBilling []DetailBilling `json:"detail_billing"`
	}
	VaResp struct {
		VaNo string `json:"vaNo"`
		Amount string `json:"amount"`
	}
	DetailBilling struct{
		ID        int       `json:"id"`
		IDBilling       int       `json:"id_billing"`
		IDItem        int       `json:"id_item"`
	}
	Item struct{
		ID        int       `json:"id"`
		Name        string       `json:"name"`
		AddressSuplier        string       `json:"address_suplier"`
		Price int `json:"price"`
		Img        string       `json:"img"`
		Stock        int       `json:"stock"`
		NamaToko        string       `json:"nama_toko"`
		Category        string       `json:"category"`
		City        string       `json:"city"`
		Description        string       `json:"description"`
		Status string `json:"status"`
		ExpiredDate time.Time `json:"expired_date"`
	}
	PageFilterItem struct{
		Page int `json:"page"`
		Total int `json:"total"`
		Name string `json:"name"`
		Category string `json:"category"`
		Order string `json:"order"`
	}

)

var (
	db *sql.DB
	dsn = fmt.Sprintf("%v:%v@%v", username, password, database)
)

// HubToMySQL
func MySQL() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		panic(err)
		return nil, err
	}
	if(db==nil){
		fmt.Println("Db Nil")
	}
	return db, nil
}
func InsertItem(ctx context.Context, item Item) error {
	db, _ = MySQL()
	queryText := fmt.Sprintf("INSERT INTO %v ( name, address_supplier, price, img,stock,nama_toko,city) values('%v','%v',%v,'%v',%v,'%v','%v')", "item",
		item.Name,
		item.AddressSuplier,
		item.Price,
		item.Img,
		item.Stock,item.NamaToko,item.City)
	fmt.Println(queryText)
	_, err := db.ExecContext(ctx, queryText)
	if err != nil {
		db.Close()
		panic(err)
		return err
	}
	defer db.Close()
	return nil
}
func GetItemById(id int,db *sql.DB)Item{
	var i Item
	var querystr = fmt.Sprintf("SELECT id, name,address_supplier,price,img,stock,nama_toko" +
		",city,description FROM item where id = %v",id)
	fmt.Println(querystr)
	db.QueryRow(querystr).Scan(&i.ID, &i.Name,&i.AddressSuplier,&i.Price,&i.Img,&i.Stock,&i.NamaToko,&i.City,&i.Description)
	log.Println(i.Name)
	log.Println(i.Price)
	return i
}
func CountItemByFilter(page PageFilterItem,db *sql.DB)int64{
	var orderby string
	var filter string

	var total int64
	filter = ""
	if(len(page.Name)>0){
		words := strings.Fields(page.Name)
		filter = " AND ("
		for i,v := range words{
			if(i>0){
				filter = filter + " OR "
			}
			filter = filter+ "  LOWER(name) like lower('%"+v+"%') OR LOWER(nama_toko) like lower('%"+v+"%') " +
				"OR LOWER(city) like lower('%"+v+"%') OR LOWER(description) like lower('%"+v+"%') "
		}
		filter=filter+") "
	}
	if(len(page.Category)>0){
		filter = filter + " AND lower(category) = lower('"+page.Category+"')"
	}
	var lmt1 int
	var lmt2 int
	lmt2 = page.Total
	if(page.Page>0){
		lmt1 = (page.Page-1)*lmt2
	}else{
		lmt1 =page.Page
	}
	if(len(page.Order)>0){
		orderby = "order by " + page.Order
	}
	fmt.Println(lmt1)

	var querystr = fmt.Sprintf("SELECT count(*) as total " +
		" FROM item where 1 = 1 %v %v ",filter,orderby)
	fmt.Println(querystr)
	rows, err := db.Query(querystr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	for rows.Next() {
		rows.Scan(&total)
	}
	defer rows.Close()
	return total
}

func GetItemList(page PageFilterItem,db *sql.DB)[]Item{
	var limit string
	var orderby string
	var filter string
	if(len(page.Name)>0){
		filter = " AND ("
		words := strings.Fields(page.Name)
		for i,v := range words{
			if(i>0){
				filter = filter + " OR "
			}
			filter = filter+ "  LOWER(name) like lower('%"+v+"%') OR LOWER(nama_toko) like lower('%"+v+"%') " +
				"OR LOWER(city) like lower('%"+v+"%') OR LOWER(description) like lower('%"+v+"%') "
		}
		filter=filter+") "
	}
	if(len(page.Category)>0){
		filter = filter + " AND lower(category) = lower('"+page.Category+"')"
	}
	var lmt1 int
	var lmt2 int
	lmt2 = page.Total
	if(page.Page>0){
		lmt1 = (page.Page-1)*lmt2
	}else{
		lmt1 =page.Page
	}
	limit = fmt.Sprintf(" limit %d,%d ",lmt1,lmt2)

	if(len(page.Order)>0){
		orderby = "order by " + page.Order
	}

	var querystr = fmt.Sprintf("SELECT id, name,address_supplier,price,img,stock,nama_toko" +
		",city,description FROM item where 1 = 1 %v %v %v",filter,orderby,limit)
	fmt.Println(querystr)
	rows, err := db.Query(querystr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer rows.Close()

	var tw []Item
	count := 0
	for rows.Next() {
		var restw Item
		err := rows.Scan(&restw.ID,&restw.Name,&restw.AddressSuplier,&restw.Price,&restw.Img,&restw.Stock,&restw.NamaToko,&restw.City,&restw.Description)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		tw = append(tw,restw)
		count++;
	}
	return tw
}
func InsertBilling(ctx context.Context, b Billing) string {
	var vaAmount string
	var vaDesc string
	var vaNo string
	totalxp := b.ExpiredDate
	currentTime := time.Now()
	for i := 0 ; i < totalxp ; i++ {
		currentTime = currentTime.Add(24 * time.Hour)
	}
	var expiredDate = currentTime.Format(time.RFC3339)
	db, _ = MySQL()
	a := b.DetailBilling
	TotalAmount := 0
	vaDesc = "Pembelian "
	for _, s := range a {
		it := GetItemById(s.IDItem,db)
		TotalAmount = TotalAmount+it.Price
		vaDesc = vaDesc + "-"+it.Name+" "
	}
	vaAmount = fmt.Sprintf("%v",TotalAmount)


	queryText := fmt.Sprintf("INSERT INTO %v (expired_date,amount,status_payment,payment_tools" +
		",address,phone_number,email,name) values('%v',%v,'%v','%v','%v','%v','%v','%v'" +
		");","billing", expiredDate,
		TotalAmount,
		b.StatusPayment,
		b.PaymentTools,
		b.Address,
		b.PhoneNumber,
		b.Email,
		b.Name)

	var concat ="";

	fmt.Println(queryText)
	res, err := db.ExecContext(ctx, queryText)

	if err != nil {
		db.Close()
		panic(err)
		return ""
	}

	for _, s := range a {
		idheader,_:=res.LastInsertId()
		vaNo = fmt.Sprintf("100010%v",idheader)
		concat = concat + fmt.Sprintf("INSERT INTO %v (id_billing,id_item) values(%v,%v);", "detailbilling",
			idheader,
			s.IDItem)
		fmt.Println(concat)
		_, err = db.ExecContext(ctx, concat)
		if err != nil {
			db.Close()
			panic(err)
			return ""
		}
		concat = ""
	}
	vaName := "an "+b.Name+" - " +time.Now().Format(time.RFC822)
	var x = GeneratReqBodyForCreateVA(vaAmount,vaDesc ,vaName,vaNo)
	result := CreateVA(x)
	if!(extractValue(result,"message")=="Success Saving VA"){
		idh,_:= res.LastInsertId()
		queryText := fmt.Sprintf("DELETE FROM detailbilling where id_billing = %v",idh)

		fmt.Println(queryText)
		_, err := db.ExecContext(ctx, queryText)

		if err != nil {
			db.Close()
			panic(err)
			return ""
		}
		queryText = fmt.Sprintf("DELETE FROM billing where id = %v",idh)

		fmt.Println(queryText)
		_, err = db.ExecContext(ctx, queryText)

		if err != nil {
			db.Close()
			panic(err)
			return ""
		}
		defer db.Close()
		return extractValue(result,"message")
	}

	defer db.Close()

	return vaNo
}
func UpdatePayment(ctx context.Context, vaNo, amountva string) error {
	db, _ = MySQL()
	for i:=0 ;i<6;i++{
		vaNo = trimFirstRuneRep(vaNo)
	}
	queryText := fmt.Sprintf("Update billing set status_payment = 'P' , payment_date = current_timestamp() where id = '%v' and amount = %v",vaNo,amountva)
	fmt.Println(queryText)
	_, err := db.ExecContext(ctx, queryText)

	if err != nil {
		panic(err)
		db.Close()
		return err
	}
	queryGetListDetail := fmt.Sprintf("select id,id_billing,id_item from detailbilling where id_billing = %v", vaNo)
	fmt.Println(queryGetListDetail)
	rows, err := db.Query(queryGetListDetail)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer rows.Close()

	var tw []DetailBilling
	count := 0
	for rows.Next() {
		var restw DetailBilling
		err := rows.Scan(&restw.ID,&restw.IDBilling,&restw.IDItem)
		if err != nil {
			fmt.Println(err)
			return nil
		}
		tw = append(tw,restw)
		count++;
	}
	for _, s := range tw{
		queryText := fmt.Sprintf("update item set stock=(select stock from item where id = %v)-1 where id = %v ",s.IDItem,s.IDItem)
		fmt.Println(queryText)
		_, err := db.ExecContext(ctx, queryText)

		if err != nil {
			panic(err)
			db.Close()
			return err
		}
	}

	defer db.Close()
	return nil
}
func trimFirstRuneRep(s string) string {
	_, i := utf8.DecodeRuneInString(s)
	return s[i:]
}

func extractValue(body string, key string) string {
	keystr := "\"" + key + "\":[^,;\\]}]*"
	r, _ := regexp.Compile(keystr)
	match := r.FindString(body)
	keyValMatch := strings.Split(match, ":")
	return strings.ReplaceAll(keyValMatch[1], "\"", "")
}
func GetItemByOrderId(db *sql.DB,id string)[]Item{


	var querystr = fmt.Sprintf("SELECT i.id, i.name,i.address_supplier,i.price,i.img,i.stock,i.nama_toko" +
		",i.city,i.description,b.status_payment,b.expired_date FROM detailbilling d inner join item i on d.id_item = i.id  inner join billing b on b.id = d.id_billing  where d.id_billing = %v",id)
	fmt.Println(querystr)
	rows, err := db.Query(querystr)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer rows.Close()

	var tw []Item
	count := 0
	for rows.Next() {
		var restw Item
		err := rows.Scan(&restw.ID,&restw.Name,&restw.AddressSuplier,&restw.Price,&restw.Img,&restw.Stock,&restw.NamaToko,&restw.City,&restw.Description,&restw.Status,&restw.ExpiredDate)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		tw = append(tw,restw)
		count++;
	}
	return tw
}
