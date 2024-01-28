package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// define the path of the wkhtmltopdf
const path = "/usr/local/bin/wkhtmltopdf"

type AccountStatements struct {
	Date               string `json:"date"`
	PaymentTypeDetails string `json:"paymentTypeDetails"`
	PaidEst            string `json:"paidEst"`
	PaidDue            string `json:"paidDue"`
	Balance            string `json:"balance"`
}

type Person struct {
	FirstName string
	LastName  string
	JobTitle  string
	Skills    []string
}

func GeneratePDF(htmlTemplate string, person *Person) ([]byte, error) {
	var htmlBuffer bytes.Buffer
	err := template.Must(template.New("statement").
		Parse(htmlTemplate)).
		Execute(&htmlBuffer, person)
	if err != nil {
		return nil, err
	}

	// set the predefined path in the wkhtmltopdf's global state
	wkhtmltopdf.SetPath(path)

	// Initialize the converter
	pdfGen, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}

	page := wkhtmltopdf.NewPageReader(bytes.NewReader(htmlBuffer.Bytes()))
	page.EnableLocalFileAccess.Set(true)
	pdfGen.AddPage(page)

	// Set PDF page attributes
	// Refer to http://wkhtmltopdf.org/usage/wkhtmltopdf.txt for more information
	pdfGen.MarginLeft.Set(10)
	pdfGen.MarginRight.Set(10)
	pdfGen.Dpi.Set(300)
	pdfGen.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfGen.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	// Generate the PDF
	err = pdfGen.Create()
	if err != nil {
		return nil, err
	}
	return pdfGen.Bytes(), nil
}

func main() {
	person := &Person{
		FirstName: "Simon Peter",
		LastName:  "Damian",
		JobTitle:  "Software Developer",

		Skills: []string{
			"Go",
			"Ruby",
			"Python",
			"JavaScript",
		},
	}
	htmlTemplate := `
		<!DOCTYPE html>
<html lang="en">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
        }

        .letter-container {
            max-width: 600px;
            margin: auto;
            margin-top: 20px;
            padding: 20px;
        }
       
        .statement{
            color: red;
        }

        .address {
            float: left;
            width: 45%;
            margin-bottom: 20px;
        }

        .text {
            clear: both;
            margin-top: 20px;
        }
        .table-one {
            width: 100%;
            border-collapse: collapse;
            margin-top: 20px;
            background-color: rgb(202, 195, 195);
        }

        th, td {
            border: 1px solid #ccc;
            text-align: left;
        }
        h6{
            margin-top: 20px;
        }
        .account{
            font-size: 10px;
        }
       .code{
        display: flex;
        flex-direction: row;
        justify-content: space-between;
       }
       strong{
        font-size: 12px;
       }
       h1{
        font-size: 40px;
       }
       img{
        width: 100px;
        height: 100px;
        margin-left: 400px;
       
       }
       .left-address{
        margin-top: 5rem;
       }
    </style>
</head>

<body>
    <img src="png.jpeg" >
    <div class="letter-container">

        <div class="address left-address">
            
            Mrs I P Jeses<br>
            April Cottage<br>
            18 East street <br>
            Crawley<br>
            West Success<br>
            PH10D 4PU
        </div>

        <div class="address">
            <strong>HBSC Advance:</strong><br>
            Contact tel 03457404404<br>
            Text phone 03457125563<br>
            www.hbsc.co.uk
            <h3 class="statement">Your Statement</h3>
            <div>
                <table class="table-one">
                    <thead>
                        <tr>
                            <th>Summary</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr>
                            <td>Opening Balance</td>
                            <td>$100.00</td>
                        </tr>
                        <tr>
                            <td>Payments in</td>
                            <td>$200.00</td>
                        </tr>
                        <tr>
                            <td>Payments out</td>
                            <td>$200.00</td>
                        </tr>
                        <tr>
                            <td>Closing Balance</td>
                            <td>$200.00</td>
                        </tr>
                        <tr>
                            <td>Arrage withdrawal limit</td>
                            <td>$2000.00</td>
                        </tr>
                    </tbody>
                </table>
            </div>
        </div>

        <div class="letter-container">
            <div class="address">
                <h6>15 July to 14 August 2019</h6>
                <strong>Account name</strong><br>
                Mrs Lisa Patricia Jesus<br>
                
            </div>   
            <div class="address">
                <strong class="account">International Bank Account Number</strong>
                GB3659FJFKJFGJKG<br>
                <strong class="account">Branch Identifier Code</strong> <br>
                HBUK3455545
                <div class="code">
                    <div>
<strong>Sortcode </strong><br>
47847<br>
                    </div>   
                    <div>
                        <strong>Account Number</strong><br>
                        47847<br>
                    </div>
                <div>
                    <strong>Sheet Number</strong><br>
                    47847<br>
                </div>  
                </div>
        </div>
    </div>
    <div>
        <table class="table-two">
            <h1>Your HBSC Advance details</h1>
            <thead>
                <tr>
                    <th>Date</th>
                    <th>Payment ups and details</th>
                    <th>Paid est</th>
                    <th>Paid due</th>
                    <th>Balance</th>
                </tr>
				
            </thead>
             <tbody>
               <tr>
                    <td>{{.date}}</td>
                    <td>{{.paymentTypeDetails}}</td>
                    <td>{{.paidEst}}</td>
                    <td>{{.paidDue}}</td>
                    <td>{{.balance}}</td>
                </tr>
               
            </tbody>
        </table>
    </div>
</body>

</html>

	`

	htmlTemplate2 := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
		</head>
		<body>
			<h1>{{.FirstName}} {{.LastName}}</h1>
			<p>{{.JobTitle}}</p>

			<h2>Skills</h2>
			<ul>
				{{range .Skills}}
				<li>{{.}}</li>
				{{end}}
			</ul>

			<a href="/download">Download PDF</a>
		</body>
		</html>
	`

	// render the HTML template on the index route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := template.Must(template.New("profile").
			Parse(htmlTemplate2)).
			Execute(w, person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// serve the PDF on the download route
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {

		fileContent, err := os.Open("account-statement.json")

		if err != nil {
			log.Fatal(err)
			return
		}

		fmt.Println("The File is opened successfully...")

		defer fileContent.Close()

		byteResult, _ := io.ReadAll(fileContent)

		var account []AccountStatements

		err2 := json.Unmarshal(byteResult, &account)

		if err2 != nil {
			fmt.Println("Error JSON Unmarshalling")
			fmt.Println(err2.Error())
		}

		for _, x := range account {
			fmt.Printf("%s %s \n", x.Date, x.Balance)
		}

		pdf, err := GeneratePDF(htmlTemplate, person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=account-statement.pdf")
		w.Write(pdf)
	})

	println("listening on port", os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)

}
