package webhook

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

)

type MutationAdmissionReviewResult struct {
	ApiVersion string           `json:"apiVersion"`
	Kind       string           `json:"kind"`
	Response   MutationResponse `json:"response"`
}

type MutationResponse struct {
	Allowed   bool   `json:"allowed"`
	Uid       string `json:"uid"`
	Patchtype string `json:"patchType"`
	Patch     string `json:"patch"`
}

type PatchOperation struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value,omitempty"`
}

type AdmissionRequest struct {
	Request struct {
		Uid    string
		Object struct {
			Metadata struct {
				Name      string
				Namespace string
			}
		}
	}
}

type ResponseMessage struct {
	Message string `json:"message"`
}


func Mutate_handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "Mutate server, Hello")
	} else if r.Method == http.MethodPost {

		bufbody := new(bytes.Buffer)
		bufbody.ReadFrom(r.Body)
		body := bufbody.String()

		var mutate_request AdmissionRequest
		if err := json.Unmarshal([]byte(body), &mutate_request); err != nil {
			fmt.Println(err)
			return
		}

		uid := mutate_request.Request.Uid
		is_valid := true

		var nodeselector map[string]string = map[string]string{"namespace": mutate_request.Request.Object.Metadata.Namespace}
		patchOperation, _ := json.Marshal(&PatchOperation{
			Op:    "add",
			Path:  "/spec/nodeSelector",
			Value: nodeselector,
		})

		str := "[" + string(patchOperation) + "]"
		patchOperation_base64 := base64.StdEncoding.EncodeToString([]byte(str))

		result := &MutationAdmissionReviewResult{
			ApiVersion: "admission.k8s.io/v1",
			Kind:       "AdmissionReview",
			Response: MutationResponse{
				Allowed:   is_valid,
				Uid:       uid,
				Patchtype: "JSONPatch",
				Patch:     patchOperation_base64,
			},
		}

		marshaled_result, _ := json.Marshal(result)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(marshaled_result)

		return
	}
}