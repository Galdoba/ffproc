package sourcefiles

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/Galdoba/ffproc/internal/pkg/define"
)

type Sourcefile struct {
	Name string
	Dir  string
	Tags map[string]string
	err  error
}

func Inspect(paths []string) ([]*Sourcefile, error) {
	if len(paths) < 1 {
		return nil, fmt.Errorf("no source files provided")
	}
	sourcefiles := []*Sourcefile{}
	for _, path := range paths {
		sourcefiles = append(sourcefiles, NewSource(path))
	}
	for _, sf := range sourcefiles {
		if sf.err != nil {
			return sourcefiles, sf.err
		}
	}
	return sourcefiles, nil
}

func NewSource(name string) *Sourcefile {
	dir := filepath.Dir(name) + string(filepath.Separator)
	dir = filepath.ToSlash(dir)
	base := filepath.Base(name)
	source := Sourcefile{}
	source.Name = base
	source.Dir = dir
	source.Tags = make(map[string]string)
	parts := []string{}
	parts = strings.Split(source.Name, "--")
	switch len(parts) {
	case 0, 1:
		source.err = errors.New("no tags")
	default:
		for _, tag := range parts[0 : len(parts)-1] {
			switch {
			case isProjectType(tag):
				source.Tags[define.TAG_PROJ_TYPE] = tag
			case isEpisodeTag(tag):
				source.Tags[define.TAG_PROJ_EPISODE] = tag
				source.Tags[define.TAG_PROJ_SEASON] = strings.Split(tag, "e")[0]
			case tag == source.Name:
			default:
				source.Tags[define.TAG_PROJ_BASE] = tag
			}
		}
	}
	source.err = checkTags(&source)
	return &source
}

func checkTags(sf *Sourcefile) error {
	switch sf.Tags[define.TAG_PROJ_TYPE] {
	case define.PROJ_TYPE_FLM, define.PROJ_TYPE_TRL:
	case define.PROJ_TYPE_SER:
		if sf.Tags[define.TAG_PROJ_EPISODE] == "" {
			return fmt.Errorf("sourcefile: %v: no episode tag provided", sf.Name)
		}
	case "":
		return fmt.Errorf("sourcefile: %v: no project type tag provided", sf.Name)
	default:
		return fmt.Errorf("sourcefile: %v: unknown project type tag provided", sf.Name)
	}

	return nil
}

func ProjectTypesMatch(tgs ...*Sourcefile) bool {
	ptTags := []string{}
	for _, tg := range tgs {
		ptTags = append(ptTags, tg.Tags[define.TAG_PROJ_TYPE])
	}
	return sliceMatch(ptTags)
}

func ProjectBaseMatch(tgs ...*Sourcefile) bool {
	ptTags := []string{}
	for _, tg := range tgs {
		ptTags = append(ptTags, tg.Tags[define.TAG_PROJ_BASE])
	}
	return sliceMatch(ptTags)
}

func sliceMatch(sl []string) bool {
	if len(sl) == 0 {
		return false
	}
	first := sl[0]
	for _, s := range sl {
		if s != first {
			return false
		}
	}
	return true
}

func isProjectType(tag string) bool {
	switch tag {
	default:
		return false
	case define.PROJ_TYPE_FLM, define.PROJ_TYPE_SER, define.PROJ_TYPE_TRL:
		return true
	}
}

func isEpisodeTag(tag string) bool {
	re := regexp.MustCompile(`s(\d+)e(\d+)`)
	match := re.FindStringSubmatch(tag)
	switch len(match) {
	case 0:
		return false
	default:
		return true
	}
}

func SplitByKeys(sources []*Sourcefile) map[string][]*Sourcefile {
	output := make(map[string][]*Sourcefile)
	for _, src := range sources {
		// key := src.Tags[define.TAG_PROJ_TYPE] + "--"
		// if val, ok := src.Tags[define.TAG_PROJ_EPISODE]; ok {
		// 	key += val + "--"
		// }
		// key += src.Tags[define.TAG_PROJ_BASE]
		key := src.Key()
		output[key] = append(output[key], src)
	}
	return output
}

func ProcessType(projData []*Sourcefile) string {
	if len(projData) < 1 {
		return ""
	}
	etalon := projData[0].Tags[define.TAG_PROJ_TYPE]
	for _, data := range projData {
		if data.Tags[define.TAG_PROJ_TYPE] != etalon {
			return ""
		}
	}
	return etalon
}

func appendUnique(sl []string, str string) []string {
	for _, s := range sl {
		if s == str {
			return sl
		}
	}
	sl = append(sl, str)
	return sl
}

func (src *Sourcefile) Key() string {
	base := src.Tags[define.TAG_PROJ_BASE]
	switch src.Tags[define.TAG_PROJ_TYPE] {
	case define.PROJ_TYPE_SER:
		season := src.Tags[define.TAG_PROJ_SEASON]
		base = base + "_" + strings.TrimPrefix(season, "s")
	}
	return base
}

func (src *Sourcefile) EpisodeTag() string {
	episode := ""
	switch src.Tags[define.TAG_PROJ_TYPE] {
	case define.PROJ_TYPE_SER:
		if ep, ok := src.Tags[define.TAG_PROJ_EPISODE]; ok {
			episode = ep
		}
	}
	return episode
}

func SortedByEpisode(sources []*Sourcefile) [][]*Sourcefile {
	eps := []string{}
	srcNames := []string{}
	for _, src := range sources {
		eps = appendUnique(eps, src.EpisodeTag())
		srcNames = append(srcNames, src.Name)
	}
	sort.Strings(eps)
	sort.Strings(srcNames)
	sortedSrcs := [][]*Sourcefile{}

	for _, ep := range eps {
		localN := 0

		namedSrc := []*Sourcefile{}
		for _, src := range sources {
			if ep == src.EpisodeTag() {
				namedSrc = append(namedSrc, src)
				localN++
			}

		}
		sortedSrcs = append(sortedSrcs, namedSrc)

	}
	return sortedSrcs
}
