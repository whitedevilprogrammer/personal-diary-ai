// package models

// import (
// 	"encoding/json"
// 	"net/http"
// )

// type Payload struct {
// 	Status  string `json:"status,omitempty"`
// 	Message string `json:"message,omitempty"`
// 	Data    any    `json:"data"`
// }

// type Response struct {
// 	Status  string `json:"status,omitempty"`
// 	Message string `json:"message,omitempty"`
// 	Data    any    `json:"data,omitempty"`
// }

// func NewPayload() *Payload {
// 	return &Payload{}
// }

// func NewResponse() *Response {
// 	return &Response{}
// }

// // DecodePayload decodes the request body into the provided data structure
// func (p *Payload) DecodePayload(r *http.Request, data any) error {

// 	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
// 		return err
// 	}

// 	p.Data = *&data

// 	return nil
// }

// func (r *Response) SuccessResponse(w http.ResponseWriter, message string) error {

// 	r.Status = "success"
// 	if message == "" {
// 		message = "success"
// 	}
// 	r.Message = message
// 	if err := json.NewEncoder(w).Encode(r); err != nil {
// 		return err
// 	}

// 	return nil
// }

// // ErrorResponse sends an error response with the given message and status code
// // It also encodes the response in JSON format
// func (r *Response) ErrorResponse(w http.ResponseWriter, message string) error {

// 	r.Status = "error"
// 	if message == "" {
// 		message = "error"
// 	}
// 	r.Message = message
// 	if err := json.NewEncoder(w).Encode(r); err != nil {
// 		return err
// 	}

// 	return nil
// }

// func (p *Payload) SetStatus(status string) {
// 	p.Status = status
// }
// func (p *Payload) SetMessage(message string) {
// 	p.Message = message
// }
// func (p *Payload) SetData(data any) {
// 	p.Data = data
// }
// func (p *Payload) GetStatus() string {
// 	return p.Status
// }

// func (p *Payload) GetMessage() string {
// 	return p.Message
// }
// func (p *Payload) GetData() any {
// 	return p.Data
// }

// func (p *Response) SetStatus(status string) {
// 	p.Status = status
// }
// func (p *Response) SetMessage(message string) {
// 	p.Message = message
// }
// func (p *Response) SetData(data any) {
// 	p.Data = data
// }
// func (p *Response) GetStatus() string {
// 	return p.Status
// }

// func (p *Response) GetMessage() string {
// 	return p.Message
// }
// func (p *Response) GetData() any {
// 	return p.Data
// }

package models

import (
	"encoding/json"
	"io"
	"net/http"
	"personal-diary/utils"
)

const IsEncrypted = true // Set to true if you want to encrypt the payload

type Payload struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data"`
}

type Response struct {
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
}

func NewPayload() *Payload {
	return &Payload{}
}

func NewResponse() *Response {
	return &Response{}
}

// Decrypt incoming request body before decoding into `data`
func (p *Payload) DecodePayload(r *http.Request, data any) error {

	if IsEncrypted {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		// Step 1: Unmarshal raw JSON string (removes outer quotes)
		var encryptedString string
		if err := json.Unmarshal(bodyBytes, &encryptedString); err != nil {
			return err
		}

		// Step 2: Decrypt the AES-encrypted string
		decryptedText, err := utils.DecryptAES(encryptedString)
		if err != nil {
			return err
		}
		// Step 3: Unmarshal the decrypted JSON into the provided struct
		if err := json.Unmarshal([]byte(decryptedText), &data); err != nil {
			return err
		}
	} else {
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			return err
		}
	}

	p.Data = data
	return nil
}

func (r *Response) SuccessResponse(w http.ResponseWriter, message string) error {
	r.Status = "success"
	if message == "" {
		message = "success"
	}
	r.Message = message

	if IsEncrypted {
		responseBytes, err := json.Marshal(r)
		if err != nil {
			return err
		}
		encrypted, err := utils.EncryptAES(string(responseBytes))
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(encrypted))
		return err
	} else {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(r); err != nil {
			return err
		}
	}

	return nil

}

func (r *Response) ErrorResponse(w http.ResponseWriter, message string) error {
	r.Status = "error"
	if message == "" {
		message = "error"
	}
	r.Message = message

	if IsEncrypted {
		responseBytes, err := json.Marshal(r)
		if err != nil {
			return err
		}
		encrypted, err := utils.EncryptAES(string(responseBytes))
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write([]byte(encrypted))
		return err
	} else {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(r); err != nil {
			return err
		}
	}

	return nil
}

func (p *Payload) SetStatus(status string) {
	p.Status = status
}
func (p *Payload) SetMessage(message string) {
	p.Message = message
}
func (p *Payload) SetData(data any) {
	p.Data = data
}
func (p *Payload) GetStatus() string {
	return p.Status
}

func (p *Payload) GetMessage() string {
	return p.Message
}
func (p *Payload) GetData() any {
	return p.Data
}

func (r *Response) SetStatus(status string) {
	r.Status = status
}
func (r *Response) SetMessage(message string) {
	r.Message = message
}
func (r *Response) SetData(data any) {
	r.Data = data
}
func (r *Response) GetStatus() string {
	return r.Status
}

func (r *Response) GetMessage() string {
	return r.Message
}
func (r *Response) GetData() any {
	return r.Data
}
