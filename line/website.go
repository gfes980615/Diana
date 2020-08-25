package line
//
//import "github.com/gfes980615/Diana/repository/website"
//
//func SaveWebsite(web []string) error {
//	var err error
//	switch len(web) {
//	case 1:
//		err = saveToDefaultTag(web[0])
//	case 2:
//		err = saveToTag(web[0], web[1])
//	}
//	return err
//}
//
//func saveToDefaultTag(url string) error {
//	webSite := website.WebsiteRepository{}
//	return webSite.Insert(url, "未分類")
//}
//
//func saveToTag(tag, url string) error {
//	webSite := website.WebsiteRepository{}
//	return webSite.Insert(url, tag)
//}
