package ordergenerator

import (
	"bytes"
	"log"
	"strconv"
	"time"

	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
)

type Address struct {
	Street   string `fake:"{street}"`
	City     string `fake:"{city}"`
	State    string `fake:"{state}"`
	PostCode string `fake:"{zip}"`
}

type CompanyDetails struct {
	CompanyName string `fake:"{company}"`
	Phone       string `fake:"{phoneformatted}"`
}

type Buyer struct {
	CompanyDetails CompanyDetails
	Address        Address
}

type Seller struct {
	CompanyDetails CompanyDetails
	Address        Address
}

type OrderLine struct {
	Code        string
	Category    string
	Name        string `fake:"{beername}"`
	Description string
	Quantity    int `fake:"{number:1,99}"`
}

type OrderArgs struct {
	Buyer      Buyer
	Seller     Seller
	OrderLines []OrderLine
}

type Generator struct {
}

func New() *Generator {
	return &Generator{}
}

var commonTextProps = props.Text{
	Color:       color.NewBlack(),
	Size:        9,
	Align:       consts.Left,
	Extrapolate: false,
	Family:      consts.Arial,
}

var commonBoldTextProps = props.Text{
	Top:         3,
	Color:       color.NewBlack(),
	Size:        9,
	Align:       consts.Left,
	Extrapolate: false,
	Family:      consts.Arial,
	Style:       consts.Bold,
}

var centeredTextProps = props.Text{
	Color:       color.NewBlack(),
	Size:        9,
	Align:       consts.Center,
	Extrapolate: false,
	Family:      consts.Arial,
}
var rowHeight = 5.0

func (g *Generator) Generate(orderArgs OrderArgs) (bytes.Buffer, error) {
	start := time.Now()

	basePdf := pdf.NewMaroto(consts.Portrait, consts.A4)
	basePdf.SetPageMargins(10, 15, 10)
	basePdf.SetCreator("Solution GmbH", false)
	basePdf.SetAuthor("Solution GmbH", false)
	basePdf.SetTitle("Order", false)

	order := orderPdfBuilder{basePdf}

	order.Row(rowHeight, func() {
		order.Col(12, func() {
			order.Text("My ugly prototype order", centeredTextProps)
		})
	})

	order.addBoldTextRow("Seller")
	order.addTextRow(orderArgs.Seller.CompanyDetails.CompanyName)
	order.addTextRow(orderArgs.Seller.CompanyDetails.Phone)
	order.addAddress(orderArgs.Seller.Address)

	order.addBoldTextRow("Buyer")
	order.addTextRow(orderArgs.Buyer.CompanyDetails.CompanyName)
	order.addTextRow(orderArgs.Buyer.CompanyDetails.Phone)
	order.addAddress(orderArgs.Buyer.Address)

	order.Row(rowHeight, func() {})

	orderLines := [][]string{}
	for _, orderLine := range orderArgs.OrderLines {
		orderLines = append(orderLines, []string{
			orderLine.Code,
			orderLine.Category,
			orderLine.Name,
			orderLine.Description,
			strconv.Itoa(orderLine.Quantity),
		})
	}
	order.TableList(
		[]string{"Code", "Category", "Name", "Description", "Quantity"},
		orderLines,
		props.TableList{
			HeaderProp: props.TableListContent{
				Family:    consts.Arial,
				Size:      9,
				GridSizes: []uint{2, 2, 3, 4, 1},
			},
			ContentProp: props.TableListContent{
				Family: consts.Arial, Size: 9,
				GridSizes: []uint{2, 2, 3, 4, 1},
			},
			Line: true,
			LineProp: props.Line{
				Style: consts.Dotted,
				Width: 0.2,
			},
		})

	elapsed := time.Since(start)
	log.Printf("Generate took %s", elapsed)

	return order.Output()
}

type orderPdfBuilder struct {
	pdf.Maroto
}

func (o *orderPdfBuilder) addAddress(address Address) {
	o.addTextRow(address.Street)
	o.addTextRow(address.PostCode + " " + address.City)
	o.addTextRow(address.State)
}

func (o *orderPdfBuilder) addBoldTextRow(value string) {
	o.Row(rowHeight+commonBoldTextProps.Top, func() {
		o.Col(12, func() {
			o.Text(value, commonBoldTextProps)
		})
	})
}

func (o *orderPdfBuilder) addTextRow(value string) {
	o.Row(rowHeight, func() {
		o.Col(12, func() {
			o.Text(value, commonTextProps)
		})
	})
}
