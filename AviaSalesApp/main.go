package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "postgres"
	dbname   = "aviadb"
)

type Booking struct {
	Id       int
	Booked   bool
	FlightId int64
}

type Flight struct {
	Id            int
	AirportIdFrom int64
	AirportIdTo   int64
	DepartureDate time.Time
	ArrivalDate   time.Time
	Competed      bool
}

// Simply did not have enought time to complete the full excersise :(
func main() {
	// connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	// open database
	db, err := sql.Open("postgres", psqlconn)
	CheckError(err)

	// close database
	defer db.Close()

	// check db
	err = db.Ping()
	CheckError(err)

	fmt.Println("Connected!")
}

func CheckError(err error) {
	if err != nil {
		panic(err)
	}
}

func BookTicket(db *sql.DB, id int64) {
	stmt, e := db.Prepare("INSERT INTO public.\"Bookings\"(\"Booked\", \"FlightId\")VALUES (?, ?);")
	ErrorCheck(e)

	res, e := stmt.Exec(true, id)
	ErrorCheck(e)

	id, e = res.LastInsertId()
	ErrorCheck(e)

	fmt.Println("Booked ticket: id", id)
}

func UnbookTicket(db *sql.DB, id int64) {
	stmt, e := db.Prepare("UPDATE public.\"Bookings\" SET \"Booked\" = false WHERE \"Id\" = ?;")
	ErrorCheck(e)

	res, e := stmt.Exec(true, id)
	ErrorCheck(e)

	id, e = res.LastInsertId()
	ErrorCheck(e)

	fmt.Println("Unbooked ticket: id", id)
}

func ChangeTicketFlightDates(db *sql.DB, id int64, newDate time.Time) {
	bookingRow, e := db.Query(fmt.Sprintf("SELECT * FROM public.\"Bookings\" WHERE \"Id\"=%d;", id))
	ErrorCheck(e)

	var booking = Booking{}
	e = bookingRow.Scan(&booking.Id, &booking.Booked, &booking.FlightId)
	ErrorCheck(e)
	fmt.Printf("Changing flight dates for booking %d. \n", booking.Id)

	initialFlightRow, e := db.Query(fmt.Sprintf("SELECT * FROM public.\"Flights\" WHERE \"Id\"=%d;", booking.FlightId))
	ErrorCheck(e)
	var initialFlight = Flight{}
	e = initialFlightRow.Scan(&initialFlight.Id, &initialFlight.AirportIdFrom, &initialFlight.AirportIdTo, &initialFlight.DepartureDate, &initialFlight.ArrivalDate, &initialFlight.Competed)
	ErrorCheck(e)
	fmt.Printf("Initial flight was for dates from %s to %s. \n", initialFlight.DepartureDate, initialFlight.ArrivalDate)

	// see if any ticket is available
	newFlightRow, e := db.Query(fmt.Sprintf("SELECT * FROM public.\"Flights\" WHERE \"DepartureDate\">=%s AND \"AirportIdFrom\"=%d AND \"AirportIdTo\"=%d;", newDate, initialFlight.AirportIdFrom, initialFlight.AirportIdTo))
	ErrorCheck(e)

	var newFlight = Flight{}
	e = newFlightRow.Scan(&newFlight.Id, &newFlight.AirportIdFrom, &newFlight.AirportIdTo, &newFlight.DepartureDate, &newFlight.ArrivalDate, &newFlight.Competed)
	ErrorCheck(e)

	stmt, e := db.Prepare(fmt.Sprintf("UPDATE public.\"Bookings\" SET \"FlightId\"=%d false WHERE \"Id\"=%d;", newFlight.Id, id))
	ErrorCheck(e)

	res, e := stmt.Exec(true, id)
	ErrorCheck(e)

	id, e = res.LastInsertId()
	ErrorCheck(e)

	fmt.Println("Unbooked ticket: id", id)
}

func ErrorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}
