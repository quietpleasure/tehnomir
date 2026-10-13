package tehnomir

import "testing"

func TestIsPartnerWarehouse(t *testing.T) {
	cfg := Config{PartnerWarehouses: []string{"DDSH", " ffgh "}}
	for _, logo := range []string{"DDSH", "ddsh", "FFGH"} {
		if !cfg.IsPartnerWarehouse(logo) {
			t.Errorf("%q должен считаться партнёрским складом", logo)
		}
	}
	if cfg.IsPartnerWarehouse("ATG") {
		t.Error("ATG не должен считаться партнёрским складом")
	}

	// Пустой список — фолбэк на склад по умолчанию: configs/ не в репозитории,
	// и без конфига подсказка всё равно должна работать.
	var empty Config
	if !empty.IsPartnerWarehouse(DEFAULT_PARTNER_WAREHOUSE) {
		t.Errorf("без конфига %q должен считаться партнёрским складом", DEFAULT_PARTNER_WAREHOUSE)
	}
	if empty.IsPartnerWarehouse("ATG") {
		t.Error("без конфига посторонний склад не должен считаться партнёрским")
	}
}
