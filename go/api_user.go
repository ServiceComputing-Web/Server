/*
 * Minimalist Blog
 *
 * This is a sample Blog server.
 *
 * API version: 1.0.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */

package swagger

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Response struct {
	Data string `json:"data"`
}

type Token struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func LogOut(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func SignIn(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "testuser:123@tcp(mysql:3306)/?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var user User

	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		response := ErrorResponse{err.Error()}
		JsonResponse(response, w, http.StatusBadRequest)
		return
	}

	query, err := db.Query("select * from test.User where username='" + user.Username + "'")
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()
	v, err := getJSON(query)
	if err != nil {
		log.Fatal(err)
	}
	if string(v) == "[]" {
		reponse := ErrorResponse{"Invalid username/password supplied"}
		JsonResponse(reponse, w, http.StatusNotFound)
		return
	}

	var userQuery User
	v = v[1 : len(v)-1]
	_ = json.Unmarshal(v, &userQuery)

	if userQuery.Password != user.Password {
		response := ErrorResponse{"Invalid username/password supplied"}
		JsonResponse(response, w, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func SignUp(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "testuser:123@tcp(mysql:3306)/?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	var user User
	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil || user.Password == "" || user.Username == "" {
		response := ErrorResponse{"Invalid username/password supplied"}
		JsonResponse(response, w, http.StatusBadRequest)
		return
	}

	query, err := db.Query("select * from test.User where username='" + user.Username + "'")
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	if query.Next() {
		response := ErrorResponse{"User Exists"}
		JsonResponse(response, w, http.StatusBadRequest)
		return
	}

	query, err = db.Query("INSERT INTO `test`.`User` (`username`, `password`) VALUES ('" + user.Username + "', '" + user.Password + "')")
	if err != nil {
		log.Fatal(err)
	}
	defer query.Close()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
}

func JsonResponse(response interface{}, w http.ResponseWriter, code int) {
	json, err := json.Marshal(response)
	if err != nil {
		log.Fatal(err)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	w.Write(json)
}

func getJSON(rows *sql.Rows) ([]byte, error) {
	columns, err := rows.Columns()
	if err != nil {
		return []byte(""), err
	}
	count := len(columns)

	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return []byte(""), err
	}
	return jsonData, nil
}
