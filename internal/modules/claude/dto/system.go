package dto

import "encoding/json"

type AnthropicSystemContent struct {
    Text string
}

func (s *AnthropicSystemContent) UnmarshalJSON(data []byte) error {

    // string
    var str string
    if err := json.Unmarshal(data, &str); err == nil {
        s.Text = str
        return nil
    }

    // Claude Code format
    var arr []struct {
        Type string `json:"type"`
        Text string `json:"text"`
    }

    if err := json.Unmarshal(data, &arr); err == nil {

        for _, v := range arr {

            if v.Type == "text" {
                if s.Text != "" {
                    s.Text += "\n"
                }
                s.Text += v.Text
            }

        }

        return nil
    }

    return nil
}