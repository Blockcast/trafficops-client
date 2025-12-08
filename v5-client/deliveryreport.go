package client

/*

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/Blockcast/trafficops-client/go-tc"
	"github.com/Blockcast/trafficops-client/toclientlib"
)

// apiDeliveryReports is the API version-relative path to the /deliveryreports API endpoint.
const apiDeliveryReports = "/deliveryreports"

// CreateDeliveryReport creates a delivery report from JSON data.
func (to *Session) CreateDeliveryReport(report interface{}, opts RequestOptions) (tc.DeliveryReportResponse, toclientlib.ReqInf, error) {
	var resp tc.DeliveryReportResponse
	var reqInf toclientlib.ReqInf

	// Ensure Content-Type is set for JSON
	if opts.Header == nil {
		opts.Header = make(http.Header)
	}
	if opts.Header.Get("Content-Type") == "" {
		opts.Header.Set("Content-Type", "application/json")
	}

	reqInf, err := to.post(apiDeliveryReports, opts, report, &resp)
	return resp, reqInf, err
}

// CreateDeliveryReportXML creates a delivery report from XML data.
func (to *Session) CreateDeliveryReportXML(xmlData string, opts RequestOptions) (tc.DeliveryReportResponse, toclientlib.ReqInf, error) {
	var resp tc.DeliveryReportResponse
	var reqInf toclientlib.ReqInf

	// Ensure Content-Type is set for XML
	if opts.Header == nil {
		opts.Header = make(http.Header)
	}
	contentType := opts.Header.Get("Content-Type")
	if contentType == "" {
		opts.Header.Set("Content-Type", "application/xml")
	}

	// Convert XML string to bytes and use raw request to avoid JSON marshaling
	xmlBytes := []byte(xmlData)

	// Build the full path with API base (like reqAPI middleware does)
	path := strings.TrimSuffix(to.TOClient.APIBase(), "/") + "/" + strings.TrimPrefix(apiDeliveryReports, "/")
	if len(opts.QueryParameters) > 0 {
		path += "?" + opts.QueryParameters.Encode()
	}

	// Use RawRequestWithHdr to send XML directly without JSON marshaling
	httpResp, _, err := to.TOClient.RawRequestWithHdr(http.MethodPost, path, xmlBytes, opts.Header)
	if err != nil {
		return resp, reqInf, err
	}
	defer httpResp.Body.Close()

	reqInf.StatusCode = httpResp.StatusCode
	reqInf.RespHeaders = httpResp.Header.Clone()

	// Read the response body
	body, readErr := io.ReadAll(httpResp.Body)
	if readErr != nil {
		return resp, reqInf, fmt.Errorf("failed to read response: %v", readErr)
	}

	// Handle error status codes
	if httpResp.StatusCode >= 400 {
		var alerts tc.Alerts
		if jsonErr := json.Unmarshal(body, &alerts); jsonErr == nil {
			errStr := alerts.ErrorString()
			if errStr != "" {
				err = fmt.Errorf("error requesting Traffic Ops: HTTP error %d %s - error-level alerts: %s",
					httpResp.StatusCode, httpResp.Status, errStr)
			} else {
				err = fmt.Errorf("error requesting Traffic Ops: HTTP error %d %s",
					httpResp.StatusCode, httpResp.Status)
			}
		} else {
			err = fmt.Errorf("error requesting Traffic Ops: HTTP error %d %s",
				httpResp.StatusCode, httpResp.Status)
		}
	}

	// Try to decode the response
	if decodeErr := json.Unmarshal(body, &resp); decodeErr != nil && err == nil {
		err = fmt.Errorf("failed to decode response: %v", decodeErr)
	}

	return resp, reqInf, err
}

// GetDeliveryReports returns delivery reports based on the provided query parameters.
func (to *Session) GetDeliveryReports(opts RequestOptions) (tc.DeliveryReportsResponse, toclientlib.ReqInf, error) {
	var data tc.DeliveryReportsResponse
	reqInf, err := to.get(apiDeliveryReports, opts, &data)
	return data, reqInf, err
}

// DeleteDeliveryReport deletes a delivery report based on the provided query parameters.
// Required parameters: deliveryMethodId, sessionStartTime, hostName
func (to *Session) DeleteDeliveryReport(opts RequestOptions) (tc.DeleteResponse, toclientlib.ReqInf, error) {
	var alerts tc.DeleteResponse

	// Validate required parameters
	if opts.QueryParameters == nil {
		return alerts, toclientlib.ReqInf{}, fmt.Errorf("missing required query parameters")
	}

	required := []string{"deliveryMethodId", "sessionStartTime", "hostName"}
	missing := []string{}

	for _, param := range required {
		if opts.QueryParameters.Get(param) == "" {
			missing = append(missing, param)
		}
	}

	if len(missing) > 0 {
		return alerts, toclientlib.ReqInf{}, fmt.Errorf("missing required parameters: %s", strings.Join(missing, ", "))
	}

	reqInf, err := to.del(apiDeliveryReports, opts, &alerts)
	return alerts, reqInf, err
}
