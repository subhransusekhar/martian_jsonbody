package jsonbody

import (
  "bytes"
  "encoding/json"
  "github.com/google/martian"
  martianparser "github.com/google/martian/parse"
  "io/ioutil"
  "net/http"
)

func init() {
  martianparser.Register("jsonbody.Modifier", modifierFromJSON)
}

type JsonBodyModifierPayload struct {
  Body  string                       `json:"body"`
  Scope []martianparser.ModifierType `json:"scope"`
}

// ModifyRequest modifies the query string of the request with the given key and value.
func (m *JsonBodyModifier) ModifyRequest(req *http.Request) error {
  newRequestBody := make(map[string]interface{})
  json.Unmarshal([]byte(m.Body), &newRequestBody)
  if len(newRequestBody) == 0 {
    return nil
  }

  buf, err := ioutil.ReadAll(req.Body)
  if err != nil {
    return err
  }
  request := ioutil.NopCloser(bytes.NewBuffer(buf))
  req.Body.Close()

  decoder := json.NewDecoder(request)
  var requestBody map[string]interface{}
  err = decoder.Decode(&requestBody)
  if len(requestBody) == 0 {
    requestBody = make(map[string]interface{})
  }

  for k, v := range newRequestBody {
    requestBody[k] = v
  }

  body, err := json.Marshal(requestBody)

  if err != nil {
    return err
  }
  req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
  req.ContentLength = int64(len(body))
  return nil
}


type JsonBodyModifier struct {
  Body string
}

func NewJsonBodyModifier(body string) martian.RequestModifier {
  return &JsonBodyModifier{
    Body: body,
  }
}

func modifierFromJSON(b []byte) (*martianparser.Result, error) {
  msg := &JsonBodyModifierPayload{}
  if err := json.Unmarshal(b, msg); err != nil {
    return nil, err
  }
  return martianparser.NewResult(NewJsonBodyModifier(msg.Body), msg.Scope)
}
