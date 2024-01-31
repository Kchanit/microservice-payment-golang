package utils

import (
	"log"
	"strconv"
	"time"

	"github.com/Kchanit/microservice-payment-golang/internal/core/domain"
	generator "github.com/angelodlfrtr/go-invoice-generator"
	"github.com/joho/godotenv"
)

func LoadSecret() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

}
func GenerateInvoice(customer domain.User, products []domain.Product, ref string) (string, error) {
	doc, _ := generator.New(generator.Invoice, &generator.Options{
		TextTypeInvoice:  "PIXELMAN",
		AutoPrint:        true,
		CurrencySymbol:   " ",
		CurrencyThousand: ",",
	})

	doc.SetHeader(&generator.HeaderFooter{
		Text:       "<center>PIXELMAN INVOICE.</center>",
		Pagination: true,
	})

	doc.SetFooter(&generator.HeaderFooter{
		Text:       "<center>Cupcake ipsum dolor sit amet bonbon. I love croissant cotton candy. Carrot cake sweet I love sweet roll cake powder.</center>",
		Pagination: true,
	})

	doc.SetRef(ref)
	// doc.SetVersion("someversion")

	doc.SetDescription("A description")
	doc.SetNotes("I love croissant cotton candy. Carrot cake sweet I love sweet roll cake powder! I love croissant cotton candy. Carrot cake sweet I love sweet roll cake powder! I love croissant cotton candy. Carrot cake sweet I love sweet roll cake powder! I love croissant cotton candy. Carrot cake sweet I love sweet roll cake powder! ")

	doc.SetDate(time.Now().Format("02/01/2006"))
	// doc.SetPaymentTerm("02/04/2021")

	// logoBytes, err := ioutil.ReadFile("./example_logo.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }

	doc.SetCompany(&generator.Contact{
		Name: "Pixelman Company",
		// Logo: logoBytes,
		Address: &generator.Address{
			Address:    "89 Rue de Brest",
			Address2:   "Appartement 2",
			PostalCode: "10110",
			City:       "Bangkok",
			Country:    "Thailand",
		},
	})

	doc.SetCustomer(&generator.Contact{
		Name: customer.Name,
		Address: &generator.Address{
			Address:    customer.Addresses[0].Address,
			PostalCode: customer.Addresses[0].PostalCode,
			City:       customer.Addresses[0].City,
			Country:    customer.Addresses[0].Country,
		},
	})

	for _, product := range products {
		doc.AppendItem(&generator.Item{
			Name:        product.Name,
			Description: product.Description,
			UnitCost:    strconv.Itoa(int(product.Price)),
			Quantity:    strconv.Itoa(product.Quantity),
		})
	}

	doc.SetDefaultTax(&generator.Tax{
		Percent: "7",
	})

	// doc.SetDiscount(&generator.Discount{
	// 	Percent: "10",
	// })

	pdf, err := doc.Build()
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	// Implement upload to minio
	outputName := "out.pdf"
	err = pdf.OutputFileAndClose(outputName)

	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return outputName, nil
}
