package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) testDatabase(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Println("Testing the database...")
	fmt.Println("Test 1: Insertion")
	_, err := rt.db.Insert("users", "('Pippo', 'Paperino', Null)")
	if err != nil {
		insertion_error := BackendError{
			Affinity: "Database testing",
			Message:  "Database insertion test has failed. There might be a problem with string formatting, or insertion formatting",
			OG_error: err,
		}
		fmt.Println(insertion_error.Error())
	}

	fmt.Println("Test 2: Updating values")
	_, err = rt.db.Update("users", "id = 'Pluto'", "1=1")
	if err != nil {
		update_error := BackendError{
			Affinity: "Database testing",
			Message:  "Database update test has failed. There might be a problem with string formatting, or update formatting",
			OG_error: err,
		}
		fmt.Println(update_error.Error())
	}

	fmt.Println("Test 3: Selection")
	selection, err := rt.db.Select("*", "users", "1=1")
	if err != nil {
		selection_error := BackendError{
			Affinity: "Database testing",
			Message:  "Database selection test has failed. There might be a problem with string formatting, or selecton formatting",
			OG_error: err,
		}
		fmt.Println(selection_error.Error())
	}
	answer, err := UsersRowReading(selection)
	if err != nil {
		reading_error := BackendError{
			Affinity: "Database testing",
			Message:  "Database reading test has failed. There might be a problem with sql row reading",
			OG_error: err,
		}
		fmt.Println(reading_error.Error())
	} else {
		fmt.Println("Selected elements are the following:")
		for _, el := range answer {
			fmt.Println(el)
		}
	}

	w.Header().Set("content-type", "text/plain")
	w.WriteHeader(http.StatusNoContent)
	_, _ = w.Write([]byte("Tests executed!"))
}
