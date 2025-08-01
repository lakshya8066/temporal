package searchattribute

import (
	"context"
	"maps"
	"sync"
	"sync/atomic"
	"time"

	enumspb "go.temporal.io/api/enums/v1"
	"go.temporal.io/api/serviceerror"
	persistencespb "go.temporal.io/server/api/persistence/v1"
	"go.temporal.io/server/common/clock"
	"go.temporal.io/server/common/dynamicconfig"
	"go.temporal.io/server/common/headers"
	"go.temporal.io/server/common/persistence"
)

const (
	cacheRefreshInterval              = 60 * time.Second
	cacheRefreshIfUnavailableInterval = 20 * time.Second
)

type (
	managerImpl struct {
		timeSource             clock.TimeSource
		clusterMetadataManager persistence.ClusterMetadataManager
		forceRefresh           dynamicconfig.BoolPropertyFn

		cacheUpdateMutex sync.Mutex
		cache            atomic.Value // of type cache
	}

	cache struct {
		// indexName -> NameTypeMap
		searchAttributes map[string]NameTypeMap
		dbVersion        int64
		expireOn         time.Time
	}
)

var _ Manager = (*managerImpl)(nil)

func NewManager(
	timeSource clock.TimeSource,
	clusterMetadataManager persistence.ClusterMetadataManager,
	forceRefresh dynamicconfig.BoolPropertyFn,
) *managerImpl {

	var saCache atomic.Value
	saCache.Store(cache{
		searchAttributes: map[string]NameTypeMap{},
		dbVersion:        0,
		expireOn:         time.Time{},
	})

	return &managerImpl{
		timeSource:             timeSource,
		cache:                  saCache,
		clusterMetadataManager: clusterMetadataManager,
		forceRefresh:           forceRefresh,
	}
}

// GetSearchAttributes returns all search attributes (including system and build-in) for specified index.
// indexName can be an empty string for backward compatibility.
func (m *managerImpl) GetSearchAttributes(
	indexName string,
	forceRefreshCache bool,
) (NameTypeMap, error) {

	now := m.timeSource.Now()
	saCache := m.cache.Load().(cache)

	if m.needRefreshCache(saCache, forceRefreshCache, now) {
		m.cacheUpdateMutex.Lock()
		saCache = m.cache.Load().(cache)
		if m.needRefreshCache(saCache, forceRefreshCache, now) {
			var err error
			saCache, err = m.refreshCache(saCache, now)
			if err != nil {
				m.cacheUpdateMutex.Unlock()
				return NameTypeMap{}, err
			}
		}
		m.cacheUpdateMutex.Unlock()
	}

	result := NameTypeMap{}
	indexSearchAttributes, ok := saCache.searchAttributes[indexName]
	if ok {
		result.customSearchAttributes = maps.Clone(indexSearchAttributes.customSearchAttributes)
	}

	// TODO (rodrigozhou): remove following block for v1.21.
	// Try to look for the empty string indexName for backward compatibility: up to v1.19,
	// empty string was used when Elasticsearch was not configured.
	// If there's a value, merging with current index name value. This is to avoid handling
	// all code references to GetSearchAttributes.
	if indexName != "" {
		indexSearchAttributes, ok = saCache.searchAttributes[""]
		if ok {
			if result.customSearchAttributes == nil {
				result.customSearchAttributes = maps.Clone(indexSearchAttributes.customSearchAttributes)
			} else {
				maps.Copy(result.customSearchAttributes, indexSearchAttributes.customSearchAttributes)
			}
		}
	}
	return result, nil
}

func (m *managerImpl) needRefreshCache(saCache cache, forceRefreshCache bool, now time.Time) bool {
	return forceRefreshCache || saCache.expireOn.Before(now) || m.forceRefresh()
}

func (m *managerImpl) refreshCache(saCache cache, now time.Time) (cache, error) {
	// TODO: specify a timeout for the context
	ctx := headers.SetCallerInfo(
		context.TODO(),
		headers.SystemBackgroundHighCallerInfo,
	)

	clusterMetadata, err := m.clusterMetadataManager.GetCurrentClusterMetadata(ctx)
	if err != nil {
		switch err.(type) {
		case *serviceerror.NotFound:
			// NotFound means cluster metadata was never persisted and custom search attributes are not defined.
			// Ignore the error.
			saCache.expireOn = now.Add(cacheRefreshInterval)
		case *serviceerror.Unavailable:
			// If persistence is Unavailable, ignore the error and use existing cache for cacheRefreshIfUnavailableInterval.
			saCache.expireOn = now.Add(cacheRefreshIfUnavailableInterval)
		default:
			return saCache, err
		}
		m.cache.Store(saCache)
		return saCache, nil
	}

	// clusterMetadata.Version <= saCache.dbVersion means DB is not changed.
	if clusterMetadata.Version <= saCache.dbVersion {
		saCache.expireOn = now.Add(cacheRefreshInterval)
		m.cache.Store(saCache)
		return saCache, nil
	}

	saCache = cache{
		searchAttributes: buildIndexNameTypeMap(clusterMetadata.GetIndexSearchAttributes()),
		expireOn:         now.Add(cacheRefreshInterval),
		dbVersion:        clusterMetadata.Version,
	}
	m.cache.Store(saCache)
	return saCache, nil
}

// SaveSearchAttributes saves search attributes to cluster metadata.
// indexName can be an empty string when Elasticsearch is not configured.
func (m *managerImpl) SaveSearchAttributes(
	ctx context.Context,
	indexName string,
	newCustomSearchAttributes map[string]enumspb.IndexedValueType,
) error {

	clusterMetadataResponse, err := m.clusterMetadataManager.GetCurrentClusterMetadata(ctx)
	if err != nil {
		return err
	}

	clusterMetadata := clusterMetadataResponse.ClusterMetadata
	if clusterMetadata.IndexSearchAttributes == nil {
		clusterMetadata.IndexSearchAttributes = map[string]*persistencespb.IndexSearchAttributes{indexName: nil}
	}
	clusterMetadata.IndexSearchAttributes[indexName] = &persistencespb.IndexSearchAttributes{CustomSearchAttributes: newCustomSearchAttributes}
	_, err = m.clusterMetadataManager.SaveClusterMetadata(ctx, &persistence.SaveClusterMetadataRequest{
		ClusterMetadata: clusterMetadata,
		Version:         clusterMetadataResponse.Version,
	})
	// Flush local cache, even if there was an error, which is most likely version mismatch (=stale cache).
	m.cache.Store(cache{})

	return err
}
