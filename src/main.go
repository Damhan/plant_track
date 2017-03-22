package main

import (
  "bufio"
  "os"
  "fmt"
  "strings"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  )

//also simple input/trim function for reading input and stripping newlines
func readTrim(r *bufio.Reader) string {
  line, _ := r.ReadString('\n')
  line = strings.Replace(line,"\n","",-1)
  return line
}

//simple error checking function that outputs any issues
func checkErr(err error) {
  if err != nil {
    panic(err)
  }
}

//creates database if it doesn't exist and connects to it.
func createConnect() *sql.DB {
  db, err := sql.Open("sqlite3", "./plants.db")
  checkErr(err)
  return db
}

//function that creates required tables if they don't exist
func createTables(db *sql.DB) {
  _, err := db.Exec("CREATE TABLE IF NOT EXISTS vegPlants(ID TEXT, Plant TEXT,PlantDate TEXT)")
  checkErr(err)
}

//function that handles inserting values into a dynamic table of choice from main.
func table_insertion(db *sql.DB, table string, insert_info string) {
  var plant_info = strings.Split(insert_info,"/") //split our info to be inserted into an array
  stmt, err:= db.Prepare("INSERT INTO vegPlants(ID,Plant,PlantDate) VALUES(?,?,?)")
  checkErr(err)
  _,err = stmt.Exec(plant_info[0],plant_info[1],plant_info[2])
  checkErr(err)
}

func tableViewing(db *sql.DB,table_choice string) {
  rows, err := db.Query("SELECT ID,Plant,PlantDate,Cast((JulianDay()-JulianDay(PlantDate)) AS INT) FROM vegPlants")
  checkErr(err)

  var ID string
  var Plant string
  var PlantDate string
  var Age string

  fmt.Println("ID|PlantName|PlantDate|PlantAge")

  for rows.Next() {
    err = rows.Scan(&ID,&Plant,&PlantDate,&Age)
    checkErr(err)

  fmt.Print(ID + " | ")
  fmt.Print(Plant + " | ")
  fmt.Print(PlantDate + " | ")
  fmt.Print(Age)
  fmt.Println(" ")
  }
}
//main
func main() {
  //setup
  db := createConnect()
  createTables(db)
  reader := bufio.NewReader(os.Stdin)

  //program loop
  for {

    fmt.Println("Press i to insert values, press q to quit or v to view tables.")
    fmt.Print("What would you like to do? ")
    choice := readTrim(reader)

    //insert if i
    if strings.Compare(choice,"i") == 0{
      fmt.Print("What table would you like to insert to? ")
      table_choice := readTrim(reader)
      fmt.Println("Please insert in the form; 'id/plant/plantdate' followed by an enter ")
      fmt.Println("Please make sure the plant date is in the form YYYY-MM-DD")
      insert_info := readTrim(reader)
      table_insertion(db,table_choice,insert_info)
    }

    //quit if q
    if strings.Compare(choice,"q") == 0 {
      break
    }

    if strings.Compare(choice,"v") == 0 {
      fmt.Print("What table would you like to view? ")
      table_choice := readTrim(reader)
      tableViewing(db,table_choice)

    }

  }
}
