package tehnomir

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"avtodetali/pkg/utilits"
)

var ErrBadResponse error = fmt.Errorf("bad response")

const (
	EUR Currency = "EUR"
	USD Currency = "USD"

	TestConnect apiPath = "test/connect"

	PriceSearch   apiPath = "price/search"
	GetStockPrice apiPath = "price/getStockPrice"

	GetUnloads    apiPath = "unload/search"
	GetUnloadData apiPath = "unload/getData"
	GetBoxesReady apiPath = "unload/getBoxesReadyToSend"

	GetSuppliers        apiPath = "info/getSuppliers"
	GetBrands           apiPath = "info/getBrands"
	GetBrandGroups      apiPath = "info/getBrandGroups"
	GetProductInfo      apiPath = "info/getProductInfo"
	GetCurrencies       apiPath = "info/getCurrencies"
	GetBrandsByCode     apiPath = "info/getBrandsByCode"
	GetPositionStatuses apiPath = "info/getPositionStatuses"
	InfoProfile         apiPath = "info/infoProfile"

	BasketAdd            apiPath = "basket/add"
	GetBasketPositions   apiPath = "basket/getPositions"
	BasketDeletePosition apiPath = "basket/delete"
	BasketClear          apiPath = "basket/clear"

	GetPositionInfo           apiPath = "order/getPositionInfo"
	OrderCreate               apiPath = "order/create"
	GetActiveOrders           apiPath = "order/getActive"
	OrderSearch               apiPath = "order/search"
	GetChangedPositions       apiPath = "order/getChangedPositions"
	GetOrderPositions         apiPath = "order/getOrderPositions"
	GetOrderPositionsByStatus apiPath = "order/getOrderPositionsByStatus"
)

func makeError(body io.ReadCloser) error {
	var resp ResponseError
	defer body.Close()
	if err := json.NewDecoder(body).Decode(&resp); err != nil {
		return err
	}
	return fmt.Errorf("%d:%s - %s", resp.Data.Status, resp.Data.Name, resp.Data.Message)
}

type Client struct {
	cfg    Config
	client *http.Client
}

func New(cfg Config) *Client {
	cfg = overrideConfig(cfg)
	return &Client{
		cfg: cfg,
		client: &http.Client{
			// Раньше здесь стояло жёстко зашитое 10s, а cfg.Timeout не
			// использовался вовсе. По «широкому» артикулу (сотни аналогов)
			// ответ price/search — сотни килобайт, и поставщик думает над ним
			// заметно дольше, чем над обычным запросом.
			Timeout: cfg.Timeout,
			Transport: &http.Transport{
				DialContext:         (&net.Dialer{Timeout: 2 * time.Second}).DialContext,
				TLSHandshakeTimeout: time.Second,
				// 3s не хватало: на тяжёлом поиске ТМ отдаёт первый байт позже,
				// и запрос падал по таймауту вместо того, чтобы дождаться ответа.
				ResponseHeaderTimeout: 15 * time.Second,
				MaxIdleConns:          100,
				MaxIdleConnsPerHost:   100,
				IdleConnTimeout:       60 * time.Second,
			},
		}}
}

func (c *Client) newRequest(path apiPath, body ...any) (*http.Response, error) {
	return c.newRequestCtx(context.Background(), path, body...)
}

func (c *Client) newRequestCtx(ctx context.Context, path apiPath, body ...any) (*http.Response, error) {
	var reqbody any
	if body == nil {
		reqbody = TokenRequestBody{
			Token: c.cfg.Token,
		}
	} else {
		reqbody = c.makeRequestBody(path, body[0])
	}
	buff := new(bytes.Buffer)
	if err := json.NewEncoder(buff).Encode(reqbody); err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.makeApiPath(path),
		buff)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, makeError(res.Body)
	}
	return res, nil
}

func (c *Client) makeApiPath(path apiPath) string {
	// u := url.URL{
	// 	Host:   c.cfg.Host,
	// 	Path:   string(path),
	// }
	return c.cfg.Host + "/" + string(path)
}

func (c *Client) makeRequestBody(path apiPath, body any) any {
	switch path {
	case TestConnect:
		b, ok := body.(*TestRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case PriceSearch:
		b, ok := body.(*PriceSearchRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetProductInfo:
		b, ok := body.(*ProductInfoRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetUnloads:
		b, ok := body.(*GetUnloadsRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetUnloadData:
		b, ok := body.(*GetUnloadDataRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case BasketAdd:
		b, ok := body.(*BasketAddRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetPositionInfo:
		b, ok := body.(*PositionInfoRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case BasketDeletePosition:
		b, ok := body.(*BasketDeletePositionRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case OrderCreate:
		b, ok := body.(*OrderCreateRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case OrderSearch:
		b, ok := body.(*OrderSearchRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetChangedPositions:
		b, ok := body.(*GetChangedPositionsRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetOrderPositions:
		b, ok := body.(*GetOrderPositionsRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	case GetOrderPositionsByStatus:
		b, ok := body.(*GetOrderPositionsByStatusRequestBody)
		if ok {
			b.Token = c.cfg.Token
		}
	}
	return body
}

func (c *Client) requestAndDecode(request apiPath, response any, requestbody ...any) error {
	return c.requestAndDecodeCtx(context.Background(), request, response, requestbody...)
}

func (c *Client) requestAndDecodeCtx(ctx context.Context, request apiPath, response any, requestbody ...any) error {
	var (
		body any
		err  error
		resp *http.Response
	)
	if requestbody != nil {
		body = requestbody[0]
		resp, err = c.newRequestCtx(
			ctx,
			request,
			body,
		)
	} else {
		resp, err = c.newRequestCtx(ctx, request)
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(response)
}

func (c *Client) TestConnect(s ...string) error {
	var phrase string
	if s != nil {
		phrase = s[0]
	}
	var res TestConnectResponse
	if err := c.requestAndDecode(TestConnect, &res, &TestRequestBody{Phrase: phrase}); err != nil {
		return err
	}
	if s != nil && res.Data.TestString != s[0] {
		return ErrBadResponse
	}
	return nil
}

func (c *Client) priceSearch(ctx context.Context, code string, currency Currency, brand int, analog bool) (*PriceSearchResponse, error) {
	var res PriceSearchResponse
	if err := c.requestAndDecodeCtx(ctx, PriceSearch, &res,
		&PriceSearchRequestBody{
			BrandID:     brand,
			Code:        code,
			ShowAnalogs: utilits.BoolToInt(analog),
			Currency:    string(currency),
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) SearchWithoutBrand(code string, currency ...Currency) (*PriceSearchResponse, error) {
	return c.SearchWithoutBrandCtx(context.Background(), code, currency...)
}

// SearchWithoutBrandCtx — вариант с контекстом: страница поиска отменяет
// запрос, когда пользователь ушёл со страницы или истёк таймаут хендлера,
// вместо того чтобы качать многосотенный килобайт «в никуда».
func (c *Client) SearchWithoutBrandCtx(ctx context.Context, code string, currency ...Currency) (*PriceSearchResponse, error) {
	cur := USD
	if currency != nil {
		cur = currency[0]
	}
	return c.priceSearch(ctx, code, cur, 0, false)
}

func (c *Client) SearchByBrand(code string, brand int, withAnalog bool, currency ...Currency) (*PriceSearchResponse, error) {
	return c.SearchByBrandCtx(context.Background(), code, brand, withAnalog, currency...)
}

func (c *Client) SearchByBrandCtx(ctx context.Context, code string, brand int, withAnalog bool, currency ...Currency) (*PriceSearchResponse, error) {
	cur := USD
	if currency != nil {
		cur = currency[0]
	}
	return c.priceSearch(ctx, code, cur, brand, withAnalog)
}

func (c *Client) SearchByBrandWithoutAnalogs(code string, brand int, currency ...Currency) (*PriceSearchResponse, error) {
	cur := USD
	if currency != nil {
		cur = currency[0]
	}
	return c.priceSearch(context.Background(), code, cur, brand, false)
}

func (c *Client) GetSuppliers() (*SuppliersResponse, error) {
	var res SuppliersResponse
	if err := c.requestAndDecode(GetSuppliers, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetBrands() (*BrandsResponse, error) {
	var res BrandsResponse
	if err := c.requestAndDecode(GetBrands, &res); err != nil {
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetBrandGroups() (*BrandGroupsResponse, error) {
	var res BrandGroupsResponse
	if err := c.requestAndDecode(GetBrandGroups, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetProductInfo(code string, brand int) (*ProductInfoResponse, error) {
	var res ProductInfoResponse
	if err := c.requestAndDecode(GetProductInfo, &res, &ProductInfoRequestBody{
		BrandID: brand,
		Code:    code,
	}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetUnloads(from, to time.Time) (*UnloadsResponse, error) {
	timeFormat := "2006-01-02"
	var res UnloadsResponse
	if err := c.requestAndDecode(GetUnloads, &res, &GetUnloadsRequestBody{
		FromDate: from.Format(timeFormat),
		ToDate:   to.Format(timeFormat),
	}); err != nil {
		return nil, err
	}
	return &res, nil
}

// get completed unload positions error="CustomFloat64: UnmarshalJSON: invalid character ',' after top-level value"
func (c *Client) GetUnloadData(unloadID int) (*UnloadResponse, error) {
	var res UnloadResponse
	if err := c.requestAndDecode(GetUnloadData, &res,
		&GetUnloadDataRequestBody{
			UnloadID: unloadID,
		}); err != nil {
		return nil, fmt.Errorf("unload id: %d: %s", unloadID, err)
	}
	return &res, nil
}

func (c *Client) GetBoxesReady() (*BoxesReadyToSendResponse, error) {
	var res BoxesReadyToSendResponse
	if err := c.requestAndDecode(GetBoxesReady, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) BasketAdd(prodid int64, priceLogo string, quantity int, reference string, comment ...string) (*BasketAddResponse, error) {
	var com string
	if comment != nil {
		com = comment[0]
	}
	var res BasketAddResponse
	if err := c.requestAndDecode(BasketAdd, &res,
		&BasketAddRequestBody{
			ProductID: prodid,
			PriceLogo: priceLogo,
			Quantity:  quantity,
			Reference: reference,
			Comment:   com,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetPositionInfo(reference string) (*PositionsInfoResponse, error) {
	var res PositionsInfoResponse
	if err := c.requestAndDecode(GetPositionInfo, &res,
		&PositionInfoRequestBody{
			Reference: reference,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) GetBasketPositions() (*BasketPositionsResponse, error) {
	var res BasketPositionsResponse
	if err := c.requestAndDecode(GetBasketPositions, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) BasketDeletePosition(basketid int) error {
	var res SuccessResponse
	if err := c.requestAndDecode(BasketDeletePosition, &res,
		&BasketDeletePositionRequestBody{
			BasketID: basketid,
		}); err != nil {
		return err
	}
	if !res.Success {
		return ErrBadResponse
	}
	return nil
}

func (c *Client) BasketClear() error {
	var res SuccessResponse
	if err := c.requestAndDecode(BasketClear, &res); err != nil {
		return err
	}
	if !res.Success {
		return ErrBadResponse
	}
	return nil
}

func (c *Client) GetCurrencies() (*CurrenciesResponse, error) {
	var res CurrenciesResponse
	if err := c.requestAndDecode(GetCurrencies, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) BrandsByCode(code string) (*BrandsByCodeResponse, error) {
	var res BrandsByCodeResponse
	if err := c.requestAndDecode(GetBrandsByCode, &res,
		&BrandsByCodeRequestBody{
			Code: code,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) PositionStatuses() (*PositionStatusesResponse, error) {
	var res PositionStatusesResponse
	if err := c.requestAndDecode(GetPositionStatuses, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) OrderCreate(basketIds []int, ordernum ...string) (*OrderResponse, error) {
	var num string
	if ordernum != nil {
		num = ordernum[0]
	}
	var res OrderResponse
	if err := c.requestAndDecode(OrderCreate, &res,
		&OrderCreateRequestBody{
			OrderNumber: num,
			BasketIds:   basketIds,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) ActiveOrders() (*OrdersResponse, error) {
	var res OrdersResponse
	if err := c.requestAndDecode(GetActiveOrders, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) OrderSearchByDate(fromDate, toDate time.Time) (*OrdersResponse, error) {
	timeFormat := "2006-01-02"
	var res OrdersResponse
	if err := c.requestAndDecode(OrderSearch, &res,
		&OrderSearchRequestBody{
			FromDate: fromDate.Format(timeFormat),
			ToDate:   toDate.Format(timeFormat),
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) OrderSearchByNumber(ordernum string) (*OrdersResponse, error) {
	var res OrdersResponse
	if err := c.requestAndDecode(OrderSearch, &res,
		&OrderSearchRequestBody{
			OrderNum: ordernum,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) ChangedPositions(fromDateTime time.Time) (*PositionsInfoResponse, error) {
	timeFormat := "2006-01-02 15:04:05"
	var res PositionsInfoResponse
	if err := c.requestAndDecode(GetChangedPositions, &res,
		&GetChangedPositionsRequestBody{
			FromDate: fromDateTime.Format(timeFormat),
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) OrderPositions(orderid int) (*PositionsInfoResponse, error) {
	var res PositionsInfoResponse
	if err := c.requestAndDecode(GetOrderPositions, &res,
		&GetOrderPositionsRequestBody{
			OrderID: orderid,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

// Отличаются статусы: В работе - StatusID: 2 и Закрыт - StatusID: 4
func (c *Client) OrderPositionsByStatus(statusid int) (*PositionsInfoResponse, error) {
	var res PositionsInfoResponse
	if err := c.requestAndDecode(GetOrderPositionsByStatus, &res,
		&GetOrderPositionsByStatusRequestBody{
			StatusID: statusid,
		}); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) StockPrice() (*StockPriceResponse, error) {
	var res StockPriceResponse
	if err := c.requestAndDecode(GetStockPrice, &res); err != nil {
		return nil, err
	}
	return &res, nil
}

func (c *Client) InfoProfile() (*FinanceInfo, error) {
	var res FinanceInfo
	if err := c.requestAndDecode(InfoProfile, &res); err != nil {
		return nil, err
	}
	return &res, nil
}
