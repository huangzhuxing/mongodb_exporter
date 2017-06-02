package collector

import (
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	count = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "collection",
		Name:      "total_objects",
		Help:      "The number of objects or documents in this collection",
	}, []string{"ns"})

	size = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "collection",
		Name:      "size_bytes",
		Help:      "The total size in memory of all records in a collection",
	}, []string{"ns"})

	avgObjSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "collection",
		Name:      "avg_objsize_bytes",
		Help:      "The average size of an object in the collection",
	}, []string{"ns"})

	storageSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "collection",
		Name:      "storage_size_bytes",
		Help:      "The total amount of storage allocated to this collection for document storage",
	}, []string{"ns"})

	collIndexSize = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: Namespace,
		Subsystem: "collection",
		Name:      "index_size_bytes",
		Help:      "The total size of all indexes",
	}, []string{"ns"})
)

type CollectionStatus struct {
	Name        string `bson:"ns"`
	Count       int32  `bson:"count"`
	Size        int32  `bson:"size"`
	AvgSize     int32  `bson:"avgObjSize"`
	StorageSize int32  `bson:"storageSize"`
	IndexSize   int32  `bson:"totalIndexSize"`
}

func (collStatus *CollectionStatus) Export(ch chan<- prometheus.Metric) {
	count.WithLabelValues(collStatus.Name).Set(float64(collStatus.Count))
	size.WithLabelValues(collStatus.Name).Set(float64(collStatus.Size))
	avgObjSize.WithLabelValues(collStatus.Name).Set(float64(collStatus.AvgSize))
	storageSize.WithLabelValues(collStatus.Name).Set(float64(collStatus.StorageSize))
	collIndexSize.WithLabelValues(collStatus.Name).Set(float64(collStatus.IndexSize))

	count.Collect(ch)
	size.Collect(ch)
	avgObjSize.Collect(ch)
	storageSize.Collect(ch)
	collIndexSize.Collect(ch)

	count.Reset()
	size.Reset()
	avgObjSize.Reset()
	storageSize.Reset()
	collIndexSize.Reset()
}

func (collStatus *CollectionStatus) Describe(ch chan<- *prometheus.Desc) {
	count.Describe(ch)
	size.Describe(ch)
	avgObjSize.Describe(ch)
	storageSize.Describe(ch)
	collIndexSize.Describe(ch)
}

func GetCollectionStatus(session *mgo.Session, db string, collection string) *CollectionStatus {
	var collStatus CollectionStatus
	err := session.DB(db).Run(bson.D{{"collStats", collection}, {"scale", 1}}, &collStatus)
	if err != nil {
		glog.Error(err)
		return nil
	}

	return &collStatus
}
