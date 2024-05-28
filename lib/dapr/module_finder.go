package dapr

//type moduleFinder interface {
//	FindInvocationModulePath(skipDirs ...string) (string, error)
//	FindEventModulePath(skipDirs ...string) (string, error)
//}
//
//
//type regexModuleFinder struct {
//	sourceDir string
//}
//
//const (
//	regexInvocationModulePath = `.+dapr.InvocationModule$`
//	regexEventModulePath      = `.+dapr.EventModule$`
//)
//
//func newRegexModuleFinder(sourceDir string) moduleFinder {
//	return &regexModuleFinder{sourceDir: sourceDir}
//}
//
//func (r regexModuleFinder) FindInvocationModulePath(skipDirs ...string) (string, error) {
//	return r.findModulePath(regexInvocationModulePath, skipDirs...)
//}
//
//func (r regexModuleFinder) FindEventModulePath(skipDirs ...string) (string, error) {
//	return r.findModulePath(regexEventModulePath, skipDirs...)
//}
//
//// FindInvocationModulePath 尝试找到
//func (r regexModuleFinder) findModulePath(regex string, skipDirs ...string) (string, error) {
//	st, err := os.Stat(r.sourceDir)
//	if err != nil {
//		return "", err
//	}
//
//	if !st.IsDir() {
//		return "", fmt.Errorf("invalid dir, dir: %s", r.sourceDir)
//	}
//
//	var found string
//	match, _ := regexp.Compile(regex)
//	_ = filepath.Walk(r.sourceDir, func(path string, info os.FileInfo, err error) error {
//		if info.IsDir() && pie.Contains(skipDirs, info.Name()) {
//			return filepath.SkipDir
//		}
//
//		if !info.IsDir() && strings.HasSuffix(path, ".go") {
//			f, err := os.Open(path)
//			if err != nil {
//				return err
//			}
//			defer f.Close()
//
//			scanner := bufio.NewScanner(f)
//			for scanner.Scan() {
//				d := scanner.Text()
//				if match.MatchString(d) {
//					parentDir, _ := filepath.Split(path)
//					relDir, _ := filepath.Rel(r.sourceDir, parentDir)
//					found = relDir
//					break
//				}
//			}
//		}
//		return nil
//	})
//	return found, nil
//}
