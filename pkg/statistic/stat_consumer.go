package statistic

import (
	"fmt"
	"net/url"
	"shorturl/app/models/short"
	"shorturl/app/models/statistic"
	"shorturl/pkg/logger"
)

const KAFKA_TOPIC = "shorturl_stat"
const KAFKA_COUSUMER_GROUP = "shorturl_stat_group_1"

func Handle(shortStr string) error {
	model := short.GetByShort(shortStr)
	if model.ID == 0 {
		logger.WarnString("statistic", "handle", fmt.Sprintf("找不到数据：%s", shortStr))
		return nil
	}

	parseUrl, err := url.Parse(model.Long)
	if err != nil {
		return err
	}
	host := parseUrl.Host

	statModel := statistic.GetByShort(shortStr)
	if statModel.ID == 0 {
		// 不存在数据，新建
		statModel = statistic.Statistic{
			Host:    host,
			Long:    model.Long,
			Short:   model.Short,
			Counter: 1,
		}
		statModel.Create()
	} else {
		statistic.IncreaseCounterByShort(shortStr)
	}
	return nil
}
