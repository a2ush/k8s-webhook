package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"log"
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

type ResponseMessage struct {
	Message string `json:"message"`
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

func validate_handler(w http.ResponseWriter, r *http.Request) {
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

func mutate_handler(w http.ResponseWriter, r *http.Request) {
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

func main() {
	log.Print("Server is running...")

	mux := http.NewServeMux()
	mux.HandleFunc("/validate", validate_handler)
	mux.HandleFunc("/mutate", mutate_handler)

	err := http.ListenAndServeTLS(":8080", "/tls/tls.crt", "/tls/tls.key", mux)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
