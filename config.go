package tehnomir

import (
	"strings"
	"time"
)

const (
	// PROTO_TM             = "https"
	URL_API_TM           = "https://api.tehnomir.com.ua"
	PRICE_SEA    float64 = 5.0
	PRICE_AVIA   float64 = 10.0
	PRICE_VOLUME float64 = 15.0
)

// DEFAULT_PARTNER_WAREHOUSE — склад ТМ, с которым у нас есть прямой договор.
// Используется, когда в configs/tehnomir.yaml список partner_warehouses не задан:
// configs/ не в репозитории, поэтому дефолт должен быть самодостаточным.
const DEFAULT_PARTNER_WAREHOUSE = "DDSH"

type (
	Currency string
	apiPath  string
)

// Теги именно yaml: конфиг разбирается по ним. С прежними cfg-тегами ключи
// price_avia/price_sea не сопоставлялись с полями, ставки доставки оставались
// нулевыми, и PriceForUser не добавлял стоимость доставки по весу.
type Config struct {
	Token string `yaml:"token"`
	// Proto       string        `yaml:"proto"`
	Host        string        `yaml:"host"`
	Timeout     time.Duration `yaml:"timeout"`
	PriceAvia   float64       `yaml:"price_avia"`
	PriceSea    float64       `yaml:"price_sea"`
	PriceVolume float64       `yaml:"price_volume"`
	// PartnerWarehouses — priceLogo складов, у которых мы можем заказать деталь
	// напрямую, минуя Техномир. Админу такие позиции в корзине ТМ подсвечиваются
	// подсказкой, чтобы он решил, оформлять через ТМ или напрямую.
	PartnerWarehouses []string `yaml:"partner_warehouses"`
}

func DefaultConfig() *Config {
	return &Config{
		// Proto:       PROTO_TM,
		Host: URL_API_TM,
		// Таймаут на весь запрос вместе с чтением тела: выдача по «широкому»
		// артикулу вроде W71295 — сотни килобайт. configs/ не в репозитории,
		// поэтому дефолт должен быть самодостаточным.
		Timeout:           30 * time.Second,
		PriceAvia:         PRICE_AVIA,
		PriceSea:          PRICE_SEA,
		PriceVolume:       PRICE_VOLUME,
		PartnerWarehouses: []string{DEFAULT_PARTNER_WAREHOUSE},
	}
}

// WithDefaults проставляет значения по умолчанию незаданным параметрам.
// Вызывается при проверке общего конфига приложения: роутер и фоновые сервисы
// считают цену через cfg.Tehnomir.PriceForUser, а не через копию в Client.
func (c Config) WithDefaults() Config {
	return overrideConfig(c)
}

func overrideConfig(cfg Config) Config {
	deffault := DefaultConfig()

	if cfg.Host == "" {
		cfg.Host = deffault.Host
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = deffault.Timeout
	}
	if cfg.PriceAvia == 0 {
		cfg.PriceAvia = deffault.PriceAvia
	}
	if cfg.PriceSea == 0 {
		cfg.PriceSea = deffault.PriceSea
	}
	if cfg.PriceVolume == 0 {
		cfg.PriceVolume = deffault.PriceVolume
	}
	if len(cfg.PartnerWarehouses) == 0 {
		cfg.PartnerWarehouses = deffault.PartnerWarehouses
	}
	return cfg
}

// IsPartnerWarehouse сообщает, является ли склад партнёрским, то есть возможен ли
// прямой заказ у него в обход Техномир. Сравнение без учёта регистра: priceLogo
// приходит от ТМ и в конфиге может быть записан иначе.
func (c Config) IsPartnerWarehouse(priceLogo string) bool {
	logos := c.PartnerWarehouses
	if len(logos) == 0 {
		logos = []string{DEFAULT_PARTNER_WAREHOUSE}
	}
	priceLogo = strings.TrimSpace(priceLogo)
	for _, l := range logos {
		if strings.EqualFold(strings.TrimSpace(l), priceLogo) {
			return true
		}
	}
	return false
}

// PriceForUser считает конечную цену для пользователя с учётом стоимости
// доставки (авиа/море, по весу) и наценки (markup).
func (c Config) PriceForUser(price, weight, markup float64, deliveryTypeID int) float64 {
	var costWithShipping float64
	switch deliveryTypeID {
	case 1:
		costWithShipping = c.PriceAvia*weight + price
	case 2:
		costWithShipping = c.PriceSea*weight + price
	default:
		costWithShipping = price
	}
	return costWithShipping * markup
}