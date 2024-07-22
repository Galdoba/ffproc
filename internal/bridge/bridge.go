package bridge

import (
	"fmt"

	"github.com/Galdoba/ffproc/internal/mfscan"
	"github.com/Galdoba/ffproc/internal/pkg/sourcefiles"
	"github.com/Galdoba/ffproc/internal/pkg/table"
)

type Bridge struct {
	Key              string
	ProcessType      string
	SourceTableMap   map[*sourcefiles.Sourcefile]table.Entry
	TableSourceMap   map[table.Entry][]*sourcefiles.Sourcefile
	SourceStreamMap  map[*sourcefiles.Sourcefile]map[string]*mfscan.Stream
	EpisodeSourceMap map[string][]*sourcefiles.Sourcefile
}

func New(entry table.Entry, allsources []*sourcefiles.Sourcefile) *Bridge {
	b := Bridge{}
	b.Key = entry.ProcessType + "--" + table.Key(entry)
	b.SourceTableMap = make(map[*sourcefiles.Sourcefile]table.Entry)
	b.TableSourceMap = make(map[table.Entry][]*sourcefiles.Sourcefile)
	b.SourceStreamMap = make(map[*sourcefiles.Sourcefile]map[string]*mfscan.Stream)
	b.TableSourceMap[entry] = allsources

	for _, episodeSrc := range sourcefiles.SortedByEpisode(allsources) {
		for i, src := range episodeSrc {

			b.SourceTableMap[src] = entry
			scanRep := mfscan.NewProfile()
			scanRep.ConsumeFile(src.Dir + src.Name)
			strMap := scanRep.MapStreams()
			b.SourceStreamMap[src] = make(map[string]*mfscan.Stream)
			bufferMap := make(map[string]*mfscan.Stream)
			for key, stream := range strMap {
				// fmt.Println("AAAADDD", fmt.Sprintf("%v%v %v", i, key, e))
				bufferMap[fmt.Sprintf("%v%v", i, key)] = stream
				// localSourceNum++
			}
			b.SourceStreamMap[src] = bufferMap
		}
	}

	b.ProcessType = entry.ProcessType
	return &b
}

/*
bridgeMap[key] = bridge
bridge.SourceTableMap[entry]

*/
