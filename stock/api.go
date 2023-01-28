package stock

import (
	"bytes"
	"errors"
	"fmt"
	financeChart "github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/equity"
	"github.com/piquette/finance-go/quote"
	"github.com/wcharczuk/go-chart/v2"
	"time"
)

func getDescription(symbol string) (string, error) {
	q, err := quote.Get(symbol)
	if err != nil {
		return "", err
	}

	s := fmt.Sprintf("%s: $%.2f (%.2f%%)", symbol, q.RegularMarketPrice, q.RegularMarketChangePercent)

	return s, nil
}

func getDetails(symbol string) (string, error) {
	q, err := quote.Get(symbol)
	if err != nil {
		return "", err
	}
	e, err := equity.Get(symbol)
	if err != nil {
		return "", err
	}

	// Title info
	s := fmt.Sprintf("*%s*\n", q.ShortName)
	s += fmt.Sprintf("_%s_\n", q.FullExchangeName)
	s += "\n"

	// Market info
	s += "*Market Info*\n"
	s += fmt.Sprintf("_Market Cap_: $%.2fB\n", float64(e.MarketCap)/1000000000)
	s += fmt.Sprintf("_P/E_: %.2f\n", e.TrailingPE)
	s += fmt.Sprintf("_EPS_: $%.2f\n", e.EpsTrailingTwelveMonths)
	s += fmt.Sprintf("_Dividend Yield_: %.2f%%\n", e.TrailingAnnualDividendYield*100)
	s += "\n"

	// 52 week info
	s += "*Annual Data*\n"
	s += fmt.Sprintf("_52 Week High_: $%.2f (%.2f%%)\n", q.FiftyTwoWeekHigh, q.FiftyTwoWeekHighChangePercent*100)
	s += fmt.Sprintf("_52 Week Low_: $%.2f (%.2f%%)\n", q.FiftyTwoWeekLow, q.FiftyTwoWeekLowChangePercent*100)
	s += fmt.Sprintf("_50 Day Average_: $%.2f (%.2f%%)\n", q.FiftyDayAverage, q.FiftyDayAverageChangePercent*100)
	s += "\n"

	return s, nil
}

func getChart(symbol string, start *time.Time, end *time.Time, interval datetime.Interval) (*bytes.Buffer, error) {
	var xValues []time.Time
	var yValues []float64

	// fetch financial data
	params := &financeChart.Params{
		Symbol:   symbol,
		Interval: interval,
	}

	if start != nil && end != nil {
		params.Start = datetime.New(start)
		params.End = datetime.New(end)
	} else if (start == nil && end != nil) || (start != nil && end == nil) {
		return nil, errors.New("both start and end must be provided")
	}

	iter := financeChart.Get(params)

	// populate axes with data
	for iter.Next() {
		b := iter.Bar()
		if b.High.IsZero() {
			return nil, errors.New("no trading data for this time period")
		}
		yValues = append(yValues, b.Close.InexactFloat64())
		t := time.Unix(int64(b.Timestamp), 0)
		xValues = append(xValues, t)
	}

	if err := iter.Err(); err != nil {
		return nil, err
	}

	// create graph
	graph := chart.Chart{
		XAxis: chart.XAxis{
			ValueFormatter: chart.TimeValueFormatterWithFormat("3:04pm"),
		},
		Series: []chart.Series{
			chart.TimeSeries{
				XValues: xValues,
				YValues: yValues,
			},
		},
	}

	// write graph
	buff := bytes.NewBuffer([]byte{})
	if err := graph.Render(chart.PNG, buff); err != nil {
		return nil, err
	}

	return buff, nil
}
