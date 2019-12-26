package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

func main() {
	e := echo.New()

	// Middleware
	// e.Use(middleware.Logger())
	// e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Routes
	e.POST("/employees/add", createUser)
	e.GET("/employees/name/:name", getName)
	e.GET("/employees/:id", getUser)
	e.GET("/employees", getUserall)
	e.PUT("/employees/:id", updateUser)
	e.DELETE("/employees/:id", deleteUser)

	e.POST("/member/add", createMember)
	e.GET("/member/:id", getMember)
	e.GET("/member", getMemberAll)
	e.PUT("/member/:id", updateMember)
	e.DELETE("/member/:id", deleteMember)
	e.POST("/member/login", login)

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}

func login(c echo.Context) error {

	mem := new(Member)

	fmt.Println("Username " + mem.Username)
	if err := c.Bind(mem); err != nil { //ตัวเก็บค่า
		return err
	}
	fmt.Println("Username " + mem.Username)

	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		fmt.Println(err) //ไลบารี่ปริ้นแสดงคอนโซ
		panic(err)
	}
	err = db.Ping() //เช็คการเชื่อมต่อ
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Succesfully connected! DB")

	var list Member

	sqlStatment := `SELECT id,username , password FROM member WHERE username = $1 AND password = $2;`
	err = db.QueryRow(sqlStatment, mem.Username, mem.Password).Scan(&list.ID, &list.Username, &list.Password)

	if err != nil {
		fmt.Println(err)

		return c.JSON(http.StatusOK, "Login failed!")
	}

	if list.Username == mem.Username && list.Password == mem.Password {
		fmt.Println("Login succesful!")

		return c.JSON(http.StatusOK, list)

	} else {

		fmt.Println("Login failed!")

		return c.JSON(http.StatusOK, "Login failed!")
	}

}

func createMember(c echo.Context) error {
	mem := new(Member)
	if err := c.Bind(mem); err != nil {
		return err
	}

	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		fmt.Println(err) //ไลบารี่ปริ้นแสดงคอนโซ
		panic(err)
	}
	err = db.Ping() //เช็คการเชื่อมต่อ
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	sqlStatment := `INSERT INTO member (username, password, name, email, telephone, role) VALUES ($1,$2,$3,$4,$5,$6);`

	_, err = db.Exec(sqlStatment, mem.Username, mem.Password, mem.Name, mem.Email, mem.Telephone, mem.Role)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer db.Close()

	return c.JSON(http.StatusOK, "Create Succes")

}

func getMember(c echo.Context) error {

	id := c.Param("id")

	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		fmt.Println(err) //ไลบารี่ปริ้นแสดงคอนโซ
		panic(err)
	}
	err = db.Ping() //เช็คการเชื่อมต่อ
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Succesfully connected! DB")

	var list Member

	sqlStatment := `SELECT * FROM member WHERE id=$1;`
	err = db.QueryRow(sqlStatment, id).Scan(&list.ID, &list.Username, &list.Password, &list.Name, &list.Email, &list.Telephone, &list.Role)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer db.Close()

	return c.JSON(http.StatusOK, list)

}

func getMemberAll(c echo.Context) error {
	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Println("Succesfully connected! DB")

	sqlStatment := `SELECT * FROM member;`
	result, err := db.Query(sqlStatment)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	var list []Member

	for result.Next() {
		var data Member
		err := result.Scan(
			&data.ID,
			&data.Username,
			&data.Password,
			&data.Name,
			&data.Email,
			&data.Telephone,
			&data.Role)
		if err != nil {
			return err
		}
		list = append(list, data)
	}

	defer result.Close()
	defer db.Close()

	return c.JSON(http.StatusOK, list)
}

func updateMember(c echo.Context) error {
	id := c.Param("id")
	mem := new(Member)
	if err := c.Bind(mem); err != nil {
		return err
	}

	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		fmt.Println(err) //ไลบารี่ปริ้นแสดงคอนโซ
		panic(err)
	}
	err = db.Ping() //เช็คการเชื่อมต่อ
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	sqlStatment := `UPDATE member SET username=$1, password=$2, name=$3, email=$4, telephone=$5, role =$6 WHERE id=$7;`
	_, err = db.Exec(sqlStatment, mem.Username, mem.Password, mem.Name, mem.Email, mem.Telephone, mem.Role, id)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer db.Close()

	return c.JSON(http.StatusOK, "UPDATE Succes")

}

func deleteMember(c echo.Context) error {
	id := c.Param("id")

	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		fmt.Println(err) //ไลบารี่ปริ้นแสดงคอนโซ
		panic(err)
	}
	err = db.Ping() //เช็คการเชื่อมต่อ
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	sqlStatment := `DELETE FROM member WHERE id=$1 ;`
	_, err = db.Exec(sqlStatment, id)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer db.Close()

	return c.JSON(http.StatusOK, "Delete Succes")

}

type (
	//Member
	Member struct {
		ID        int    `json:"id"`
		Username  string `json:"username"`
		Password  string `json:"password"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		Telephone string `json:"telephone"`
		Role      string `json:"role"`
	}
)

func createUser(c echo.Context) error {
	emp := new(employee)
	if err := c.Bind(emp); err != nil {
		return err
	}

	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		fmt.Println(err) //ไลบารี่ปริ้นแสดงคอนโซ
	}
	err = db.Ping() //เช็คการเชื่อมต่อ
	if err != nil {
		fmt.Println(err)
	}

	sqlStatment := `INSERT INTO employee (name, email) VALUES ($1, $2);`
	_, err = db.Exec(sqlStatment, emp.Name, emp.Email)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	return c.JSON(http.StatusOK, "Create Succes")

}

func getName(c echo.Context) error {
	name := c.Param("name")
	fmt.Println("name=" + name)

	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		fmt.Println(err) //ไลบารี่ปริ้นแสดงคอนโซ
	}
	err = db.Ping() //เช็คการเชื่อมต่อ
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Succes connected! DB")

	var employee employee

	sqlStatment := `SELECT id,name,email FROM employee WHERE name=$1;`
	err = db.QueryRow(sqlStatment, name).Scan(&employee.ID, &employee.Name, &employee.Email)

	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()
	fmt.Println("name=" + employee.Name)
	fmt.Println("email=" + employee.Email)
	return c.JSON(http.StatusOK, employee)

}

func getUser(c echo.Context) error {
	id := c.Param("id")

	//คอนฟิก ดาต้าเบส
	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		fmt.Println(err) //ไลบารี่ปริ้นแสดงคอนโซ
	}
	err = db.Ping() //เช็คการเชื่อมต่อ
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Succes connected! DB")

	var employee employee

	sqlStatment := `SELECT * FROM employee WHERE id=$1;`
	err = db.QueryRow(sqlStatment, id).Scan(&employee.ID, &employee.Name, &employee.Email)

	if err != nil {
		fmt.Println(err)
	}

	defer db.Close()

	return c.JSON(http.StatusOK, employee)

}

func getUserall(c echo.Context) error {
	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		fmt.Println(err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Succesfully connected! DB")

	sqlStatment := `SELECT * FROM employee;`
	result, err := db.Query(sqlStatment)
	if err != nil {
		fmt.Println(err)
	}
	var employees []employee

	for result.Next() {
		var data employee
		err := result.Scan(
			&data.ID,
			&data.Name,
			&data.Email)
		if err != nil {
			return err
		}
		employees = append(employees, data)
	}

	defer result.Close()
	defer db.Close()

	return c.JSON(http.StatusOK, employees)
}

func updateUser(c echo.Context) error {
	id := c.Param("id")
	emp := new(employee)
	if err := c.Bind(emp); err != nil {
		return err
	}

	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		fmt.Println(err) //ไลบารี่ปริ้นแสดงคอนโซ
	}
	err = db.Ping() //เช็คการเชื่อมต่อ
	if err != nil {
		fmt.Println(err)
	}

	sqlStatment := `UPDATE employee SET name =$1, email=$2 WHERE id =$3;`
	_, err = db.Exec(sqlStatment, emp.Name, emp.Email, id)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	return c.JSON(http.StatusOK, "UPDATE Succes")

}

func deleteUser(c echo.Context) error {
	id := c.Param("id")

	psql := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psql)
	if err != nil {
		fmt.Println(err) //ไลบารี่ปริ้นแสดงคอนโซ
	}
	err = db.Ping() //เช็คการเชื่อมต่อ
	if err != nil {
		fmt.Println(err)
	}

	sqlStatment := `DELETE FROM employee WHERE id=$1 ;`
	_, err = db.Exec(sqlStatment, id)
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	return c.JSON(http.StatusOK, "Delete Succes")

}

type (
	employee struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	// Member struct {
	// 	ID        int    `json:"id"`
	// 	Username  string `json:"username"`
	// 	Password  string `json:"password"`
	// 	Name      string `json:"name"`
	// 	Email     string `json:"email"`
	// 	Telephone string `json:"telephone"`
	// 	Role      string `json:"role"`
	// }
)

const (
	host     = "34.68.7.41"
	port     = 81
	user     = "postgres"
	password = "1234qwer"
	dbname   = "employeedatabase"
)
