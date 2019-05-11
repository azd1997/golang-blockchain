package blockchain

import "os"

const (
	dbPath = "./tmp/blocks/blocks_%s"
)

/*检查数据库是否存在*/
func DbExists(path string) bool {
	if _, err := os.Stat(path + "/MANIFEST"); os.IsNotExist(err) {
		return false
	}
	return true
}
