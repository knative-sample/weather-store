package tablestore

import (
	"os"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	"github.com/golang/glog"
	"github.com/knative-sample/weather-store/pkg/weather"
	uuid "github.com/satori/go.uuid"
)

type TableClient struct {
	tableName string
	client    *tablestore.TableStoreClient
}

func InitClient() *TableClient {
	endpoint := os.Getenv("OTS_TEST_ENDPOINT")
	tableName := os.Getenv("TABLE_NAME")
	instanceName := os.Getenv("OTS_TEST_INSTANCENAME")
	accessKeyId := os.Getenv("OTS_TEST_KEYID")
	accessKeySecret := os.Getenv("OTS_TEST_SECRET")
	client := tablestore.NewClient(endpoint, instanceName, accessKeyId, accessKeySecret)
	return &TableClient{
		tableName: tableName,
		client:    client,
	}
}

func (tableClient *TableClient) Store(forecast weather.Forecast) error {
	for _, cast := range forecast.Casts {
		putRowRequest := new(tablestore.PutRowRequest)
		putRowChange := new(tablestore.PutRowChange)
		putRowChange.TableName = tableClient.tableName
		putPk := &tablestore.PrimaryKey{}

		updateRowRequest := new(tablestore.UpdateRowRequest)
		updateRowChange := new(tablestore.UpdateRowChange)
		updateRowChange.TableName = tableClient.tableName

		getRowRequest := new(tablestore.GetRowRequest)
		criteria := new(tablestore.SingleRowQueryCriteria)

		putPk.AddPrimaryKeyColumn("adcode", forecast.Adcode)
		putPk.AddPrimaryKeyColumn("date", cast.Date)

		putRowChange.PrimaryKey = putPk
		criteria.PrimaryKey = putPk
		updateRowChange.PrimaryKey = putPk

		getRowRequest.SingleRowQueryCriteria = criteria
		getRowRequest.SingleRowQueryCriteria.TableName = tableClient.tableName
		getRowRequest.SingleRowQueryCriteria.MaxVersion = 1
		getResp, err := tableClient.client.GetRow(getRowRequest)
		if err != nil {
			glog.Errorf("BatchWriteRow failed with error: %s", err.Error())
			continue
		}
		if len(getResp.Columns) == 0 {
			uid, _ := uuid.NewV4()
			putRowChange.AddColumn("id", uid.String())
			putRowChange.AddColumn("city", forecast.City)
			putRowChange.AddColumn("province", forecast.Province)
			putRowChange.AddColumn("reporttime", forecast.Reporttime)
			putRowChange.AddColumn("week", cast.Week)
			putRowChange.AddColumn("dayweather", cast.Dayweather)
			putRowChange.AddColumn("nightweather", cast.Nightweather)
			putRowChange.AddColumn("daytemp", cast.Daytemp)
			putRowChange.AddColumn("nighttemp", cast.Nighttemp)
			putRowChange.AddColumn("daywind", cast.Daywind)
			putRowChange.AddColumn("nightwind", cast.Nightwind)
			putRowChange.AddColumn("daypower", cast.Daypower)
			putRowChange.AddColumn("nightpower", cast.Nightpower)
			putRowChange.SetCondition(tablestore.RowExistenceExpectation_IGNORE)
			putRowRequest.PutRowChange = putRowChange
			_, err = tableClient.client.PutRow(putRowRequest)
			if err != nil {
				glog.Errorf("PutRow failed with error: %s", err.Error())
				continue
			}
		} else {
			updateRowChange.PutColumn("reporttime", forecast.Reporttime)
			updateRowChange.PutColumn("week", cast.Week)
			updateRowChange.PutColumn("dayweather", cast.Dayweather)
			updateRowChange.PutColumn("nightweather", cast.Nightweather)
			updateRowChange.PutColumn("daytemp", cast.Daytemp)
			updateRowChange.PutColumn("nighttemp", cast.Nighttemp)
			updateRowChange.PutColumn("daywind", cast.Daywind)
			updateRowChange.PutColumn("nightwind", cast.Nightwind)
			updateRowChange.PutColumn("daypower", cast.Daypower)
			updateRowChange.PutColumn("nightpower", cast.Nightpower)
			updateRowChange.SetCondition(tablestore.RowExistenceExpectation_IGNORE)
			updateRowRequest.UpdateRowChange = updateRowChange
			_, err = tableClient.client.UpdateRow(updateRowRequest)
			if err != nil {
				glog.Errorf("UpdateRow failed with error: %s", err.Error())
				continue
			}
		}

	}

	return nil
}
