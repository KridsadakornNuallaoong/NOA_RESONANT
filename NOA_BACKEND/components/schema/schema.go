package schema

import "time"

type GyroData struct {
	DeviceID  string         `json:"deviceID"`
	UserID    string         `json:"userID"`
	DateTime  string         `json:"Datetime"`
	TimeStamp int64          `json:"TimeStamp"`
	Data      GyroDataDetail `json:"data"`
}

type GyroDataDetail struct {
	DeviceAddress string   `json:"DeviceAddress"`
	X             AxisData `json:"X"`
	Y             AxisData `json:"Y"`
	Z             AxisData `json:"Z"`
	Temperature   float64  `json:"Temperature"`
}

type AxisData struct {
	Acceleration          float64 `json:"Acceleration"`
	VelocityAngular       float64 `json:"VelocityAngular"`
	VibrationSpeed        float64 `json:"VibrationSpeed"`
	VibrationAngle        float64 `json:"VibrationAngle"`
	VibrationDisplacement float64 `json:"VibrationDisplacement"`
	Frequency             float64 `json:"Frequency"`
}

type GyroStruct struct {
	Acceleration          float32 `json:"Acceleration"`
	VelocityAngular       float32 `json:"VelocityAngular"`
	VibrationSpeed        float32 `json:"VibrationSpeed"`
	VibrationAngle        float32 `json:"VibrationAngle"`
	VibrationDisplacement float32 `json:"VibrationDisplacement"`
	Frequency             float32 `json:"Frequency"`
}

type MQTTData struct {
	DeviceAddress string     `json:"DeviceAddress"`
	Timestamp     int64      `json:"Timestamp"`
	X             GyroStruct `json:"X"`
	Y             GyroStruct `json:"Y"`
	Z             GyroStruct `json:"Z"`
	Temperature   float32    `json:"Temperature"`
}

type MQTTDataWithData struct {
	DeviceID string   `json:"deviceID"`
	UserID   string   `json:"userID"`
	Data     GyroData `json:"data"`
}

type PasswordRequest struct {
	Password string `json:"Password"`
	CFP      string `json:"CFP"`
}

type User struct {
	ID       string `bson:"userID"`   // User ID
	Username string `bson:"username"` // User name
	Email    string `bson:"email"`    // User email
	Password string `bson:"password"` // User password
}

type Account struct {
	ID  string `bson:"id,omitempty"` // User ID
	OTP string `bson:"otp"`          // OTP
}

type Device struct {
	ID       string `bson:"deviceID"` // Device ID
	Email    string `bson:"email"`    // User email
	Password string `bson:"password"` // User password
}

type GetDevice struct {
	UserID      string    `bson:"userID"`
	DeviceName  string    `bson:"deviceName"`
	DeviceID    string    `bson:"deviceID"`
	CreateDate  time.Time `bson:"createDate"`
	CurrentDate time.Time `bson:"currentDate"`
	Bookmark    bool      `bson:"bookmark"`
	Usage       int       `bson:"usage"`
	Status      bool      `bson:"status"`
}

type DataPayload struct {
	DataX []float32 `json:"dataX"`
	DataY []float32 `json:"dataY"`
	DataZ []float32 `json:"dataZ"`
}

// DataPayload represents the structure of the incoming data
type Data struct {
	DeviceID string `json:"deviceID"`
	UserID   string `json:"userID"`
	Data     struct {
		DeviceAddress string `json:"DeviceAddress"`
		X             struct {
			Acceleration          float64 `json:"Acceleration"`
			VelocityAngular       float64 `json:"VelocityAngular"`
			VibrationSpeed        float64 `json:"VibrationSpeed"`
			VibrationAngle        float64 `json:"VibrationAngle"`
			VibrationDisplacement float64 `json:"VibrationDisplacement"`
			Frequency             float64 `json:"Frequency"`
		} `json:"X"`
		Y struct {
			Acceleration          float64 `json:"Acceleration"`
			VelocityAngular       float64 `json:"VelocityAngular"`
			VibrationSpeed        float64 `json:"VibrationSpeed"`
			VibrationAngle        float64 `json:"VibrationAngle"`
			VibrationDisplacement float64 `json:"VibrationDisplacement"`
			Frequency             float64 `json:"Frequency"`
		} `json:"Y"`
		Z struct {
			Acceleration          float64 `json:"Acceleration"`
			VelocityAngular       float64 `json:"VelocityAngular"`
			VibrationSpeed        float64 `json:"VibrationSpeed"`
			VibrationAngle        float64 `json:"VibrationAngle"`
			VibrationDisplacement float64 `json:"VibrationDisplacement"`
			Frequency             float64 `json:"Frequency"`
		} `json:"Z"`
		Temperature float64 `json:"Temperature"`
	} `json:"data"`
}
