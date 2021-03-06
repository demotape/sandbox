package sandbox

import (
	"fmt"
	"path"
	"reflect"
)

type BindVolume struct {
	DotSsh   string
	Tutorial string
}

func (b *BindVolume) VolumeMappings() []string {
	var vols []string
	val := reflect.ValueOf(b).Elem()

	for i := 0; i < val.NumField(); i++ {
		vols = append(vols, fmt.Sprintf("%v", val.Field(i)))
	}

	return vols
}

func NewBindVolumes(runtimeEnv *RuntimeEnv) BindVolume {

	return BindVolume{
		DotSsh:   runtimeEnv.SshEnv.DotSshDir + ":/root/.ssh",
		Tutorial: path.Join(runtimeEnv.RuntimeDir, WelcomeDir) + ":" + Container.WelcomeDir,
	}
}
