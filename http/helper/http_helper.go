package helper

import (
	"finance-api/config/helper/str"
	"finance-api/http/response"

	"math"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/universal-translator"
	"github.com/labstack/echo"
	validator "gopkg.in/go-playground/validator.v9"
)

const (
	textError             string = `error`
	textOk                string = `ok`
	codeSuccess           int    = 200
	codeDenied            int    = 202
	codeBadRequestError   int    = 400
	codeUnauthorizedError int    = 401
	codeDatabaseError     int    = 402
	codeValidationError   int    = 403
	codeNotFound          int    = 404
)

// ResponseHelper ...
type ResponseHelper struct {
	C        echo.Context
	Status   string
	Message  string
	Data     interface{}
	Code     int // not the http code
	CodeType string
}

// HTTPHelper ...
type HTTPHelper struct {
	Validate   *validator.Validate
	Translator ut.Translator
}

func (u *HTTPHelper) getTypeData(i interface{}) string {
	v := reflect.ValueOf(i)
	v = reflect.Indirect(v)

	return v.Type().String()
}

// GetStatusCode ...
func (u *HTTPHelper) GetStatusCode(err error) int {
	statusCode := http.StatusOK
	if err != nil {
		switch u.getTypeData(err) {
		case "models.ErrorUnauthorized":
			statusCode = http.StatusUnauthorized
		case "models.ErrorNotFound":
			statusCode = http.StatusNotFound
		case "models.ErrorConflict":
			statusCode = http.StatusConflict
		case "models.ErrorInternalServer":
			statusCode = http.StatusInternalServerError
		default:
			statusCode = http.StatusInternalServerError
		}
	}

	return statusCode
}

// SetResponse ...
// Set response data.
func (u *HTTPHelper) SetResponse(c echo.Context, status string, message string, data interface{}, code int, codeType string) ResponseHelper {
	return ResponseHelper{c, status, message, data, code, codeType}
}

// SendError ...
// Send error response to consumers.
func (u *HTTPHelper) SendError(c echo.Context, message string, data interface{}, code int, codeType string) error {
	res := u.SetResponse(c, `error`, message, data, code, codeType)

	return u.SendResponse(res)
}

// SendBadRequest ...
// Send bad request response to consumers.
func (u *HTTPHelper) SendBadRequest(c echo.Context, message string, data interface{}) error {
	res := u.SetResponse(c, `error`, message, data, codeBadRequestError, `badRequest`)

	return u.SendResponse(res)
}

// SendBadRequest ...
// Send bad request response to consumers.
func (u *HTTPHelper) SendDeniedRequest(c echo.Context, message string, data interface{}) error {
	res := u.SetResponse(c, `denied`, message, data, codeBadRequestError, `badRequest`)

	return u.SendResponse(res)
}

// SendValidationError ...
// Send validation error response to consumers.
func (u *HTTPHelper) SendValidationError(c echo.Context, validationErrors validator.ValidationErrors) error {
	errorResponse := map[string][]string{}
	errorTranslation := validationErrors.Translate(u.Translator)
	for _, err := range validationErrors {
		errKey := str.Underscore(err.StructField())
		errorResponse[errKey] = append(errorResponse[errKey], errorTranslation[err.Namespace()])
	}

	return c.JSON(400, map[string]interface{}{
		"code":         codeValidationError,
		"code_type":    "[Gateway] validationError",
		"code_message": errorResponse,
		"data":         u.EmptyJsonMap(),
	})
}

// SendDatabaseError ...
// Send database error response to consumers.
func (u *HTTPHelper) SendDatabaseError(c echo.Context, message string, data interface{}) error {
	return u.SendError(c, message, data, codeDatabaseError, `databaseError`)
}

// SendUnauthorizedError ...
// Send unauthorized response to consumers.
func (u *HTTPHelper) SendUnauthorizedError(c echo.Context, message string, data interface{}) error {
	return u.SendError(c, message, data, codeUnauthorizedError, `unAuthorized`)
}

// SendNotFoundError ...
// Send not found response to consumers.
func (u *HTTPHelper) SendNotFoundError(c echo.Context, message string, data interface{}) error {
	return u.SendError(c, message, data, codeNotFound, `notFound`)
}

// SendSuccess ...
// Send success response to consumers.
func (u *HTTPHelper) SendSuccess(c echo.Context, message string, data interface{}) error {
	res := u.SetResponse(c, `ok`, message, data, codeSuccess, `success`)

	return u.SendResponse(res)
}

// SendResponse ...
// Send response
func (u *HTTPHelper) SendResponse(res ResponseHelper) error {
	if len(res.Message) == 0 {
		res.Message = `success`
	}

	var resCode int
	if res.Code != 200 {
		resCode = http.StatusBadRequest
	} else {
		resCode = http.StatusOK
	}

	return res.C.JSON(resCode, map[string]interface{}{
		"code":         res.Code,
		"code_type":    res.CodeType,
		"code_message": res.Message,
		"data":         res.Data,
	})
}

// EmptyJsonMap ...
// just return empty array instead of null
func (u *HTTPHelper) EmptyJsonMap() map[string]interface{} {
	return make(map[string]interface{})
}

//get pagination URL
func (u *HTTPHelper) GetPagingUrl(c echo.Context, page, limit int) string {

	r := c.Request()
	currentURL := c.Scheme() + "://" + r.Host + r.URL.Path + "?page={page}&limit={limit}"

	defaultLinkReplacer := strings.NewReplacer("{page}", strconv.Itoa(page), "{limit}", strconv.Itoa(limit)).Replace(currentURL)

	return defaultLinkReplacer
}

//Set paginantion response
func (u *HTTPHelper) GeneratePaging(c echo.Context, prev, next, limit, page, totalRecord int) map[string]interface{} {

	prevURL, nextURL, firstURL, lastURL := "", "", "", ""
	paramPrevURL, paramNextURL, paramFirstURL, paramLastURL := "", "", "", ""

	totalPages := int(math.Ceil(float64(totalRecord) / float64(limit)))

	if page >= 1 {
		prev = page - 1
		if page < totalPages {
			next = page + 1
		} else {
			next = totalPages
		}
	}

	if totalPages >= page && page > 1 {
		prevURL = u.GetPagingUrl(c, prev, limit)
		paramPrevURL = "?page=" + strconv.Itoa(prev) + "&limit=" + strconv.Itoa(limit)
	}

	if totalPages > page {
		nextURL = u.GetPagingUrl(c, next, limit)
		paramNextURL = "?page=" + strconv.Itoa(next) + "&limit=" + strconv.Itoa(limit)
	}

	if totalPages >= page && page > 1 {
		firstURL = u.GetPagingUrl(c, 1, limit)
		paramFirstURL = "?page=1" + "&limit=" + strconv.Itoa(limit)
	}

	if totalPages >= page && totalPages != page {
		lastURL = u.GetPagingUrl(c, totalPages, limit)
		paramLastURL = "?page=" + strconv.Itoa(totalPages) + "&limit=" + strconv.Itoa(limit)
	}

	links := map[string]interface{}{
		"previous": prevURL,
		"next":     nextURL,
		"first":    firstURL,
		"last":     lastURL,
	}

	linkParameter := map[string]interface{}{
		"previous": paramPrevURL,
		"next":     paramNextURL,
		"first":    paramFirstURL,
		"last":     paramLastURL,
	}

	pagination := map[string]interface{}{
		"total_records":  totalRecord,
		"per_page":       limit,
		"current_page":   page,
		"total_pages":    totalPages,
		"links":          links,
		"link_parameter": linkParameter,
	}

	return pagination
}

//get pagination MPC Area
func (u *HTTPHelper) GetPagingMPCArea(c echo.Context, page, limit int) string {

	currentURL := "?page={page}&limit={limit}"

	defaultLinkReplacer := strings.NewReplacer("{page}", strconv.Itoa(page), "{limit}", strconv.Itoa(limit)).Replace(currentURL)

	return defaultLinkReplacer
}

//Set paginantion response for MPC
func (u *HTTPHelper) GeneratePagingMPC(c echo.Context, prev, next, limit, page, totalRecord int, cityOrigin, cityDestination, serviceID, timeID, sizePriceID, code, isForBackup string) map[string]interface{} {

	prevURL, nextURL, firstURL, lastURL := "", "", "", ""
	paramPrevURL, paramNextURL, paramFirstURL, paramLastURL := "", "", "", ""

	totalPages := int(math.Ceil(float64(totalRecord) / float64(limit)))

	if page >= 1 {
		prev = page - 1
		if page < totalPages {
			next = page + 1
		} else {
			next = totalPages
		}
	}

	if totalPages >= page && page > 1 {
		prevURL = u.GetPagingUrl(c, prev, limit)
		if prevURL != "" {
			prevURL = prevURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramPrevURL = u.GetPagingMPCArea(c, prev, limit)
		if paramPrevURL != "" {
			paramPrevURL = paramPrevURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	if totalPages > page {
		nextURL = u.GetPagingUrl(c, next, limit)
		if nextURL != "" {
			nextURL = nextURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramNextURL = u.GetPagingMPCArea(c, next, limit)
		if paramNextURL != "" {
			paramNextURL = paramNextURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	if totalPages >= page && page > 1 {
		firstURL = u.GetPagingUrl(c, 1, limit)
		if firstURL != "" {
			firstURL = firstURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramFirstURL = u.GetPagingMPCArea(c, 1, limit)
		if paramFirstURL != "" {
			paramFirstURL = paramFirstURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	if totalPages >= page && totalPages != page {
		lastURL = u.GetPagingUrl(c, totalPages, limit)
		if lastURL != "" {
			lastURL = lastURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramLastURL = u.GetPagingMPCArea(c, totalPages, limit)
		if paramLastURL != "" {
			paramLastURL = paramLastURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	currentParameter := "city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
		"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup

	links := map[string]interface{}{
		"previous": prevURL,
		"next":     nextURL,
		"first":    firstURL,
		"last":     lastURL,
	}

	linkParameter := map[string]interface{}{
		"previous": paramPrevURL,
		"next":     paramNextURL,
		"first":    paramFirstURL,
		"last":     paramLastURL,
	}

	pagination := map[string]interface{}{
		"total_records":     totalRecord,
		"per_page":          limit,
		"current_page":      page,
		"total_pages":       totalPages,
		"links":             links,
		"link_parameter":    linkParameter,
		"current_parameter": currentParameter,
	}

	return pagination
}

//Set paginantion response for MPC area
func (u *HTTPHelper) GeneratePagingMPCArea(c echo.Context, prev, next, limit, page, totalRecord int, cityOrigin, cityDestination, clusterOrigin, clusterDestination,
	areaOrigin, areaDestination, phOrigin, phDestination, ppOrigin, ppDestination, serviceID, timeID, sizePriceID, code, isForBackup string) map[string]interface{} {

	prevURL, nextURL, firstURL, lastURL := "", "", "", ""
	paramPrevURL, paramNextURL, paramFirstURL, paramLastURL := "", "", "", ""

	totalPages := int(math.Ceil(float64(totalRecord) / float64(limit)))

	if page >= 1 {
		prev = page - 1
		if page < totalPages {
			next = page + 1
		} else {
			next = totalPages
		}
	}

	if totalPages >= page && page > 1 {
		prevURL = u.GetPagingUrl(c, prev, limit)
		if prevURL != "" {
			prevURL = prevURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramPrevURL = u.GetPagingMPCArea(c, prev, limit)
		if paramPrevURL != "" {
			paramPrevURL = paramPrevURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	if totalPages > page {
		nextURL = u.GetPagingUrl(c, next, limit)
		if nextURL != "" {
			nextURL = nextURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramNextURL = u.GetPagingMPCArea(c, next, limit)
		if paramNextURL != "" {
			paramNextURL = paramNextURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	if totalPages >= page && page > 1 {
		firstURL = u.GetPagingUrl(c, 1, limit)
		if firstURL != "" {
			firstURL = firstURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramFirstURL = u.GetPagingMPCArea(c, 1, limit)
		if paramFirstURL != "" {
			paramFirstURL = paramFirstURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	if totalPages >= page && totalPages != page {
		lastURL = u.GetPagingUrl(c, totalPages, limit)
		if lastURL != "" {
			lastURL = lastURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramLastURL = u.GetPagingMPCArea(c, totalPages, limit)
		if paramLastURL != "" {
			paramLastURL = paramLastURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	currentParameter := "city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
		"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
		"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
		"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup

	links := map[string]interface{}{
		"previous": prevURL,
		"next":     nextURL,
		"first":    firstURL,
		"last":     lastURL,
	}

	linkParameter := map[string]interface{}{
		"previous": paramPrevURL,
		"next":     paramNextURL,
		"first":    paramFirstURL,
		"last":     paramLastURL,
	}

	pagination := map[string]interface{}{
		"total_records":     totalRecord,
		"per_page":          limit,
		"current_page":      page,
		"total_pages":       totalPages,
		"links":             links,
		"link_parameter":    linkParameter,
		"current_parameter": currentParameter,
	}

	return pagination
}

func (u *HTTPHelper) GeneratePagination(c echo.Context, prev, next, limit, page, totalRecord int) response.PaginationResponse {

	prevURL, nextURL, firstURL, lastURL := "", "", "", ""
	paramPrevURL, paramNextURL, paramFirstURL, paramLastURL := "", "", "", ""

	totalPages := int(math.Ceil(float64(totalRecord) / float64(limit)))

	if page >= 1 {
		prev = page - 1
		if page < totalPages {
			next = page + 1
		} else {
			next = totalPages
		}
	}

	if totalPages >= page && page > 1 {
		prevURL = u.GetPagingUrl(c, prev, limit)
		paramPrevURL = "?page=" + strconv.Itoa(prev) + "&limit=" + strconv.Itoa(limit)
	}

	if totalPages > page {
		nextURL = u.GetPagingUrl(c, next, limit)
		paramNextURL = "?page=" + strconv.Itoa(next) + "&limit=" + strconv.Itoa(limit)
	}

	if totalPages >= page && page > 1 {
		firstURL = u.GetPagingUrl(c, 1, limit)
		paramFirstURL = "?page=1" + "&limit=" + strconv.Itoa(limit)
	}

	if totalPages >= page && totalPages != page {
		lastURL = u.GetPagingUrl(c, totalPages, limit)
		paramLastURL = "?page=" + strconv.Itoa(totalPages) + "&limit=" + strconv.Itoa(limit)
	}

	links := response.LinksResponse{
		First:    firstURL,
		Previous: prevURL,
		Next:     nextURL,
		Last:     lastURL,
	}

	linkParameter := response.LinksResponse{
		First:    paramFirstURL,
		Previous: paramPrevURL,
		Next:     paramNextURL,
		Last:     paramLastURL,
	}

	pagination := response.PaginationResponse{
		CurrentPage:   page,
		PerPage:       limit,
		TotalPages:    totalPages,
		TotalRecords:  totalRecord,
		Links:         links,
		LinkParameter: linkParameter,
	}

	return pagination
}

func (u *HTTPHelper) GenerateMPCClusterPagination(c echo.Context, prev, next, limit, page, totalRecord int, cityOrigin, cityDestination, clusterOrigin, clusterDestination,
	areaOrigin, areaDestination, phOrigin, phDestination, ppOrigin, ppDestination, serviceID, timeID, sizePriceID, code, isForBackup string) response.PaginationResponse {

	prevURL, nextURL, firstURL, lastURL := "", "", "", ""
	paramPrevURL, paramNextURL, paramFirstURL, paramLastURL := "", "", "", ""

	totalPages := int(math.Ceil(float64(totalRecord) / float64(limit)))

	if page >= 1 {
		prev = page - 1
		if page < totalPages {
			next = page + 1
		} else {
			next = totalPages
		}
	}

	if totalPages >= page && page > 1 {
		prevURL = u.GetPagingUrl(c, prev, limit)
		if prevURL != "" {
			prevURL = prevURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramPrevURL = u.GetPagingMPCArea(c, prev, limit)
		if paramPrevURL != "" {
			paramPrevURL = paramPrevURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	if totalPages > page {
		nextURL = u.GetPagingUrl(c, next, limit)
		if nextURL != "" {
			nextURL = nextURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramNextURL = u.GetPagingMPCArea(c, next, limit)
		if paramNextURL != "" {
			paramNextURL = paramNextURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	if totalPages >= page && page > 1 {
		firstURL = u.GetPagingUrl(c, 1, limit)
		if firstURL != "" {
			firstURL = firstURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramFirstURL = u.GetPagingMPCArea(c, 1, limit)
		if paramFirstURL != "" {
			paramFirstURL = paramFirstURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	if totalPages >= page && totalPages != page {
		lastURL = u.GetPagingUrl(c, totalPages, limit)
		if lastURL != "" {
			lastURL = lastURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramLastURL = u.GetPagingMPCArea(c, totalPages, limit)
		if paramLastURL != "" {
			paramLastURL = paramLastURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
				"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
				"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	currentParameter := "city_origin=" + cityOrigin + "&city_destination=" + cityDestination +
		"&cluster_origin=" + clusterOrigin + "&cluster_destination=" + clusterDestination + "&area_origin=" + areaOrigin + "&area_destination=" + areaDestination +
		"&ph_origin=" + phOrigin + "&ph_destination=" + phDestination + "&pp_origin=" + ppOrigin + "&pp_destination=" + ppDestination + "&service_id=" + serviceID +
		"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup

	links := response.LinksResponse{
		First:    firstURL,
		Previous: prevURL,
		Next:     nextURL,
		Last:     lastURL,
	}

	linkParameter := response.LinksResponse{
		First:    paramFirstURL,
		Previous: paramPrevURL,
		Next:     paramNextURL,
		Last:     paramLastURL,
	}

	pagination := response.PaginationResponse{
		CurrentPage:      page,
		PerPage:          limit,
		TotalPages:       totalPages,
		TotalRecords:     totalRecord,
		Links:            links,
		LinkParameter:    linkParameter,
		CurrentParameter: currentParameter,
	}

	return pagination
}

func (u *HTTPHelper) GenerateMPCPagination(c echo.Context, prev, next, limit, page, totalRecord int, cityOrigin, cityDestination, serviceID, timeID, sizePriceID, code, isForBackup string) response.PaginationResponse {

	prevURL, nextURL, firstURL, lastURL := "", "", "", ""
	paramPrevURL, paramNextURL, paramFirstURL, paramLastURL := "", "", "", ""

	totalPages := int(math.Ceil(float64(totalRecord) / float64(limit)))

	if page >= 1 {
		prev = page - 1
		if page < totalPages {
			next = page + 1
		} else {
			next = totalPages
		}
	}

	if totalPages >= page && page > 1 {
		prevURL = u.GetPagingUrl(c, prev, limit)
		if prevURL != "" {
			prevURL = prevURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramPrevURL = u.GetPagingMPCArea(c, prev, limit)
		if paramPrevURL != "" {
			paramPrevURL = paramPrevURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	if totalPages > page {
		nextURL = u.GetPagingUrl(c, next, limit)
		if nextURL != "" {
			nextURL = nextURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramNextURL = u.GetPagingMPCArea(c, next, limit)
		if paramNextURL != "" {
			paramNextURL = paramNextURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	if totalPages >= page && page > 1 {
		firstURL = u.GetPagingUrl(c, 1, limit)
		if firstURL != "" {
			firstURL = firstURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramFirstURL = u.GetPagingMPCArea(c, 1, limit)
		if paramFirstURL != "" {
			paramFirstURL = paramFirstURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	if totalPages >= page && totalPages != page {
		lastURL = u.GetPagingUrl(c, totalPages, limit)
		if lastURL != "" {
			lastURL = lastURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}

		paramLastURL = u.GetPagingMPCArea(c, totalPages, limit)
		if paramLastURL != "" {
			paramLastURL = paramLastURL + "&city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
				"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup
		}
	}

	currentParameter := "city_origin=" + cityOrigin + "&city_destination=" + cityDestination + "&service_id=" + serviceID +
		"&time_id=" + timeID + "&size_price_id=" + sizePriceID + "&code=" + code + "&is_for_backup=" + isForBackup

	links := response.LinksResponse{
		First:    firstURL,
		Previous: prevURL,
		Next:     nextURL,
		Last:     lastURL,
	}

	linkParameter := response.LinksResponse{
		First:    paramFirstURL,
		Previous: paramPrevURL,
		Next:     paramNextURL,
		Last:     paramLastURL,
	}

	pagination := response.PaginationResponse{
		CurrentPage:      page,
		PerPage:          limit,
		TotalPages:       totalPages,
		TotalRecords:     totalRecord,
		Links:            links,
		LinkParameter:    linkParameter,
		CurrentParameter: currentParameter,
	}

	return pagination
}

func (u *HTTPHelper) SendResponseWithPagination(res ResponseHelper, pagination interface{}) error {
	if len(res.Message) == 0 {
		res.Message = `success`
	}

	var resCode int
	if res.Code != 200 {
		resCode = http.StatusBadRequest
	} else {
		resCode = http.StatusOK
	}

	return res.C.JSON(resCode, map[string]interface{}{
		"code":         res.Code,
		"code_type":    res.CodeType,
		"code_message": res.Message,
		"data":         res.Data,
		"pagination":   pagination,
	})
}

// SendSuccess with Pagination ...
// Send success response to consumers.
func (u *HTTPHelper) SendSuccessWithPagination(c echo.Context, message string, data, pagination interface{}) error {
	res := u.SetResponse(c, `ok`, message, data, codeSuccess, `success`)

	return u.SendResponseWithPagination(res, pagination)
}
