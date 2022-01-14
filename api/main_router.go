package api

import (
	"encoding/json"
	"firstProject/config"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/itrepablik/itrlog"
)

// MainRouters are the collection of all URLs for the Main App.
func MainRouters(r *mux.Router) {
	r.HandleFunc("/", Home).Methods("GET")
	r.HandleFunc("/convert", Convert).Methods("POST")
}

type contextData map[string]interface{}

func Convert(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusOK)

	body, errBody := ioutil.ReadAll(r.Body)
	if errBody != nil {
		itrlog.Error(errBody)
		panic(errBody.Error())
	}

	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)

	amount := keyVal["value"]
	source := keyVal["source"]
	dest := keyVal["dest"]

	const kenyanShellingRate = 113.39
	const ghanaCedisRate = 6.2
	const nigerianNiaraRate = 414.10

	var sourceRate = 0.0
	var destRate = 0.0

	if len(amount) == 0 {
		w.Write([]byte(`{ "isSuccess": "false", "AlertTitle": "Amount is Required", "AlertMsg": "Please enter amount.", "AlertType": "error" }`))
		return
	}

	i, _ := strconv.ParseFloat(amount, 64)
	if i == 0 {
		w.Write([]byte(`{ "isSuccess": "false", "AlertTitle": "Amount is Required", "AlertMsg": "Please enter amount greater than zero.", "AlertType": "error" }`))
		return
	}
	if len(source) == 0 {
		w.Write([]byte(`{ "isSuccess": "false", "AlertTitle": "Source Currency is Required", "AlertMsg": "Please Select Source Currency.", "AlertType": "error" }`))
		return
	}

	if len(dest) == 0 {
		w.Write([]byte(`{ "isSuccess": "false", "AlertTitle": "Destination Currency is Required", "AlertMsg": "Please Select Destination Currency.", "AlertType": "error" }`))
		return
	}

	if source == dest {
		w.Write([]byte(`{ "isSuccess": "false", "AlertTitle": "Destination Currency cannot be the same as Source Currency", "AlertMsg": "Please Select different Destination Currency.", "AlertType": "error" }`))
		return
	}
	var source_description string
	var dest_description string

	switch source {
	case "NGN":
		source_description = "Nigerian Niara (NGN)"
		sourceRate = nigerianNiaraRate

	case "GHS":
		source_description = "Ghana Cedis (GHS)"
		sourceRate = ghanaCedisRate
	case "KES":
		source_description = "Kenyan Shellings (KSH)"
		sourceRate = kenyanShellingRate
	default:
		source_description = "Unknown"
	}

	switch dest {
	case "NGN":
		dest_description = "Nigerian Niara (NGN)"
		destRate = nigerianNiaraRate
	case "GHS":
		dest_description = "Ghana Cedis (GHS)"
		destRate = ghanaCedisRate
	case "KES":
		dest_description = "Kenyan Shellings (KSH)"
		destRate = kenyanShellingRate
	default:
		dest_description = "Unknown"
	}

	// url := "https://free-currency-converter.herokuapp.com/list/convert?source=" + source + "&destination=" + dest + "&price=" + amount
	// resp, err := http.Get(url)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// defer resp.Body.Close()
	// bodyBytes, _ := ioutil.ReadAll(resp.Body)

	// type Response struct {
	// 	Success         bool
	// 	Message         string
	// 	Source          string
	// 	Destination     string
	// 	Price           string
	// 	Converted_value float64
	// }

	// var responseObject Response
	// json.Unmarshal(bodyBytes, &responseObject)

	amount_to_convert, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Fatalln(err)
	}
	calculated_amount := (destRate / sourceRate) * amount_to_convert

	w.Write([]byte(`{ "isSuccess": "true","converted_amount":"` + strconv.FormatFloat(calculated_amount, 'f', 6, 64) + `","value": "` + amount + `","sourceDescription": "` + source_description + `","destDescription": "` + dest_description + `", "alertMsg": "Amount converted successfully",
		"alertType": "success"}`))
}

// Home function is to render the homepage page.
func Home(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(config.SiteRootTemplate+"dashboard/index.html", config.SiteHeaderTemplate, config.SiteFooterTemplate))

	data := contextData{
		"PageTitle":    "Convertor Web App",
		"PageMetaDesc": config.SiteSlogan,
		"CanonicalURL": r.RequestURI,
		"CsrfToken":    csrf.Token(r),
	}
	tmpl.Execute(w, data)
}
