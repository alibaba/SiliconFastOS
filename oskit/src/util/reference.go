package util

import "strings"

const DOCKER_REPO = "reg.docker.alibaba-inc.com"

func ReferenceExpand(ref string) string {
	parts := strings.Split(ref, "/")
	switch len(parts) {
	case 1:
		return DOCKER_REPO + "/silicon_fast/" + ref
	case 2:
		return DOCKER_REPO + "/" + ref
	default:
		return ref
	}
}
