package bridge

import (
	"github.com/Galdoba/ffproc/internal/pkg/sourcefiles"
	"github.com/Galdoba/ffproc/internal/pkg/table"
)

type Bridge struct {
	Key            string
	ProcessType    string
	SourceTableMap map[*sourcefiles.Sourcefile]table.Entry
	TableSourceMap map[table.Entry][]*sourcefiles.Sourcefile
}

func New(entry table.Entry, sources []*sourcefiles.Sourcefile) *Bridge {
	b := Bridge{}
	b.Key = entry.ProcessType + "--" + table.Key(entry)
	b.SourceTableMap = make(map[*sourcefiles.Sourcefile]table.Entry)
	b.TableSourceMap = make(map[table.Entry][]*sourcefiles.Sourcefile)
	b.TableSourceMap[entry] = sources
	for _, src := range sources {
		b.SourceTableMap[src] = entry
	}
	b.ProcessType = entry.ProcessType
	return &b
}

/*
bridgeMap[key] = bridge
bridge.SourceTableMap[entry]

*/
