package parsing

import (
	"regexp"

	"github.com/samber/lo"
	"github.com/surkovvs/kpxc-exp/internal/entity"
)

type key string

const (
	notesKey string = "Notes"

	passwordKey key = `Password`
	userNameKey key = `UserName`
	urlKey      key = `URL`
)

var (
	defRexes                    map[key]*regexp.Regexp
	rexFlag, rexAdd, rexAddName *regexp.Regexp
)

func init() {
	defRexes = make(map[key]*regexp.Regexp, 3)
	defRexes[passwordKey] = regexp.MustCompile(`#(?i:Password)={[a-zA-Z_][a-zA-Z0-9_]+}`)
	defRexes[userNameKey] = regexp.MustCompile(`#(?i:UserName)={[a-zA-Z_][a-zA-Z0-9_]+}`)
	defRexes[urlKey] = regexp.MustCompile(`#(?i:URL)={[a-zA-Z_][a-zA-Z0-9_]+}`)

	rexFlag = regexp.MustCompile(`#(?i:exportable)`)
	rexAdd = regexp.MustCompile(`#(?i:add)={[a-zA-Z_][a-zA-Z0-9_]+}='.*'{1}`)
	rexAddName = regexp.MustCompile(`={[a-zA-Z_][a-zA-Z0-9_]+}=`)
}

// root num==0
func ParseGroups(head entity.Group, num int) entity.EnvGroup {
	var envEntrys []entity.EnvEntry
	for i, he := range head.Entrys {
		envs := parseEntys(he)
		if envs == nil {
			continue
		}
		envEntrys = append(envEntrys, entity.EnvEntry{
			Num:       i + 1,
			GroupName: head.Name,
			Envs:      envs,
		})
	}
	if len(envEntrys) == 0 && len(head.Groups) == 0 {
		return entity.EnvGroup{}
	}

	envGroup := entity.EnvGroup{
		Num:    num,
		Name:   head.Name,
		Note:   head.Notes,
		Entrys: envEntrys,
	}

	for i, headSub := range head.Groups {
		parsed := ParseGroups(headSub, i+1)
		if parsed.Entrys == nil && parsed.SubGroups == nil {
			continue
		}
		envGroup.SubGroups = append(envGroup.SubGroups, parsed)
	}

	return envGroup
}

func parseEntys(entry entity.Entry) []entity.Env {
	var envs []entity.Env
	posMap := lo.SliceToMap(entry.Positions, func(pos entity.Position) (string, string) {
		return pos.Key, pos.Value
	})

	if !isExportable(posMap) {
		return nil
	}

	envs = append(envs, parseDefEnvs(passwordKey, posMap)...)
	envs = append(envs, parseDefEnvs(userNameKey, posMap)...)
	envs = append(envs, parseDefEnvs(urlKey, posMap)...)
	envs = append(envs, parseAddEnvs(posMap[notesKey])...)

	return envs
}

func isExportable(poss map[string]string) bool {
	notesVal, ok := poss[notesKey]
	return ok && rexFlag.MatchString(notesVal)
}

func parseDefEnvs(k key, poss map[string]string) []entity.Env {
	envNames := defRexes[k].FindAllString(poss[notesKey], -1)
	if len(envNames) == 0 {
		return nil
	}

	envVal, ok := poss[string(k)]
	if !ok {
		return nil
	}

	var envs []entity.Env
	for _, name := range envNames {
		envs = append(envs, entity.Env{
			Name:  name[3+len(k) : len(name)-1],
			Value: envVal,
		})
	}
	return envs
}

func parseAddEnvs(toParse string) []entity.Env {
	var envs []entity.Env
	rawAdds := rexAdd.FindAllString(toParse, -1)
	for _, rawAdd := range rawAdds {
		rawName := rexAddName.FindString(rawAdd)
		envs = append(envs, entity.Env{
			Name:  rawName[2 : len(rawName)-2],
			Value: rawAdd[5+len(rawName) : len(rawAdd)-1],
		})
	}
	return envs
}
