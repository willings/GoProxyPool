package provider

type Table struct {
	Tr []Tr `xml:"tr"`
}

type Tbody struct {
	Tr []Tr `xml:"tr"`
}

type Tr struct {
	Td []string `xml:"td"`
}
