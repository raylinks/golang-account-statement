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

func GeneratePDF(htmlTemplate string, statement *AccountStatements) ([]byte, error) {
	var htmlBuffer bytes.Buffer
	err := template.Must(template.New("statement").
		Parse(htmlTemplate)).
		Execute(&htmlBuffer, statement)
	if err != nil {
		return nil, err
	}

	// Initialize the converter
	pdfGen, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(htmlBuffer.Bytes()))
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
	person := &AccountStatements{
		Date:               "01/04/2023",
		PaymentTypeDetails: "VOUCHER",
		PaidEst:            "2000",
		PaidDue:            "5400",
		Balance:            "7234",
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
    <img src="../bird.svg" alt="SVG Image">
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
            </tbody>
        </table>
    </div>
</body>

</html>

	`

	// render the HTML template on the index route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := template.Must(template.New("statement").
			Parse(htmlTemplate)).
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
		w.Header().Set("Content-Disposition", "attachment; filename=statement.pdf")
		w.Write(pdf)
	})
	// println("listening on port", os.Getenv("PORT"))
	// http.ListenAndServe(":"+os.Getenv("PORT"), nil)
	println("listening on port", 8080)
	http.ListenAndServe("localhost:8080", nil)
}

//err := router.Run("localhost:8080")
// 	if err != nil {
// 		log.Fatal(err)
// 	}