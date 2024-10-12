package companies

//go:generate go run github.com/dmarkham/enumer -type=CompanyType -transform=title
type CompanyType int

const (
	Corporation CompanyType = iota
	NonProfit
	Cooperative
	SoleProprietorship
)
