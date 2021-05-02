package webhook

import (
	"bytes"
	"fmt"
	"net/http"
	"encoding/json"

)

type ValidationAdmissionReviewResult struct {
	ApiVersion string             `json:"apiVersion"`
	Kind       string             `json:"kind"`
	Response   ValidationResponse `json:"response"`
}

type ValidationResponse struct {
	Allowed bool            `json:"allowed"`
	Uid     string          `json:"uid"`
	Status  ResponseMessage `json:"status"`
}

func Validate_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "Validate server, Hello")

	} else if r.Method == http.MethodPost {

		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(r.Body)
		body := bufbody.String()

		var validate_request AdmissionRequest
		if err := json.Unmarshal([]byte(body), &validate_request); err != nil {
			fmt.Println(err)
			return
		}

		uid := validate_request.Request.Uid
		is_valid := true

		result := &ValidationAdmissionReviewResult{
			ApiVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
			Response: ValidationResponse{
				Allowed: is_valid,
				Uid:     uid,
				Status: ResponseMessage{
					Message: "result for valitate",
				},
			},
		}

		marshaled_result, _ := json.Marshal(result)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(marshaled_result)

		return
	}
}