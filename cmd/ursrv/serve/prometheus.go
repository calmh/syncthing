package serve

import (
	"database/sql"
	"log"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/syncthing/syncthing/lib/ur/contract"
)

const namePrefix = "syncthing_usage_"

type metricsSet struct {
	db *sql.DB

	gauges         map[string]prometheus.Gauge
	gaugeVecs      map[string]*prometheus.GaugeVec
	gaugeVecLabels map[string][]string
	summaries      map[string]*metricSummary
}

func newMetricsSet(db *sql.DB) *metricsSet {
	s := &metricsSet{
		db:             db,
		gauges:         make(map[string]prometheus.Gauge),
		gaugeVecs:      make(map[string]*prometheus.GaugeVec),
		gaugeVecLabels: make(map[string][]string),
		summaries:      make(map[string]*metricSummary),
	}

	var initForType func(reflect.Type)
	initForType = func(t reflect.Type) {
		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)
			if field.Type.Kind() == reflect.Struct {
				initForType(field.Type)
				continue
			}
			name, typ, label := fieldNameTypeLabel(field)
			sname, labels := nameConstLabels(name)
			switch typ {
			case "gauge":
				s.gauges[name] = prometheus.NewGauge(prometheus.GaugeOpts{
					Name:        namePrefix + sname,
					ConstLabels: labels,
				})
			case "summary":
				s.summaries[name] = newMetricSummary(namePrefix+sname, nil, labels)
			case "gaugeVec":
				s.gaugeVecLabels[name] = append(s.gaugeVecLabels[name], label)
			case "summaryVec":
				s.summaries[name] = newMetricSummary(namePrefix+sname, []string{label}, labels)
			}
		}
	}
	initForType(reflect.ValueOf(contract.Report{}).Type())

	for name, labels := range s.gaugeVecLabels {
		s.gaugeVecs[name] = prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: namePrefix + name,
		}, labels)
	}

	return s
}

func fieldNameTypeLabel(rf reflect.StructField) (string, string, string) {
	metric := rf.Tag.Get("metric")
	name, typ, ok := strings.Cut(metric, ",")
	if !ok {
		return "", "", ""
	}
	gv, label, ok := strings.Cut(typ, ":")
	if ok {
		typ = gv
	}
	return name, typ, label
}

func nameConstLabels(name string) (string, prometheus.Labels) {
	if name == "-" {
		return "", nil
	}
	name, labels, ok := strings.Cut(name, "{")
	if !ok {
		return name, nil
	}
	lls := strings.Split(labels[:len(labels)-1], ",")
	m := make(map[string]string)
	for _, l := range lls {
		k, v, _ := strings.Cut(l, "=")
		m[k] = v
	}
	return name, m
}

func (s *metricsSet) addReport(r *contract.Report) {
	gaugeVecs := make(map[string][]string)
	s.addReportStruct(reflect.ValueOf(r).Elem(), gaugeVecs)
	for name, lv := range gaugeVecs {
		s.gaugeVecs[name].WithLabelValues(lv...).Add(1)
	}
}

func (s *metricsSet) addReportStruct(v reflect.Value, gaugeVecs map[string][]string) {
	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			s.addReportStruct(field, gaugeVecs)
			continue
		}

		name, typ, label := fieldNameTypeLabel(t.Field(i))
		switch typ {
		case "gauge":
			switch v := field.Interface().(type) {
			case int:
				s.gauges[name].Add(float64(v))
			case string:
				s.gaugeVecs[name].WithLabelValues(v).Add(1)
			case bool:
				if v {
					s.gauges[name].Add(1)
				}
			}
		case "gaugeVec":
			var labelValue string
			switch v := field.Interface().(type) {
			case string:
				labelValue = v
			case int:
				labelValue = strconv.Itoa(v)
			case map[string]int:
				for k, v := range v {
					labelValue = k
					field.SetInt(int64(v))
					break
				}
			}
			if _, ok := gaugeVecs[name]; !ok {
				gaugeVecs[name] = make([]string, len(s.gaugeVecLabels[name]))
			}
			for i, l := range s.gaugeVecLabels[name] {
				if l == label {
					gaugeVecs[name][i] = labelValue
					break
				}
			}
		case "summary", "summaryVec":
			switch v := field.Interface().(type) {
			case int:
				s.summaries[name].Observe("", float64(v))
			case float64:
				s.summaries[name].Observe("", v)
			case []int:
				for _, v := range v {
					s.summaries[name].Observe("", float64(v))
				}
			case map[string]int:
				for k, v := range v {
					s.summaries[name].Observe(k, float64(v))
				}
			}
		}
	}
}

func (s *metricsSet) Describe(c chan<- *prometheus.Desc) {
	for _, g := range s.gauges {
		g.Describe(c)
	}
	for _, g := range s.gaugeVecs {
		g.Describe(c)
	}
	for _, g := range s.summaries {
		g.Describe(c)
	}
}

func (s *metricsSet) Collect(c chan<- prometheus.Metric) {
	for _, g := range s.gauges {
		g.Set(0)
	}
	for _, g := range s.gaugeVecs {
		g.Reset()
	}
	for _, g := range s.summaries {
		g.Reset()
	}

	allReports(s.db)(func(r *contract.Report) bool {
		r.Version = transformVersion(r.Version)
		r.OS, r.Arch, _ = strings.Cut(r.Platform, "-")
		s.addReport(r)
		return true
	})

	for _, g := range s.gauges {
		c <- g
	}
	for _, g := range s.gaugeVecs {
		g.Collect(c)
	}
	for _, g := range s.summaries {
		g.Collect(c)
	}
}

type metricSummary struct {
	name   string
	values map[string][]float64
	zeroes map[string]int

	qDesc     *prometheus.Desc
	countDesc *prometheus.Desc
	sumDesc   *prometheus.Desc
	zDesc     *prometheus.Desc
}

func newMetricSummary(name string, labels []string, constLabels prometheus.Labels) *metricSummary {
	return &metricSummary{
		name:      name,
		values:    make(map[string][]float64),
		zeroes:    make(map[string]int),
		qDesc:     prometheus.NewDesc(name, "", append(labels, "quantile"), constLabels),
		countDesc: prometheus.NewDesc(name+"_nonzero_count", "", labels, constLabels),
		sumDesc:   prometheus.NewDesc(name+"_sum", "", labels, constLabels),
		zDesc:     prometheus.NewDesc(name+"_zero_count", "", labels, constLabels),
	}
}

func (q *metricSummary) Observe(labelValue string, v float64) {
	if v == 0 {
		q.zeroes[labelValue]++
		return
	}
	q.values[labelValue] = append(q.values[labelValue], v)
}

func (q *metricSummary) Describe(c chan<- *prometheus.Desc) {
	c <- q.qDesc
	c <- q.countDesc
	c <- q.sumDesc
	c <- q.zDesc
}

func (q *metricSummary) Collect(c chan<- prometheus.Metric) {
	for lv, vs := range q.values {
		var labelVals []string
		if lv != "" {
			labelVals = []string{lv}
		}

		c <- prometheus.MustNewConstMetric(q.countDesc, prometheus.GaugeValue, float64(len(vs)), labelVals...)
		c <- prometheus.MustNewConstMetric(q.zDesc, prometheus.GaugeValue, float64(q.zeroes[lv]), labelVals...)

		var sum float64
		for _, v := range vs {
			sum += v
		}
		c <- prometheus.MustNewConstMetric(q.sumDesc, prometheus.GaugeValue, sum, labelVals...)

		if len(vs) == 0 {
			return
		}

		slices.Sort(vs)
		c <- prometheus.MustNewConstMetric(q.qDesc, prometheus.GaugeValue, vs[0], append(labelVals, "0")...)
		c <- prometheus.MustNewConstMetric(q.qDesc, prometheus.GaugeValue, vs[len(vs)*5/100], append(labelVals, "0.05")...)
		c <- prometheus.MustNewConstMetric(q.qDesc, prometheus.GaugeValue, vs[len(vs)/2], append(labelVals, "0.5")...)
		c <- prometheus.MustNewConstMetric(q.qDesc, prometheus.GaugeValue, vs[len(vs)*9/10], append(labelVals, "0.9")...)
		c <- prometheus.MustNewConstMetric(q.qDesc, prometheus.GaugeValue, vs[len(vs)*95/100], append(labelVals, "0.95")...)
		c <- prometheus.MustNewConstMetric(q.qDesc, prometheus.GaugeValue, vs[len(vs)-1], append(labelVals, "1")...)
	}
}

func (q *metricSummary) Reset() {
	clear(q.values)
	clear(q.zeroes)
}

func allReports(db *sql.DB) func(func(*contract.Report) bool) {
	rows, err := db.Query(`SELECT Received, Report FROM ReportsJson WHERE Received > now() - '1 day'::INTERVAL LIMIT 1000`)
	if err != nil {
		log.Println("sql:", err)
		return nil
	}

	return func(fn func(*contract.Report) bool) {
		defer rows.Close()
		for rows.Next() {
			var rep contract.Report
			err := rows.Scan(&rep.Received, &rep)
			if err != nil {
				log.Println("sql:", err)
				return
			}
			if !fn(&rep) {
				return
			}
		}
	}
}
