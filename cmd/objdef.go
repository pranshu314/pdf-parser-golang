package main

type PDFObject interface {
	GetType() string
}

type PDFBoolean struct {
	Value bool
}

func (b PDFBoolean) GetType() string {
	return "Boolean"
}

type PDFRealNumber struct {
	Value float64
}

func (r PDFRealNumber) GetType() string {
	return "RealNumber"
}

type PDFInteger struct {
	Value int
}

func (i PDFInteger) GetType() string {
	return "Integer"
}

type PDFString struct {
	Value string
}

func (s PDFString) GetType() string {
	return "String"
}

type PDFName struct {
	Value string
}

func (n PDFName) GetType() string {
	return "Name"
}

type PDFArray []PDFObject

func (a PDFArray) GetType() string {
	return "Array"
}

type PDFDictionary map[string]PDFObject

func (d PDFDictionary) GetType() string {
	return "Dictionary"
}

type PDFStream struct {
	Dictionary PDFDictionary
	Content    []byte
}

func (s PDFStream) GetType() string {
	return "Stream"
}

type PDFNull struct{}

func (n PDFNull) GetType() string {
	return "Null"
}

type PDFIndirectObject struct {
	ObjectNumber     int
	GenerationNumber int
	Object           PDFObject
}

func (i PDFIndirectObject) GetType() string {
	return "IndirectObject"
}

type PDFHeader struct {
	Version string
}

type PDFCrossReferenceTable struct {
	Offsets []int
}

type PDFFooter struct {
	StartXRef int
}

type PDFDocument struct {
	Header              PDFHeader
	Objects             []PDFObject
	CrossReferenceTable PDFCrossReferenceTable
	Trailer             PDFDictionary
	Footer              PDFFooter
}
