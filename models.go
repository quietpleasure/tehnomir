package tehnomir

import (
	"github.com/quietpleasure/utilits"
)

/*
	Requests body structs
*/

type TokenRequestBody struct {
	Token string `json:"apiToken"`
}
type TestRequestBody struct {
	TokenRequestBody
	Phrase string `json:"string"`
}

type PriceSearchRequestBody struct {
	TokenRequestBody
	BrandID     int    `json:"brandId"`
	Code        string `json:"code"`
	ShowAnalogs int    `json:"isShowAnalogs"`
	Currency    string `json:"currency"`
}

type ProductInfoRequestBody struct {
	TokenRequestBody
	BrandID int    `json:"brandId"`
	Code    string `json:"code"`
}

type GetUnloadsRequestBody struct {
	TokenRequestBody
	FromDate string `json:"fromDate"` //"2006-02-01"
	ToDate   string `json:"toDate"`   //"2006-02-01"
}

type GetUnloadDataRequestBody struct {
	TokenRequestBody
	UnloadID int `json:"unloadId"`
}

type BasketAddRequestBody struct {
	TokenRequestBody
	ProductID int64  `json:"productId"`
	PriceLogo string `json:"priceLogo"`
	Quantity  int    `json:"quantity"`
	Reference string `json:"reference"`
	Comment   string `json:"comment"`
}

type PositionInfoRequestBody struct {
	TokenRequestBody
	Reference string `json:"reference"`
}

type BasketDeletePositionRequestBody struct {
	TokenRequestBody
	BasketID int `json:"basketId"`
}

type BrandsByCodeRequestBody struct {
	TokenRequestBody
	Code string `json:"code"`
}

type OrderCreateRequestBody struct {
	TokenRequestBody
	OrderNumber string `json:"orderNumber"`
	BasketIds   []int  `json:"basketIds"`
}

type OrderSearchRequestBody struct {
	TokenRequestBody
	FromDate string `json:"fromDate"` //"2006-01-02"
	ToDate   string `json:"toDate"`   //"2006-01-02"
	OrderNum string `json:"orderNumber"`
}

type GetChangedPositionsRequestBody struct {
	TokenRequestBody
	FromDate string `json:"fromDateTime"` //"2006-01-02"
}

type GetOrderPositionsRequestBody struct {
	TokenRequestBody
	OrderID int `json:"orderId"`
}

type GetOrderPositionsByStatusRequestBody struct {
	TokenRequestBody
	StatusID int `json:"statusId"`
}

/*
	Response structs
*/

type SuccessResponse struct {
	Success bool `json:"success"`
}

type ResponseError struct {
	SuccessResponse
	Data struct {
		Name       string `json:"name"`
		Status     int    `json:"status"`
		Message    string `json:"message"`
		TestString string `json:"testString,omitempty"`
	} `json:"data"`
}

type TestConnectResponse struct {
	Success bool `json:"success"`
	Data    struct {
		TestString string `json:"testString"`
	} `json:"data"`
}

type PriceSearchResponse struct {
	SuccessResponse
	Details []FoundDetail `json:"data"`
}

type FoundDetail struct {
	ProductID          int                   `json:"productId"`
	BrandID            utilits.CustomInt     `json:"brandId"`
	BrandGroupID       utilits.CustomInt     `json:"brandGroupId"`
	Brand              string                `json:"brand"`
	Code               string                `json:"code"`
	DescriptionRus     string                `json:"descriptionRus"`
	DescriptionUa      string                `json:"descriptionUa"`
	Weight             utilits.CustomFloat64 `json:"weight"`
	IsOriginal         utilits.CustomBool    `json:"isOriginal"`
	IsExistProductInfo utilits.CustomBool    `json:"isExistProductInfo"`
	Stocks             []OfferSupplier       `json:"rests"`
}

type OfferSupplier struct {
	PriceLogo       string              `json:"priceLogo"`
	Price           float64             `json:"price"`
	Currency        string              `json:"currency"`
	Quantity        utilits.CustomInt64 `json:"quantity"`
	QuantityType    string              `json:"quantityType"`
	Multiplicity    utilits.CustomInt   `json:"multiplicity"`
	PriceQuality    float64             `json:"priceQuality"`
	DeliveryTypeID  utilits.CustomInt   `json:"deliveryTypeId"`
	DeliveryType    string              `json:"deliveryType"`
	DeliveryDays    int                 `json:"deliveryTime"`
	DeliveryDate    utilits.CustomTime  `json:"deliveryDate"`
	DeliveryPercent int                 `json:"deliveryPercent"`
	PriceChangeDate utilits.CustomTime  `json:"priceChangeDate"`
	IsReturn        utilits.CustomBool  `json:"isReturn"`
	IsPriceFinal    utilits.CustomBool  `json:"isPriceFinal"`
}

type SuppliersResponse struct {
	SuccessResponse
	Suppliers []Supplier `json:"data"`
}

type Supplier struct {
	PriceLogo      string              `json:"priceLogo"`
	DeliveryTypeID utilits.CustomInt   `json:"deliveryTypeId"`
	DeliveryType   string              `json:"deliveryType"`
	DeliveryDays   int                 `json:"deliveryTime"`
	DeliveryHours  utilits.CustomInt64 `json:"deliveryTimeHours"`
	DeliveryDate   utilits.CustomTime  `json:"deliveryDate"`
	Region         string              `json:"region"`
	RegionEn       string              `json:"regionEn"`
	RegionUa       string              `json:"regionUa"`
	IsReturn       utilits.CustomBool  `json:"isReturnFlag"`
	IsPriceFinal   utilits.CustomBool  `json:"isPriceFinalFlag"`
}

type BrandGroupsResponse struct {
	SuccessResponse
	BrandGroups []BrandGroup `json:"data"`
}

type BrandGroup struct {
	GroupID   utilits.CustomInt   `json:"groupId"`
	GroupName string              `json:"group"`
	BrandIds  []utilits.CustomInt `json:"brandIds"`
}

type BrandsResponse struct {
	Success bool    `json:"success"`
	Brands  []Brand `json:"data"`
}

type Brand struct {
	BrandID    utilits.CustomInt  `json:"brandId"`
	BrandName  string             `json:"brand"`
	IsOriginal utilits.CustomBool `json:"isOriginal"`
}

type UnloadsResponse struct {
	SuccessResponse
	Unloads []Unload `json:"data"`
}

type Unload struct {
	UnloadID       utilits.CustomInt     `json:"unloadId"`
	CreateTime     utilits.CustomTime    `json:"createTime"` //"2006-01-02 15:04:05.000000",
	BoxQuantity    utilits.CustomInt     `json:"boxQuantity"`
	SumPositions   utilits.CustomFloat64 `json:"sumPositions"`
	SumWorks       utilits.CustomInt     `json:"sumWorks"`
	SumDelivery    utilits.CustomFloat64 `json:"sumDelivery"`
	SumTotal       utilits.CustomFloat64 `json:"sumTotal"`
	Carrier        string                `json:"carrier"`
	CarrierWaybill interface{}           `json:"carrierWaybill"` // maybe null
}

type UnloadResponse struct {
	SuccessResponse
	Unload UnloadData `json:"data"`
}
type UnloadData struct {
	Boxes     []UnloadBox      `json:"boxes"`
	Positions []UnloadPosition `json:"positions"`
}

type UnloadBox struct {
	BoxID        utilits.CustomInt     `json:"boxId"`
	SumPositions utilits.CustomFloat64 `json:"sumPositions"`
	SumWorks     utilits.CustomFloat64 `json:"sumWorks"` //стоимость доставки
	Length       utilits.CustomFloat64 `json:"length"`
	Width        utilits.CustomFloat64 `json:"width"`
	Height       utilits.CustomFloat64 `json:"height"`
	Weight       utilits.CustomFloat64 `json:"weight"`
}

type UnloadPosition struct {
	BoxID           utilits.CustomInt     `json:"boxId"`
	OrderID         utilits.CustomInt     `json:"orderId"`     //номер заказа в ТМ в цифровом типе
	OrderNumber     string                `json:"orderNumber"` //номер заказа в ТМ в строковом типе
	OrderPositionID utilits.CustomInt     `json:"orderPositionId"`
	PriceLogo       string                `json:"priceLogo"`
	Brand           string                `json:"brand"`
	BrandID         utilits.CustomInt     `json:"brandId"`
	Code            string                `json:"code"`
	DescriptionRus  string                `json:"descriptionRus"`
	DescriptionUa   string                `json:"descriptionUa"`
	Quantity        utilits.CustomInt     `json:"quantity"`
	Price           utilits.CustomFloat64 `json:"price"`
	PriceFinal      utilits.CustomFloat64 `json:"priceFinal"`
	Currency        string                `json:"currency"`
	Reference       string                `json:"reference"` //!отдает в верхнем регистре
	Comment         string                `json:"comment"`
	AdminComment    string                `json:"adminComment"`
	Weight          utilits.CustomFloat64 `json:"weight"`
	Sticker         string                `json:"sticker"`
}

type BoxesReadyToSendResponse struct {
	SuccessResponse
	ReadyBoxes []UnloadBox `json:"data"`
}

type ProductInfoResponse struct {
	SuccessResponse
	Data struct {
		ProductID      utilits.CustomInt     `json:"productId"`
		Brand          string                `json:"brand"`
		Code           string                `json:"code"`
		DescriptionRus string                `json:"descriptionRus"`
		DescriptionUa  string                `json:"descriptionUa"`
		Weight         utilits.CustomFloat64 `json:"weight"`
		Volume         utilits.CustomFloat64 `json:"volume"`
		Images         []struct {
			Image string `json:"image"`
		} `json:"images"`
		Properties []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"properties"`
		Analogs []struct {
			ProductID      utilits.CustomInt     `json:"productId"`
			Brand          string                `json:"brand"`
			Code           string                `json:"code"`
			DescriptionRus string                `json:"descriptionRus"`
			DescriptionUa  string                `json:"descriptionUa"`
			Weight         utilits.CustomFloat64 `json:"weight"`
			Volume         utilits.CustomFloat64 `json:"volume"`
		} `json:"analogs"`
	} `json:"data"`
}

type BasketAddResponse struct {
	SuccessResponse
	Data struct {
		BasketID int `json:"basketId"`
	} `json:"data"`
}

type PositionsInfoResponse struct {
	SuccessResponse
	Positions []Position `json:"data"` //всегда 1 позиция
}

type Position struct {
	OrderID         utilits.CustomInt     `json:"orderId,omitempty"`
	OrderNumber     string                `json:"orderNumber,omitempty"`
	OrderPositionID utilits.CustomInt     `json:"orderPositionId,omitempty"`
	PriceLogo       string                `json:"priceLogo"`
	BrandID         utilits.CustomInt     `json:"brandId"`
	Brand           string                `json:"brand"`
	Code            string                `json:"code"`                  //приходит неочищенный
	ReplaceCode     string                `json:"replaceCode,omitempty"` //если в тм меняют номер при заказе
	DescriptionRus  string                `json:"descriptionRus,omitempty"`
	DescriptionUa   string                `json:"descriptionUa,omitempty"`
	Price           utilits.CustomFloat64 `json:"price"`
	Currency        string                `json:"currency"`
	Reference       string                `json:"reference"` //приходит в верхнем регистре
	Comment         string                `json:"comment"`
	AdminComment    string                `json:"adminComment,omitempty"`
	States          []StatePosition       `json:"states,omitempty"` //может быть несколько с отказными
}

type StatePosition struct {
	Quantity          utilits.CustomInt  `json:"quantity"`
	StatusID          utilits.CustomInt  `json:"statusId"`
	Status            string             `json:"status"`
	StatusChangedDate utilits.CustomTime `json:"statusChangedDate"` //"2006-01-02 15:04:05.000000"
}

type BasketPositionsResponse struct {
	SuccessResponse
	Positions []struct {
		BasketID     utilits.CustomInt     `json:"basketId"`
		PriceLogo    string                `json:"priceLogo"`
		BrandID      utilits.CustomInt     `json:"brandId"`
		Brand        string                `json:"brand"`
		Code         string                `json:"code"`
		Quantity     utilits.CustomInt     `json:"quantity"`
		Price        utilits.CustomFloat64 `json:"price"`
		Multiplicity utilits.CustomInt     `json:"multiplicity"`
		Currency     string                `json:"currency"`
		Reference    string                `json:"reference"`
		Comment      string                `json:"comment"`
	} `json:"data"`
}

type CurrenciesResponse struct {
	SuccessResponse
	Currencies []struct {
		Currency string  `json:"currency"`
		Rate     float64 `json:"rate"`
	} `json:"data"`
}

type BrandsByCodeResponse struct {
	SuccessResponse
	Data []struct {
		BrandID        int    `json:"brandId"`
		Brand          string `json:"brand"`
		DescriptionRus string `json:"descriptionRus"`
		DescriptionUa  string `json:"descriptionUa"`
		OffersCount    int    `json:"offersCount"`
		BrandGroupID   int    `json:"brandGroupId"`
	} `json:"data"`
}

type PositionStatusesResponse struct {
	Success  bool `json:"success"`
	Statuses []struct {
		StatusID    int    `json:"statusId"`
		Status      string `json:"status"`
		Description string `json:"description"`
	} `json:"data"`
}

type Order struct {
	OrderID     int                   `json:"orderId"`
	OrderNumber string                `json:"orderNumber"`
	Sum         utilits.CustomFloat64 `json:"sum"`
	StatusID    utilits.CustomInt     `json:"statusId"`
	Status      string                `json:"status"`
	CreateTime  utilits.CustomTime    `json:"createTime"` //"2023-11-09 13:40:39"
}

type OrderResponse struct {
	Success bool  `json:"success"`
	Order   Order `json:"data"`
}

type OrdersResponse struct {
	Success bool    `json:"success"`
	Orders  []Order `json:"data"`
}

type StockPriceResponse struct {
	Success  bool           `json:"success"`
	Products []StockProduct `json:"data"`
}

type StockProduct struct {
	ProductID      int                   `json:"productId"`
	Brand          string                `json:"brand"`
	Code           string                `json:"code"`
	DescriptionRus string                `json:"descriptionRus"`
	Quantity       int                   `json:"quantity"`
	Price          utilits.CustomFloat64 `json:"price"`
	Currency       string                `json:"currency"`
	CodePrinted    string                `json:"codePrinted"`
	PriceForRemote utilits.CustomFloat64 `json:"priceForRemote"`
}

type FinanceInfo struct {
	Success bool `json:"success"`
	Data    struct {
		Account struct {
			CurrencyCode              string  `json:"currencyCode"`
			Balance                   float64 `json:"balance"`
			ReadyToShipping           float64 `json:"readyToShipping"`
			NeedPayForUnload          float64 `json:"needPayForUnload"`
			DebtDebit                 int     `json:"debtDebit"`
			ExpectedPositionsSum      float64 `json:"expectedPositionsSum"`
			PrepaymentForNewOrdersSum int     `json:"prepaymentForNewOrdersSum"`
			ActivePositionsSum        float64 `json:"activePositionsSum"`
		} `json:"account"`
		FinanceConditions struct {
			Orders struct {
				Type    string `json:"type"`
				Value   int    `json:"value"`
				Unit    string `json:"unit"`
				Comment string `json:"comment"`
			} `json:"orders"`
			Shipment struct {
				Type    string `json:"type"`
				Value   int    `json:"value"`
				Unit    string `json:"unit"`
				Comment string `json:"comment"`
			} `json:"shipment"`
			CreditTerm int `json:"creditTerm"`
		} `json:"financeConditions"`
	} `json:"data"`
}
