package redirect

import (
	"LPManager/model"
	"LPManager/db"
	"log"
	"net/http"
)

func CheckRedirect(req *http.Request) (bool, string, error) {
	var result *model.RedirectType

	url := req.URL.String()
	domain := req.Host

	result, err := db.GetRedirect(url)
	if err != nil {
		log.Fatal(err)
		return false,nil,err
	}

	if result.Domain != domain {
		return false,nil,err
	}

	return true,result.To,err
}

/*func loopDetection(r model.RedirectType) (bool, error) {
	results, err := db.GetRedirects()
	if err != nil {
		return false,err
	}

	redMap := map[string]string{}

	for _, r := range *results {
		redMap[r.From] = r.To
	}

}*/
//054-9579308