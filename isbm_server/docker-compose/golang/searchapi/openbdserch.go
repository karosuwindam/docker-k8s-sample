package searchapi

import (
	"fmt"

	"github.com/seihmd/openbd"
)

type NameType struct {
	Isbn     string `json:isbn`
	Title    string `json:title`
	Writer   string `json:writer`
	Brand    string `json:brand`
	Synopsis string `json:synopsis`
	Ext      string `json:ext`
	Image    string `json:image`
}

func OnixTitle(onix openbd.Onix) string {
	var output string
	output = onix.DescriptiveDetail.TitleDetail.TitleElement.TitleText.Content
	output += "(" + onix.DescriptiveDetail.TitleDetail.TitleElement.TitleText.Collationkey + ")"
	return output
}
func OnixBrand(onix openbd.Onix) string {
	var output string
	if len(onix.DescriptiveDetail.Collection.TitleDetail.TitleElement) > 0 {
		output = onix.DescriptiveDetail.Collection.TitleDetail.TitleElement[0].TitleText.Content
		output += "(" + onix.DescriptiveDetail.Collection.TitleDetail.TitleElement[0].TitleText.CollationKey + ")"
	}
	if output == "" {
		output = onix.PublishingDetail.Imprint.ImprintName
	}
	return output
}
func OnixWriter(onix openbd.Onix) (string, string) {
	var out1, out2 string
	for num, tmp := range onix.DescriptiveDetail.Contributor {
		if num == 0 {
			out1 = tmp.PersonName.Content
			out1 += "(" + tmp.PersonName.Collationkey + ")"
		} else {
			if out2 == "" {
				out2 = tmp.PersonName.Content
				out2 += "(" + tmp.PersonName.Collationkey + ")"
			} else {
				out2 += "," + tmp.PersonName.Content
				out2 += "(" + tmp.PersonName.Collationkey + ")"
			}
		}
	}
	return out1, out2
}
func OnixSynopsis(onix openbd.Onix) string {
	var output string
	for _, tmp := range onix.CollateralDetail.TextContent {
		if output == "" {
			output = tmp.Text
		} else {
			output += "\n" + tmp.Text
		}
	}
	return output
}

func GetOpenBdData(isbn string) NameType {
	var output NameType
	o := openbd.New()
	output.Isbn = isbn
	book, err := o.Get(output.Isbn)
	if err != nil {
		fmt.Println(err.Error())
		return output
	}
	output.Title = OnixTitle(book.Onix)
	if output.Title == "" {
		output.Title = book.GetTitle()
	}
	output.Synopsis = OnixSynopsis(book.Onix)
	output.Brand = OnixBrand(book.Onix)
	output.Writer, output.Ext = OnixWriter(book.Onix)
	if output.Writer == "" {
		output.Writer = book.GetSeries()
	}
	output.Image = book.GetCover()
	return output
}
