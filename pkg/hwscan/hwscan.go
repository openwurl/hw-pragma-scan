package hwscan

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
)

// Scanner encapsulates the scanning operation
type Scanner struct {
	Target string
	Pragma *Pragma
}

// Report displays the data in a meaningful way
func (s *Scanner) Report() error {
	fmt.Printf("\n\n\n")
	fmt.Printf("Scanned: %s\n", s.Target)

	data := [][]string{
		[]string{
			"CDN Cache",
			fmt.Sprintf("%d", s.Pragma.XHWCacheTTL),
			fmt.Sprintf("%d", s.Pragma.XHWCacheTTL/60),
			fmt.Sprintf("%d", s.Pragma.XHWCacheTTL/60/60),
		},
	}
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Field", "TTL (seconds)", "TTL (minutes)", "TTL (hours)"})
	//table.SetFooter([]string{"", "", "", ""})
	table.SetRowLine(true)
	table.AppendBulk(data)
	table.Render()

	now := time.Now()
	lastModified := now.Sub(time.Unix(int64(s.Pragma.XHWCacheLastModified), 0))
	lastRequest := now.Sub(time.Unix(int64(s.Pragma.XHWCacheLastRequest), 0))
	lastRefresh := now.Sub(time.Unix(int64(s.Pragma.XHWCacheLastRefresh), 0))
	originated := now.Sub(time.Unix(int64(s.Pragma.XHWCacheOriginated), 0))

	headerData := [][]string{
		[]string{"File Size", fmt.Sprintf("%d", s.Pragma.XHWCacheFileSize)},
		[]string{"Access-Control-Allow-Origin", strings.Join(s.Pragma.AccessControlAllowOrigin, "")},
		[]string{"Cache-Control", strings.Join(s.Pragma.CacheControl, ",")},
		[]string{"Content-Type", strings.Join(s.Pragma.ContentType, "")},
		[]string{"X-HW-Cache-Compressed-Size", fmt.Sprintf("%v", s.Pragma.XHWCacheCompressedSize)},
		[]string{"X-HW-Cache-Behavior", strings.Join(s.Pragma.XHWCacheBehavior, "")},
		[]string{"X-HW-Cache-Last-Modified", lastModified.String()},
		[]string{"X-HW-Cache-Originated", originated.String()},
		[]string{"X-HW-Cache-Last-Refresh", lastRefresh.String()},
		[]string{"X-HW-Cache-Last-Request", lastRequest.String()},
	}
	table2 := tablewriter.NewWriter(os.Stdout)
	table2.SetHeader([]string{"Field", "Value"})
	table2.SetRowLine(true)
	table2.AppendBulk(headerData)
	table2.Render()

	fmt.Printf("X-HW: %s\n", strings.Join(s.Pragma.XHW, ","))

	//fmt.Printf("CDN Cache: [%d] seconds ( [%d] minutes or [%d] hours )\n", s.Pragma.XHWCacheTTL, s.Pragma.XHWCacheTTL/60, s.Pragma.XHWCacheTTL/60/60)
	//fmt.Printf("Last Refreshed: [%d] seconds ( [%d] minutes or [%d] hours ) ago\n", s.Pragma.XHWCacheLastRefresh, s.Pragma.XHWCacheLastRefresh/60, s.Pragma.XHWCacheLastRefresh/60/60)
	//fmt.Printf("Reported browser cache: ( %s )\n", strings.Join(s.Pragma.CacheControl, ", "))
	//fmt.Printf("Content Length: %d ", s.Pragma.ContentLength)
	//if s.Pragma.ContentLength > 0 {
	//	fmt.Printf("indicates content returned\n")
	//} else {
	//	fmt.Printf("\n")
	//}
	return nil
}

// Scan scans the target host
func (s *Scanner) Scan() error {
	c := &http.Client{}
	req, err := http.NewRequest("GET", s.Target, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Pragma", "X-HW-Cache-All")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	errs := s.Pack(resp)
	if errs != nil {
		log.Println("Errors packing headers:")
		for _, thisErr := range errs {
			log.Println(thisErr)
		}
	}

	return nil
}

// Pack packs the pragma struct with header data from a response
func (s *Scanner) Pack(r *http.Response) []error {
	s.Pragma = &Pragma{}
	prepackStorage := make(map[string]interface{})
	reflection := reflect.TypeOf(*s.Pragma)

	//spew.Dump(r.Header)

	// Extract all headers that correlate to a field
	for k, v := range r.Header {
		//fmt.Println("==================")
		//fmt.Printf("===== Key: %s\n", k)
		for x := 0; x < reflection.NumField(); x++ {
			//fmt.Println("=========")
			//fmt.Printf("Iter %d of %d\n", x, reflection.NumField())
			thisField := reflection.Field(x)
			if val, ok := thisField.Tag.Lookup("pragma"); ok {
				//fmt.Printf("found | %s\n", val)
				//fmt.Printf("compare | %s\n", k)
				if strings.ToLower(k) == strings.ToLower(val) {
					//fmt.Printf("Found key: %v\n", k)
					prepackStorage[thisField.Name] = v
				}
			}
		}
	}

	err := s.Pragma.Pack(prepackStorage)
	if err != nil {
		return err
	}

	return nil
}

// Pragma captures the pragma headers
type Pragma struct {
	Connection                 []string    `pragma:"Connection"`
	AcceptRanges               []string    `pragma:"Accept-Ranges"`
	ContentLength              int         `pragma:"Content-Length"`
	ContentType                []string    `pragma:"Content-Type"`
	ContentMD5                 []string    `pragma:"Content-MD5"`
	CacheControl               []string    `pragma:"Cache-Control"`
	AccessControlAllowHeaders  []string    `pragma:"Access-Control-Allow-Headers"`
	AccessControlExposeHeaders []string    `pragma:"Access-Control-Expose-Headers"`
	AccessControlAllowMethods  []string    `pragma:"Access-Control-Allow-Methods"`
	AccessControlAllowOrigin   []string    `pragma:"Access-Control-Allow-Origin"`
	XHW                        []string    `pragma:"X-HW"`
	XHWCacheKey                []string    `pragma:"X-HW-Cache-Key"`
	XHWCacheFileName           []string    `pragma:"X-HW-Cache-File-Name"`
	XHWCacheTTL                int         `pragma:"X-HW-Cache-TTL"`
	XHWCacheCRC                int         `pragma:"X-HW-Cache-CRC"`
	XHWCacheLastModified       int         `pragma:"X-HW-Cache-Last-Modified"`
	XHWCacheOriginated         int         `pragma:"X-HW-Cache-Originated"`
	XHWCacheLastRefresh        int         `pragma:"X-HW-Cache-Last-Refresh"`
	XHWCacheLastRequest        int         `pragma:"X-HW-Cache-Last-Request"`
	XHWCacheMimeType           []string    `pragma:"X-HW-Cache-Mime-Type"`
	XHWCacheHeaders            []string    `pragma:"X-HW-Cache-Headers"` // *2
	XHWCacheFileSize           int         `pragma:"X-HW-Cache-File-Size"`
	XHWCacheCompressedSize     interface{} `pragma:"X-HW-Cache-Compressed-Size"`
	XHWCacheControl            []string    `pragma:"X-HW-Cache-Control"`
	XHWCacheBehavior           []string    `pragma:"X-HW-Cache-Behavior"`
}

// Pack loads the struct with prepack contents from reflection
func (p *Pragma) Pack(contents map[string]interface{}) []error {
	var ok bool
	var errs []error
	p.XHW = contents["XHW"].([]string)
	p.XHWCacheKey = contents["XHWCacheKey"].([]string)
	p.XHWCacheFileName = contents["XHWCacheFileName"].([]string)
	p.XHWCacheMimeType = contents["XHWCacheMimeType"].([]string)
	p.XHWCacheHeaders = contents["XHWCacheHeaders"].([]string)
	p.XHWCacheCompressedSize = contents["XHWCacheCompressedSize"].(interface{})
	p.XHWCacheControl = contents["XHWCacheControl"].([]string)
	p.XHWCacheBehavior = contents["XHWCacheBehavior"].([]string)

	p.Connection, ok = contents["Connection"].([]string)
	if !ok {
		errs = append(errs, errConversionFailed("Connection"))
	}
	p.AcceptRanges, ok = contents["AcceptRanges"].([]string)
	if !ok {
		errs = append(errs, errConversionFailed("AcceptRanges"))
	}
	p.ContentType, ok = contents["ContentType"].([]string)
	if !ok {
		errs = append(errs, errConversionFailed("ContentType"))
	}
	p.ContentMD5, ok = contents["ContentMD5"].([]string)
	if !ok {
		errs = append(errs, errConversionFailed("ContentMD5"))
	}
	p.CacheControl, ok = contents["CacheControl"].([]string)
	if !ok {
		errs = append(errs, errConversionFailed("CacheControl"))
	}
	p.AccessControlAllowHeaders, ok = contents["AccessControlAllowHeaders"].([]string)
	if !ok {
		errs = append(errs, errConversionFailed("AccessControlAllowHeaders"))
	}
	p.AccessControlExposeHeaders, ok = contents["AccessControlExposeHeaders"].([]string)
	if !ok {
		errs = append(errs, errConversionFailed("AccessControlExposeHeaders"))
	}
	p.AccessControlAllowMethods, ok = contents["AccessControlAllowMethods"].([]string)
	if !ok {
		errs = append(errs, errConversionFailed("AccessControlAllowMethods"))
	}
	p.AccessControlAllowOrigin, ok = contents["AccessControlAllowOrigin"].([]string)
	if !ok {
		errs = append(errs, errConversionFailed("AccessControlAllowOrigin"))
	}

	// ints
	ContentLength, err := extractIntegerFromStringSlice(contents["ContentLength"].([]string))
	if err != nil {
		errs = append(errs, err)
	}
	p.ContentLength = ContentLength

	XHWCacheTTL, err := extractIntegerFromStringSlice(contents["XHWCacheTTL"].([]string))
	if err != nil {
		errs = append(errs, err)
	}
	p.XHWCacheTTL = XHWCacheTTL

	XHWCacheCRC, err := extractIntegerFromStringSlice(contents["XHWCacheCRC"].([]string))
	if err != nil {
		errs = append(errs, err)
	}
	p.XHWCacheCRC = XHWCacheCRC

	XHWCacheLastModified, err := extractIntegerFromStringSlice(contents["XHWCacheLastModified"].([]string))
	if err != nil {
		errs = append(errs, err)
	}
	p.XHWCacheLastModified = XHWCacheLastModified

	XHWCacheOriginated, err := extractIntegerFromStringSlice(contents["XHWCacheOriginated"].([]string))
	if err != nil {
		errs = append(errs, err)
	}
	p.XHWCacheOriginated = XHWCacheOriginated

	XHWCacheLastRefresh, err := extractIntegerFromStringSlice(contents["XHWCacheLastRefresh"].([]string))
	if err != nil {
		errs = append(errs, err)
	}
	p.XHWCacheLastRefresh = XHWCacheLastRefresh

	XHWCacheFileSize, err := extractIntegerFromStringSlice(contents["XHWCacheFileSize"].([]string))
	if err != nil {
		errs = append(errs, err)
	}
	p.XHWCacheFileSize = XHWCacheFileSize

	XHWCacheLastRequest, err := extractIntegerFromStringSlice(contents["XHWCacheLastRequest"].([]string))
	if err != nil {
		errs = append(errs, err)
	}
	p.XHWCacheLastRequest = XHWCacheLastRequest

	return errs
}

func extractHeadersToMap(r *http.Response) map[string]interface{} {
	headers := make(map[string]interface{})

	for k, v := range r.Header {
		headers[strings.ToLower(k)] = string(v[0])
	}

	return headers
}

func extractIntegerFromStringSlice(s []string) (int, error) {
	i, err := strconv.Atoi(s[0])
	return i, err
}

func errConversionFailed(field string) error {
	return fmt.Errorf("failed to convert %s", field)
}

// Pack packs the struct with data from a scan

////////////////////////////////////////////////
/*

	*2
		X-HW-Cache-Headers: Content-Length=411720;Content-MD5=mtz16PLObXslE3A3NXhhwA==;Expires=Mon, 16 Dec 2019 18:41:33 GMT;Cache-Control=max-age=0, no-cache, no-store;Pragma=no-cache;Date=Mon, 16 Dec 2019 18:41:33 GMT;Connection=keep-alive;Access-Control-Allow-Headers=origin,range,hdntl,hdnts;Access-Control-Expose-Headers=Server,range,hdntl,hdnts;Access-Control-Allow-Methods=GET, HEAD, OPTIONS;Access-Control-Allow-Origin=*;Content-Type=video/MP2T
*/
