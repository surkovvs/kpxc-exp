package layout

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/surkovvs/kpxc-exp/internal/entity"
)

func GetLayout(envGroup entity.EnvGroup) (map[string]entity.EnvEntry, string) {
	layout := &strings.Builder{}
	em := make(map[string]entity.EnvEntry)
	collectEntrysFromGroup(envGroup, "", em, layout)
	return em, layout.String()
}

func collectEntrysFromGroup(envGroup entity.EnvGroup, parentName string, entryMap map[string]entity.EnvEntry, layout *strings.Builder) {
	if parentName == "" && envGroup.Num == 0 {
		parentName = "R"
	}
	layout.WriteString(fmt.Sprintf("--- Group %s: %s ---\nEntrys:\n", envGroup.Name, parentName))

	for _, entry := range envGroup.Entrys {
		entryCode := parentName + ":" + strconv.Itoa(entry.Num)
		entryMap[entryCode] = entry
		layout.WriteString(fmt.Sprintf("%s\tenvs:\n", entryCode))
		for _, env := range entry.Envs {
			layout.WriteString("\t" + env.String() + "\n")
		}
	}

	for _, group := range envGroup.SubGroups {
		subname := strconv.Itoa(group.Num)
		if parentName != "R" {
			subname = parentName + "." + subname
		}

		collectEntrysFromGroup(group, subname, entryMap, layout)
	}
}
