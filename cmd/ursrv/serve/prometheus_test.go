package serve

// func TestAddReport(t *testing.T) {
// 	r := &contract.Report{
// 		UniqueID:       "123",
// 		Version:        "1.0.0",
// 		Platform:       "linux",
// 		NumFolders:     42,
// 		FolderMaxFiles: 100,
// 	}
// 	s := newMetricsSet()
// 	s.addReport(r)

// 	m := make(chan prometheus.Metric)
// 	go func() {
// 		for metric := range m {
// 			var dm dto.Metric
// 			metric.Write(&dm)
// 			fmt.Println(&dm)
// 		}
// 	}()
// 	s.Collect(m)
// 	close(m)
// }
