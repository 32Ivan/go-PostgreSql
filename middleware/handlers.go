package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-PostgresSql/models"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type response struct {
	ID      int64  `json:"id,omitempty"`
	Massage string `json:"message.omitempty"`
}

func createConnection() *sql.DB {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")

	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))

	if err != nil {
		panic(err)
	}

	err = db.Ping()

	if err != nil {
		panic(err)
	}

	fmt.Println("Successfulluy connected to postgres")
	return db
}

func CreateStock(w http.ResponseWriter, r *http.Request) {
	var stock models.Stock

	err := json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)

	}

	insertID := insertStock(stock)

	res := response{
		ID:      insertID,
		Massage: "sock created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

func GetStock(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert the strnig into int. %v", err)
	}

	stock, err := getStock(int64(id))

	if err != nil {
		log.Fatalf("Unable to get stock. %v", err)

	}

	json.NewEncoder(w).Encode(stock)

}

func GetAllStock(w http.ResponseWriter, r *http.Request) {

	stocks, err := getAllStock()

	if err != nil {
		log.Fatalf("Unable to get all stocks %v", err)
	}

	json.NewEncoder(w).Encode(stocks)
}

func UpdateStock(w http.ResponseWriter, r *http.Request) {

	parmas := mux.Vars(r)

	id, err := strconv.Atoi(parmas["id"])

	if err != nil {
		log.Fatalf("Unable to convert the string into int %v", err)

	}

	var stock models.Stock

	err = json.NewDecoder(r.Body).Decode(&stock)

	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	updatedRows := updateStock(int64(id), stock)

	msg := fmt.Sprintf("Stock updated successfully. Total rows/records affected %v", updatedRows)

	res := response{
		ID:      int64(id),
		Massage: msg,
	}

	json.NewEncoder(w).Encode(res)

}

func DeleteStock(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)

	id, err := strconv.Atoi(params["id"])

	if err != nil {
		log.Fatalf("Unable to convert string to int. %v", err)
	}

	deleteRows := deleteStock(int64(id))

	msg := fmt.Sprintf("Stock deleted successfully. total rows/records %v", deleteRows)

	res := response{
		ID:      int64(id),
		Massage: msg,
	}

	json.NewEncoder(w).Encode(res)
}

func insertStock(stock models.Stock) int64 {

	db := createConnection()

	defer db.Close()
	sqlStatment := `INSERT INTO stocks(name, price, company) VALUES ($1, $2, $3) RETURNING stockid`

	var id int64

	err := db.QueryRow(sqlStatment, stock.Name, stock.Price, stock.Company).Scan(&id)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single record. %v", err)
	return id
}

func getStock(id int64) (models.Stock, error) {

	db := createConnection()

	defer db.Close()

	var stock models.Stock

	sqlStatment := `SELECT * FROM stocks WHERE stockid = $1`

	row := db.QueryRow(sqlStatment, id)

	err := row.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("Now rows were retuyrned !")
		return stock, nil
	case nil:
		return stock, nil
	default:
		log.Fatalf("Unable to scan the ro. %v", err)
	}
	return stock, err

}

func getAllStock() ([]models.Stock, error) {
	db := createConnection()

	defer db.Close()

	var stocks []models.Stock

	sqlStatment := `SELECT * FROM stocks`

	rows, err := db.Query(sqlStatment)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)

	}

	defer rows.Close()

	for rows.Next() {
		var stock models.Stock

		err = rows.Scan(&stock.StockID, &stock.Name, &stock.Price, &stock.Company)

		if err != nil {
			log.Fatalf("Unable to scan the row %v", err)
		}

		stocks = append(stocks, stock)

	}
	return stocks, err
}

func updateStock(id int64, stock models.Stock) int64 {

	db := createConnection()

	defer db.Close()

	sqlStatment := `UPDATE stocks SET name=$2, price=$3, company=$4 WHERE stockid = $1`

	res, err := db.Exec(sqlStatment, id, stock.Name, stock.Price, stock.Company)

	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows, %v", err)
	}

	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected
}

func deleteStock(id int64) int64 {
	db := createConnection()

	defer db.Close()

	sqlStatment := `DELETE FROM stocks WHERE stockid=$1`

	res, err := db.Exec(sqlStatment, id)

	if err != nil {
		log.Fatalf("Unable to execute the query %v", err)
	}

	rowsAffected, err := res.RowsAffected()

	if err != nil {
		log.Fatalf("Error while checking the affected rows, %v", err)
	}

	fmt.Printf("Total rows/records affected %v", rowsAffected)
	return rowsAffected

}
