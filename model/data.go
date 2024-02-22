package model

type Data struct {
	Unique   string   `json:"unique"`
	Info     Info     `json:"info"`
	Workflow Workflow `json:"workflow"`
}

type Info struct {
	Image struct {
		LowResolution  Resolution `json:"low_resolution"`
		HighResolution Resolution `json:"high_resolution"`
	} `json:"image"`
	Event struct {
		Event     string                 `json:"event"`
		EventType string                 `json:"eventType"`
		EventInfo map[string]interface{} `json:"eventInfo"`
		AccountID string                 `json:"accountid"`
		CameraID  string                 `json:"cameraid"`
		Timestamp string                 `json:"timestamp"`
		EventID   string                 `json:"eventid"`
	} `json:"event"`
	Settings struct {
		CustomAI struct {
			Enable bool `json:"enable"`
		} `json:"customAI"`
		CloudLPR       map[string]interface{} `json:"cloudLPR"`
		CloudInference bool                   `json:"cloudInference"`
	} `json:"settings"`
}

type Resolution struct {
	MimeType   string `json:"mimetype"`
	Resolution struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"resolution"`
	Framing struct {
		Axis struct {
			Pan  int `json:"pan"`
			Tilt int `json:"tilt"`
		} `json:"axis"`
		Fov struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"fov"`
	} `json:"framing"`
	URL string `json:"url"`
}

type Workflow struct {
	WorkFlowId    string            `json:"workFlowId"`
	CurrentWorker string            `json:"currentWorker"`
	CurrentPath   []interface{}     `json:"currentPath"`
	Workers       map[string]Worker `json:"workers"`
}

type Worker struct {
	Topic string `json:"topic"`
	Emits []struct {
		Filter map[string]interface{} `json:"filter"`
		Worker string                 `json:"worker"`
	} `json:"emits"`
}
